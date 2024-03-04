package app

import (
	"github.com/augustjourney/urlshrt/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	BadRequest(ctx *fiber.Ctx) error
	CreateURL(ctx *fiber.Ctx) error
	APICreateURL(ctx *fiber.Ctx) error
	GetURL(ctx *fiber.Ctx) error
}

func New(c Controller) *fiber.App {
	app := fiber.New()

	app.Use(middleware.RequestCompress)
	app.Use(middleware.RequestLogger)

	app.Post("/", c.CreateURL)
	app.Post("/api/shorten", c.APICreateURL)
	app.Get("/:short", c.GetURL)
	app.Use("/*", c.BadRequest)

	return app
}
