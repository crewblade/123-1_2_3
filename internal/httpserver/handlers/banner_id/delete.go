package banner_id

import "github.com/crewblade/banner-management-service/internal/lib/api/response"

type ResponseDelete struct {
	response.Response
}

type BannerDeleter interface {
	DeleteBanner()
}

func DeleteBanner() {

}
