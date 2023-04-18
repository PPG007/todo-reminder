package cron

import (
	"context"
	"todo-reminder/gocq"
	"todo-reminder/model"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

func init() {
	registerCronTask("@every 20s", Remind, false)
}

func Remind() {
	ctx := context.Background()
	records, err := model.CTodoRecord.ListNeedRemindOnes(ctx)
	if err != nil {
		return
	}
	var succeedIds []bsoncodec.ObjectId
	var succeedTodoIds []bsoncodec.ObjectId
	for _, record := range records {
		if err := notify(ctx, record); err == nil {
			succeedIds = append(succeedIds, record.Id)
			succeedTodoIds = append(succeedTodoIds, record.TodoId)
		}
	}
	model.CTodoRecord.MarkAsReminded(ctx, succeedIds)
	//for _, id := range succeedTodoIds {
	//	model.CTodo.GenNextRecord(ctx, id, false)
	//}
}

func notify(ctx context.Context, record model.TodoRecord) error {
	err := gocq.GetGocqInstance().SendPrivateStringMessage(ctx, record.Content, record.UserId)
	if err != nil {
		return err
	}
	for _, image := range record.Images {
		url, err := util.MinioClient.SignObjectUrl(ctx, image)
		if err == nil {
			err = gocq.GetGocqInstance().SendPrivateImageMessage(ctx, record.UserId, image, url)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
