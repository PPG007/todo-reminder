package util

import "time"

var (
	weekdayMap = map[time.Weekday]int{
		time.Monday:    1,
		time.Tuesday:   2,
		time.Wednesday: 3,
		time.Thursday:  4,
		time.Friday:    5,
		time.Saturday:  6,
		time.Sunday:    7,
	}
)

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

func PtrValue[T any](value T) *T {
	v := new(T)
	v = &value
	return v
}

func GetNextWeekday(t time.Time, weekday int) time.Time {
	for weekdayMap[t.Weekday()] != weekday {
		t = t.AddDate(0, 0, 1)
	}
	return t
}
