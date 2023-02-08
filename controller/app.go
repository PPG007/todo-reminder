package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	"todo-reminder/model"
)

func init() {
	registerApi(ReminderApi{
		Endpoint: "/appVersion/latest",
		Method:   http.MethodGet,
		Handler:  GetLatestAppVersion,
		NoAuth:   true,
	})
}

func GetLatestAppVersion(ctx *gin.Context) {
	appVersion, err := model.CAppVersion.GetLatestAppVersion(ctx)
	if err != nil {
		ReturnError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, map[string]string{
		"appVersionId": appVersion.Id.Hex(),
		"releasedAt":   appVersion.CreatedAt.Format(time.RFC3339),
		"url":          appVersion.URL,
	})
}
