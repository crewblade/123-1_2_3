package user_banner

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
	GetUserBanner(ctx context.Context)
}
type UserProvider interface {
	IsAdmin(ctx context.Context, token string) (bool, error)
}

func GetUserBanner(
	log *slog.Logger,
	userBannerGetter UserBannerGetter,
	userProvider UserProvider,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.httpserver.handlers.user_banner.GetUserBanner"

		log = log.With("op", op)
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		var req RequestGet
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))

			return
		}

		token := r.Header.Get("token")
		log.With("token", token)

		_, err = userProvider.IsAdmin(r.Context(), token)
		if err != nil {
			log.Error("Invalid token: ", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusUnauthorized, "User is not authorized"))
			return
		}
	}

}
