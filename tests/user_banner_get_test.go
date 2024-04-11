package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/crewblade/banner-management-service/internal/app"
	"github.com/crewblade/banner-management-service/internal/cache"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/httpserver/handlers/user_banner"
	"github.com/crewblade/banner-management-service/internal/repo/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
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

	err := godotenv.Load("../.env")
	if err != nil {
		s.FailNow("failed reading .env file %w", err)
	}
	storagePath := os.Getenv("PG_URL")
	fmt.Println("storagePath:", storagePath)
	storage, err := postgres.New(storagePath)
	if err != nil {
		s.FailNow("failed to create DB connection: %w", err)
	}

	cache := cache.NewBannerCacheImpl(5*time.Minute, 10*time.Minute)
	log := app.SetupLogger("local")

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Get("/user_banner", user_banner.GetUserBanner(log, storage, storage, cache))

	s.handler = router
	s.cache = cache
	s.repo = storage
	s.ctx = context.Background()

	if err = s.clearData(); err != nil {
		s.FailNow("Failed to clear DB data", err)
	}

	if err = s.loadDbData(); err != nil {
		s.FailNow("Failed to load DB data", err)
	}

	if err = s.loadCacheData(); err != nil {
		s.FailNow("Failed to load cache data", err)
	}

}

func (s *Suite) clearData() error {
	err := s.repo.ClearData(s.ctx)
	if err != nil {
		return fmt.Errorf("failed to clear user data: %w", err)
	}
	return nil
}

// every 2nd banner is not active
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

func (s *Suite) TestHappyGetUserBannerFromDB() {
	req := httptest.NewRequest("GET", "/user_banner?tag_id=2&feature_id=1&use_last_revision=true", nil)
	req.Header.Set("token", "admin_token")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	content := response["content"].(map[string]interface{})

	r.Equal(200, int(response["status"].(float64)))
	r.Equal("some_title1", content["title"])
	r.Equal("some_text1", content["text"])
	r.Equal("some_url1", content["url"])
}

func (s *Suite) TestGetUserBannerNotAuthorized() {
	req := httptest.NewRequest("GET", "/user_banner?tag_id=4&feature_id=2&use_last_revision=true", nil)
	req.Header.Set("token", "unknown_token")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(http.StatusUnauthorized, int(response["status"].(float64)))
}

func (s *Suite) TestGetUserBannerHaveNoAccess() {
	req := httptest.NewRequest("GET", "/user_banner?tag_id=4&feature_id=2&use_last_revision=true", nil)
	req.Header.Set("token", "user_token")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(http.StatusForbidden, int(response["status"].(float64)))
}

func (s *Suite) TestHappyGetUserBannerFromCache() {
	req := httptest.NewRequest("GET", "/user_banner?tag_id=2&feature_id=1&use_last_revision=false", nil)
	req.Header.Set("token", "admin_token")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	r.NoError(err)

	content := response["content"].(map[string]interface{})
	r.Equal(http.StatusOK, int(response["status"].(float64)))
	r.Equal("some_title1", content["title"])
	r.Equal("some_text1", content["text"])
	r.Equal("some_url1", content["url"])
}

func (s *Suite) TestGetUserBannerNotFoundFeature() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=1&feature_id=0&use_last_revision=true", nil)
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()

	s.handler.ServeHTTP(w, req)

	r := s.Require()
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	r.NoError(err)

	r.Equal(http.StatusNotFound, int(response["status"].(float64)))
}

func (s *Suite) TestGetUserBannerNotFoundTag() {

	req := httptest.NewRequest("GET", "/user_banner?tag_id=0&feature_id=1&use_last_revision=true", nil)
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()

	s.handler.ServeHTTP(w, req)

	r := s.Require()
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	r.NoError(err)

	r.Equal(http.StatusNotFound, int(response["status"].(float64)))
}

func (s *Suite) TestGetUserBannerBadRequest() {
	req := httptest.NewRequest("GET", "/user_banner?tag_id=bad&feature_id=bad&use_last_revision=bad", nil)
	req.Header.Set("token", "admin_token")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	r.NoError(err)

	r.Equal(http.StatusBadRequest, int(response["status"].(float64))) // Проверяем статус код в JSON-ответе
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
