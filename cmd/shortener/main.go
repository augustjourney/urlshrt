package main

import (
	"github.com/augustjourney/urlshrt/internal/app"
	"github.com/augustjourney/urlshrt/internal/controller"
	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/augustjourney/urlshrt/internal/storage/inmemory"
)

func main() {

	repo := inmemory.New()
	service := service.New(&repo)
	c := controller.New(&service)
	server := app.New(&c)

	server.Listen(":8080")
}
