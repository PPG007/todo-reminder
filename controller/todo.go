package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo-reminder/model"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

func init() {
	registerApi(ReminderApi{
		Endpoint: "/todos/upsert",
		Method:   http.MethodPost,
		Handler:  UpsertTodo,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/:id",
		Method:   http.MethodDelete,
		Handler:  DeleteTodo,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/uploadUrl",
		Method:   http.MethodGet,
		Handler:  GenUploadUrl,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/object/:name",
		Method:   http.MethodGet,
		Handler:  GetObject,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/object/uploadProxy",
		Method:   http.MethodPost,
		Handler:  UploadProxy,
	})
}

type UpsertTodoRequest struct {
	Id               string   `json:"id"`
	NeedRemind       bool     `json:"needRemind"`
	Content          string   `json:"content" binding:"required"`
	RemindAt         string   `json:"remindAt"`
	IsRepeatable     bool     `json:"isRepeatable"`
	RepeatType       string   `json:"repeatType"`
	RepeatDateOffset int      `json:"repeatDateOffset"`
	Images           []string `json:"images"`
}

type TodoDetail struct {
	Id          string `json:"id"`
	CreatedAt   string `json:"createdAt"`
	HasBeenDone bool   `json:"hasBeenDone"`
	NeedRemind  bool   `json:"needRemind"`
	DoneAt      string `json:"doneAt"`
	Content     string `json:"content"`
	RemindAt    string `json:"remindAt"`
	RepeatType  string `json:"repeatType"`
	Month       int    `json:"month"`
	Weekday     int    `json:"weekday"`
	Day         int    `json:"day"`
}

type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type SearchTodoRequest struct {
	HasBeenDone bool       `json:"hasBeenDone"`
	OrderBy     []string   `json:"orderBy"`
	RemindAt    *TimeRange `json:"remindAt"`
}

type SearchTodoResponse struct {
	Total int          `json:"total"`
	Items []TodoDetail `json:"items"`
}

func UpsertTodo(ctx *gin.Context) {
	req := UpsertTodoRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	var (
		remindAt time.Time
	)
	if req.NeedRemind {
		remindAt, err = util.TransTimeStrToTime(req.RemindAt)
		if err != nil {
			ReturnError(ctx, err)
			return
		}
	}
	if req.IsRepeatable && !util.StrInArray(req.RepeatType, &[]string{
		model.REPEAT_TYPE_DAILY,
		model.REPEAT_TYPE_WEEKLY,
		model.REPEAT_TYPE_MONTHLY,
		model.REPEAT_TYPE_YEARLY,
		model.REPEAT_TYPE_WORKING_DAY,
		model.REPEAT_TYPE_HOLIDAY,
	}) {
		ReturnError(ctx, errors.New("invalid repeat type"))
		return
	}
	todo := model.Todo{
		NeedRemind: req.NeedRemind,
		Content:    req.Content,
		UserId:     util.ExtractUserId(ctx),
		RemindSetting: model.RemindSetting{
			RemindAt:     remindAt,
			IsRepeatable: req.IsRepeatable,
			RepeatSetting: model.RepeatSetting{
				Type:       req.RepeatType,
				DateOffset: req.RepeatDateOffset,
			},
		},
		Images: req.Images,
	}
	if req.Id != "" {
		if bsoncodec.IsObjectIdHex(req.Id) {
			todo.Id = bsoncodec.ObjectIdHex(req.Id)
		} else {
			ReturnError(ctx, errors.New("invalid todo id"))
			return
		}
	} else {
		todo.Id = bsoncodec.NewObjectId()
	}
	err = todo.Upsert(ctx)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func DeleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ReturnError(ctx, errors.New("invalid todo id"))
		return
	}
	err := model.CTodo.DeleteById(ctx, bsoncodec.ObjectIdHex(id))
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func GenUploadUrl(ctx *gin.Context) {
	fileName := ctx.Query("fileName")
	uniqueName := fmt.Sprintf("%s_%s", bsoncodec.NewObjectId().Hex(), fileName)
	url, err := util.MinioClient.GetPreSignPutObjectUrl(ctx, uniqueName)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{
		"url":  url,
		"name": uniqueName,
	})
}

func GetObject(ctx *gin.Context) {
	fileName := ctx.Param("name")
	url, err := util.MinioClient.SignObjectUrl(ctx, fileName)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{
		"url": url,
	})
}

func UploadProxy(ctx *gin.Context) {
	fileName := ctx.Query("fileName")
	if fileName == "" {
		ctx.JSON(http.StatusBadRequest, map[string]string{
			"message": "need file name",
		})
		return
	}
	file, _, err := ctx.Request.FormFile(fileName)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	uniqueName := fmt.Sprintf("%s_%s", bsoncodec.NewObjectId().Hex(), fileName)
	url, err := util.MinioClient.PutObject(ctx, uniqueName, file)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{
		"name": uniqueName,
		"url":  url,
	})
}
