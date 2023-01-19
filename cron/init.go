package cron

import "github.com/robfig/cron/v3"

var (
	c *cron.Cron
)

func init() {
	RefreshHoliday()
	SyncUser()
	c = cron.New()
	_, err := c.AddFunc("@every 1m", remind)
	if err != nil {
		panic(err)
	}
	_, err = c.AddFunc("@weekly", RefreshHoliday)
	if err != nil {
		panic(err)
	}
}

func Start() {
	c.Start()
}
