package user_banner

import (
	"context"
	"errors"
	"github.com/crewblade/banner-management-service/internal/lib/api/errs"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
	"github.com/crewblade/banner-management-service/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type ResponseGet struct {
	response.Response
	Content string `json:"content"`
}

type UserBannerGetter interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int) (string, bool, error)
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

		tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
		if err != nil {
			log.Error("error converting tagID", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}
		featureID, err := strconv.Atoi(r.URL.Query().Get("feature_id"))
		if err != nil {
			log.Error("error converting featureID", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}
		useLastRevision := false
		useLastRevisionStr := r.URL.Query().Get("use_last_revision")
		if useLastRevisionStr == "true" {
			useLastRevision = true
		} else if useLastRevisionStr != "false" {
			log.Error("Incorrect data")
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}

		token := r.Header.Get("token")
		log.With("token", token)

		isAdmin, err := userProvider.IsAdmin(r.Context(), token)
		if err != nil {
			log.Error("In: ", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusUnauthorized, "User is not authorized"))
			return
		}
		var bannerContent string
		var bannerIsActive bool
		if !useLastRevision {
			//banner, err = cache.GetUserBanner()
		} else {
			bannerContent, bannerIsActive, err = userBannerGetter.GetUserBanner(r.Context(), tagID, featureID)
			if err != nil {
				if errors.Is(err, errs.ErrBannerNotFound) {
					log.Error("Banner is not found", sl.Err(err))
					render.JSON(w, r, response.NewError(http.StatusNotFound, "Banner is not found"))
					return
				}
				log.Error("Internal error", sl.Err(err))
				render.JSON(w, r, response.NewError(http.StatusInternalServerError, "Intrenal error"))
				return
			}

		}
		if !isAdmin && !bannerIsActive {
			log.Error("User have no access to inactive banner")
			render.JSON(w, r, response.NewError(http.StatusForbidden, errs.ErrUserDoesNotHaveAccess.Error()))
			return
		}
		log.Info("Successful respnose:", slog.String("banner content", bannerContent))
		render.JSON(w, r, ResponseGet{
			response.NewSuccess(200),
			bannerContent,
		})
	}

}
