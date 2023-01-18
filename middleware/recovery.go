package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Recovery(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, map[string]interface{}{
				"message": err,
			})
		}
	}()
	ctx.Next()
}
