package utils

import "encoding/json"

func DecodeJSON[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return result, err
}

func EncodeJSON[T any](data T) ([]byte, error) {
	return json.Marshal(data)
}
