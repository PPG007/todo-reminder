package middleware

import (
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"todo-reminder/constant"
)

func init() {
	registerMiddleware(requestId, 1)
}

func requestId(ctx *gin.Context) {
	reqId := uuid.NewV4().String()
	ctx.Request.Header.Set(constant.HEADER_REQUEST_ID, reqId)
	ctx.Header(constant.HEADER_REQUEST_ID, reqId)
	ctx.Next()
}
