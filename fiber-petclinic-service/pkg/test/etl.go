package test

import jsoniter "github.com/json-iterator/go"

func JsonString(obj interface{}) string {
	b, _ := jsoniter.Marshal(obj)
	return string(b)
}
