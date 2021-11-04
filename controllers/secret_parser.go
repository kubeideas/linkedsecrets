package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"unicode/utf8"
)

// DELIMITER used to split data
const DELIMITER = "="

// parse plain data and split key/value using delimiter
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

				// trim leading and trailing whitespaces
				key := string(bytes.TrimSpace(splitDatakv[0]))
				//value := bytes.TrimSpace(splitDatakv[1])

				//check for base64 encoded file idenfication
				value, err := base64.StdEncoding.DecodeString(string(bytes.TrimSpace(splitDatakv[1])))
				if err != nil {
					// fallback to raw data
					value = bytes.TrimSpace(splitDatakv[1])
				} else {
					// Decoded value must be utf8
					if !utf8.Valid(value) {
						// fallback to raw data
						value = bytes.TrimSpace(splitDatakv[1])
					}
				}

				// add raw key and data to map
				result[key] = value

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

		return result, &InvalidCloudSecret{}
	}

	for k, v := range dat {
		// skip key without value and value without key
		if k != "" && v != "" {

			//add to map
			//result[k] = []byte(fmt.Sprintf("%v", v))

			//check for base64 encoded file content
			value, err := base64.StdEncoding.DecodeString(fmt.Sprintf("%v", v))
			if err != nil {
				// fallback to raw data
				value = []byte(fmt.Sprintf("%v", v))
			} else {
				// Decoded value must be utf8
				if !utf8.Valid(value) {
					// fallback to raw data
					value = []byte(fmt.Sprintf("%v", v))
				}
			}

			// add to map
			result[k] = value

		}
	}

	return result, nil
}
