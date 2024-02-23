package logger

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

func New() *logrus.Logger {
	Log = logrus.New()
	Log.SetLevel(logrus.InfoLevel)
	return Log
}

func RequestLogger(ctx *fiber.Ctx) error {
	start := time.Now()
	result := ctx.Next()
	duration := time.Since(start).Milliseconds()
	log := fmt.Sprintf("HTTP Request – Method: %s, Path: %s, Duration: %dms, Status: %d, ContentLength: %d", ctx.Method(), ctx.OriginalURL(), duration, ctx.Response().StatusCode(), len(ctx.Response().Body()))
	Log.Info(log)
	return result
}
