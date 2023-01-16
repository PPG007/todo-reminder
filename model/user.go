package model

import (
	"context"
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

type User struct {
	Id        bsoncodec.ObjectId `json:"id" bson:"_id"`
	UserId    string             `json:"userId" bson:"userId"`
	Password  string             `json:"password" bson:"password"`
	CreatedAt time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time          `json:"updatedAt" bson:"updatedAt"`
	IsDeleted bool               `json:"isDeleted" bson:"isDeleted"`
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
	user := User{}
	err := repository.Mongo.FindOne(ctx, C_USER, condition, &user)
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
