package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crewblade/banner-management-service/internal/cache"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers"
	"github.com/crewblade/banner-management-service/internal/lib/logger/handlers/slogempty"
	"github.com/crewblade/banner-management-service/internal/repo/postgres"
	"github.com/go-chi/chi/v5"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"os"
	"strconv"
	"testing"
	"time"
)

type Suite struct {
	suite.Suite
	handler   *chi.Mux
	cache     *cache.BannerCacheImpl
	repo      *postgres.Storage
	ctx       context.Context
	bannerIDs []int
}

const configPath = "config/config.yaml"

// const keyDB = "PG_URL_LOCALHOST"
const keyDB = "PG_URL_TEST"

//const keyDB = "PG_URL"

func (s *Suite) SetupSuite() {

	err := godotenv.Load("../.env")
	if err != nil {
		s.FailNow("failed reading .env file %w", err)
	}

	storagePath := os.Getenv(keyDB)
	fmt.Println("storagePath:", storagePath)
	storage, err := postgres.New(storagePath, 5, 10, uint64(100))
	if err != nil {
		s.FailNow("failed to create DB connection: %w", err)
	}

	cache := cache.NewBannerCacheImpl(5*time.Minute, 10*time.Minute)
	log := slogempty.NewEmptyLogger()

	router := handlers.NewRouter(log, storage, cache)

	s.handler = router
	s.cache = cache
	s.repo = storage
	s.ctx = context.Background()

	if err = s.repo.PrepareForTest(s.ctx); err != nil {
		s.Fail("Failed to prepare db for test", err)
	}
	if err = s.clearData(); err != nil {
		s.FailNow("Failed to clear DB data", err)
	}

	if err = s.loadDbData(); err != nil {
		s.FailNow("Failed to load DB data", err)
	}

}

func (s *Suite) clearData() error {
	err := s.repo.ClearData(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to clear user data: %w", err)
	}
	return nil
}

const (
	defaultAttempts = 20
	defaultTimeout  = time.Second
)

// every 2nd banner is not active
func (s *Suite) loadDbData() error {
	s.repo.ClearData(s.ctx)
	isActive := true
	bannerIDs := make([]int, 10)
	for i := 1; i < 10; i++ {
		b := models.Banner{
			TagIDs:    []int{i * 2, i * 3, i * 4},
			FeatureID: i,
			Content:   json.RawMessage(`{"title":"some_title` + strconv.Itoa(i) + `","text":"some_text` + strconv.Itoa(i) + `","url":"some_url` + strconv.Itoa(i) + `"}`),
			IsActive:  isActive,
		}
		id, err := s.repo.SaveBanner(s.ctx, b.TagIDs, b.FeatureID, b.Content, b.IsActive)
		bannerIDs[i] = id
		if err != nil {
			return err
		}
		isActive = !isActive

	}
	adminToken := "admin_token"
	userToken := "user_token"

	err := s.repo.AddUser(s.ctx, adminToken, true)
	if err != nil {
		return err
	}
	err = s.repo.AddUser(s.ctx, userToken, false)
	if err != nil {
		return err
	}

	s.bannerIDs = bannerIDs

	return nil
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
