package logger

import (
	"github.com/sirupsen/logrus"
)

// Логгер
var Log *logrus.Logger

// Создает новый экземпляр логгера
func New() *logrus.Logger {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)
	return Log
}
