package main

import (
	"fmt"

	"github.com/augustjourney/urlshrt/internal/app"
	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/controller"
	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/augustjourney/urlshrt/internal/storage/inmemory"
)

func main() {

	config := config.New()

	repo := inmemory.New()
	service := service.New(&repo, config)
	c := controller.New(&service)
	server := app.New(&c)

	server.Listen(fmt.Sprintf(":%s", config.Port))
}
