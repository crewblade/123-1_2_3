package app

import (
	"github.com/crewblade/banner-management-service/config"
	"github.com/crewblade/banner-management-service/internal/lib/logger/sl"
	"github.com/crewblade/banner-management-service/internal/storage/postgres"
	"log"
	"log/slog"
)

func Run(configPath string) {
	const op = "internal.app.Run"
	cfg, err := config.NewConfig(configPath)

	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	log := SetupLogger(cfg.Log.Level)
	log.Info("starting app", slog.Any("cfg", cfg))

	log.Info("Initializing postgres...")
	storage, err := postgres.New(cfg.PG.URL)
	if err != nil {
		log.Error("failed to init postgres", sl.Err(err))
	}

	log.Info("Initializing handlers and routes...")
	router := v1.New
}
