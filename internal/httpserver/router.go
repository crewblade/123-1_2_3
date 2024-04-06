package httpserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)

	//router.Get("/user_banner", GetUserBanner)
	//router.Get("/banner", GetBanners)
	//router.Post("/banner", CreateBanner)
	//router.Patch("/banner/{id}", UpdateBanner)
	//router.Delete("/banner/{id}", DeleteBanner)
	return router
}
