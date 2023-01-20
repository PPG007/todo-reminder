package repository

import (
	"context"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	mgo_option "go.mongodb.org/mongo-driver/mongo/options"
	"todo-reminder/repository/bsoncodec"
)

var (
	Mongo dbRepository
)

func init() {
	ctx := context.Background()
	uri := viper.GetString("mongodb.uri")
	database := viper.GetString("mongodb.database")
	client, err := qmgo.Open(ctx, &qmgo.Config{
		Uri:      uri,
		Database: database,
	}, options.ClientOptions{
		&mgo_option.ClientOptions{
			Registry: bsoncodec.DefaultRegistry,
		},
	})
	if err != nil {
		panic(err)
	}
	Mongo = mongoRepository{
		client: client,
	}
}

type dbRepository interface {
	Insert(ctx context.Context, collection string, docs ...interface{}) error
	UpdateOne(ctx context.Context, collection string, condition bson.M, updater bson.M) error
	UpdateAll(ctx context.Context, collection string, condition, updater bsoncodec.M) (int64, error)
	FindAll(ctx context.Context, collection string, condition bson.M, result interface{}) error
	FindOne(ctx context.Context, collection string, condition bson.M, result interface{}) error
	Count(ctx context.Context, collection string, condition bson.M) (int64, error)
	FindAndApply(ctx context.Context, collection string, condition bson.M, change qmgo.Change, result interface{}) error
	CreateIndex(ctx context.Context, collection string, index options.IndexModel) error
	FindAllWithSorter(ctx context.Context, collection string, sorter []string, condition bsoncodec.M, result interface{}) error
	FindOneWithSorter(ctx context.Context, collection string, sorter []string, condition bsoncodec.M, result interface{}) error
	FindAllWithPage(ctx context.Context, collection string, sorter []string, page, perPage int64, condition bsoncodec.M, result interface{}) error
}

type mongoRepository struct {
	client *qmgo.QmgoClient
}

func (m mongoRepository) Insert(ctx context.Context, collection string, docs ...interface{}) error {
	if len(docs) == 0 {
		return nil
	}
	if len(docs) == 1 {
		_, err := m.client.Database.Collection(collection).InsertOne(ctx, docs[0])
		return err
	} else {
		_, err := m.client.Database.Collection(collection).InsertMany(ctx, docs)
		return err
	}
}

func (m mongoRepository) UpdateOne(ctx context.Context, collection string, condition bson.M, updater bson.M) error {
	return m.client.Database.Collection(collection).UpdateOne(ctx, condition, updater)
}

func (m mongoRepository) FindAll(ctx context.Context, collection string, condition bson.M, result interface{}) error {
	return m.client.Database.Collection(collection).Find(ctx, condition).All(result)
}

func (m mongoRepository) FindOne(ctx context.Context, collection string, condition bson.M, result interface{}) error {
	return m.client.Database.Collection(collection).Find(ctx, condition).One(result)
}

func (m mongoRepository) Count(ctx context.Context, collection string, condition bson.M) (int64, error) {
	return m.client.Database.Collection(collection).Find(ctx, condition).Count()
}

func (m mongoRepository) FindAndApply(ctx context.Context, collection string, condition bson.M, change qmgo.Change, result interface{}) error {
	return m.client.Database.Collection(collection).Find(ctx, condition).Apply(change, result)
}

func (m mongoRepository) CreateIndex(ctx context.Context, collection string, index options.IndexModel) error {
	return m.client.Database.Collection(collection).CreateOneIndex(ctx, index)
}

func (m mongoRepository) FindOneWithSorter(ctx context.Context, collection string, sorter []string, condition bsoncodec.M, result interface{}) error {
	return m.client.Database.Collection(collection).Find(ctx, condition).Sort(sorter...).One(result)
}

func (m mongoRepository) UpdateAll(ctx context.Context, collection string, condition, updater bsoncodec.M) (int64, error) {
	result, err := m.client.Database.Collection(collection).UpdateAll(ctx, condition, updater)
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

func (m mongoRepository) FindAllWithSorter(ctx context.Context, collection string, sorter []string, condition bsoncodec.M, result interface{}) error {
	return m.client.Database.Collection(collection).Find(ctx, condition).Sort(sorter...).All(result)
}

func (m mongoRepository) FindAllWithPage(ctx context.Context, collection string, sorter []string, page, perPage int64, condition bsoncodec.M, result interface{}) error {
	return m.client.Database.Collection(collection).Find(ctx, condition).Sort(sorter...).Skip((page - 1) * perPage).Limit(perPage).All(result)
}
