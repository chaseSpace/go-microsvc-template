package util

import (
	jsoniter "github.com/json-iterator/go"
)

func ToJson(v interface{}) []byte {
	b, _ := jsoniter.Marshal(v)
	return b
}

func ToJsonStr(v interface{}) string {
	b, _ := jsoniter.Marshal(v)
	return string(b)
}
