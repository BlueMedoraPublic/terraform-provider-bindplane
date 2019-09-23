package jsonutil

import (
	"encoding/json"
	"github.com/pkg/errors"
)

// JSONToInterface unmarshals a json string to an interface
func JSONToInterface(jsonStr string, i *interface{}) error {
	jsonByte := []byte(jsonStr)
	if err := json.Unmarshal(jsonByte, i); err != nil {
		return errors.Wrap(err, "JsonToInterface() Failed to convert json string to interface{}")
	}
	return nil
}

// InterfaceToJSONBytes returns a json []byte from an interface
func InterfaceToJSONBytes(i interface{}) ([]byte, error) {
	x, err := json.Marshal(i)
	if err != nil {
		return nil, errors.Wrap(err, "InterfaceToJson() failed to convert interface{} to json []byte")
	}
	return x, err
}
