package middleware

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"todo-reminder/constant"
)

func init() {
	registerMiddleware(responseWriter, 3)
}

type responseWriterMiddleware struct {
	gin.ResponseWriter
	ResponseBody *bytes.Buffer
}

func (rw responseWriterMiddleware) Write(data []byte) (int, error) {
	rw.ResponseBody.Write(data)
	return rw.ResponseWriter.Write(data)
}

func responseWriter(ctx *gin.Context) {
	rw := &responseWriterMiddleware{
		ResponseBody:   &bytes.Buffer{},
		ResponseWriter: ctx.Writer,
	}
	ctx.Writer = rw
	ctx.Next()
	if ctx.Writer.Status() != http.StatusOK {
		ctx.Set(constant.GIN_KEY_RESPONSE_BODY, rw.ResponseBody.String())
	}
}
