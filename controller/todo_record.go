package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo-reminder/model"
	"todo-reminder/repository/bsoncodec"
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
		Endpoint: "/todoRecord/:id/delay",
		Method:   http.MethodPost,
		Handler:  DelayTodoRecord,
	})
	registerApi(ReminderApi{
		Endpoint: "/todoRecords",
		Method:   http.MethodGet,
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
	Items []TodoRecordDetail `json:"items,omitempty"`
}
type TodoRecordDetail struct {
	Id          string `json:"id,omitempty"`
	RemindAt    string `json:"remindAt,omitempty"`
	HasBeenDone bool   `json:"hasBeenDone,omitempty"`
	Content     string `json:"content,omitempty"`
	DoneAt      string `json:"doneAt,omitempty"`
	NeedRemind  bool   `json:"needRemind,omitempty"`
}

func ListTodoRecords(ctx *gin.Context) {
	req := ListTodoRecordsRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	userId := ctx.GetString("userId")
	condition := bsoncodec.M{
		"isDeleted":   false,
		"hasBeenDone": req.HasBeenDone,
		"userId":      userId,
	}
	todoRecords, err := model.CTodoRecord.ListByCondition(ctx, condition, req.ListCondition.Page, req.ListCondition.PerPage, req.ListCondition.OrderBy)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, ListTodoRecordsResponse{
		Items: formatTodoRecordDetails(todoRecords),
	})
}

func formatTodoRecordDetail(record model.TodoRecord) TodoRecordDetail {
	return TodoRecordDetail{
		Id:          "",
		RemindAt:    "",
		HasBeenDone: false,
		Content:     "",
		DoneAt:      "",
		NeedRemind:  false,
	}
}

func formatTodoRecordDetails(records []model.TodoRecord) []TodoRecordDetail {
	details := make([]TodoRecordDetail, 0, len(records))
	for _, record := range records {
		details = append(details, formatTodoRecordDetail(record))
	}
	return details
}
