package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	_ "todo-reminder/conf"
	"todo-reminder/controller"
	"todo-reminder/cron"
	"todo-reminder/middleware"
)

func startGin() {
	e := gin.New()
	gin.SetMode(gin.ReleaseMode)
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
	port := viper.GetInt("gin.port")
	if port == 0 {
		port = 8080
	}
	err := e.Run(fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		panic(err)
	}
}

func main() {
	go cron.Start()
	startGin()
}
