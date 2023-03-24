package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"todo-reminder/log"
)

func Recovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			stack := make([]byte, 4096)
			stack = stack[:runtime.Stack(stack, false)]
			log.WarnTrace("Panic during processing", logrus.Fields{
				"error": fmt.Sprintf("%v", err),
			}, stack)
			ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
				"message": err,
			})
		}
	}()
	ctx.Next()
}
