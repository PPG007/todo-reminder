package model

import (
	"context"
	"github.com/qiniu/qmgo"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
)

const (
	C_CHINA_HOLIDAY = "C_CHINA_HOLIDAY"
)

var (
	CChinaHoliday = &ChinaHoliday{}
)

type ChinaHoliday struct {
	Id           bsoncodec.ObjectId `bson:"_id"`
	Date         string             `bson:"date"`
	IsWorkingDay bool               `bson:"isWorkingDay"`
}

func (c *ChinaHoliday) Create(ctx context.Context) error {
	condition := bsoncodec.M{
		"date":         c.Date,
		"isWorkingDay": c.IsWorkingDay,
	}
	change := qmgo.Change{
		Upsert: true,
	}
	return repository.Mongo.FindAndApply(ctx, C_CHINA_HOLIDAY, condition, change, c)
}
