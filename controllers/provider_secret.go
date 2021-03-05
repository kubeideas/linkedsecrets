package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	securityv1 "linkedsecrets/api/v1"

	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DELIMITER used to split data
const DELIMITER = "="

// parse plain data into key an value using delimiter
func parsePlainData(data []byte) (map[string][]byte, error) {

	// key/value map
	result := make(map[string][]byte)

	splitDatalines := bytes.Split(data, []byte("\n"))

	for _, line := range splitDatalines {
		// skip empty lines
		if bytes.TrimSpace(line) != nil {

			//split by delimiter
			splitDatakv := bytes.SplitN(line, []byte(DELIMITER), 2)

			// skip key without delimiter, key without value or value without key
			if len(splitDatakv) == 2 && bytes.TrimSpace(splitDatakv[0]) != nil && bytes.TrimSpace(splitDatakv[1]) != nil {

				//trim leading and trailing whitespaces
				key := bytes.TrimSpace(splitDatakv[0])
				value := bytes.TrimSpace(splitDatakv[1])

				// add to map
				result[string(key)] = value
			}

		}
	}

	//check for invalid data
	if len(result) == 0 {
		return result, &InvalidCloudSecret{}
	}

	return result, nil
}

// parse json format
func parseJSON(data []byte) (map[string][]byte, error) {
	// key/value map
	result := make(map[string][]byte)

	var dat map[string]interface{}

	// parse json
	if err := json.Unmarshal(data, &dat); err != nil {
		return result, err
	}

	for k, v := range dat {
		// skip key without value and value without key
		if k != "" && v != "" {
			//add to map
			result[k] = []byte(fmt.Sprintf("%v", v))
		}
	}

	return result, nil
}

//check provider and return filled secret object
func getProviderSecret(linkedsecret *securityv1.LinkedSecret) (corev1.Secret, error) {

	var err error
	data := []byte{}
	secret := corev1.Secret{}
	secretMap := make(map[string][]byte)

	if linkedsecret.Spec.Provider == GOOGLE {

		// ######## begin fake data ########
		// if linkedsecret.Name == "google-secret-one" {
		// 	data = []byte("password=senha1")
		// } else {
		// 	data = []byte("password=senha2")
		// }
		// err = nil
		// ######## end fake data ########

		//Restore code after all linkedsecret crud is finished.
		data, err = GetGCPSecret(linkedsecret)

		if err != nil {
			return corev1.Secret{}, err
		}
	}

	// create key/value map based on choosen format
	if linkedsecret.Spec.ProviderDataFormat == JSONFORMAT {
		secretMap, err = parseJSON(data)
	} else {
		secretMap, err = parsePlainData(data)
	}

	if err != nil {
		return secret, err
	}

	// create new secret object and add data
	secret = corev1.Secret{
		TypeMeta: v1.TypeMeta{
			Kind:       "Secret",
			APIVersion: corev1.SchemeGroupVersion.String(),
		},
		ObjectMeta: v1.ObjectMeta{
			Name:      linkedsecret.Spec.SecretName,
			Namespace: linkedsecret.Namespace,
		},
		Data: secretMap,
		Type: "Opaque",
	}

	return secret, nil

}
