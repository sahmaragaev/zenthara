package main

import (
	"log"
	"os"

	"zenthara/internal/app"
	"zenthara/internal/config"
	"zenthara/internal/shared/logger"
)

func main() {
	cfg, err := config.Load("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log := logger.NewLogger(logger.Config{
		Environment: cfg.App.Environment,
		LogLevel:    cfg.Logger.Level,
		Format:      cfg.Logger.Format,
	})

	application := app.New(cfg, log)
	if err := application.Setup(); err != nil {
		log.Fatal().Err(err).Msg("Failed to setup application")
		os.Exit(1)
	}

	if err := application.Run(); err != nil {
		log.Fatal().Err(err).Msg("Application error")
		os.Exit(1)
	}
}
