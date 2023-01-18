package cron

import (
	"context"
	"todo-reminder/gocq"
	"todo-reminder/model"
)

func SyncUser() {
	ctx := context.Background()
	userIds, err := gocq.GoCq.ListFriends(ctx)
	if err != nil {
		return
	}
	for _, id := range userIds {
		user := model.User{
			UserId: id,
		}
		user.UpsertWithoutPassword(ctx)
	}
}
