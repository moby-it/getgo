package registry_ops

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

const stable = "stable"

var cli *client.Client
var ctx context.Context

func HandleContainerPush(w http.ResponseWriter, r *http.Request) {
	var res HookResponse
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	serviceName := strings.Trim(strings.TrimPrefix(r.URL.Path, "/deploy/"), " ")
	if len(serviceName) <= 0 {
		errorMessage := "Invalid CD endpoint. Pushed a service with no name."
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(errorMessage))
		log.Println(errorMessage)
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
	if res.PushData.Tag == stable {
		log.Println("Tag pushed is stable. Starting build of", fmt.Sprintf("%v:%v", serviceName, stable))
		err := updateContainer(res, serviceName)
		if err != nil {
			message := "Failed to update container"
			log.Println(message, serviceName, err)
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(message))
			return
		} else {
			log.Println("Succesfully updated container", serviceName)
		}
	}
	w.WriteHeader(http.StatusOK)
	defer r.Body.Close()
}
func updateContainer(hookResponse HookResponse, serviceName string) error {
	ctx = context.Background()
	_cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	cli = _cli
	defer cli.Close()
	if err != nil {
		return err
	}
	out, err := cli.ImagePull(ctx, hookResponse.Repository.RepoName+fmt.Sprintf(":%v", stable), types.ImagePullOptions{})
	if err != nil {
		return err
	}
	container, err := destroyContainer(serviceName)
	if err != nil {
		return err
	}
	portMap, err := extractPortMapFromContainer(container)
	if err != nil {
		return err
	}
	err = recreateContainer(serviceName, hookResponse.Repository.RepoName, portMap)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}
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
func recreateContainer(containerName string, imageName string, portMap nat.PortMap) error {
	container, err := cli.ContainerCreate(ctx, &container.Config{Image: fmt.Sprintf("%v:stable", imageName)}, &container.HostConfig{PortBindings: portMap}, &network.NetworkingConfig{}, &v1.Platform{}, containerName)
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