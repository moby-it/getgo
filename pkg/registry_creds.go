package registry_ops

import (
	"encoding/base64"
	"encoding/json"
	"os"

	"github.com/docker/docker/api/types"
)

type RegistryCreds struct {
	Username string
	Password string
}

func getRegistryCredsFromEnv() (string, error) {
	username := os.Getenv("DOCKER_USERNAME")
	password := os.Getenv("DOCKER_PASSWORD")
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
