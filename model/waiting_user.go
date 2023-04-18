package model

import (
	"context"
	"github.com/qiniu/qmgo"
	"time"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
)

const (
	C_WAITING_USER = "waitingUser"
)

var (
	CWaitingUser = &WaitingUser{}
)

type WaitingUser struct {
	Id        bsoncodec.ObjectId `bson:"_id"`
	Approved  bool               `bson:"approved"`
	Notified  bool               `bson:"notified"`
	UserId    string             `bson:"userId"`
	Flag      string             `bson:"flag"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func (*WaitingUser) Upsert(ctx context.Context, userId, flag string) (WaitingUser, error) {
	condition := bsoncodec.M{
		"userId": userId,
		"flag":   flag,
	}
	updater := bsoncodec.M{
		"$setOnInsert": bsoncodec.M{
			"createdAt": time.Now(),
			"notified":  false,
		},
	}
	change := qmgo.Change{
		Upsert:    true,
		ReturnNew: true,
		Update:    updater,
	}
	var result WaitingUser
	err := repository.Mongo.FindAndApply(ctx, C_WAITING_USER, condition, change, &result)
	return result, err
}

func (w *WaitingUser) MarkAsNotified(ctx context.Context) error {
	condition := bsoncodec.M{
		"_id": w.Id,
	}
	updater := bsoncodec.M{
		"notified": true,
	}
	return repository.Mongo.UpdateOne(ctx, C_WAITING_USER, condition, updater)
}

func (*WaitingUser) GetWaitingOne(ctx context.Context, userId string) (WaitingUser, error) {
	condition := bsoncodec.M{
		"userId":   userId,
		"approved": false,
	}
	result := WaitingUser{}
	err := repository.Mongo.FindOne(ctx, C_WAITING_USER, condition, &result)
	return result, err
}

func (*WaitingUser) HandleFriendAdded(ctx context.Context, userId string) error {
	condition := bsoncodec.M{
		"userId": userId,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"approved": true,
		},
	}
	_, err := repository.Mongo.UpdateAll(ctx, C_WAITING_USER, condition, updater)
	return err
}
