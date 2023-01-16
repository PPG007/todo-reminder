package model

import (
	"context"
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
	Id              bsoncodec.ObjectId `json:"id" bson:"_id"`
	IsDeleted       bool               `json:"isDeleted" bson:"isDeleted"`
	CreatedAt       time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time          `json:"updatedAt" bson:"updatedAt"`
	HasBeenDone     bool               `json:"hasBeenDone" bson:"hasBeenDone"`
	HasBeenReminded bool               `json:"hasBeenReminded" bson:"hasBeenReminded"`
	NeedRemind      bool               `json:"needRemind" bson:"needRemind"`
	Content         string             `json:"content" bson:"content"`
	UserId          string             `json:"userId" bson:"userId"`
	RemindSetting   RemindSetting      `json:"remindSetting" bson:"remindSetting"`
	DoneAt          time.Time          `json:"doneAt,omitempty" bson:"doneAt,omitempty"`
}

func (t *Todo) Create(ctx context.Context) error {
	t.Id = bsoncodec.NewObjectId()
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	t.IsDeleted = false
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

func (*Todo) Done(ctx context.Context, id bsoncodec.ObjectId) error {
	condition := bsoncodec.M{
		"_id":         id,
		"isDeleted":   false,
		"hasBeenDone": false,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"hasBeenDone": true,
			"doneAt":      time.Now(),
			"updatedAt":   time.Now(),
		},
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
	return repository.Mongo.UpdateOne(ctx, C_TODO, condition, updater)
}

func (*Todo) ListByCondition(ctx context.Context, condition bsoncodec.M) ([]Todo, error) {
	var todos []Todo
	err := repository.Mongo.FindAll(ctx, C_TODO, condition, &todos)
	return todos, err
}

func (*Todo) RollbackTodo(ctx context.Context, id bsoncodec.ObjectId) error {
	condition := bsoncodec.M{
		"_id": id,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"updatedAt":       time.Now(),
			"hasBeenDone":     false,
			"hasBeenReminded": false,
		},
		"$unset": bsoncodec.M{
			"doneAt": "",
		},
	}
	return repository.Mongo.UpdateOne(ctx, C_TODO, condition, updater)
}

func (t *Todo) Notify(ctx context.Context) {
	err := t.UpdateById(ctx, t.Id, bsoncodec.M{
		"$set": bsoncodec.M{
			"hasBeenReminded": true,
		},
	})
	if err != nil {
		return
	}
	// TODO: notify user
}

func (t *Todo) RenewRemindAt(ctx context.Context) {
	if t.RemindSetting.RepeatSetting.Type == REPEAT_TYPE_NONE {
		return
	}
	nextRemindAt := t.RemindSetting.GetNextRemindAt()
	t.UpdateById(ctx, t.Id, bsoncodec.M{
		"$set": bsoncodec.M{
			"remindSetting.remindAt": nextRemindAt,
			"hasBeenReminded":        false,
		},
	})
}
