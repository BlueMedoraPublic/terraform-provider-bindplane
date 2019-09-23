package collector

import (
	"reflect"
	"testing"
)

func TestCollectorTimeout(t *testing.T) {
	x := reflect.TypeOf(collectorTimeout).Kind()
	if x != reflect.Int64 {
		t.Errorf("Expected constant 'collectorTimeout' to be an int64")
	}
	if collectorTimeout != 300 {
		t.Errorf("Expected constant 'collectorTimeout' to be '240'")
	}
}
