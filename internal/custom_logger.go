// /
package custom_logger

import (
	"fmt"
	"log"
	"os"
)

// / Currently configured to works exclusively with a simple logsPath/log file.
func SetLogger(logsPath string) (*os.File, error) {
	os.MkdirAll(logsPath, os.ModePerm)
	file, err := os.OpenFile(fmt.Sprintf("%v/log", logsPath), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	log.SetOutput(file)
	log.Println("Service Started. v1.3.0")
	return file, nil
}
