package docker_ops

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
)

type RegistryCreds struct {
	Username string
	Password string
}

// / Gets the docker registry credentials from environment. Returns an error if any of them are not found.
func getRegistryCredsFromEnv() (string, error) {
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
	if len(username) <= 0 {
		return "", fmt.Errorf("did not find docker username")
	}
	if len(password) <= 0 {
		return "", fmt.Errorf("did not find docker password")
	}
	authConfig := types.AuthConfig{
		Username: username,
		Password: password,
	}
	encodedJSON, err := json.Marshal(authConfig)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(encodedJSON), nil
}
