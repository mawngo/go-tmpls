package internal

import "encoding/json"

// toPrettyJSON encodes an item into a pretty (indented) JSON string.
func toPrettyJSON(v any) string {
	output, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return ""
	}
	return string(output)
}

// toJSON encodes an item into a JSON string.
func toJSON(v any) string {
	output, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(output)
}
