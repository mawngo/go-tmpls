package internal

import "encoding/json"

// toPrettyJson encodes an item into a pretty (indented) JSON string.
func toPrettyJson(v interface{}) string {
	output, _ := json.MarshalIndent(v, "", "  ")
	return string(output)
}

// toJson encodes an item into a JSON string.
func toJson(v any) string {
	output, _ := json.Marshal(v)
	return string(output)
}
