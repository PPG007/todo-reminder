package cron

import "github.com/robfig/cron/v3"

var (
	c *cron.Cron
)

func init() {
	//c = cron.New()
	//_, err := c.AddFunc("@every 1m", remind)
	//if err != nil {
	//	panic(err)
	//}
	//RefreshHoliday()
	//_, err = c.AddFunc("@every 1d", RefreshHoliday)
	//if err != nil {
	//	panic(err)
	//}
}

func Start() {
	c.Start()
}
