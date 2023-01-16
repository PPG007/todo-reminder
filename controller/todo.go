package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
	"todo-reminder/model"
	"todo-reminder/repository/bsoncodec"
	"todo-reminder/util"
)

type handler = func(c *gin.Context)

type ReminderApi struct {
	Endpoint string
	Handler  handler
	Method   string
}

func init() {
	registerApi(ReminderApi{
		Endpoint: "/todos",
		Method:   http.MethodPost,
		Handler:  AddTodo,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/:id",
		Method:   http.MethodDelete,
		Handler:  DeleteTodo,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/:id",
		Method:   http.MethodPut,
		Handler:  EditTodo,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/:id/done",
		Method:   http.MethodPost,
		Handler:  DoneTodo,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/search",
		Method:   http.MethodPost,
		Handler:  SearchTodos,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/:id",
		Method:   http.MethodGet,
		Handler:  GetTodo,
	})
	registerApi(ReminderApi{
		Endpoint: "/todos/:id/rollback",
		Method:   http.MethodPost,
		Handler:  RollbackTodo,
	})
}

type AddTodoRequest struct {
	NeedRemind    bool                 `json:"needRemind"`
	Content       string               `json:"content" binding:"required"`
	RemindSetting *model.RemindSetting `json:"remindSetting"`
	HasBeenDone   bool                 `json:"hasBeenDone"`
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

func AddTodo(ctx *gin.Context) {
	req := AddTodoRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
	if req.NeedRemind && req.RemindSetting == nil {
		ctx.Error(errors.New("empty remind setting"))
		return
	}
	todo := model.Todo{
		HasBeenDone: false,
		NeedRemind:  req.NeedRemind,
		Content:     req.Content,
		UserId:      ctx.GetString("userId"),
	}
	if req.RemindSetting != nil {
		todo.RemindSetting = *req.RemindSetting
	}
	err = todo.Create(ctx)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func DeleteTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ctx.Error(errors.New("invalid todo id"))
		return
	}
	err := model.CTodo.DeleteById(ctx, bsoncodec.ObjectIdHex(id))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func EditTodo(ctx *gin.Context) {
	req := AddTodoRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ctx.Error(errors.New("invalid todo id"))
		return
	}
	setter := bson.M{
		"needRemind":  req.NeedRemind,
		"content":     req.Content,
		"hasBeenDone": req.HasBeenDone,
	}
	if req.RemindSetting != nil {
		setter["remindSetting"] = *req.RemindSetting
	}
	err = model.CTodo.UpdateById(ctx, bsoncodec.ObjectIdHex(id), bsoncodec.M{
		"$set": setter,
	})
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func DoneTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ctx.Error(errors.New("invalid todo id"))
		return
	}
	err := model.CTodo.Done(ctx, bsoncodec.ObjectIdHex(id))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func SearchTodos(ctx *gin.Context) {
	req := SearchTodoRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
	condition := bson.M{
		"isDeleted":   false,
		"userId":      ctx.GetString("userId"),
		"hasBeenDone": req.HasBeenDone,
	}
	if req.RemindAt != nil {
		condition["remindSetting.remindAt"] = bsoncodec.M{}
		if req.RemindAt.Start != "" {
			condition["remindSetting.remindAt"].(bsoncodec.M)["$gte"] = util.MustTransTimeStrToTime(req.RemindAt.Start)
		}
		if req.RemindAt.End != "" {
			condition["remindSetting.remindAt"].(bsoncodec.M)["$lte"] = util.MustTransTimeStrToTime(req.RemindAt.End)
		}
	}
	todos, err := model.CTodo.ListByCondition(ctx, condition)
	if err != nil {
		ctx.Error(err)
		return
	}
	resp := SearchTodoResponse{
		Total: len(todos),
		Items: formatTodoDetails(todos),
	}
	ctx.JSON(http.StatusOK, resp)
}

func GetTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ctx.Error(errors.New("invalid todo id"))
		return
	}
	todo, err := model.CTodo.GetById(ctx, bsoncodec.ObjectIdHex(id))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, formatTodoDetail(todo))
}

func RollbackTodo(ctx *gin.Context) {
	id := ctx.Param("id")
	if !bsoncodec.IsObjectIdHex(id) {
		ctx.Error(errors.New("invalid todo id"))
		return
	}
	err := model.CTodo.RollbackTodo(ctx, bsoncodec.ObjectIdHex(id))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func formatTodoDetails(todos []model.Todo) []TodoDetail {
	details := make([]TodoDetail, 0, len(todos))
	for _, todo := range todos {
		details = append(details, formatTodoDetail(todo))
	}
	return details
}

func formatTodoDetail(todo model.Todo) TodoDetail {
	return TodoDetail{
		Id:          todo.Id.Hex(),
		CreatedAt:   util.TransTimeToRFC3339(todo.CreatedAt),
		DoneAt:      util.TransTimeToRFC3339(todo.DoneAt),
		HasBeenDone: todo.HasBeenDone,
		NeedRemind:  todo.NeedRemind,
		Content:     todo.Content,
		RemindAt:    util.TransTimeToRFC3339(todo.RemindSetting.RemindAt),
		RepeatType:  todo.RemindSetting.RepeatSetting.Type,
		Month:       todo.RemindSetting.RepeatSetting.Month,
		Weekday:     todo.RemindSetting.RepeatSetting.Weekday,
		Day:         todo.RemindSetting.RepeatSetting.Day,
	}
}
