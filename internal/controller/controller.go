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

type APICreateURLBody struct {
	URL string `json:"url"`
}

type APICreateURLResult struct {
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
	short, err := c.service.Shorten(originalURL)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Response
	return ctx.Status(http.StatusCreated).SendString(short)
}

func (c *Controller) APICreateURL(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "application/json")

	if ctx.Method() != http.MethodPost {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	var body APICreateURLBody

	err := json.Unmarshal(ctx.Body(), &body)

	if err != nil || body.URL == "" {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Make a short url
	result, err := c.service.Shorten(body.URL)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	response, err := json.Marshal(APICreateURLResult{
		Result: result,
	})

	if err != nil {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	return ctx.Status(http.StatusCreated).Send(response)
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
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	if originalURL == "" {
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
