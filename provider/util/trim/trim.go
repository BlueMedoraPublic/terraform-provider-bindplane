package trim

import (
	"bytes"
	"encoding/json"
	"regexp"
)

// Trim strips all whitespace and new lines frmo a json string
func Trim(s string) (string, error) {
	s, err := removeWhiteSpace(s)
	if err != nil {
		return "", err
	}

	s = removeNewLine(s)

	return s, nil
}

func removeWhiteSpace(s string) (string, error) {
	buffer := new(bytes.Buffer)
	err := json.Compact(buffer, []byte(s))
	if err != nil {
		return "", err
	}
	return string(buffer.Bytes()), nil
}

func removeNewLine(s string) string {
	re := regexp.MustCompile(`\r?\n`)
	return re.ReplaceAllString(s, "")
}
