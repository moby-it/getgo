// / GetGo is a tool that aims to help you deploy your Dockerhub Repositories to your virtual machine. It is based on Dockerhub Webhooks.
package main

import (
	"fmt"
	"log"
	"net/http"

	custom_logger "moby-it/getgo/internal"
	docker_ops "moby-it/getgo/pkg"
)

const port = 32041

const logPath = "/apps/getgo"

// const logPath = "./logs"

func main() {

	file, err := custom_logger.SetLogger(logPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	http.HandleFunc("/deploy/", docker_ops.HandleContainerPush)
	err = http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
	if err != nil {
		log.Fatalln("App crash.", err)
	}
}
