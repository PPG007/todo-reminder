package model

import (
	"context"
	"github.com/qiniu/qmgo"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
)

const (
	C_CHINA_HOLIDAY = "chinaHoliday"
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
		"date": c.Date,
	}
	change := qmgo.Change{
		Upsert:    true,
		ReturnNew: true,
		Update: bsoncodec.M{
			"$set": bsoncodec.M{
				"isWorkingDay": c.IsWorkingDay,
			},
		},
	}
	return repository.Mongo.FindAndApply(ctx, C_CHINA_HOLIDAY, condition, change, c)
}
