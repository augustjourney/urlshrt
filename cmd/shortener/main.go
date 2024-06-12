package main

import (
	"context"
	"fmt"
	"github.com/augustjourney/urlshrt/internal/app"
	"github.com/augustjourney/urlshrt/internal/config"
	"github.com/augustjourney/urlshrt/internal/controller"
	"github.com/augustjourney/urlshrt/internal/infra"
	"github.com/augustjourney/urlshrt/internal/logger"
	"github.com/augustjourney/urlshrt/internal/service"
	"github.com/augustjourney/urlshrt/internal/storage"
	"github.com/augustjourney/urlshrt/internal/storage/infile"
	"github.com/augustjourney/urlshrt/internal/storage/postgres"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	config := config.New()
	logger.New()

	logger.Log.Printf("Build version: %v\n", buildVersion)
	logger.Log.Printf("Build date: %v\n", buildDate)
	logger.Log.Printf("Build commit: %v\n", buildCommit)

	db, err := infra.InitPostgres(config)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	var repo storage.IRepo

	repo, err = postgres.New(context.Background(), db)
	if err != nil {
		logger.Log.Error("Could not connect to postgres, using in-file storage")
		repo = infile.New(config)
	}
	service := service.New(repo, config)
	controller := controller.New(&service)
	app := app.New(&controller, db)

	logger.Log.Info(fmt.Sprintf("Launching on %s", config.ServerAddress))

	err = app.Listen(config.ServerAddress)
	if err != nil {
		panic(err)
	}
}
