package model

import (
	"context"
	"time"
)

const (
	REPEAT_TYPE_DAILY       = "daily"
	REPEAT_TYPE_WEEKLY      = "weekly"
	REPEAT_TYPE_MONTHLY     = "monthly"
	REPEAT_TYPE_YEARLY      = "yearly"
	REPEAT_TYPE_WORKING_DAY = "workingDay"
	REPEAT_TYPE_HOLIDAY     = "holiday"
)

type RemindSetting struct {
	RemindAt      time.Time     `json:"remindAt" bson:"remindAt"`
	LastRemindAt  time.Time     `json:"lastRemindAt" bson:"lastRemindAt"`
	IsRepeatable  bool          `json:"isRepeatable" bson:"isRepeatable"`
	RepeatSetting RepeatSetting `json:"repeatSetting" bson:"repeatSetting"`
}

type RepeatSetting struct {
	Type       string `json:"type" bson:"type"`
	DateOffset int    `json:"dateOffset" bson:"dateOffset"`
}

type dateOffset struct {
	year  int
	month int
	week  int
	day   int
}

func (r *RemindSetting) GetNextRemindAt(ctx context.Context) time.Time {
	if !r.IsRepeatable {
		return r.RemindAt
	}
	now := time.Now()
	var (
		offset     dateOffset
		targetDate time.Time
	)
	switch r.RepeatSetting.Type {
	case REPEAT_TYPE_DAILY:
		offset.day = r.RepeatSetting.DateOffset
	case REPEAT_TYPE_WEEKLY:
		offset.week = r.RepeatSetting.DateOffset
	case REPEAT_TYPE_MONTHLY:
		offset.month = r.RepeatSetting.DateOffset
	case REPEAT_TYPE_YEARLY:
		offset.year = r.RepeatSetting.DateOffset
	}
	if r.RepeatSetting.Type == REPEAT_TYPE_HOLIDAY {
		temp := r.LastRemindAt
		// 第一次提醒
		if r.LastRemindAt.Unix() < 0 {
			temp = r.RemindAt
			if now.After(r.RemindAt) {
				temp = r.RemindAt.AddDate(0, 0, 1)
			}
		} else {
			temp = temp.AddDate(0, 0, 1)
		}
		nextRemindAt, err := CChinaHoliday.GetNextHoliday(ctx, temp)
		if err != nil {
			return time.Time{}
		}
		r.LastRemindAt = nextRemindAt
		targetDate = nextRemindAt
	} else if r.RepeatSetting.Type == REPEAT_TYPE_WORKING_DAY {
		temp := r.LastRemindAt
		// 第一次提醒
		if r.LastRemindAt.Unix() < 0 {
			temp = r.RemindAt
			if now.After(r.RemindAt) {
				temp = r.RemindAt.AddDate(0, 0, 1)
			}
		} else {
			temp = temp.AddDate(0, 0, 1)
		}
		nextRemindAt, err := CChinaHoliday.GetNextWorkingDay(ctx, temp)
		if err != nil {
			return time.Time{}
		}
		r.LastRemindAt = nextRemindAt
		targetDate = nextRemindAt
	} else {
		targetDate := r.LastRemindAt
		// 从未提醒过
		if r.LastRemindAt.Unix() < 0 {
			targetDate = r.RemindAt
			if now.After(r.RemindAt) {
				targetDate = r.RemindAt.AddDate(offset.year, offset.month, offset.week*7+offset.day)
			}
		} else {
			targetDate = r.LastRemindAt.AddDate(offset.year, offset.month, offset.week*7+offset.day)
		}
		r.LastRemindAt = targetDate
		return targetDate
	}
	return targetDate
}
