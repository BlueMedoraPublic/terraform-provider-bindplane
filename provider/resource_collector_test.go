package provider

import (
	"reflect"
	"testing"
)

// This test is to ensure that const collectorTimeout is an int64
// set to 300. It should not be changed unless we have a good reason
// as this timeout is in the provider's documentation
func TestCollectorTimeout(t *testing.T) {
	x := reflect.TypeOf(collectorTimeout).Kind()
	if x != reflect.Int64 {
		t.Errorf("Expected constant 'collectorTimeout' to be an int64")
	}
	if collectorTimeout != 300 {
		t.Errorf("Expected constant 'collectorTimeout' to be '300'")
	}
}
