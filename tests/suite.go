package tests

import (
	"context"
	"encoding/json"
	"github.com/crewblade/banner-management-service/config"
	"github.com/crewblade/banner-management-service/internal/app"
	"github.com/crewblade/banner-management-service/internal/cache"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/user_banner"
	"github.com/crewblade/banner-management-service/internal/repo/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"strconv"
	"time"
)

type Suite struct {
	suite.Suite
	handler *chi.Mux
	cache   *cache.BannerCacheImpl
	repo    *postgres.Storage
	ctx     context.Context
}

const configPath = "config/config.yaml"

func (s *Suite) SetupSuite() {
	err := godotenv.Load(".env")
	if err != nil {
		s.FailNow("failed to load .env file: %v", err)
	}

	cfg, err := config.NewConfig(configPath)

	if err != nil {
		s.FailNow("failed to init config: %v", err)
	}

	storage, err := postgres.New(cfg.PG.URL)
	if err != nil {
		s.FailNow("failed to create DB connection: %v", err)
	}

	cache := cache.NewBannerCacheImpl(5*time.Minute, 10*time.Minute)
	log := app.SetupLogger(cfg.Level)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Get("/user_banner", user_banner.GetUserBanner(log, storage, storage, cache))

	s.handler = router
	s.cache = cache
	s.repo = storage
	s.ctx = context.Background()

	if err = s.loadDbData(); err != nil {
		s.FailNow("Failed to load DB data", err)
	}

	if err = s.loadCacheData(); err != nil {
		s.FailNow("Failed to load cache data", err)
	}

}

func (s *Suite) loadDbData() error {
	isActive := true
	for i := 1; i < 10; i++ {
		b := models.Banner{
			TagIDs:    []int{i * 2, i * 3, i * 4},
			FeatureID: i,
			Content:   json.RawMessage(`{"title":"some_title` + strconv.Itoa(i) + `","text":"some_text` + strconv.Itoa(i) + `","url":"some_url` + strconv.Itoa(i) + `"}`),
			IsActive:  isActive,
		}
		_, err := s.repo.SaveBanner(s.ctx, b.TagIDs, b.FeatureID, b.Content, b.IsActive)
		if err != nil {
			return err
		}
		isActive = !isActive
	}
	return nil
}

func (s *Suite) loadCacheData() error {
	var active = true
	for i := 1; i < 10; i++ {
		b := models.Banner{
			TagIDs:    []int{i * 2, i * 3, i * 4},
			FeatureID: i,
			Content:   json.RawMessage(`{"title":"some_title` + strconv.Itoa(i) + `","text":"some_text` + strconv.Itoa(i) + `","url":"some_url` + strconv.Itoa(i) + `"}`),
			IsActive:  active,
		}
		for _, id := range b.TagIDs {
			err := s.cache.SetBanner(s.ctx, id, b.FeatureID, &models.BannerForUser{b.Content, b.IsActive})
			if err != nil {
				return err
			}
		}
		active = !active
	}
	return nil
}
