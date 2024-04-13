package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (s *Suite) TestSaveUserBanner() {

	requestBody := []byte(`{
        "tag_ids": [123, 456],
        "feature_id": 789,
        "content": {
            "title": "test title",
            "text": "test text",
            "url": "test url"
        },
        "is_active": true
    }`)
	req := httptest.NewRequest("POST", "/banner", bytes.NewBuffer(requestBody))
	req.Header.Set("token", "admin_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(s.bannerIDs[9]+1, int(response["banner_id"].(float64)))
	r.Equal(http.StatusCreated, int(response["status"].(float64)))
}

func (s *Suite) TestSaveUserBannerAlreadyExists() {
	requestBody := []byte(`{
        "tag_ids": [1, 2, 3, 4],
        "feature_id": 2,
        "content": {
            "title": "test title",
            "text": "test text",
            "url": "test url"
        },
        "is_active": true
    }`)
	req := httptest.NewRequest("POST", "/banner", bytes.NewBuffer(requestBody))
	req.Header.Set("token", "admin_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(http.StatusInternalServerError, int(response["status"].(float64)))
}

func (s *Suite) TestSaveUserBannerBadRequest() {
	requestBody := []byte(`{
        "tag_ids": asd123,
        "feature_id": 123sdf,
        "content": {
            "title": "test title",
            "text": "test text",
            "url": "test url"
        },
        "is_active": t?rue
    }`)
	req := httptest.NewRequest("POST", "/banner", bytes.NewBuffer(requestBody))
	req.Header.Set("token", "admin_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(http.StatusBadRequest, int(response["status"].(float64)))
}

func (s *Suite) TestSaveUserBannerNoAccess() {
	requestBody := []byte(`{
        "tag_ids": [15, 25, 35, 45],
        "feature_id": 2212,
        "content": {
            "title": "test title",
            "text": "test text",
            "url": "test url"
        },
        "is_active": true
    }`)
	req := httptest.NewRequest("POST", "/banner", bytes.NewBuffer(requestBody))
	req.Header.Set("token", "user_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	r.NoError(err)

	r.Equal(http.StatusForbidden, int(response["status"].(float64)))
}

func (s *Suite) TestSaveUserBannerNotAuthorized() {

	requestBody := []byte(`{
        "tag_ids": [11, 22, 33, 44],
        "feature_id": 2020,
        "content": {
            "title": "test title",
            "text": "test text",
            "url": "test url"
        },
        "is_active": true
    }`)
	req := httptest.NewRequest("POST", "/banner", bytes.NewBuffer(requestBody))
	req.Header.Set("token", "unknown_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	s.handler.ServeHTTP(w, req)

	r := s.Require()
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	r.NoError(err)

	r.Equal(http.StatusUnauthorized, int(response["status"].(float64)))
}
