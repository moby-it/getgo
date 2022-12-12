// / Docker operations for handling deployments.
package docker_ops

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

var cli *client.Client
var ctx context.Context

func HandleContainerPush(w http.ResponseWriter, r *http.Request) {
	var res HookResponse
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("failed to read request body")
	}
	err = json.Unmarshal(b, &res)
	if err != nil {
		log.Fatalln(err)
	}
	ctx = r.Context()
	_cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		log.Fatalln(err)
	}
	cli = _cli
	defer cli.Close()
	containerName := fmt.Sprintf("%v-%v", res.Repository.Name, res.PushData.Tag)
	if !containerRunning(containerName) {
		message := fmt.Sprintf("Container with name %v not running. Please start the service before wiring up your webhooks.", containerName)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return
	}
	log.Println("Starting build of", containerName)
	err = updateContainer(res, containerName)
	if err != nil {
		message := "Failed to update container"
		log.Println(message, containerName, "Reason:", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(message))
		return
	} else {
		log.Println("Succesfully updated container", containerName)
		w.WriteHeader(http.StatusOK)
	}

}

// / If there is any container with a name that matches containerName, it destroys and recreates it with the same network configuration.
func updateContainer(hookResponse HookResponse, containerName string) error {
	creds, err := getRegistryCredsFromEnv()
	if err != nil {
		return err
	}
	image := fmt.Sprintf("%v:%v", hookResponse.Repository.RepoName, hookResponse.PushData.Tag)
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{
		RegistryAuth: creds,
	})
	io.Copy(os.Stdout, out)
	defer out.Close()

	if err != nil {
		return err
	}
	container, err := destroyContainer(containerName)
	if err != nil {
		return err
	}
	portMap, err := extractPortMapFromContainer(container)
	if err != nil {
		return err
	}
	err = recreateContainer(containerName, image, portMap)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Destroys a container with the given name. Throws an error if no active containers are found.
func destroyContainer(containerName string) (*types.Container, error) {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return nil, err
	}
	if len(containers) == 0 {
		return nil, fmt.Errorf("no active containers")
	}
	for _, c := range containers {
		if strings.Contains(c.Names[0], containerName) {
			stopTimeout := 2 * time.Minute
			cli.ContainerStop(ctx, c.ID, &stopTimeout)
			cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
			return &c, nil
		}
	}
	return nil, fmt.Errorf("failed to destroy container. No container found with name %v", containerName)

}

// / Given an image name, and a portMap, it recreates the container with the given image name, appending the stable tag.
func recreateContainer(containerName string, image string, portMap nat.PortMap) error {
	container, err := cli.ContainerCreate(ctx, &container.Config{Image: image}, &container.HostConfig{PortBindings: portMap}, &network.NetworkingConfig{}, &v1.Platform{}, containerName)
	if err != nil {
		return err
	}
	err = cli.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
	if err != nil {
		return err
	}
	return nil
}
func extractPortMapFromContainer(container *types.Container) (nat.PortMap, error) {
	portMap := make(nat.PortMap, 1)
	if len(container.Ports) <= 0 {
		return nil, fmt.Errorf("container has no open ports")
	}
	firstPort := container.Ports[0]
	publicPort := strconv.Itoa(int(firstPort.PublicPort))
	portMap[nat.Port(fmt.Sprintf("%v/%v", firstPort.PrivatePort, firstPort.Type))] = []nat.PortBinding{{HostIP: firstPort.IP, HostPort: publicPort}}
	return portMap, nil
}

// / Check if there is any container that matches the serviceName.
func containerRunning(serviceName string) bool {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		log.Fatalln(err)
		return false
	}
	for _, c := range containers {
		if strings.Contains(c.Names[0], serviceName) {
			return true
		}
	}
	return false
}
