package middleware

import (
	"github.com/gin-gonic/gin"
	"sort"
)

var middlewares []middleware

type middleware struct {
	f     gin.HandlerFunc
	order int
}

func registerMiddleware(f gin.HandlerFunc, order int) {
	middlewares = append(middlewares, middleware{
		f,
		order,
	})
}

func Init(e *gin.Engine) {
	sort.SliceStable(middlewares, func(i, j int) bool {
		return middlewares[i].order < middlewares[j].order
	})
	for _, m := range middlewares {
		e.Use(m.f)
	}
}
