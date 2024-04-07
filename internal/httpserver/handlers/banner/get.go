package banner

import (
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
	"github.com/crewblade/banner-management-service/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
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

func GetBanners(log *slog.Logger, bannersGetter BannersGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.httpserver.handlers.banner.GetBanners"

		log = log.With("op", op)
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		var req RequestGet
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))

			return
		}

		log.Info("request body decoded", slog.Any("request", req))

	}

}
