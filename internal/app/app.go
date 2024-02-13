package app

import (
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	BadRequest(ctx *fiber.Ctx) error
	CreateURL(ctx *fiber.Ctx) error
	GetURL(ctx *fiber.Ctx) error
}

func New(c Controller) *fiber.App {
	app := fiber.New()

	app.Post("/", c.CreateURL)
	app.Get("/:short", c.GetURL)
	app.Use("/*", c.BadRequest)

	return app
}
