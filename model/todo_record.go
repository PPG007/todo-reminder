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
	C_TODO_RECORD = "todoRecord"
)

var (
	CTodoRecord = &TodoRecord{}
)

func init() {
	repository.Mongo.CreateIndex(context.Background(), C_TODO_RECORD, options.IndexModel{
		Key: []string{"isDeleted", "todoId", "hasBeenDone"},
		IndexOptions: &mgo_option.IndexOptions{
			Background: util.PtrValue[bool](true),
		},
	})
	repository.Mongo.CreateIndex(context.Background(), C_TODO_RECORD, options.IndexModel{
		Key: []string{"isDeleted", "hasBeenDone", "userId", "remindAt"},
		IndexOptions: &mgo_option.IndexOptions{
			Background: util.PtrValue[bool](true),
		},
	})
	repository.Mongo.CreateIndex(context.Background(), C_TODO_RECORD, options.IndexModel{
		Key: []string{"isDeleted", "remindAt", "needRemind", "hasBeenDone", "hasBeenReminded"},
		IndexOptions: &mgo_option.IndexOptions{
			Background: util.PtrValue[bool](true),
		},
	})
}

type TodoRecord struct {
	Id               bsoncodec.ObjectId `bson:"_id"`
	IsDeleted        bool               `bson:"isDeleted"`
	CreatedAt        time.Time          `bson:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt"`
	RemindAt         time.Time          `bson:"remindAt,omitempty"`
	HasBeenDone      bool               `bson:"hasBeenDone"`
	Content          string             `bson:"content"`
	TodoId           bsoncodec.ObjectId `bson:"todoId"`
	DoneAt           time.Time          `bson:"doneAt,omitempty"`
	NeedRemind       bool               `bson:"needRemind"`
	UserId           string             `bson:"userId"`
	HasBeenReminded  bool               `bson:"hasBeenReminded"`
	IsRepeatable     bool               `bson:"isRepeatable"`
	RepeatType       string             `bson:"repeatType"`
	RepeatDateOffset int                `bson:"repeatDateOffset"`
	Images           []string           `bson:"images,omitempty"`
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
		"isDeleted":       false,
		"todoId":          todoId,
		"hasBeenDone":     false,
		"hasBeenReminded": false,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"isDeleted": true,
		},
	}
	_, err := repository.Mongo.UpdateAll(ctx, C_TODO_RECORD, condition, updater)
	return err
}

func (*TodoRecord) Done(ctx context.Context, id bsoncodec.ObjectId) error {
	condition := bsoncodec.M{
		"_id": id,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"hasBeenDone": true,
			"doneAt":      time.Now(),
		},
	}
	change := qmgo.Change{
		Upsert:    false,
		ReturnNew: true,
		Update:    updater,
	}
	r := TodoRecord{}
	err := repository.Mongo.FindAndApply(ctx, C_TODO_RECORD, condition, change, &r)
	if err != nil {
		return err
	}
	go func() {
		CTodo.GenNextRecord(ctx, r.TodoId, false)
	}()
	return nil
}

func (*TodoRecord) Undo(ctx context.Context, id bsoncodec.ObjectId) error {
	condition := bsoncodec.M{
		"_id": id,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"hasBeenDone": false,
		},
		"$unset": bsoncodec.M{
			"doneAt": "",
		},
	}
	return repository.Mongo.UpdateOne(ctx, C_TODO_RECORD, condition, updater)
}

func (*TodoRecord) Delete(ctx context.Context, id bsoncodec.ObjectId) error {
	condition := bsoncodec.M{
		"_id": id,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"isDeleted": true,
		},
	}
	return repository.Mongo.UpdateOne(ctx, C_TODO_RECORD, condition, updater)
}

func (*TodoRecord) Delay(ctx context.Context, id bsoncodec.ObjectId, delayDuration time.Duration) error {
	condition := bsoncodec.M{
		"_id": id,
	}
	r := TodoRecord{}
	err := repository.Mongo.FindOne(ctx, C_TODO_RECORD, condition, &r)
	if err != nil {
		return err
	}
	r.RemindAt = r.RemindAt.Add(delayDuration)
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"remindAt":  r.RemindAt,
			"updatedAt": time.Now(),
		},
	}
	return repository.Mongo.UpdateOne(ctx, C_TODO_RECORD, condition, updater)
}

func (*TodoRecord) UpdateById(ctx context.Context, id bsoncodec.ObjectId, updater bsoncodec.M) error {
	condition := bsoncodec.M{
		"_id": id,
	}
	return repository.Mongo.UpdateOne(ctx, C_TODO_RECORD, condition, updater)
}

func (*TodoRecord) ListNeedRemindOnes(ctx context.Context) ([]TodoRecord, error) {
	condition := bsoncodec.M{
		"isDeleted": false,
		"remindAt": bsoncodec.M{
			"$lte": time.Now(),
		},
		"needRemind":      true,
		"hasBeenDone":     false,
		"hasBeenReminded": false,
	}
	var records []TodoRecord
	err := repository.Mongo.FindAll(ctx, C_TODO_RECORD, condition, &records)
	if err != nil {
		return nil, err
	}
	return records, nil
}

func (*TodoRecord) MarkAsReminded(ctx context.Context, ids []bsoncodec.ObjectId) error {
	condition := bsoncodec.M{
		"_id": bsoncodec.M{
			"$in": ids,
		},
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"hasBeenReminded": true,
		},
	}
	_, err := repository.Mongo.UpdateAll(ctx, C_TODO_RECORD, condition, updater)
	return err
}

func (*TodoRecord) ListByPagination(ctx context.Context, condition bsoncodec.M, page, perPage int64, orderBy []string) (int64, []TodoRecord, error) {
	var r []TodoRecord
	total, err := repository.Mongo.FindAllWithPage(ctx, C_TODO_RECORD, orderBy, page, perPage, condition, &r)
	return total, r, err
}

func (*TodoRecord) GetById(ctx context.Context, id bsoncodec.ObjectId) (TodoRecord, error) {
	condition := bsoncodec.M{
		"_id": id,
	}
	r := TodoRecord{}
	err := repository.Mongo.FindOne(ctx, C_TODO_RECORD, condition, &r)
	return r, err
}

func (*TodoRecord) DeleteUndoneOnesByTodoId(ctx context.Context, todoId bsoncodec.ObjectId) error {
	condition := bsoncodec.M{
		"todoId":      todoId,
		"hasBeenDone": false,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"isDeleted": true,
		},
	}
	if _, err := repository.Mongo.UpdateAll(ctx, C_TODO_RECORD, condition, updater); err != nil && err != qmgo.ErrNoSuchDocuments {
		return err
	}
	return nil
}

func (*TodoRecord) CountNotDoneRecordsByTodoId(ctx context.Context, todoId bsoncodec.ObjectId) (int64, error) {
	condition := bsoncodec.M{
		"isDeleted":   false,
		"todoId":      todoId,
		"hasBeenDone": false,
	}
	return repository.Mongo.Count(ctx, C_TODO_RECORD, condition)
}
