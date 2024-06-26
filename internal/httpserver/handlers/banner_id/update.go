package banner_id

import (
	"context"
	"encoding/json"
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

type RequestUpdate struct {
	TagIDs    []int           `json:"tag_ids"`
	FeatureID int             `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
}

type ResponseUpdate struct {
	response.Response
}

type BannerUpdater interface {
	UpdateBanner(
		ctx context.Context,
		bannerID int,
		tagIDs []int,
		featureID int,
		content json.RawMessage,
		isActive bool,
	) error
}

func UpdateBanner(
	log *slog.Logger,
	bannerUpdater BannerUpdater,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.httpserver.handlers.banner_id.UpdateBanner"

		log = log.With("op", op)
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		var req RequestUpdate
		err := render.DecodeJSON(r.Body, &req)

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))

			return
		}

		bannerID, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Error("error converting bannerID from URLParam", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}

		isAdmin := r.Context().Value("isAdmin").(bool)

		if !isAdmin {
			log.Error("User have no access")
			render.JSON(w, r, response.NewError(http.StatusForbidden, "User have no access"))
			return
		}

		err = bannerUpdater.UpdateBanner(
			r.Context(),
			bannerID,
			req.TagIDs,
			req.FeatureID,
			req.Content,
			req.IsActive,
		)

		if err != nil {
			if errors.Is(err, errs.ErrBannerNotFound) {
				log.Error("Banner is not found", sl.Err(err))
				render.JSON(w, r, response.NewError(http.StatusNotFound, "Banner is not found"))
				return
			} else {
				log.Error("Internal error", sl.Err(err))
				render.JSON(w, r, response.NewError(http.StatusInternalServerError, "Internal error"))
				return
			}
		}

		log.Info("Successful update", slog.Any("Updated banner", req))
		render.JSON(w, r, response.NewSuccess(http.StatusOK))
	}

}
