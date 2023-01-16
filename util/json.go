package util

import (
	"encoding/json"
	"time"
)

func MarshalToJson(in interface{}) string {
	bytes, _ := json.Marshal(in)
	return string(bytes)
}

func UnmarshalFromJson[T any](text string) (*T, error) {
	result := new(T)
	err := json.Unmarshal([]byte(text), result)
	return result, err
}

func MustUnmarshalFromJson[T any](text string) T {
	result, _ := UnmarshalFromJson[T](text)
	return *result
}

func TransTimeStrToTime(str string) (time.Time, error) {
	return time.Parse(time.RFC3339, str)
}

func MustTransTimeStrToTime(str string) time.Time {
	t, _ := TransTimeStrToTime(str)
	return t
}

func TransTimeToRFC3339(t time.Time) string {
	return t.Format(time.RFC3339)
}
