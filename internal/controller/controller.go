package controller

import (
	"net/http"

	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service service.IService
}

func (c *Controller) BadRequest(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/plain")
	return ctx.SendStatus(400)
}

func (c *Controller) CreateURL(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/plain")
	if ctx.Method() != http.MethodPost {
		return ctx.SendStatus(400)
	}
	// Getting url from text plain body

	originalURL := string(ctx.Body())

	if originalURL == "" {
		return ctx.SendStatus(400)
	}

	// Make a short url
	short := c.service.Shorten(originalURL)

	// Response
	return ctx.Status(201).SendString(short)
}

func (c *Controller) GetURL(ctx *fiber.Ctx) error {

	ctx.Set("Content-type", "text/plain")

	if ctx.Method() != http.MethodGet {
		return ctx.SendStatus(400)
	}

	// Parse short url
	short := ctx.Params("short")

	// Find original
	originalURL, err := c.service.FindOriginal(short)

	if err != nil || originalURL == "" {
		return ctx.SendStatus(400)
	}

	// Response
	ctx.Location(originalURL)
	return ctx.Status(307).SendString(originalURL)
}

func New(service service.IService) Controller {
	return Controller{
		service: service,
	}
}
