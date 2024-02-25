package controller

import (
	"encoding/json"
	"net/http"

	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	service service.IService
}

type ApiCreateURLBody struct {
	URL string `json:"url"`
}

type ApiCreateURLResult struct {
	Result string `json:"result"`
}

func (c *Controller) BadRequest(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/plain")
	return ctx.SendStatus(http.StatusBadRequest)
}

func (c *Controller) CreateURL(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "text/plain")
	if ctx.Method() != http.MethodPost {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	// Getting url from text plain body

	originalURL := string(ctx.Body())

	if originalURL == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Make a short url
	short := c.service.Shorten(originalURL)

	// Response
	return ctx.Status(201).SendString(short)
}

func (c *Controller) ApiCreateURL(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "application/json")
	if ctx.Method() != http.MethodPost {
		return ctx.SendStatus(http.StatusBadRequest)
	}
	// Getting url from json body

	var body ApiCreateURLBody

	err := json.Unmarshal(ctx.Body(), &body)

	if err != nil || body.URL == "" {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Make a short url
	short := c.service.Shorten(body.URL)

	resp, err := json.Marshal(ApiCreateURLResult{
		Result: short,
	})

	if err != nil {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Response
	return ctx.Status(201).Send(resp)
}

func (c *Controller) GetURL(ctx *fiber.Ctx) error {

	ctx.Set("Content-type", "text/plain")

	if ctx.Method() != http.MethodGet {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Parse short url
	short := ctx.Params("short")

	// Find original
	originalURL, err := c.service.FindOriginal(short)

	if err != nil || originalURL == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Response
	ctx.Location(originalURL)
	return ctx.Status(http.StatusTemporaryRedirect).SendString(originalURL)
}

func New(service service.IService) Controller {
	return Controller{
		service: service,
	}
}
