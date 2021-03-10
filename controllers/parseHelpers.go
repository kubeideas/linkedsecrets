package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
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

		return result, &InvalidCloudSecret{}
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
