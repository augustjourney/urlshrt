package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// Middleware — который создает сжатие входящих http-запросов
func RequestCompress(ctx *fiber.Ctx) error {

	// Processing request compress
	requestGzipped := strings.Contains(ctx.Get(fiber.HeaderContentEncoding), "gzip")

	if requestGzipped && ctx.Method() != fiber.MethodGet {
		buf := bytes.NewBuffer(ctx.Request().Body())
		gzipWriter, err := gzip.NewReader(buf)

		if err != nil {
			if err != gzip.ErrHeader {
				return err
			}
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

		defer gzipWriter.Close()

		body, err := io.ReadAll(gzipWriter)
		if err != nil {
			return err
		}

		ctx.Request().Header.Del(fiber.HeaderContentEncoding)

		ctx.Request().SetBodyStream(io.NopCloser(bytes.NewReader(body)), len(body))
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

	body := ctx.Response().Body()

	var bodyBytes bytes.Buffer

	gzipWriter, err := gzip.NewWriterLevel(&bodyBytes, gzip.BestSpeed)

	if err != nil {
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
