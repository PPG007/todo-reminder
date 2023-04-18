package controller

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo-reminder/model"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

func init() {
	registerApi(ReminderApi{
		Endpoint: "/todoRecord/:id/done",
		Method:   http.MethodPost,
		Handler:  DoneTodoRecord,
	})
	registerApi(ReminderApi{
		Endpoint: "/todoRecord/:id/undo",
		Method:   http.MethodPost,
		Handler:  UndoTodoRecord,
	})
	registerApi(ReminderApi{
		Endpoint: "/todoRecord/:id",
		Method:   http.MethodDelete,
		Handler:  DeleteOneRecord,
	})
	registerApi(ReminderApi{
		Endpoint: "/todoRecord/:id",
		Method:   http.MethodGet,
		Handler:  GetTodoRecordById,
	})
	registerApi(ReminderApi{
		Endpoint: "/todoRecord/:id/delay",
		Method:   http.MethodPost,
		Handler:  DelayTodoRecord,
	})
	registerApi(ReminderApi{
		Endpoint: "/todoRecords/search",
		Method:   http.MethodPost,
		Handler:  ListTodoRecords,
	})
}

func DoneTodoRecord(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ReturnError(ctx, errors.New("invalid id"))
		return
	}
	err := model.CTodoRecord.Done(ctx, bsoncodec.ObjectIdHex(id))
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func UndoTodoRecord(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ReturnError(ctx, errors.New("invalid id"))
		return
	}
	model.CTodoRecord.Undo(ctx, bsoncodec.ObjectIdHex(id))
}

func DeleteOneRecord(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ReturnError(ctx, errors.New("invalid id"))
		return
	}
	err := model.CTodoRecord.Delete(ctx, bsoncodec.ObjectIdHex(id))
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

type DelayTodoRecordRequest struct {
	Second int64 `json:"second" binding:"required"`
}

func DelayTodoRecord(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ReturnError(ctx, errors.New("invalid id"))
		return
	}
	req := DelayTodoRecordRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	err = model.CTodoRecord.Delay(ctx, bsoncodec.ObjectIdHex(id), time.Second*time.Duration(req.Second))
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

type ListTodoRecordsRequest struct {
	HasBeenDone   bool          `json:"hasBeenDone"`
	ListCondition ListCondition `json:"listCondition"`
}

type ListCondition struct {
	OrderBy []string `json:"orderBy"`
	Page    int64    `json:"page"`
	PerPage int64    `json:"perPage"`
}

type ListTodoRecordsResponse struct {
	Items []TodoRecordDetail `json:"items"`
	Total int64              `json:"total"`
}
type TodoRecordDetail struct {
	Id               string  `json:"id"`
	RemindAt         string  `json:"remindAt"`
	HasBeenDone      bool    `json:"hasBeenDone"`
	Content          string  `json:"content"`
	DoneAt           string  `json:"doneAt"`
	NeedRemind       bool    `json:"needRemind"`
	IsRepeatable     bool    `json:"isRepeatable"`
	RepeatType       string  `json:"repeatType"`
	RepeatDateOffset int     `json:"repeatDateOffset"`
	TodoId           string  `json:"todoId"`
	Images           []Image `json:"images"`
}

type Image struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

func formatListCondition(listCondition ListCondition) ListCondition {
	if listCondition.Page == 0 {
		listCondition.Page = 1
	}
	if listCondition.PerPage == 0 {
		listCondition.PerPage = 100
	}
	return listCondition
}

func ListTodoRecords(ctx *gin.Context) {
	req := ListTodoRecordsRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	userId := util.ExtractUserId(ctx)
	condition := bsoncodec.M{
		"isDeleted":   false,
		"hasBeenDone": req.HasBeenDone,
		"userId":      userId,
	}
	req.ListCondition = formatListCondition(req.ListCondition)
	if !req.HasBeenDone {
		req.ListCondition.OrderBy = []string{"remindAt"}
	}
	total, todoRecords, err := model.CTodoRecord.ListByPagination(ctx, condition, req.ListCondition.Page, req.ListCondition.PerPage, req.ListCondition.OrderBy)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, ListTodoRecordsResponse{
		Items: formatTodoRecordDetails(ctx, todoRecords),
		Total: total,
	})
}

func GetTodoRecordById(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ReturnError(ctx, errors.New("invalid id"))
		return
	}
	record, err := model.CTodoRecord.GetById(ctx, bsoncodec.ObjectIdHex(id))
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, formatTodoRecordDetail(ctx, record))
}

func formatTodoRecordDetail(ctx context.Context, record model.TodoRecord) TodoRecordDetail {
	return TodoRecordDetail{
		Id:               record.Id.Hex(),
		RemindAt:         util.TransTimeToRFC3339(record.RemindAt),
		HasBeenDone:      record.HasBeenDone,
		Content:          record.Content,
		DoneAt:           util.TransTimeToRFC3339(record.DoneAt),
		NeedRemind:       record.NeedRemind,
		IsRepeatable:     record.IsRepeatable,
		RepeatType:       record.RepeatType,
		RepeatDateOffset: record.RepeatDateOffset,
		TodoId:           record.TodoId.Hex(),
		Images: func() []Image {
			result := make([]Image, 0, len(record.Images))
			for _, image := range record.Images {
				url, _ := util.MinioClient.SignObjectUrl(ctx, image)
				result = append(result, Image{
					Name: image,
					Url:  url,
				})
			}
			return result
		}(),
	}
}

func formatTodoRecordDetails(ctx context.Context, records []model.TodoRecord) []TodoRecordDetail {
	details := make([]TodoRecordDetail, 0, len(records))
	for _, record := range records {
		details = append(details, formatTodoRecordDetail(ctx, record))
	}
	return details
}
