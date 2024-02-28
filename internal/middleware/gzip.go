package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"

	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/gofiber/fiber/v2"
)

func RequestCompress(ctx *fiber.Ctx) error {

	// Processing request compress
	requestGzipped := strings.Contains(ctx.Get(fiber.HeaderContentEncoding), "gzip")

	if requestGzipped && ctx.Method() != fiber.MethodGet {
		bodyBytes := bytes.NewBuffer(ctx.Request().Body())

		gzipWriter, err := gzip.NewReader(bodyBytes)

		if err != nil {
			return err
		}

		defer gzipWriter.Close()

		body, err := io.ReadAll(gzipWriter)
		if err != nil {
			return err
		}

		ctx.Request().SetBody(body)
	}

	result := ctx.Next()

	// Processing response compress
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

	ctx.Response().SetBody(bodyBytes.Bytes())

	return result
}
