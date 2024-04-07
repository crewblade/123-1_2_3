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

type RequestGet struct {
	TagID           int  `json:"tag_id,omitempty"`
	FeatureID       int  `json:"feature_id,omitempty"`
	UseLastRevision bool `json:"use_last_revision,omitempty"`
	Limit           int  `json:"limit,omitempty"`
	Offset          int  `json:"offset,omitempty"`
}
type ResponseGet struct {
	response.Response
	Banners []models.Banner `json:"items"`
}

type BannersGetter interface {
	GetBanners(ctx context.Context) ([]models.Banner, error)
}

type UserProvider interface {
	IsAdmin(ctx context.Context, token string) (bool, error)
}

func GetBanners(
	log *slog.Logger,
	bannersGetter BannersGetter,
	userProvider UserProvider,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "internal.httpserver.handlers.banner.GetBanners"
		log = log.With("op", op)
		log = log.With("request_id", middleware.GetReqID(r.Context()))

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

		var req RequestGet
		err = render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))

			return
		}

		//banners, err := bannersGetter.GetBanners(r.Context(), )
		//if err != nil {
		//	log.Error("Internal error:", sl.Err(err))
		//	render.JSON(w, r, response.NewError(http.StatusInternalServerError, "Internal error"))
		//
		//	return
		//
		//}
		//
		//log.Info("get banners:", banners)
		//render.JSON(w, r, ResponseGet{
		//	response.NewSuccess(200),
		//	banners,
		//})
		render.JSON(w, r, response.NewSuccess(200))
	}

}
