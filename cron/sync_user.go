package cron

import (
	"context"
	"todo-reminder/gocq"
	"todo-reminder/log"
	"todo-reminder/model"
)

func SyncUser() {
	ctx := context.Background()
	userIds, err := gocq.GoCq.ListFriends(ctx)
	if err != nil {
		return
	}
	for _, id := range userIds {
		err := UpsertUser(ctx, id)
		if err != nil {
			log.Warn("Failed to sync user", map[string]interface{}{
				"userId": id,
			})
		}
	}
}

func UpsertUser(ctx context.Context, userId string) error {
	user := model.User{
		UserId: userId,
	}
	return user.UpsertWithoutPassword(ctx)
}
