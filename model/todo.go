package model

import (
	"context"
	"errors"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	mgo_option "go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

const (
	C_TODO = "todo"
)

var (
	CTodo = &Todo{}
)

func init() {
	repository.Mongo.CreateIndex(context.Background(), C_TODO, options.IndexModel{
		Key: []string{"isDeleted", "userId"},
		IndexOptions: &mgo_option.IndexOptions{
			Background: util.PtrValue[bool](true),
		},
	})
	repository.Mongo.CreateIndex(context.Background(), C_TODO, options.IndexModel{
		Key: []string{"isDeleted", "needRemind", "remindSetting.isRepeatable"},
		IndexOptions: &mgo_option.IndexOptions{
			Background: util.PtrValue[bool](true),
		},
	})
}

type Todo struct {
	Id            bsoncodec.ObjectId `json:"id" bson:"_id"`
	IsDeleted     bool               `json:"isDeleted" bson:"isDeleted"`
	CreatedAt     time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time          `json:"updatedAt" bson:"updatedAt"`
	NeedRemind    bool               `json:"needRemind" bson:"needRemind"`
	Content       string             `json:"content" bson:"content"`
	UserId        string             `json:"userId" bson:"userId"`
	RemindSetting RemindSetting      `json:"remindSetting" bson:"remindSetting"`
	Images        []string           `json:"images" bson:"images,omitempty"`
}

func (t *Todo) Create(ctx context.Context) error {
	t.Id = bsoncodec.NewObjectId()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.IsDeleted = false
	err := repository.Mongo.Insert(ctx, C_TODO, t)
	if err != nil {
		return err
	}
	return t.GenNextRecord(ctx, t.Id, true)
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
				"images":        t.Images,
			},
			"$setOnInsert": bsoncodec.M{
				"isDeleted": false,
				"createdAt": time.Now(),
			},
		},
	}

	err := repository.Mongo.FindAndApply(ctx, C_TODO, condition, change, t)
	if err != nil {
		return err
	}
	go func() {
		CTodoRecord.DeleteUndoneOnesByTodoId(ctx, t.Id)
		t.GenNextRecord(ctx, t.Id, true)
	}()
	return nil
}

func (*Todo) GenNextRecord(ctx context.Context, id bsoncodec.ObjectId, isFirst bool) error {
	condition := bsoncodec.M{
		"_id": id,
	}
	t := Todo{}
	err := repository.Mongo.FindOne(ctx, C_TODO, condition, &t)
	if err != nil {
		return err
	}
	if (!t.NeedRemind || !t.RemindSetting.IsRepeatable) && !isFirst {
		return nil
	}
	if total, err := CTodoRecord.CountNotDoneRecordsByTodoId(ctx, id); err == nil {
		if total >= 1 {
			return nil
		}
	} else {
		return err
	}
	r := TodoRecord{
		UserId:     t.UserId,
		NeedRemind: t.NeedRemind,
		Content:    t.Content,
		TodoId:     t.Id,
		Images:     t.Images,
	}
	if t.NeedRemind && t.RemindSetting.IsRepeatable {
		r.IsRepeatable = true
		r.RepeatType = t.RemindSetting.RepeatSetting.Type
		r.RepeatDateOffset = t.RemindSetting.RepeatSetting.DateOffset
	}
	if t.NeedRemind {
		remindAt := t.RemindSetting.GetNextRemindAt(ctx)
		if remindAt.Unix() < 0 {
			return errors.New("failed to gen remind at")
		}
		updater := bsoncodec.M{
			"$set": bsoncodec.M{
				"remindSetting": t.RemindSetting,
			},
		}
		err = t.UpdateById(ctx, t.Id, updater)
		if err != nil {
			return err
		}
		r.RemindAt = remindAt
	}
	return r.Create(ctx)
}

func (*Todo) ListByIds(ctx context.Context, ids []bsoncodec.ObjectId) ([]Todo, error) {
	condition := bsoncodec.M{
		"_id": bsoncodec.M{
			"$in": ids,
		},
	}
	var todos []Todo
	err := repository.Mongo.FindAll(ctx, C_TODO, condition, &todos)
	return todos, err
}
