package banner_id

import "github.com/crewblade/banner-management-service/internal/domain/models"

type RequestUpdate struct {
	TagIDs    []int                `json:"tag_ids"`
	FeatureID int                  `json:"feature_id"`
	Content   models.BannerContent `json:"content"`
}

type ResponseUpdate struct {
}

type BannerUpdater interface {
	UpdateBanner()
}

func UpdateBanner() {

}
