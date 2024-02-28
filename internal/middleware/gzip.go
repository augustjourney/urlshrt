package middleware

import (
	"bytes"
	"compress/gzip"
	"strings"

	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/gofiber/fiber/v2"
)

func RequestCompress(ctx *fiber.Ctx) error {
	result := ctx.Next()

	supportsGzip := strings.Contains(ctx.Get(fiber.HeaderAcceptEncoding), "gzip")

	if !supportsGzip {
		return result
	}

	contentType := string(ctx.Response().Header.ContentType())

	if contentType != "application/json" && contentType != "plain/html" {
		return result
	}

	ctx.Response().Header.Set(fiber.HeaderContentEncoding, "gzip")

	body := ctx.Response().Body()

	var bodyBytes bytes.Buffer

	gzipWriter, err := gzip.NewWriterLevel(&bodyBytes, gzip.BestSpeed)

	if err != nil {
		logger.Log.Error("Err1 ", err)
		return err
	}

	if _, err = gzipWriter.Write(body); err != nil {
		return err
	}

	if err = gzipWriter.Close(); err != nil {
		return err
	}

	ctx.Response().Header.Set("Content-Encoding", "gzip")

	ctx.Response().SetBody(bodyBytes.Bytes())

	return result
}
