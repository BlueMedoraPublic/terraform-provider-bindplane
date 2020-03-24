package compare

// MapStringInterface returns true if each map[string]interface{}
// has identical key value pairs
func MapStringInterface(a, b map[string]interface{}) bool {
	if len(a) != len(b) {
		return false
	}
	return compareKeyValues(a, b)
}

func compareKeyValues(a, b map[string]interface{}) bool {
	for aKey, aValue := range a {
		if bValue, ok := b[aKey]; ok {
			if aValue != bValue {
				return false
			}
		}
	}
	return true
}
