package banner

import (
	"context"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
	"github.com/crewblade/banner-management-service/internal/lib/logger/sl"
	"github.com/crewblade/banner-management-service/internal/lib/utils"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type ResponseDelete struct {
	response.Response
	Count int `json:"count"`
}

type BannersDeleter interface {
	DeleteBanners(ctx context.Context, featureID, tagID *int) (int, error)
}

func DeleteBanners(
	log *slog.Logger,
	bannersDeleter BannersDeleter,
	userProvider UserProvider,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "internal.httpserver.handlers.banner.DeleteBanners"
		log = log.With("op", op)
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		token := r.Header.Get("token")
		log.With("token", token)

		tagID, err := utils.StrToIntPtr(r.URL.Query().Get("tag_id"), log)
		if err != nil {
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}

		featureID, err := utils.StrToIntPtr(r.URL.Query().Get("feature_id"), log)
		if err != nil {
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}

		isAdmin, err := userProvider.IsAdmin(r.Context(), token)
		if err != nil {
			log.Error("Invalid token: ", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusUnauthorized, "User is not authorized"))
			return
		}

		if !isAdmin {
			log.Error("User have no access")
			render.JSON(w, r, response.NewError(http.StatusForbidden, "User have no access"))
			return
		}

		cnt, err := bannersDeleter.DeleteBanners(r.Context(), featureID, tagID)

		if err != nil {
			log.Error("Internal error:", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusInternalServerError, "Internal error"))
			return
		}

		if cnt == 0 {
			log.Error("No banner for delete")
			render.JSON(w, r, response.NewError(http.StatusNotFound, "Banners for delete not found"))
			return
		}

		log.Info("successful delete banners:")
		render.JSON(w, r, ResponseDelete{
			response.NewSuccess(200),
			cnt,
		})
	}

}
