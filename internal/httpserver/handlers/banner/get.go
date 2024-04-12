package banner

import (
	"context"
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
	Banners []models.Banner `json:"items"`
}

type BannersGetter interface {
	GetBanners(ctx context.Context, featureID, tagID, limit, offset *int) ([]models.Banner, error)
}

type UserProvider interface {
	IsAdmin(ctx context.Context, token string) (bool, error)
}

type BannersCache interface {
	GetBanners(ctx context.Context, featureID, tagID, limit, offset *int) ([]models.Banner, error)
	SetBanners(ctx context.Context, featureID, tagID, limit, offset *int, banners []models.Banner) error
}

func GetBanners(
	log *slog.Logger,
	bannersGetter BannersGetter,
	userProvider UserProvider,
	bannersCache BannersCache,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		const op = "internal.httpserver.handlers.banner.GetBanners"
		log = log.With("op", op)
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		token := r.Header.Get("token")
		log.With("token", token)

		tagID, err := strToIntPtr(r.URL.Query().Get("tag_id"), log)
		if err != nil {
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}

		featureID, err := strToIntPtr(r.URL.Query().Get("feature_id"), log)
		if err != nil {
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}

		limit, err := strToIntPtr(r.URL.Query().Get("limit"), log)
		if err != nil || (limit != nil && *limit < 0) {
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}
		if limit == nil {
			defaultLimit := 100
			limit = &defaultLimit
		}

		offset, err := strToIntPtr(r.URL.Query().Get("offset"), log)
		if err != nil || (offset != nil && *offset < 0) {
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}
		if offset == nil {
			zero := 0
			offset = &zero
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

		var banners []models.Banner
		isCacheUsed := false

		banners, err = bannersCache.GetBanners(r.Context(), featureID, tagID, limit, offset)
		if err != nil {
			log.Error("Error fetching banners from cache", sl.Err(err))
		} else {
			log.Info("Data from cache:", slog.Any("banners", banners))

			isCacheUsed = true
		}

		if !isCacheUsed {
			banners, err = bannersGetter.GetBanners(r.Context(), featureID, tagID, limit, offset)
			if err != nil {
				log.Error("Internal error:", sl.Err(err))
				render.JSON(w, r, response.NewError(http.StatusInternalServerError, "Internal error"))

				return
			}
			err := bannersCache.SetBanners(r.Context(), featureID, tagID, limit, offset, banners)
			if err != nil {
				log.Error("Error setting banner content in cache", sl.Err(err))
			}
		}

		if len(banners) == 0 {
			log.Error("no banners found")
			render.JSON(w, r, response.NewError(http.StatusNotFound, errs.ErrNoBannersFound.Error()))
			return
		}

		log.Info("successful get banners:", slog.Any("banners", banners))
		render.JSON(w, r, ResponseGet{
			response.NewSuccess(200),
			banners,
		})
	}

}

func strToIntPtr(str string, log *slog.Logger) (*int, error) {
	if str == "" {
		return nil, nil
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Error("error converting value", sl.Err(err))
		return nil, err
	}
	return &val, nil
}
