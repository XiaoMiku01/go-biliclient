package utils

import "encoding/json"

func AnyToJSON(a any) string {
	s, _ := json.Marshal(a)
	return string(s)
}
