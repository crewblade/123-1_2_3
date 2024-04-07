package banner

import (
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
)

type RequestSave struct {
	TagIDs    []int                `json:"tag_ids"`
	FeatureID int                  `json:"feature_id"`
	Content   models.BannerContent `json:"content"`
	IsActive  bool                 `json:"is_active"`
}

type ResponseSave struct {
	response.Response
	BannerID int `json:"banner_id"`
}

type BannerSaver interface {
	SaveBanner()
}

func SaveBanner() {

}
