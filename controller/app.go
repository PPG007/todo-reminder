package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo-reminder/model"
	"todo-reminder/util"
)

func init() {
	registerApi(ReminderApi{
		Endpoint: "/app/latest",
		Handler:  GetLatestApp,
		Method:   http.MethodGet,
		NoAuth:   true,
	})
	registerApi(ReminderApi{
		Endpoint: "/app/:version/url",
		Handler:  GetAppUrl,
		Method:   http.MethodGet,
		NoAuth:   true,
	})
}

type AppVersionResponse struct {
	Id        string `json:"id"`
	Version   string `json:"version"`
	CreatedAt string `json:"createdAt"`
	FileName  string `json:"fileName"`
}

func GetLatestApp(ctx *gin.Context) {
	latestVersion, err := model.CAppVersion.GetLatestVersion(ctx)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, AppVersionResponse{
		Id:        latestVersion.Id.Hex(),
		Version:   latestVersion.Version,
		CreatedAt: latestVersion.CreatedAt.Format(time.RFC3339),
		FileName:  latestVersion.FileName,
	})
}

func GetAppUrl(ctx *gin.Context) {
	version := ctx.Param("version")
	appVersion, err := model.CAppVersion.GetByVersion(ctx, version)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	objectUrl, err := util.MinioClient.SignObjectUrl(ctx, appVersion.FileName)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{
		"url": objectUrl,
	})
}
