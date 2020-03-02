package trim

import (
	"testing"
)

func TestTrim(t *testing.T) {
	// Trim should remove the single space
	str := "{\"abc\": 20}"
	newStr, err := Trim(str)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if len(str)-len(newStr) != 1 {
		t.Errorf("Expected the length of trimmed string to be one characters less than origonal string")
	}

	// Trim should remove a single space and the new line
	str = "{\"abc\": 20\n}"
	newStr, err = Trim(str)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	if len(str)-len(newStr) != 2 {
		t.Errorf("Expected the length of trimmed string to be two characters less than origonal string")
	}
}

func TestTrimEmptyStr(t *testing.T) {
	s, err := Trim("")
	if err != nil {
		t.Errorf("Expected Trim() to return an empty string when given a string of length zery, got an error: " + err.Error())
	}

	if s != "" {
		t.Errorf("Expected Trim() to return an empty string when given an empty string, got: '" + s + "'")
	}
}
