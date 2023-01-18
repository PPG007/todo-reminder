package controller

import "github.com/gin-gonic/gin"

type handler = func(c *gin.Context)

type ReminderApi struct {
	Endpoint string
	Handler  handler
	Method   string
	NoAuth   bool
}

var (
	APIs []ReminderApi

	NoAuthPath map[string][]string
)

func init() {
	NoAuthPath = map[string][]string{}
}

func registerApi(api ReminderApi) {
	APIs = append(APIs, api)
	if api.NoAuth {
		registerNoAuthPath(api.Method, api.Endpoint)
	}
}

func registerNoAuthPath(method, path string) {
	if len(NoAuthPath[method]) == 0 {
		NoAuthPath[method] = []string{path}
		return
	}
	NoAuthPath[method] = append(NoAuthPath[method], path)
}
