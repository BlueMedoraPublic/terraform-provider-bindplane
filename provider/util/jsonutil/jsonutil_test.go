package jsonutil

import (
	"encoding/json"
	"testing"
)

// TestJSONToInterface takes a json string, converts it to
// an interface and then converts it back to a string and compares
// the two strings. If they are different, something went wrong
func TestJSONToInterface(t *testing.T) {
	jsonStr := "{\"abc\":\"xyz\",\"num\":5}"

	var i interface{}
	if err := JSONToInterface(jsonStr, &i); err != nil {
		t.Errorf("Expected JSONToInterface to return nil, got error: " + err.Error())
		return
	}

	jsonBytes, err := json.Marshal(i)
	if err != nil {
		t.Errorf("Expected json.Marshal to return nil, when marshalling the interface returned from JSONToInterface")
		return
	}

	if string(jsonBytes) != jsonStr {
		t.Errorf("Expected the result of JSONToInterface to marshal into a []byte and then convert to a string, and be identical to the origonal string. Origonal: " + jsonStr + ", result: " + string(jsonBytes))
	}
}

// TestInterfaceToJSONBytes
func TestInterfaceToJSONBytes(t *testing.T) {
	jsonStr := "{\"abc\":\"xyz\",\"num\":5}"

	var i interface{}
	if err := JSONToInterface(jsonStr, &i); err != nil {
		t.Errorf("Expected JSONToInterface to return nil, got error: " + err.Error())
		return
	}

	jsonBytes, err := InterfaceToJSONBytes(i)
	if err != nil {
		t.Errorf("Expected InterfaceToJSONBytes to return a nil error, got: " + err.Error())
		return
	}

	if string(jsonBytes) != jsonStr {
		t.Errorf("Expected the result of InterfaceToJSONBytes to be identical to the origonal string. Origonal: " + jsonStr + ", result: " + string(jsonBytes))
	}
}
