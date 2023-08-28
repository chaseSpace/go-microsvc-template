package util

import "encoding/json"

func ToJson(v interface{}) []byte {
	b, _ := json.Marshal(v)
	return b
}

func ToJsonStr(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
