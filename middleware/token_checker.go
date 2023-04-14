package middleware

import (
	"github.com/gin-gonic/gin"
	"todo-reminder/constant"
	"todo-reminder/controller"
	"todo-reminder/util"
)

func init() {
	registerMiddleware(checkToken, 3)
}

func checkToken(ctx *gin.Context) {
	if noAuthHandler(ctx) {
		ctx.Next()
		return
	}
	tokenStr := ctx.GetHeader(constant.HEADER_TOKEN)
	token, err := util.ParseToken(tokenStr)
	if err != nil {
		controller.ReturnError(ctx, err)
		return
	}
	ctx.Set(constant.GIN_KEY_USER_ID, token.UserId)
	ctx.Next()
}

func noAuthHandler(ctx *gin.Context) bool {
	path := ctx.FullPath()
	method := ctx.Request.Method
	paths, ok := controller.NoAuthPath[method]
	if !ok {
		return false
	}
	return util.StrInArray(path, &paths)
}
