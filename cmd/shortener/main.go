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
	controller := controller.New(&service)
	app := app.New(&controller)

	fmt.Printf("Launching on %s", config.ServerAddress)

	err := app.Listen(config.ServerAddress)
	if err != nil {
		panic(err)
	}
}
