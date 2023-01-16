package cron

import "github.com/robfig/cron/v3"

var (
	c *cron.Cron
)

func init() {
	c = cron.New()
	_, err := c.AddFunc("@every 1m", remind)
	if err != nil {
		panic(err)
	}
	_, err = c.AddFunc("@monthly", refreshHoliday)
	if err != nil {
		panic(err)
	}
}

func Start() {
	c.Start()
}
