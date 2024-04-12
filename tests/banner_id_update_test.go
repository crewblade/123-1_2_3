package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
)

func (s *Suite) TestBannerUpdate() {

	requestBody := []byte(`{
        "tag_ids": [1010, 5050],
        "feature_id": 999,
        "content": {
            "title": "new title",
            "text": "new text",
            "url": "new url"
        },
        "is_active": true
    }`)
	url := "/banner/" + strconv.Itoa(s.bannerIDs[4])
	req := httptest.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
	req.Header.Set("token", "admin_token")
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	s.handler.ServeHTTP(w, req)

	r := s.Require()

	responseBody := w.Body.Bytes()

	var response map[string]interface{}
	err := json.Unmarshal(responseBody, &response)
	r.NoError(err)

	r.Equal(http.StatusOK, int(response["status"].(float64)))
}

func (s *Suite) TestUpdateBannerAlreadyExists() {
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
	url := "/banner/" + strconv.Itoa(s.bannerIDs[5])

	req := httptest.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))

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

func (s *Suite) TestUpdateBannerBadRequest() {
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

	url := "/banner/" + strconv.Itoa(s.bannerIDs[2])

	req := httptest.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
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

func (s *Suite) TestUpdateBannerNoAccess() {
	requestBody := []byte(`{
        "tag_ids": [15, 25, 35, 45],
        "feature_id": 2212,
        "content": {
            "title": "new title",
            "text": "new text",
            "url": "new url"
        },
        "is_active": true
    }`)

	url := "/banner/" + strconv.Itoa(s.bannerIDs[7])

	req := httptest.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
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

func (s *Suite) TestUpdateBannerNotAuthorized() {

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

	url := "/banner/" + strconv.Itoa(s.bannerIDs[8])

	req := httptest.NewRequest("PATCH", url, bytes.NewBuffer(requestBody))
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
