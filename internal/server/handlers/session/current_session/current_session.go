package current_session

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/twitch-clone/internal/config"
	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/AlexMickh/twitch-clone/pkg/api"
	"github.com/AlexMickh/twitch-clone/pkg/logger"
	"github.com/go-chi/render"
)

type SessionProvider interface {
	SessionById(ctx context.Context, sessionId string) (entities.Session, error)
}

// @Summary		login user
// @Description	login user
// @Tags			session
// @Accept			json
// @Produce		json
// @Success		200	{object}	dtos.CurrentSessionResponse
// @Failure		401	{object}	api.ErrorResponse
// @Failure		404	{object}	api.ErrorResponse
// @Failure		500	{object}	api.ErrorResponse
// @Security		SessionAuth
// @Router			/session/current [get]
func New(sessionProvider SessionProvider, sessionCfg config.SessionConfig) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.auth.register.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		cookie, err := r.Cookie(sessionCfg.Name)
		if err != nil {
			log.Error("failed to get cookie", logger.Err(err))
			return api.Error("failed to get cookie", http.StatusUnauthorized)
		}

		session, err := sessionProvider.SessionById(ctx, cookie.Value)
		if err != nil {
			if errors.Is(err, errs.ErrSessionNotFound) {
				log.Error("session not found", logger.Err(err))
				return api.Error(errs.ErrSessionNotFound.Error(), http.StatusNotFound)
			}

			log.Error("failed to get session", logger.Err(err))
			return api.Error("failed to get session", http.StatusInternalServerError)
		}

		cookie.MaxAge = sessionCfg.MaxAge
		http.SetCookie(w, cookie)

		render.JSON(w, r, dtos.ToCurrentSessionResponse(session.ID, session.UserId, session.UserAgent))

		return nil
	}
}
