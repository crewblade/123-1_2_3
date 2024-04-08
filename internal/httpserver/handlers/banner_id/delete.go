package banner_id

import (
	"context"
	"errors"
	"github.com/crewblade/banner-management-service/internal/lib/api/errs"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
	"github.com/crewblade/banner-management-service/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type ResponseDelete struct {
	response.Response
}

type BannerDeleter interface {
	DeleteBanner(ctx context.Context, bannerID int) error
}
type UserProvider interface {
	IsAdmin(ctx context.Context, token string) (bool, error)
}

func DeleteBanner(
	log *slog.Logger,
	bannerDeleter BannerDeleter,
	userProvider UserProvider,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.httpserver.handlers.banner_id.DeleteBanner"

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
			log.Error("User have no access")
			render.JSON(w, r, response.NewError(http.StatusForbidden, "User have no access"))
			return
		}

		bannerID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("error converting bannerID from URLParam", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Invalid data"))
			return
		}
		err = bannerDeleter.DeleteBanner(r.Context(), bannerID)

		if err != nil {
			if errors.Is(err, errs.ErrBannerForTagNotFound) {
				log.Error("Banner is not found", sl.Err(err))
				render.JSON(w, r, response.NewError(http.StatusNotFound, "Banner for ID is not found"))
				return
			} else {
				log.Error("Internal error", sl.Err(err))
				render.JSON(w, r, response.NewError(http.StatusInternalServerError, "Internal error"))
			}
		}
		log.Info("Successful delete")
		render.JSON(w, r, response.NewSuccess(http.StatusOK))

	}

}
