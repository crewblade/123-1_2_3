package handlers

import (
	"github.com/crewblade/banner-management-service/internal/cache"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/banner"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/banner_id"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/user_banner"
	"github.com/crewblade/banner-management-service/internal/repo/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func NewRouter(log *slog.Logger, storage *postgres.Storage, cache *cache.BannerCacheImpl) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	router.Get("/user_banner", user_banner.GetUserBanner(log, storage, storage, cache))
	router.Get("/banner", banner.GetBanners(log, storage, storage, cache))
	router.Post("/banner", banner.SaveBanner(log, storage, storage))
	router.Patch("/banner/{id}", banner_id.UpdateBanner(log, storage, storage))
	router.Delete("/banner/{id}", banner_id.DeleteBanner(log, storage, storage))

	return router
}
