package cron

import (
	"context"
	"github.com/spf13/viper"
	"todo-reminder/gocq"
)

func init() {
	registerCronTask("0 */6 * * *", Ping, false)
}

func Ping() {
	admin := viper.GetString("admin.qq")
	gocq.GetGocqInstance().SendPrivateStringMessage(context.Background(), "ping message", admin)
}
