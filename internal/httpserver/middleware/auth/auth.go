package auth

import (
	"context"
	"github.com/crewblade/banner-management-service/internal/lib/api/response"
	"github.com/crewblade/banner-management-service/internal/lib/logger/sl"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type UserProvider interface {
	IsAdmin(ctx context.Context, token string) (bool, error)
}

func AuthMiddleware(log *slog.Logger, userProvider UserProvider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "internal.httpserver.middleware.AuthMiddleware.Handle"
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

			ctx := context.WithValue(r.Context(), "isAdmin", isAdmin)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
