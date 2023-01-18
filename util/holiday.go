package util

import (
	"fmt"
	"github.com/Lofanmi/chinese-calendar-golang/calendar"
	"time"
)

var (
	LunarHolidaysMap = map[string][]string{
		"12-30": {"除夕"},
		"01-01": {"春节"},
		"05-05": {"端午节"},
		"09-09": {"重阳节"},
		"07-07": {"七夕节"},
		"01-15": {"元宵节"},
		"08-15": {"中秋节"},
	}
)

// GetSolarHolidays 获取阳历节日，返回精确到天
func GetSolarHolidays(year int) map[string]time.Time {
	return map[string]time.Time{
		"元旦":         time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local),
		"情人节":        time.Date(year, time.February, 14, 0, 0, 0, 0, time.Local),
		"妇女节":        time.Date(year, time.March, 8, 0, 0, 0, 0, time.Local),
		"愚人节":        time.Date(year, time.April, 1, 0, 0, 0, 0, time.Local),
		"清明节":        time.Date(year, time.April, 5, 0, 0, 0, 0, time.Local),
		"劳动节":        time.Date(year, time.March, 1, 0, 0, 0, 0, time.Local),
		"青年节":        time.Date(year, time.March, 4, 0, 0, 0, 0, time.Local),
		"儿童节":        time.Date(year, time.June, 1, 0, 0, 0, 0, time.Local),
		"中国共产党成立纪念日": time.Date(year, time.July, 1, 0, 0, 0, 0, time.Local),
		"建军节":        time.Date(year, time.August, 1, 0, 0, 0, 0, time.Local),
		"教师节":        time.Date(year, time.September, 10, 0, 0, 0, 0, time.Local),
		"国庆节":        time.Date(year, time.October, 1, 0, 0, 0, 0, time.Local),
		"父亲节":        GetFatherDay(year),
		"母亲节":        GetMotherDay(year),
		"双十一":        time.Date(year, time.November, 11, 0, 0, 0, 0, time.Local),
		"双十二":        time.Date(year, time.December, 12, 0, 0, 0, 0, time.Local),
		"感恩节":        GetThanksgivingDay(year),
		"圣诞节":        time.Date(year, time.December, 25, 0, 0, 0, 0, time.Local),
		"万圣节":        time.Date(year, time.November, 1, 0, 0, 0, 0, time.Local),
	}
}

// GetFatherDay 获取指定年份的父亲节的日期，返回精确到天
func GetFatherDay(year int) time.Time {
	startTime := time.Date(year, time.June, 1, 0, 0, 0, 0, time.Local)
	for startTime.Weekday() != time.Sunday {
		startTime = startTime.AddDate(0, 0, 1)
	}
	// 到这里已经是六月的第一个周日了，直接加 14 天即可
	return startTime.AddDate(0, 0, 14)
}

// GetMotherDay 获取指定年份的母亲节的日期，返回精确到天
func GetMotherDay(year int) time.Time {
	startTime := time.Date(year, time.May, 1, 0, 0, 0, 0, time.Local)
	for startTime.Weekday() != time.Sunday {
		startTime = startTime.AddDate(0, 0, 1)
	}
	return startTime.AddDate(0, 0, 7)
}

// GetThanksgivingDay 获取指定年份的感恩节的日期，返回精确到天
func GetThanksgivingDay(year int) time.Time {
	startTime := time.Date(year, time.November, 1, 0, 0, 0, 0, time.Local)
	for startTime.Weekday() != time.Thursday {
		startTime = startTime.AddDate(0, 0, 1)
	}
	return startTime.AddDate(0, 0, 21)
}

// GetLunarHolidaysForSolarDate 获取公历日期这一天的农历节日
func GetLunarHolidaysForSolarDate(solar time.Time) []string {
	c := calendar.BySolar(int64(solar.Year()), int64(solar.Month()), int64(solar.Day()), 0, 0, 0)
	key := fmt.Sprintf("%02d-%02d", c.Lunar.GetMonth(), c.Lunar.GetDay())
	return LunarHolidaysMap[key]
}
