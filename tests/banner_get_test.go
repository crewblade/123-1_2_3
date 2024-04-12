package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (s *Suite) TestGetBanner() {
	req := httptest.NewRequest("GET", "/banner?tag_id=4&limit=3&offset=1", nil)
	req.Header.Set("token", "admin_token")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(200, int(response["status"].(float64)))
}

func (s *Suite) TestGetBannerNotAuthorized() {
	req := httptest.NewRequest("GET", "/banner?tag_id=4&limit=3&offset=1", nil)
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

func (s *Suite) TestGetBannerHaveNoAccess() {
	req := httptest.NewRequest("GET", "/banner?tag_id=4&limit=3&offset=1", nil)
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

func (s *Suite) TestGetBannerNotFound() {
	req := httptest.NewRequest("GET", "/banner?tag_id=4924&feature_id=90909090&limit=3&offset=1", nil)
	req.Header.Set("token", "admin_token")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(http.StatusNotFound, int(response["status"].(float64)))
}

func (s *Suite) TestGetBannerBadRequest() {
	req := httptest.NewRequest("GET", "/banner?tag_id=4abc2&feature_id=9ff0&limit=3&offset=1", nil)
	req.Header.Set("token", "admin_token")
	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(http.StatusBadRequest, int(response["status"].(float64)))
}
