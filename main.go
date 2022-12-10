package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	custom_logger "moby-it/getgo/internal"
)

const port = 32041

func main() {
	file, err := custom_logger.SetLogger("/apps/getgo")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	http.HandleFunc("/", handleContainerPush)
	http.ListenAndServe(fmt.Sprintf(":%v", port), nil)
}
func handleContainerPush(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal("failed to read request body")
	}
	log.Println(string(b))
	defer r.Body.Close()
}
