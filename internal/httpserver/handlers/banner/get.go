package banner

import (
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
)

type Banner struct {
	BannerID  int                  `json:"banner_id"`
	TagIDs    []int                `json:"tag_ids"`
	FeatureID int                  `json:"feature_id"`
	Content   models.BannerContent `json:"content"`
	IsActive  bool                 `json:"is_active"`
	CreatedAt string               `json:"created_at"`
	UpdatedAt string               `json:"updated_at"`
}

type RequestGet struct {
	TagID           int    `json:"tag_id"`
	FeatureID       int    `json:"feature_id"`
	UseLastRevision bool   `json:"use_last_revision"`
	Limit           int    `json:"limit"`
	Offset          int    `json:"offset"`
	Token           string `json:"token"`
}

type ResponseGet struct {
	response.Response
	Banners []Banner `json:"items"`
}

type BannersGetter interface {
	GetBanners()
}

func GetBanners() {

}
