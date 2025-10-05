package verify_email

import (
	"context"
	"errors"
	"net/http"

	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/AlexMickh/twitch-clone/pkg/api"
	"github.com/AlexMickh/twitch-clone/pkg/logger"
)

type EmailVerifier interface {
	VerifyEmail(ctx context.Context, req dtos.ValidateEmailRequest) error
}

// @Summary		verify user email
// @Description	verify user email
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			token	path	string	true	"token for email verification"
// @Success		204
// @Failure		400	{object}	api.ErrorResponse
// @Failure		404	{object}	api.ErrorResponse
// @Failure		500	{object}	api.ErrorResponse
// @Router			/user/verify-email/{token} [get]
func New(emailVerifier EmailVerifier) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.user.verify_email.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx).With("op", op)

		req := dtos.ValidateEmailRequest{
			Token: r.PathValue("token"),
		}

		err := req.Validate()
		if err != nil {
			log.Error("failed to validate request", logger.Err(err))
			return api.Error("failed to validate request", http.StatusBadRequest)
		}

		err = emailVerifier.VerifyEmail(ctx, req)
		if err != nil {
			if errors.Is(err, errs.ErrTokenNotFound) {
				log.Error("token not found", logger.Err(err))
				return api.Error(errs.ErrTokenNotFound.Error(), http.StatusNotFound)
			}
		}
		if err != nil {
			if errors.Is(err, errs.ErrUserNotFound) {
				log.Error("user not found", logger.Err(err))
				return api.Error(errs.ErrUserNotFound.Error(), http.StatusNotFound)
			}
		}

		w.WriteHeader(http.StatusNoContent)

		return nil
	}
}
