package cron

import (
	"context"
	"todo-reminder/model"
	"todo-reminder/repository/bsoncodec"
)

func Remind() {
	ctx := context.Background()
	records, err := model.CTodoRecord.ListNeedRemindOnes(ctx)
	if err != nil {
		return
	}
	var succeedIds []bsoncodec.ObjectId
	for _, record := range records {
		if err := record.Notify(ctx); err == nil {
			succeedIds = append(succeedIds, record.Id)
		}
	}
	model.CTodoRecord.MarkAsReminded(ctx, succeedIds)
}
