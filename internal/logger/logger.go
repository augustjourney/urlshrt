package logger

import "github.com/sirupsen/logrus"

var Log *logrus.Logger

func New() *logrus.Logger {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)
	return Log
}
