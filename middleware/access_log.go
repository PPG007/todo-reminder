package middleware

import (
	"github.com/gin-gonic/gin"
	"todo-reminder/model"
)

func init() {
	registerMiddleware(access, 2)
}

func access(ctx *gin.Context) {
	accessLog := model.CAccessLog.Init(ctx)
	ctx.Next()
	accessLog.Record(ctx)
}
