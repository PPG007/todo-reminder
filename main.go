package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	_ "todo-reminder/conf"
	"todo-reminder/controller"
	"todo-reminder/cron"
	"todo-reminder/middleware"
)

func startGin() {
	e := gin.New()
	e.Use(middleware.Recovery)
	e.Use(middleware.CheckToken)
	for _, api := range controller.APIs {
		switch api.Method {
		case http.MethodPost:
			e.POST(api.Endpoint, api.Handler)
		case http.MethodDelete:
			e.DELETE(api.Endpoint, api.Handler)
		case http.MethodPut:
			e.PUT(api.Endpoint, api.Handler)
		case http.MethodGet:
			e.GET(api.Endpoint, api.Handler)
		}
	}
	e.NoRoute(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusNotFound, map[string]string{
			"message": "not found",
		})
	})
	e.NoMethod(func(ctx *gin.Context) {
		ctx.AbortWithStatusJSON(http.StatusMethodNotAllowed, map[string]string{
			"message": "method not allowed",
		})
	})
	err := e.Run("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
}

func main() {
	go cron.Start()
	startGin()
}
