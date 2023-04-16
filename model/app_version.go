package model

import (
	"context"
	"errors"
	"github.com/qiniu/qmgo/options"
	mgo_option "go.mongodb.org/mongo-driver/mongo/options"
	"time"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

const (
	C_APP_VERSION = "appVersion"
)

var (
	CAppVersion = &AppVersion{}
)

func init() {
	repository.Mongo.CreateIndex(context.Background(), C_APP_VERSION, options.IndexModel{
		Key: []string{"version"},
		IndexOptions: &mgo_option.IndexOptions{
			Unique: util.PtrValue[bool](true),
		},
	})
}

type AppVersion struct {
	Id        bsoncodec.ObjectId `bson:"_id"`
	Version   string             `bson:"version"`
	FileName  string             `bson:"fileName"`
	CreatedAt time.Time          `bson:"createdAt"`
}

func (a *AppVersion) Create(ctx context.Context) error {
	if !util.MinioClient.IsObjectExist(ctx, a.FileName) {
		return errors.New("file not found")
	}
	a.CreatedAt = time.Now()
	a.Id = bsoncodec.NewObjectId()
	return repository.Mongo.Insert(ctx, C_APP_VERSION, a)
}

func (*AppVersion) GetLatestVersion(ctx context.Context) (*AppVersion, error) {
	appVersion := &AppVersion{}
	err := repository.Mongo.FindOneWithSorter(ctx, C_APP_VERSION, []string{"-createdAt"}, nil, appVersion)
	if err != nil {
		return nil, err
	}
	return appVersion, nil
}

func (*AppVersion) GetByVersion(ctx context.Context, version string) (*AppVersion, error) {
	appVersion := &AppVersion{}
	condition := bsoncodec.M{
		"version": version,
	}
	err := repository.Mongo.FindOne(ctx, C_APP_VERSION, condition, appVersion)
	if err != nil {
		return nil, err
	}
	return appVersion, nil
}
