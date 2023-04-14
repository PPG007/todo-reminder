package cron

import (
	"github.com/robfig/cron/v3"
	"todo-reminder/util"
)

var (
	c     *cron.Cron
	tasks []cronTask
)

type cronTask struct {
	fn             func()
	spec           string
	needRunAtStart bool
}

func registerCronTask(spec string, fn func(), needRunAtStart bool) {
	tasks = append(tasks, cronTask{
		fn:             util.FuncWithRecovery(fn),
		spec:           spec,
		needRunAtStart: needRunAtStart,
	})
}

func Start() {
	c = cron.New()
	for _, task := range tasks {
		_, err := c.AddFunc(task.spec, task.fn)
		if err != nil {
			panic(err)
		}
		if task.needRunAtStart {
			go task.fn()
		}
	}
	c.Start()
}
