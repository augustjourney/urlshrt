package app

import (
	"database/sql"

	"github.com/augustjourney/urlshrt/internal/middleware"
	"github.com/gofiber/fiber/v2"
)

type Controller interface {
	BadRequest(ctx *fiber.Ctx) error
	CreateURL(ctx *fiber.Ctx) error
	APICreateURL(ctx *fiber.Ctx) error
	GetURL(ctx *fiber.Ctx) error
	APICreateURLBatch(ctx *fiber.Ctx) error
	GetUserURLs(ctx *fiber.Ctx) error
}

// Временное решение прокидывать db в контроллер
// Позже db будет прокидываться только в repo
func New(c Controller, db *sql.DB) *fiber.App {
	app := fiber.New()

	app.Use(middleware.RequestCompress)
	app.Use(middleware.RequestLogger)

	app.Get("/ping", func(ctx *fiber.Ctx) error {
		err := db.Ping()
		if err != nil {
			return ctx.SendStatus(fiber.StatusInternalServerError)
		}
		return ctx.SendStatus(fiber.StatusOK)
	})

	app.Post("/", c.CreateURL)
	app.Post("/api/shorten", c.APICreateURL)
	app.Post("/api/shorten/batch", c.APICreateURLBatch)
	app.Get("/:short", c.GetURL)
	app.Get("/api/user/urls", c.GetUserURLs)
	app.Use("/*", c.BadRequest)

	return app
}
