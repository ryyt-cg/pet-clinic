package test

import (
	jsoniter "github.com/json-iterator/go"
	"time"
)

func JsonString(obj interface{}) string {
	b, _ := jsoniter.Marshal(obj)
	return string(b)
}

func ToDate(dateStr string) *time.Time {
	date, _ := time.Parse(time.DateOnly, dateStr)
	return &date
}
