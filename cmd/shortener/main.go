package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	urlService := service.New(repo, config)

	httpController := controller.NewHTTPController(&urlService)
	grpcController := controller.NewGrpcController(&urlService)

	httpServer := app.NewHTTPApp(httpController, db)
	grpcServer := app.NewGrpcApp(grpcController)

	go func() {
		if config.EnableHTTPS {
			err = app.RunHTTPS(httpServer, config)
			// Если происходит ошибка — просто логируем ее
			// И запускаем на http
			if err != nil {
				logger.Log.Fatal(err)
			}
		}

		err = app.RunHTTP(httpServer, config)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		err = app.RunGRPC(grpcServer, config)
		if err != nil {
			panic(err)
		}
	}()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	<-ch

	logger.Log.Info("Gracefully shutting down...")

	httpServer.ShutdownWithTimeout(10 * time.Second)
	grpcServer.GracefulStop()

	logger.Log.Info("Closing connections...")

	db.Close()

	logger.Log.Info("Server was shutdown successfully")
}
