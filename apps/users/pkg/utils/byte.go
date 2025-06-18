package utils

import "encoding/json"

func ByteToMap(data []byte) (*map[string]any, error) {
	var result *map[string]any
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func MapToByte(data map[string]any) ([]byte, error) {
	return json.Marshal(data)
}
