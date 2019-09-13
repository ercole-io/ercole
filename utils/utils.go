package utils

import "encoding/json"

//ToJSON convert v to a string containing the equivalent json rappresentaion
func ToJSON(v interface{}) string {
	raw, _ := json.Marshal(v)
	return string(raw)
}

//Intptr return a point to the int passed in the argument
func Intptr(v int64) *int64 {
	return &v
}
