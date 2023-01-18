package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-reminder/model"
)

func init() {
	registerApi(ReminderApi{
		Endpoint: "/user/login",
		Method:   http.MethodPost,
		Handler:  Login,
		NoAuth:   true,
	})
	registerApi(ReminderApi{
		Endpoint: "/user/register",
		Method:   http.MethodPost,
		Handler:  Register,
		NoAuth:   true,
	})
}

type LoginRequest struct {
	UserId   string `json:"userId" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type EmptyResponse struct {
}

func Login(ctx *gin.Context) {
	req := LoginRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	token, err := model.CUser.Login(ctx, req.UserId, req.Password)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, token)
}

func Register(ctx *gin.Context) {
	req := LoginRequest{}
	err := ctx.ShouldBind(&req)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	user, _ := model.CUser.GetByUserId(ctx, req.UserId)
	if user.Id.Valid() {
		ReturnError(ctx, errors.New("user has already registered"))
		return
	}
	err = model.CUser.Create(ctx, req.UserId, req.Password)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}
