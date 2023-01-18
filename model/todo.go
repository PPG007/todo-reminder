package model

import (
	"context"
	"github.com/qiniu/qmgo"
	"time"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
)

const (
	C_TODO = "todo"
)

var (
	CTodo = &Todo{}
)

type Todo struct {
	Id            bsoncodec.ObjectId `json:"id" bson:"_id"`
	IsDeleted     bool               `json:"isDeleted" bson:"isDeleted"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
	NeedRemind    bool               `json:"needRemind" bson:"needRemind"`
	Content       string             `json:"content" bson:"content"`
	UserId        string             `json:"userId" bson:"userId"`
	RemindSetting RemindSetting      `json:"remindSetting" bson:"remindSetting"`
}

func (t *Todo) Create(ctx context.Context) error {
	t.Id = bsoncodec.NewObjectId()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.IsDeleted = false
	record := TodoRecord{
		HasBeenDone: false,
		TodoId:      t.Id,
		NeedRemind:  t.NeedRemind,
		Content:     t.Content,
	}
	err := record.Create(ctx)
	if err != nil {
		return err
	}
	return repository.Mongo.Insert(ctx, C_TODO, t)
}

func (*Todo) Get(ctx context.Context, condition bsoncodec.M) (Todo, error) {
	result := Todo{}
	err := repository.Mongo.FindOne(ctx, C_TODO, condition, &result)
	return result, err
}

func (*Todo) GetById(ctx context.Context, id bsoncodec.ObjectId) (Todo, error) {
	result := Todo{}
	err := repository.Mongo.FindOne(ctx, C_TODO, bsoncodec.M{"_id": id}, &result)
	return result, err
}

func (*Todo) UpdateById(ctx context.Context, id bsoncodec.ObjectId, updater bsoncodec.M) error {
	condition := bsoncodec.M{
		"_id": id,
	}
	return repository.Mongo.UpdateOne(ctx, C_TODO, condition, updater)
}

func (*Todo) DeleteById(ctx context.Context, id bsoncodec.ObjectId) error {
	condition := bsoncodec.M{
		"_id": id,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"isDeleted": true,
			"updatedAt": time.Now(),
		},
	}
	err := repository.Mongo.UpdateOne(ctx, C_TODO, condition, updater)
	if err != nil {
		return err
	}
	return CTodoRecord.DeleteByTodoId(ctx, id)
}

func (*Todo) ListByCondition(ctx context.Context, condition bsoncodec.M) ([]Todo, error) {
	var todos []Todo
	err := repository.Mongo.FindAll(ctx, C_TODO, condition, &todos)
	return todos, err
}

func (t *Todo) Upsert(ctx context.Context) error {
	condition := bsoncodec.M{
		"_id": t.Id,
	}
	change := qmgo.Change{
		Upsert:    true,
		ReturnNew: true,
		Update: bsoncodec.M{
			"$set": bsoncodec.M{
				"updatedAt":     time.Now(),
				"needRemind":    t.NeedRemind,
				"content":       t.Content,
				"userId":        t.UserId,
				"remindSetting": t.RemindSetting,
			},
			"$setOnInsert": bsoncodec.M{
				"isDeleted": false,
				"createdAt": time.Now(),
			},
		},
	}
	return repository.Mongo.FindAndApply(ctx, C_TODO, condition, change, t)
}

func (*Todo) GenNextRecord(ctx context.Context, id bsoncodec.ObjectId) error {
	return nil
}
