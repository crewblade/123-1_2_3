package app

import (
	"fmt"
	"github.com/crewblade/banner-management-service/config"
	"github.com/crewblade/banner-management-service/internal/httpserver"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/banner"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/banner_id"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/user_banner"
	"github.com/crewblade/banner-management-service/internal/lib/logger/sl"
	"github.com/crewblade/banner-management-service/internal/repo/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func Run(configPath string) {
	cfg, err := config.NewConfig(configPath)
	const op = "internal.app.Run"
	if err != nil {
		fmt.Printf("Config errs: %s", err)
		os.Exit(1)
	}
	log := SetupLogger(cfg.Log.Level)
	log.Info("Starting app", slog.Any("cfg", cfg))

	log.Info("Initializing postgres...")
	storage, err := postgres.New(cfg.PG.URL)
	if err != nil {
		log.Error("failed to init postgres", sl.Err(err))
		os.Exit(1)
	}

	log.Info("Initializing handlers and routes...")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Get("/user_banner", user_banner.GetUserBanner(log, storage))
	router.Get("/banner", banner.GetBanners(log, storage))
	router.Post("/banner", banner.SaveBanner(log, storage))
	router.Patch("/banner/{id}", banner_id.UpdateBanner(log, storage))
	router.Delete("/banner/{id}", banner_id.DeleteBanner(log, storage))

	log.Info("Starting http server...", slog.String("addr", cfg.Addr))

	httpServer := httpserver.New(
		router,
		cfg.Timeout,
		cfg.Timeout,
		cfg.ShutdownTimeout,
		cfg.Addr,
	)

	log.Info("Configuring graceful shutdown...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Info(op + s.String())
	case err = <-httpServer.Notify():
		log.Error(op, sl.Err(err))
	}

	// Graceful shutdown
	log.Info("Shutting down...")
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(op, sl.Err(err))
	}

}
