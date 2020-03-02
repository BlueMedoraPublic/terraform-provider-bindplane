package compare

import (
    "testing"
)

func TestMapStringInterfaceDiffer(t *testing.T) {
    a := make(map[string]interface{})
    a["string"] = "string"
    a["int"] = 5

    b := make(map[string]interface{})
    b["string"] = "not"

    if MapStringInterface(a, b) {
        t.Errorf("expected MapStringInterface to return false when given different map[string]interface{}")
    }
}

func TestMapStringInterfaceSame(t *testing.T) {
    a := make(map[string]interface{})
    a["string"] = "string"
    a["int"] = 5
    a["float"] = 5.1
    a["string2"] = "abcd"

    b := a

    if MapStringInterface(a, b) != true {
        t.Errorf("expected MapStringInterface to return true when given identical map[string]interface{}")
    }
}
