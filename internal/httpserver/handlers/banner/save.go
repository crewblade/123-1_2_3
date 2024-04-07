package banner

import (
	"context"
	"github.com/crewblade/banner-management-service/internal/domain/models"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
	"github.com/crewblade/banner-management-service/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
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
	SaveBanner(ctx context.Context) (int, error)
}

func SaveBanner(log *slog.Logger, bannerSaver BannerSaver, userProvider UserProvider) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.httpserver.handlers.banner.SaveBanner"

		log = log.With("op", op)
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		var req RequestSave
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))

			return
		}

		token := r.Header.Get("token")
		log.With("token", token)

		isAdmin, err := userProvider.IsAdmin(r.Context(), token)
		if err != nil {
			log.Error("Invalid token: ", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusUnauthorized, "User is not authorized"))
			return
		}

		if !isAdmin {
			log.Error("User have no access", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusForbidden, "User have no access"))
			return
		}

	}

}
