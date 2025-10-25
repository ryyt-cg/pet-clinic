package test

import (
	"time"

	jsoniter "github.com/json-iterator/go"
)

func JsonString(obj interface{}) string {
	b, _ := jsoniter.Marshal(obj)
	return string(b)
}

// ToDate - convert date string to date
func ToDate(dateStr string) *time.Time {
	date, err := time.Parse(time.DateOnly, dateStr)
	if err != nil {
		return nil
	}
	return &date
}
