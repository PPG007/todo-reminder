package model

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	mgo_option "go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

const (
	C_CHINA_HOLIDAY = "chinaHoliday"
)

var (
	CChinaHoliday = &ChinaHoliday{}
)

func init() {
	repository.Mongo.CreateIndex(context.Background(), C_CHINA_HOLIDAY, options.IndexModel{
		Key: []string{"date"},
		IndexOptions: &mgo_option.IndexOptions{
			Background: util.PtrValue[bool](true),
			Unique:     util.PtrValue[bool](true),
		},
	})
}

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

func (*ChinaHoliday) IsTimeWorkingDay(ctx context.Context, t time.Time) (bool, error) {
	date := t.Format("20060102")
	condition := bsoncodec.M{
		"date": date,
	}
	holiday := ChinaHoliday{}
	err := repository.Mongo.FindOne(ctx, C_CHINA_HOLIDAY, condition, &holiday)
	if err != nil {
		return false, err
	}
	return holiday.IsWorkingDay, nil
}

func (*ChinaHoliday) GetNextWorkingDay(ctx context.Context, t time.Time) (time.Time, error) {
	date := t.Format("20060102")
	condition := bsoncodec.M{
		"date": bsoncodec.M{
			"$gte": date,
		},
		"isWorkingDay": true,
	}
	c := ChinaHoliday{}
	err := repository.Mongo.FindOneWithSorter(ctx, C_CHINA_HOLIDAY, []string{"date"}, condition, &c)
	if err != nil {
		return time.Time{}, err
	}
	return c.ParseDate(t)
}

func (*ChinaHoliday) GetNextHoliday(ctx context.Context, t time.Time) (time.Time, error) {
	date := t.Format("20060102")
	condition := bsoncodec.M{
		"date": bsoncodec.M{
			"$gte": date,
		},
		"isWorkingDay": false,
	}
	c := ChinaHoliday{}
	err := repository.Mongo.FindOneWithSorter(ctx, C_CHINA_HOLIDAY, []string{"date"}, condition, &c)
	if err != nil {
		return time.Time{}, err
	}
	return c.ParseDate(t)
}

func (c *ChinaHoliday) ParseDate(arg time.Time) (time.Time, error) {
	t, err := time.Parse("20060102", c.Date)
	if err != nil {
		return t, err
	}
	return time.Date(t.Year(), t.Month(), t.Day(), arg.Hour(), arg.Minute(), arg.Second(), arg.Nanosecond(), time.Local), nil
}
