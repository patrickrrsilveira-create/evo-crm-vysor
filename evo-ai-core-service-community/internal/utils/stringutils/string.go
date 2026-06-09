package stringutils

import (
	"encoding/json"
)

// Generic functions for JSON operations
func ToJSON[T any](v T) string {
	json, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(json)
}

func FromJSON[T any](s string, defaultValue T) T {
	var v T
	err := json.Unmarshal([]byte(s), &v)
	if err != nil {
		return defaultValue
	}
	return v
}

// JSONToStringMap converts a JSON string to map[string]string
func JSONToStringMap(s string) map[string]string {
	return FromJSON(s, map[string]string{})
}

// StringMapToJSON converts a map[string]string to JSON string
func StringMapToJSON(m map[string]string) string {
	return ToJSON(m)
}

// InterfaceMapSliceToJSON converts a slice of map[string]interface{} to JSON string
func InterfaceMapSliceToJSON(m []map[string]interface{}) string {
	return ToJSON(m)
}

// JSONToInterfaceMapSlice converts a JSON string to a slice of map[string]interface{}
func JSONToInterfaceMapSlice(s string) []map[string]interface{} {
	return FromJSON(s, []map[string]interface{}{})
}

// InterfaceMapToJSON converts a map[string]interface{} to JSON string
func InterfaceMapToJSON(m map[string]interface{}) string {
	return ToJSON(m)
}

// JSONToInterfaceMap converts a JSON string to map[string]interface{}
func JSONToInterfaceMap(s string) map[string]interface{} {
	return FromJSON(s, map[string]interface{}{})
}

// StructToString converts any struct to a JSON string
// This is kept for backward compatibility and is an alias for ToJSON
func StructToString[T any](s T) string {
	return ToJSON(s)
}

// JSONToStruct converts a JSON string
func JSONToStruct[T any](s string) []T {
	var t []T
	json.Unmarshal([]byte(s), &t)
	return t
}
