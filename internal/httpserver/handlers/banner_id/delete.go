package banner_id

type RequestDelete struct {
}

type ResponseDelete struct {
}

type BannerDeleter interface {
	DeleteBanner()
}

func DeleteBanner() {

}
