package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"net/http"
	"os"
	"strings"
	_ "todo-reminder/conf"
	"todo-reminder/controller"
	"todo-reminder/cron"
	"todo-reminder/log"
	"todo-reminder/middleware"
	"todo-reminder/model"
	"todo-reminder/util"
)

func startGin() {
	e := gin.New()
	gin.SetMode(gin.ReleaseMode)
	middleware.Init(e)
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

func releaseApp(ctx context.Context, version, appName, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	suffix := ""
	if !strings.HasSuffix(appName, ".apk") {
		suffix = ".apk"
	}
	fileName := fmt.Sprintf("%s-%s%s", appName, version, suffix)
	_, err = util.MinioClient.PutObject(ctx, fileName, file)
	if err != nil {
		return err
	}
	appVersion := &model.AppVersion{
		Version:  version,
		FileName: fileName,
	}
	return appVersion.Create(ctx)
}

func main() {
	version := flag.String("version", "", "app version to release")
	appName := flag.String("appName", "", "app name")
	filePath := flag.String("filePath", "", "app file path")
	flag.Parse()
	if version != nil && appName != nil && filePath != nil && *version != "" && *appName != "" && *filePath != "" {
		err := releaseApp(context.Background(), *version, *appName, *filePath)
		if err != nil {
			log.Error("Failed to release app", map[string]interface{}{
				"error": err.Error(),
			})
		}
		return
	}
	go cron.Start()
	startGin()
}
