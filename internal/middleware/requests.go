package middleware

import (
	"fmt"
	"time"

	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/gofiber/fiber/v2"
)

func RequestLogger(ctx *fiber.Ctx) error {
	start := time.Now()
	result := ctx.Next()
	duration := time.Since(start).Milliseconds()
	log := fmt.Sprintf("HTTP Request – Method: %s, Path: %s, Duration: %dms, Status: %d, ContentLength: %d", ctx.Method(), ctx.OriginalURL(), duration, ctx.Response().StatusCode(), len(ctx.Response().Body()))
	logger.Log.Info(log)
	return result
}
