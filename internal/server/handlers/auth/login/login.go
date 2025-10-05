package login

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AlexMickh/twitch-clone/internal/config"
	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/AlexMickh/twitch-clone/pkg/api"
	"github.com/AlexMickh/twitch-clone/pkg/logger"
	"github.com/go-chi/render"
)

type Loginer interface {
	Login(ctx context.Context, req dtos.LoginRequest, userAgent string) (string, error)
}

// @Summary		login user
// @Description	login user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			req	body	dtos.LoginRequest	true	"request"
// @Success		201
// @Failure		400	{object}	api.ErrorResponse
// @Failure		403	{object}	api.ErrorResponse
// @Failure		404	{object}	api.ErrorResponse
// @Failure		500	{object}	api.ErrorResponse
// @Router			/auth/login [post]
func New(loginer Loginer, sessionCfg config.SessionConfig) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.auth.login.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With(slog.String("op", op))

		var req dtos.LoginRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode body", logger.Err(err))
			return api.Error("failed to decode body", http.StatusBadRequest)
		}

		if err = req.Validate(); err != nil {
			log.Error("failed to validate body", logger.Err(err))
			return api.Error("failed to validate body", http.StatusBadRequest)
		}

		sessionId, err := loginer.Login(ctx, req, r.UserAgent())
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
				log.Error("user not found", logger.Err(err))
				return api.Error(errs.ErrUserNotFound.Error(), http.StatusNotFound)
			}
			if errors.Is(err, errs.ErrUserEmailNotVerify) {
				log.Error("email not verify", logger.Err(err))
				return api.Error(errs.ErrUserEmailNotVerify.Error(), http.StatusForbidden)
			}

			log.Error("failed to login user", logger.Err(err))
			return api.Error("failed to login user", http.StatusInternalServerError)
		}

		cookie := &http.Cookie{
			Name:     sessionCfg.Name,
			Value:    sessionId,
			Path:     "/",
			HttpOnly: sessionCfg.HttpOnly,
			Secure:   sessionCfg.Secure,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   sessionCfg.MaxAge,
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusCreated)

		return nil
	}
}
