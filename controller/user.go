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
		Method:   http.MethodGet,
		Handler:  GetDefaultPassword,
		NoAuth:   true,
	})
	registerApi(ReminderApi{
		Endpoint: "/user/password",
		Method:   http.MethodPut,
		Handler:  UpdatePassword,
	})
	registerApi(ReminderApi{
		Endpoint: "/user/userId",
		Method:   http.MethodGet,
		Handler:  GetCurrentUserId,
	})
	registerApi(ReminderApi{
		Endpoint: "/user/validToken",
		Method:   http.MethodPost,
		Handler:  ValidToken,
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
	ctx.JSON(http.StatusOK, map[string]string{
		"token": token,
	})
}

type UpdatePasswordRequest struct {
	NewPassword string `json:"newPassword"`
}

func UpdatePassword(ctx *gin.Context) {
	userId := util.ExtractUserId(ctx)
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
	_, err := model.CUser.GetByUserId(ctx, userId)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	password := util.GenRandomString(10)
	err = model.CUser.UpdatePassword(ctx, userId, password)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	err = gocq.GetGocqInstance().SendPrivateStringMessage(ctx, password, userId)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func ValidToken(ctx *gin.Context) {
	tokenStr := ctx.GetHeader("x-access-token")
	_, err := util.ParseToken(tokenStr)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, EmptyResponse{})
}

func GetCurrentUserId(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, map[string]string{
		"userId": util.ExtractUserId(ctx),
	})
}
