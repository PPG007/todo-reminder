package cron

import (
	"context"
	"time"
	"todo-reminder/model"
	"todo-reminder/repository/bsoncodec"
)

func remind() {
	ctx := context.Background()
	condition := bsoncodec.M{
		"isDeleted":       false,
		"hasBeenReminded": false,
		"hasBeenDone":     false,
		"needRemind":      true,
		"remindSetting.remindAt": bsoncodec.M{
			"$lte": time.Now(),
		},
	}
	todos, _ := model.CTodo.ListByCondition(ctx, condition)
	for _, todo := range todos {
		todo.Notify(ctx)
		todo.RenewRemindAt(ctx)
	}
}
