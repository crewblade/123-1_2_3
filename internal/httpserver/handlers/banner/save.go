package banner

import (
	"context"
	"encoding/json"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
	"github.com/crewblade/banner-management-service/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type RequestSave struct {
	TagIDs    []int           `json:"tag_ids"`
	FeatureID int             `json:"feature_id"`
	Content   json.RawMessage `json:"content"`
	IsActive  bool            `json:"is_active"`
}

type ResponseSave struct {
	response.Response
	BannerID int `json:"banner_id"`
}

type BannerSaver interface {
	SaveBanner(
		ctx context.Context,
		tagIDs []int,
		featureID int,
		content json.RawMessage,
		isActive bool,
	) (int, error)
}

func SaveBanner(log *slog.Logger, bannerSaver BannerSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "internal.httpserver.handlers.banner.SaveBanner"

		log = log.With("op", op)
		log = log.With("request_id", middleware.GetReqID(r.Context()))

		var req RequestSave
		err := render.DecodeJSON(r.Body, &req)
		log.Info("req", slog.Any("req", req))

		if err != nil {
			log.Error("failed to decode request body", sl.Err(err))

			render.JSON(w, r, response.NewError(http.StatusBadRequest, "Incorrect data"))
			return
		}

		isAdmin := r.Context().Value("isAdmin").(bool)

		if !isAdmin {
			log.Error("User have no access")
			render.JSON(w, r, response.NewError(http.StatusForbidden, "User have no access"))
			return
		}
		bannerID, err := bannerSaver.SaveBanner(
			r.Context(),
			req.TagIDs,
			req.FeatureID,
			req.Content,
			req.IsActive,
		)
		if err != nil {
			log.Error("Internal error", sl.Err(err))
			render.JSON(w, r, response.NewError(http.StatusInternalServerError, "Internal error"))
			return
		}

		log.Info("Successful saved")
		render.JSON(w, r, ResponseSave{
			response.NewSuccess(http.StatusCreated),
			bannerID,
		})

	}

}
