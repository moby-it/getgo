package custom_logger

import (
	"fmt"
	"log"
	"os"
)

func SetLogger(logsPath string) (*os.File, error) {
	os.MkdirAll(logsPath, os.ModePerm)
	file, err := os.OpenFile(fmt.Sprintf("%v/log", logsPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	log.Println("Service Started.")
	return file, nil
}
