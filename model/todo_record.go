package model

import (
	"context"
	"time"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
)

const (
	C_TODO_RECORD = "todoRecord"
)

var (
	CTodoRecord = &TodoRecord{}
)

type TodoRecord struct {
	Id          bsoncodec.ObjectId `bson:"_id"`
	IsDeleted   bool               `bson:"isDeleted"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
	RemindAt    time.Time          `bson:"remindAt"`
	HasBeenDone bool               `bson:"hasBeenDone"`
	Content     string             `bson:"content"`
	TodoId      bsoncodec.ObjectId `bson:"todoId"`
	DoneAt      time.Time          `bson:"doneAt"`
	NeedRemind  bool               `bson:"needRemind"`
}

func (t *TodoRecord) Create(ctx context.Context) error {
	t.Id = bsoncodec.NewObjectId()
	t.IsDeleted = false
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return repository.Mongo.Insert(ctx, C_TODO_RECORD, t)
}

func (*TodoRecord) DeleteByTodoId(ctx context.Context, todoId bsoncodec.ObjectId) error {
	condition := bsoncodec.M{
		"isDeleted": false,
		"todoId":    todoId,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"isDeleted": true,
		},
	}
	_, err := repository.Mongo.UpdateAll(ctx, C_TODO_RECORD, condition, updater)
	return err
}
