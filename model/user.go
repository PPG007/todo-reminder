package model

import (
	"context"
	"errors"
	"github.com/qiniu/qmgo"
	"github.com/qiniu/qmgo/options"
	"github.com/spf13/viper"
	mgo_option "go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"time"
	"todo-reminder/repository"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

const (
	C_USER = "user"
)

var (
	CUser = &User{}
)

func init() {
	repository.Mongo.CreateIndex(context.Background(), C_USER, options.IndexModel{
		Key: []string{"userId", "isDeleted"},
		IndexOptions: &mgo_option.IndexOptions{
			Background: util.PtrValue[bool](true),
			Unique:     util.PtrValue[bool](true),
		},
	})
}

type User struct {
	Id             bsoncodec.ObjectId `json:"id" bson:"_id"`
	UserId         string             `json:"userId" bson:"userId"`
	Nickname       string             `bson:"nickname"`
	Password       string             `json:"password" bson:"password"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsDeleted      bool               `json:"isDeleted" bson:"isDeleted"`
	IsEnabled      bool               `json:"isEnabled" bson:"isEnabled"`
	OpenAIApproved bool               `json:"openAIApproved" bson:"openAIApproved"`
}

func (*User) Create(ctx context.Context, userId, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user := User{
		Id:        bsoncodec.NewObjectId(),
		UserId:    userId,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsDeleted: false,
	}

	return repository.Mongo.Insert(ctx, C_USER, user)
}

func (*User) Login(ctx context.Context, userId, password string) (string, error) {
	condition := bsoncodec.M{
		"isDeleted": false,
		"userId":    userId,
	}
	change := qmgo.Change{
		Update: bsoncodec.M{
			"$set": bsoncodec.M{
				"isEnabled": true,
				"updatedAt": time.Now(),
			},
		},
		ReturnNew: true,
		Upsert:    false,
	}
	// TODO: FindAndApply
	user := User{}
	err := repository.Mongo.FindAndApply(ctx, C_USER, condition, change, &user)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}
	return util.GenerateToken(user.UserId), nil
}

func (*User) GetByUserId(ctx context.Context, userId string) (User, error) {
	condition := bsoncodec.M{
		"userId":    userId,
		"isDeleted": false,
	}
	user := User{}
	err := repository.Mongo.FindOne(ctx, C_USER, condition, &user)
	return user, err
}

func (u *User) UpsertWithoutPassword(ctx context.Context) error {
	condition := bsoncodec.M{
		"userId":    u.UserId,
		"isDeleted": false,
	}
	change := qmgo.Change{
		Upsert:    true,
		ReturnNew: true,
		Update: bsoncodec.M{
			"$set": bsoncodec.M{
				"updatedAt": time.Now(),
				"nickname":  u.Nickname,
			},
			"$setOnInsert": bsoncodec.M{
				"createdAt": time.Now(),
			},
		},
	}
	return repository.Mongo.FindAndApply(ctx, C_USER, condition, change, u)
}

func (*User) UpdatePassword(ctx context.Context, userId, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	condition := bsoncodec.M{
		"userId":    userId,
		"isDeleted": false,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"password":  string(hashedPassword),
			"updatedAt": time.Now(),
		},
	}
	return repository.Mongo.UpdateOne(ctx, C_USER, condition, updater)
}

func (u *User) GetAdmin(ctx context.Context) (User, error) {
	return u.GetByUserId(ctx, viper.GetString("goCq.admin"))
}

func (u *User) ApproveOpenAI(ctx context.Context, userId string) (User, error) {
	user, err := u.GetByUserId(ctx, userId)
	if err != nil {
		return User{}, err
	}
	if user.OpenAIApproved {
		return User{}, errors.New("user has already been approved openAI")
	}
	condition := bsoncodec.M{
		"_id": user.Id,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"openAIApproved": true,
		},
	}
	err = repository.Mongo.UpdateOne(ctx, C_USER, condition, updater)
	return user, err
}

func (u *User) BlockOpenAI(ctx context.Context, userId string) (User, error) {
	user, err := u.GetByUserId(ctx, userId)
	if err != nil {
		return User{}, err
	}
	if !user.OpenAIApproved {
		return User{}, errors.New("user has already been blocked openAI")
	}
	condition := bsoncodec.M{
		"_id": user.Id,
	}
	updater := bsoncodec.M{
		"$set": bsoncodec.M{
			"openAIApproved": false,
		},
	}
	err = repository.Mongo.UpdateOne(ctx, C_USER, condition, updater)
	return user, err
}

func (*User) HandleFriendAdded(ctx context.Context, userId string) error {
	user := User{
		UserId: userId,
	}
	return user.UpsertWithoutPassword(ctx)
}
