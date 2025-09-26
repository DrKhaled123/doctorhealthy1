package utils

import (
	"encoding/json"
)

// StringSliceToJSON converts string slice to JSON string
func StringSliceToJSON(slice []string) string {
	if len(slice) == 0 {
		return "[]"
	}
	jsonBytes, _ := json.Marshal(slice)
	return string(jsonBytes)
}

// JSONToStringSlice converts JSON string to string slice
func JSONToStringSlice(jsonStr string) []string {
	if jsonStr == "" || jsonStr == "[]" {
		return []string{}
	}
	var slice []string
	_ = json.Unmarshal([]byte(jsonStr), &slice)
	return slice
}
