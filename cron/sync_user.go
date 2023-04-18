package cron

import (
	"context"
	"github.com/spf13/cast"
	"todo-reminder/gocq"
	"todo-reminder/log"
	"todo-reminder/model"
)

func init() {
	registerCronTask("@every 1m", SyncUser, true)
}

func SyncUser() {
	ctx := context.Background()
	friends, err := gocq.GetGocqInstance().ListFriendsWithFullInfo(ctx)
	if err != nil {
		return
	}
	for _, friend := range friends {
		err := UpsertUser(ctx, friend.UserId, friend.Nickname)
		if err != nil {
			log.Warn("Failed to sync user", map[string]interface{}{
				"user":  friend,
				"error": err.Error(),
			})
		}
	}
}

func UpsertUser(ctx context.Context, userId int64, nickname string) error {
	user := model.User{
		UserId:   cast.ToString(userId),
		Nickname: nickname,
	}
	return user.UpsertWithoutPassword(ctx)
}
