package controllers

import (
	"encoding/base64"
	"encoding/json"
)

type authConfig struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email,omitempty"`
	Auth     string `json:"auth,omitempty"`
}

type dockerAuths struct {
	AuthConfigs map[string]authConfig `json:"auths"`
}

func inferDockerConfig(secretMap map[string][]byte) (map[string][]byte, bool) {

	// Infer docker required fields
	if checkFields(secretMap) {
		if dockerSecret, err := convertToDockerConfig(secretMap); err == nil {
			return dockerSecret, true
		}
	}

	return nil, false
}

// Convert secret map to docker config format
func convertToDockerConfig(secretMap map[string][]byte) (dockerSecret map[string][]byte, err error) {

	// encode username and password
	authPlain := string(secretMap["docker-username"]) + ":" + string(secretMap["docker-password"])
	authEnc := base64.URLEncoding.EncodeToString([]byte(authPlain))

	dockerAuths := dockerAuths{
		AuthConfigs: map[string]authConfig{
			string(secretMap["docker-server"]): {
				Username: string(secretMap["docker-username"]),
				Password: string(secretMap["docker-password"]),
				Email:    string(secretMap["docker-email"]),
				Auth:     authEnc,
			}},
	}

	// encode JSON
	encDockerAuths, err := json.Marshal(dockerAuths)
	if err != nil {
		return nil, err
	}

	dockerSecret = map[string][]byte{".dockerconfigjson": encDockerAuths}
	return dockerSecret, nil
}

// evaluated cloud secrets fields
func checkFields(secretMap map[string][]byte) bool {

	// check number of fields
	if len(secretMap) < 2 || len(secretMap) > 4 {
		return false
	}

	// check required required fields to be considered docker config
	dockerReqFields := []string{"docker-username", "docker-password"}
	for _, v := range dockerReqFields {
		if _, ok := secretMap[v]; !ok {
			return false
		}
	}

	// set default registry if necessary
	if _, ok := secretMap["docker-server"]; !ok {
		secretMap["docker-server"] = []byte("https://index.docker.io/v1/")
	}

	return true
}
