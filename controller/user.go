package controller

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-reminder/gocq"
	"todo-reminder/model"
	"todo-reminder/util"
)

func init() {
	registerApi(ReminderApi{
		Endpoint: "/user/login",
		Method:   http.MethodPost,
		Handler:  Login,
		NoAuth:   true,
	})
	registerApi(ReminderApi{
		Endpoint: "/user/:userId/genPassword",
		Method:   http.MethodPost,
		Handler:  GetDefaultPassword,
		NoAuth:   true,
	})
	registerApi(ReminderApi{
		Endpoint: "/user/password",
		Method:   http.MethodPut,
		Handler:  UpdatePassword,
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

type UpdatePasswordRequest struct {
	NewPassword string `json:"newPassword"`
}

func UpdatePassword(ctx *gin.Context) {
	userId := ctx.GetString("userId")
	req := UpdatePasswordRequest{}
	if err := ctx.ShouldBind(&req); err != nil {
		ReturnError(ctx, err)
		return
	}
	if len(req.NewPassword) < 10 {
		ReturnError(ctx, errors.New("password is to short"))
		return
	}
	err := model.CUser.UpdatePassword(ctx, userId, req.NewPassword)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func GetDefaultPassword(ctx *gin.Context) {
	userId := ctx.Param("userId")
	user, err := model.CUser.GetByUserId(ctx, userId)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	if user.IsEnabled {
		ReturnError(ctx, errors.New("invalid user"))
		return
	}
	password := util.GenRandomString(10)
	err = model.CUser.UpdatePassword(ctx, userId, password)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	err = gocq.GoCq.SendPrivateStringMessage(ctx, password, userId)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}
