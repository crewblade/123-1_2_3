package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
)

func (s *Suite) TestDeleteBannerByID() {

	url := "/banner/" + strconv.Itoa(s.bannerIDs[3])
	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Set("token", "admin_token")

	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(http.StatusOK, int(response["status"].(float64)))
}

func (s *Suite) TestDeleteBannerByIDNotFound() {

	url := "/banner/0"

	req := httptest.NewRequest("DELETE", url, nil)
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

func (s *Suite) TestDeleteBannerByIDrBadRequest() {

	url := "/banner/abc123"
	req := httptest.NewRequest("DELETE", url, nil)
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

func (s *Suite) TestDeleteBannerByIDNoAccess() {

	url := "/banner/" + strconv.Itoa(s.bannerIDs[9])
	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Set("token", "user_token")

	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	r.NoError(err)

	r.Equal(http.StatusForbidden, int(response["status"].(float64)))
}

func (s *Suite) TestDeleteBannerByIDNotAuthorized() {

	url := "/banner/" + strconv.Itoa(s.bannerIDs[9])
	req := httptest.NewRequest("DELETE", url, nil)
	req.Header.Set("token", "unknown_token")

	w := httptest.NewRecorder()

	s.handler.ServeHTTP(w, req)

	r := s.Require()
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	r.NoError(err)

	r.Equal(http.StatusUnauthorized, int(response["status"].(float64)))
}
