package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-reminder/controller"
	"todo-reminder/util"
)

func CheckToken(ctx *gin.Context) {
	if NoAuthHandler(ctx) {
		ctx.Next()
		return
	}
	tokenStr := ctx.GetHeader("x-access-token")
	token, err := util.ParseToken(ctx, tokenStr)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	ctx.Set("userId", token.UserId)
	ctx.Next()
}

func NoAuthHandler(ctx *gin.Context) bool {
	path := ctx.FullPath()
	method := ctx.Request.Method
	paths, ok := controller.NoAuthPath[method]
	if !ok {
		return false
	}
	return util.StrInArray(path, &paths)
}
