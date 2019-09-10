package utils

import "encoding/json"

func ToJson(v interface{}) string {
	raw, _ := json.Marshal(v)
	return string(raw)
}
