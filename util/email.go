package util

import (
	"context"
	"github.com/spf13/viper"
	"gopkg.in/gomail.v2"
)

func SendEmail(ctx context.Context, to string, subject, content string) error {
	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {viper.GetString("email.username")},
		"To":      {to},
		"Subject": {subject},
	})
	m.SetBody("text/html", content)
	return gomail.NewDialer(
		viper.GetString("email.server"),
		viper.GetInt("email.port"),
		viper.GetString("email.username"),
		viper.GetString("email.password")).DialAndSend(m)
}
