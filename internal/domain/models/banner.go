package models

type BannerContent struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	URL   string `json:"url"`
}

type Banner struct {
	BannerID  int           `json:"banner_id"`
	TagIDs    []int         `json:"tag_ids"`
	FeatureID int           `json:"feature_id"`
	Content   BannerContent `json:"content"`
	IsActive  bool          `json:"is_active"`
	CreatedAt string        `json:"created_at"`
	UpdatedAt string        `json:"updated_at"`
}
