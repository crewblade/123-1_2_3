package handlers

import (
	"github.com/crewblade/banner-management-service/internal/cache"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/banner"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/banner_id"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/user_banner"
	mwauth "github.com/crewblade/banner-management-service/internal/httpserver/middleware/auth"
	"github.com/crewblade/banner-management-service/internal/repo/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
)

func NewRouter(log *slog.Logger, storage *postgres.Storage, cache *cache.BannerCacheImpl) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(mwauth.AuthMiddleware(log, storage))

	router.Get("/user_banner", user_banner.GetUserBanner(log, storage, cache))
	router.Get("/banner", banner.GetBanners(log, storage, cache))
	router.Post("/banner", banner.SaveBanner(log, storage))
	router.Delete("/banner", banner.DeleteBanners(log, storage))
	router.Patch("/banner/{id}", banner_id.UpdateBanner(log, storage))
	router.Delete("/banner/{id}", banner_id.DeleteBanner(log, storage))

	return router
}
