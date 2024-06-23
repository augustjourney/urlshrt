package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

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
	server := app.New(&controller, db)

	go func() {
		if config.EnableHTTPS {
			err = app.RunHTTPS(server, config)
			// Если происходит ошибка — просто логируем ее
			// И запускаем на http
			if err != nil {
				logger.Log.Fatal(err)
			}
		}

		err = app.RunHTTP(server, config)
		if err != nil {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	signal.Notify(ch, os.Interrupt, syscall.SIGINT)
	signal.Notify(ch, os.Interrupt, syscall.SIGQUIT)

	<-ch

	logger.Log.Info("Gracefully shutting down...")

	server.Shutdown()

	logger.Log.Info("Closing connections...")

	db.Close()

	logger.Log.Info("Server was shutdown successfully")
}
