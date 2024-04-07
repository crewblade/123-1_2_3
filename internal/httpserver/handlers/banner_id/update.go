package banner_id

import (
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
)

type RequestUpdate struct {
	TagIDs    []int                `json:"tag_ids"`
	FeatureID int                  `json:"feature_id"`
	Content   models.BannerContent `json:"content"`
	IsActive  bool                 `json:"is_active"`
}

type ResponseUpdate struct {
	response.Response
}

type BannerUpdater interface {
	UpdateBanner()
}

func UpdateBanner() {

}
