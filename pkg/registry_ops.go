package registry_ops

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

const stable = "stable"

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
		err := updateRunningContainer(res, serviceName)
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
func updateRunningContainer(hookResponse HookResponse, serviceName string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()
	if err != nil {
		return err
	}
	out, err := cli.ImagePull(ctx, hookResponse.Repository.RepoName+fmt.Sprintf(":%v", stable), types.ImagePullOptions{})
	if err != nil {
		return err
	}
	err = restartContainer(serviceName, cli, ctx)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}
func restartContainer(containerName string, cli *client.Client, ctx context.Context) error {
	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return err
	}
	var container types.Container
	if len(containers) == 0 {
		return fmt.Errorf("no active containers")
	}
	for _, c := range containers {
		if strings.Contains(c.Names[0], containerName) {
			container = c
			break
		}
	}
	duration := 2 * time.Minute
	cli.ContainerRestart(ctx, container.ID, &duration)
	return nil
}
