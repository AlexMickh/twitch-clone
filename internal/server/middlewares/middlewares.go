package middlewares

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/twitch-clone/internal/config"
	"github.com/AlexMickh/twitch-clone/internal/consts"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/AlexMickh/twitch-clone/pkg/api"
	"github.com/AlexMickh/twitch-clone/pkg/logger"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type SessionValidator interface {
	ValidateSession(ctx context.Context, sessionId string) (uuid.UUID, error)
}

func Auth(sessionCfg config.SessionConfig, sessionValidator SessionValidator) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			const op = "middlewares.Auth"
			ctx := r.Context()
			log := logger.FromCtx(ctx).With(slog.String("op", op))

			cookie, err := r.Cookie(sessionCfg.Name)
			if err != nil {
				log.Error("failed to get session", logger.Err(err))
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, api.ErrorResponse{
					Error: "failed to get session",
				})
				return
			}

			userId, err := sessionValidator.ValidateSession(ctx, cookie.Value)
			if err != nil {
				if errors.Is(err, errs.ErrSessionNotFound) {
					log.Error("session not found", logger.Err(err))
					render.Status(r, http.StatusUnauthorized)
					render.JSON(w, r, api.ErrorResponse{
						Error: errs.ErrSessionNotFound.Error(),
					})
					return
				}
				log.Error("failed to validate session", logger.Err(err))
				render.Status(r, http.StatusUnauthorized)
				render.JSON(w, r, api.ErrorResponse{
					Error: "failed to validate session",
				})
				return
			}

			//nolint:staticcheck
			ctx = context.WithValue(ctx, consts.ContextUserId, userId)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
