package controller

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

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

	user, err := c.checkAuth(ctx, true)

	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	originalURL := string(ctx.Body())

	if originalURL == "" {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	result, err := c.service.Shorten(originalURL, user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Если уже существует такой url
	// То возвращаем url и статус 409
	if result.AlreadyExists {
		return ctx.Status(http.StatusConflict).SendString(result.ResultURL)
	}

	return ctx.Status(http.StatusCreated).SendString(result.ResultURL)
}

func (c *Controller) APICreateURLBatch(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "application/json")

	if ctx.Method() != http.MethodPost {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	user, err := c.checkAuth(ctx, true)

	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	var body []service.BatchURL

	err = json.Unmarshal(ctx.Body(), &body)

	if err != nil || len(body) == 0 {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	result, err := c.service.ShortenBatch(body, user)

	if err != nil {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	response, err := json.Marshal(result)

	if err != nil {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	return ctx.Status(http.StatusCreated).Send(response)
}

func (c *Controller) APIDeleteBatch(ctx *fiber.Ctx) error {
	if ctx.Method() != http.MethodDelete {
		return ctx.SendStatus(http.StatusMethodNotAllowed)
	}

	user, _ := c.checkAuth(ctx, false)

	if user == "" {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	var err error

	var shortIds []string

	err = json.Unmarshal(ctx.Body(), &shortIds)

	if err != nil {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	rctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = c.service.DeleteBatch(rctx, shortIds, user)

	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.SendStatus(http.StatusAccepted)
}

func (c *Controller) APICreateURL(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "application/json")

	if ctx.Method() != http.MethodPost {
		return ctx.SendStatus(http.StatusBadRequest)
	}

	user, err := c.checkAuth(ctx, true)

	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	var body APICreateURLBody

	err = json.Unmarshal(ctx.Body(), &body)

	if err != nil || body.URL == "" {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Make a short url
	result, err := c.service.Shorten(body.URL, user)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	response, err := json.Marshal(APICreateURLResult{
		Result: result.ResultURL,
	})

	if err != nil {
		logger.Log.Error(err)
		return ctx.SendStatus(http.StatusBadRequest)
	}

	// Если уже существует такой url
	// То возвращаем url и статус 409
	if result.AlreadyExists {
		return ctx.Status(http.StatusConflict).Send(response)
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

	// TODO: наверное, будет лучше вынести эти ошибки из сервиса
	// Куда-то в отдельный модуль со всеми ошибками
	if errors.Is(err, service.ErrIsDeleted) {
		return ctx.SendStatus(http.StatusGone)
	}

	if errors.Is(err, service.ErrNotFound) {
		// should be 404
		return ctx.SendStatus(http.StatusBadRequest)
	}

	if err != nil {
		return ctx.Status(http.StatusInternalServerError).SendString(err.Error())
	}

	// Response
	ctx.Location(originalURL)
	return ctx.Status(http.StatusTemporaryRedirect).SendString(originalURL)
}

func (c *Controller) GetUserURLs(ctx *fiber.Ctx) error {
	ctx.Set("Content-type", "application/json")

	if ctx.Method() != http.MethodGet {
		return ctx.SendStatus(http.StatusMethodNotAllowed)
	}

	user, _ := c.checkAuth(ctx, false)

	if user == "" {
		return ctx.SendStatus(http.StatusUnauthorized)
	}

	urls, err := c.service.GetUserURLs(context.Background(), user)

	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	if urls == nil || len(*urls) == 0 {
		return ctx.SendStatus(http.StatusNoContent)
	}

	response, err := json.Marshal(urls)

	if err != nil {
		return ctx.SendStatus(http.StatusInternalServerError)
	}

	return ctx.Status(http.StatusOK).Send(response)

}

func (c *Controller) checkAuth(ctx *fiber.Ctx, createIfEmpty bool) (string, error) {
	// ID пользователя может храниться
	// Либо в заголовке Authorization
	// Либо в куке user
	user := ctx.Get("Authorization")
	cookie := ctx.Cookies("user")

	if cookie != "" && user == "" {
		user = cookie
	}

	if createIfEmpty && user == "" {
		user, err := c.service.GenerateID()
		if err != nil {
			return "", err
		}

		cookie := new(fiber.Cookie)
		cookie.Name = "user"
		cookie.Value = user

		ctx.Cookie(cookie)

		ctx.Set("Authorization", user)
		return user, nil
	}

	return user, nil
}

func New(service service.IService) Controller {
	return Controller{
		service: service,
	}
}
