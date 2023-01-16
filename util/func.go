package util

import "time"

func StrInArray(str string, arr *[]string) bool {
	if arr == nil {
		return false
	}
	for _, s := range *arr {
		if str == s {
			return true
		}
	}
	return false
}

func GetStartTimeOfYear(argTime time.Time) time.Time {
	y, _, _ := argTime.Date()
	return time.Date(y, 1, 1, 0, 0, 0, 0, time.Local)
}

func GetEndTimeOfYear(argTime time.Time) time.Time {
	y, _, _ := argTime.Date()
	return time.Date(y+1, 1, 1, 23, 59, 59, 999999999, time.Local).AddDate(0, 0, -1)
}

func GetStartTimeOfMonth(argTime time.Time) time.Time {
	y, m, _ := argTime.Date()
	return time.Date(y, m, 1, 0, 0, 0, 0, time.Local)
}

func GetEndTimeOfMonth(argTime time.Time) time.Time {
	y, m, _ := argTime.Date()
	return time.Date(y, m+1, 0, 23, 59, 59, 999999999, time.Local)
}
