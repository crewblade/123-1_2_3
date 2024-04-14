package models

import (
	"encoding/json"
	"time"
)

type Banner struct {
	BannerID  int             `json:"banner_id"`
	TagIDs    []int           `json:"tag_ids"`
	FeatureID int             `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
type BannerForUser struct {
	Content  json.RawMessage `json:"content"`
	IsActive bool            `json:"is_active"`
}
