package user_banner

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/crewblade/banner-management-service/internal/domain/models"
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
	Content json.RawMessage `json:"content"`
}

type UserBannerGetter interface {
	GetUserBanner(ctx context.Context, tagID int, featureID int) (json.RawMessage, bool, error)
}

type UserProvider interface {
	IsAdmin(ctx context.Context, token string) (bool, error)
}

type BannerCache interface {
	GetBanner(ctx context.Context, tagID, featureID int) (json.RawMessage, bool, error)
	SetBanner(ctx context.Context, tagID, featureID int, banner *models.BannerForUser) error
}

func GetUserBanner(
	log *slog.Logger,
	userBannerGetter UserBannerGetter,
	userProvider UserProvider,
	bannerCache BannerCache,
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
		} else if useLastRevisionStr != "false" && useLastRevisionStr != "" {
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
		var bannerContent json.RawMessage
		var bannerIsActive bool
		isCacheUsed := false
		if !useLastRevision {
			bannerContent, bannerIsActive, err = bannerCache.GetBanner(r.Context(), tagID, featureID)
			if err != nil {
				log.Error("Error fetching banner content from cache", sl.Err(err))
			} else {
				isCacheUsed = true
			}
		}
		if useLastRevision || !isCacheUsed {
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
			err := bannerCache.SetBanner(r.Context(), tagID, featureID, &models.BannerForUser{bannerContent, bannerIsActive})
			if err != nil {
				log.Error("Error setting banner content in cache", sl.Err(err))
			} else {
				log.Info(
					"Data cached:",
					slog.Any("bannerContent", bannerContent),
					slog.Any("bannerIsActive", bannerIsActive),
					slog.Any("tagID", tagID),
					slog.Any("featureID", featureID))
			}
		}
		if !isAdmin && !bannerIsActive {
			log.Error("User have no access to inactive banner")
			render.JSON(w, r, response.NewError(http.StatusForbidden, errs.ErrUserDoesNotHaveAccess.Error()))
			return
		}
		log.Info("Successful respnose:", slog.Any("banner content", bannerContent))
		render.JSON(w, r, ResponseGet{
			response.NewSuccess(200),
			bannerContent,
		})
	}

}
