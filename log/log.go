package log

import (
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	//logrus.SetFormatter(&logrus.TextFormatter{
	//	ForceColors:               true,
	//	DisableColors:             false,
	//	ForceQuote:                false,
	//	DisableQuote:              false,
	//	EnvironmentOverrideColors: false,
	//	DisableTimestamp:          false,
	//	FullTimestamp:             false,
	//})
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.WarnLevel)
}

func Warn(message string, fields logrus.Fields) {
	logrus.WithFields(fields).Warn(message)
}

func WarnTrace(msg string, fields logrus.Fields, trace []byte) {
	logrus.WithFields(fields).WithField("backtrace", string(trace)).Warn(msg)
}

func ErrorTrace(msg string, fields logrus.Fields, trace []byte) {
	logrus.WithFields(fields).WithField("backtrace", string(trace)).Warn(msg)
}
