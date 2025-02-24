package util

// HasMapKeys checks a map for a number of keys.
//
// Returns wether all keys could be found in the map.
// If not, it returns the first missing key.
func HasMapKeys(input map[interface{}]interface{}, keys []string) (bool, string) {
	for _, key := range keys {
		if !hasMapKey(input, key) {
			return false, key
		}
	}
	return true, ""
}

// hasMapKey checks a map for a single key.
//
// Returns wether the key could be found in the map.
func hasMapKey(input map[interface{}]interface{}, key string) bool {
	_, ok := input[key]
	return ok
}
