package user_banner

import (
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
)

type RequestGet struct {
	TagID           int    `json:"tag_id"`
	FeatureID       int    `json:"feature_id"`
	UseLastRevision bool   `json:"use_last_revision"`
	Token           string `json:"token"`
}

type ResponseGet struct {
	response.Response
	models.BannerContent
}

type UserBannerGetter interface {
	GetUserBanner()
}

func GetUserBanner() {

}
