package model

import (
	"time"
)

const (
	REPEAT_TYPE_DAILY       = "daily"
	REPEAT_TYPE_WEEKLY      = "weekly"
	REPEAT_TYPE_MONTHLY     = "monthly"
	REPEAT_TYPE_YEARLY      = "yearly"
	REPEAT_TYPE_WORKING_DAY = "workingDay"
	REPEAT_TYPE_NONE        = "none"
)

type RemindSetting struct {
	RemindAt      time.Time     `json:"remindAt" bson:"remindAt"`
	RepeatSetting RepeatSetting `json:"repeatSetting" bson:"repeatSetting"`
}

type RepeatSetting struct {
	Type    string `json:"type" bson:"type"`
	Month   int    `json:"month" bson:"month"`
	Weekday int    `json:"weekday" bson:"weekday"`
	Day     int    `json:"day" bson:"day"`
}

func (r *RemindSetting) GetNextRemindAt() time.Time {
	switch r.RepeatSetting.Type {
	case REPEAT_TYPE_DAILY:
		return r.RemindAt.AddDate(0, 0, 1)
	case REPEAT_TYPE_WEEKLY:
		return r.RemindAt.AddDate(0, 0, 7)
	case REPEAT_TYPE_MONTHLY:
		return r.RemindAt.AddDate(0, 1, 0)
	case REPEAT_TYPE_YEARLY:
		return r.RemindAt.AddDate(1, 0, 0)
	case REPEAT_TYPE_WORKING_DAY:
		// TODO:
	}
	return time.Time{}
}
