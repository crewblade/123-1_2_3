package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (s *Suite) TestGetUserBannerWithFlag() {
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

func (s *Suite) TestGetUserBannerWithoutFlag() {
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
