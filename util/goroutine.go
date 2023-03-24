package util

import (
	"fmt"
	"github.com/panjf2000/ants/v2"
	"github.com/sirupsen/logrus"
	"runtime"
	"todo-reminder/log"
)

var (
	pool *ants.Pool
)

func init() {
	newPool, err := ants.NewPool(100, ants.WithPanicHandler(func(err interface{}) {
		stack := make([]byte, 4096)
		stack = stack[:runtime.Stack(stack, false)]
		log.WarnTrace("Panic in goroutine", logrus.Fields{
			"error": fmt.Sprintf("%v", err),
		}, stack)
	}))
	if err != nil {
		panic(err)
	}
	pool = newPool
}

func Submit(f func()) error {
	return pool.Submit(f)
}
