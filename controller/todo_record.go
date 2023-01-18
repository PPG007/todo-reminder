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
