package tests

//
//import (
//	"encoding/json"
//	"net/http"
//	"net/http/httptest"
//)
//
//func (s *Suite) TestDelete() {
//	req := httptest.NewRequest("DELETE", "/banner?tag_id=4", nil)
//	req.Header.Set("token", "admin_token")
//	w := httptest.NewRecorder()
//	s.handler.ServeHTTP(w, req)
//
//	r := s.Require()
//
//	responseBody := w.Body.Bytes()
//
//	var response map[string]interface{}
//	err := json.Unmarshal(responseBody, &response)
//	r.NoError(err)
//
//	r.Equal(200, int(response["status"].(float64)))
//
//}
//
//func (s *Suite) TestDeleteNotAuthorized() {
//	req := httptest.NewRequest("DELETE", "/banner?tag_id=4&feature_id=2", nil)
//	req.Header.Set("token", "unknown_token")
//	w := httptest.NewRecorder()
//	s.handler.ServeHTTP(w, req)
//
//	r := s.Require()
//
//	responseBody := w.Body.Bytes()
//
//	var response map[string]interface{}
//	err := json.Unmarshal(responseBody, &response)
//	r.NoError(err)
//
//	r.Equal(http.StatusUnauthorized, int(response["status"].(float64)))
//}
//
//func (s *Suite) TestDeleteHaveNoAccess() {
//	req := httptest.NewRequest("DELETE", "/banner?tag_id=2", nil)
//	req.Header.Set("token", "user_token")
//	w := httptest.NewRecorder()
//	s.handler.ServeHTTP(w, req)
//
//	r := s.Require()
//
//	responseBody := w.Body.Bytes()
//
//	var response map[string]interface{}
//	err := json.Unmarshal(responseBody, &response)
//	r.NoError(err)
//
//	r.Equal(http.StatusForbidden, int(response["status"].(float64)))
//}
//
//func (s *Suite) TestDeleteNotFound() {
//	req := httptest.NewRequest("DELETE", "/banner?tag_id=99999&feature_id=99999999", nil)
//	req.Header.Set("token", "admin_token")
//	w := httptest.NewRecorder()
//	s.handler.ServeHTTP(w, req)
//
//	r := s.Require()
//
//	responseBody := w.Body.Bytes()
//
//	var response map[string]interface{}
//	err := json.Unmarshal(responseBody, &response)
//	r.NoError(err)
//
//	r.Equal(http.StatusNotFound, int(response["status"].(float64)))
//}
//
//func (s *Suite) TestDeleteBadRequest() {
//	req := httptest.NewRequest("DELETE", "/banner?tag_id=4abc2&feature_id=9mf1", nil)
//	req.Header.Set("token", "admin_token")
//	w := httptest.NewRecorder()
//	s.handler.ServeHTTP(w, req)
//
//	r := s.Require()
//
//	responseBody := w.Body.Bytes()
//
//	var response map[string]interface{}
//	err := json.Unmarshal(responseBody, &response)
//	r.NoError(err)
//
//	r.Equal(http.StatusBadRequest, int(response["status"].(float64)))
//}
//
//func (s *Suite) TestDeletedAndCleared() {
//
//	req := httptest.NewRequest("DELETE", "/banner", nil)
//	req.Header.Set("token", "admin_token")
//	w := httptest.NewRecorder()
//	s.handler.ServeHTTP(w, req)
//
//	r := s.Require()
//
//	responseBody := w.Body.Bytes()
//
//	var response map[string]interface{}
//	err := json.Unmarshal(responseBody, &response)
//	r.NoError(err)
//
//	r.Equal(200, int(response["status"].(float64)))
//
//	s.repo.ClearData(s.ctx)
//
//	rows, err := s.repo.CountRows(s.ctx)
//
//	r.Equal(0, rows)
//
//}
