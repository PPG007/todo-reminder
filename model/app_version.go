package model

import (
	"context"
	"github.com/qiniu/qmgo/options"
	"time"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
)

const (
	C_APP_VERSION = "appVersion"
)

var (
	CAppVersion = &AppVersion{}
)

func init() {
	repository.Mongo.CreateIndex(context.Background(), C_CHINA_HOLIDAY, options.IndexModel{
		Key: []string{"createdAt"},
	})
}

type AppVersion struct {
	Id        bsoncodec.ObjectId `bson:"_id"`
	CreatedAt time.Time          `bson:"createdAt"`
	URL       string             `bson:"url"`
}

func (*AppVersion) GetLatestAppVersion(ctx context.Context) (AppVersion, error) {
	result := AppVersion{}
	err := repository.Mongo.FindOneWithSorter(ctx, C_APP_VERSION, []string{"-createdAt"}, bsoncodec.M{}, &result)
	return result, err
}

func (*AppVersion) CreateVersion(ctx context.Context, url string) error {
	appVersion := AppVersion{
		Id:        bsoncodec.NewObjectId(),
		CreatedAt: time.Now(),
		URL:       url,
	}
	return repository.Mongo.Insert(ctx, C_APP_VERSION, appVersion)
}
