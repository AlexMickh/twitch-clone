package register

import (
	"context"
	"errors"
	"net/http"

	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/AlexMickh/twitch-clone/pkg/api"
	"github.com/AlexMickh/twitch-clone/pkg/logger"
	"github.com/go-chi/render"
)

type Registerer interface {
	Register(ctx context.Context, req dtos.RegisterRequest) (string, error)
}

// @Summary		register user
// @Description	register user
// @Tags			auth
// @Accept			json
// @Produce		json
// @Param			req	body		dtos.RegisterRequest	true	"request"
// @Success		201	{object}	dtos.RegisterResponse
// @Failure		400	{object}	api.ErrorResponse
// @Failure		500	{object}	api.ErrorResponse
// @Router			/auth/register [post]
func New(registerer Registerer) api.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		const op = "handlers.auth.register.New"
		ctx := r.Context()
		log := logger.FromCtx(ctx)

		var req dtos.RegisterRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("failed to decode body", logger.Err(err))
			return api.Error("failed to decode body", http.StatusBadRequest)
		}

		if err = req.Validate(); err != nil {
			log.Error("failed to validate body", logger.Err(err))
			return api.Error("failed to validate body", http.StatusBadRequest)
		}

		id, err := registerer.Register(ctx, req)
		if err != nil {
			if errors.Is(err, errs.ErrUserAlreadyExists) {
				log.Error("user already exists", logger.Err(err))
				return api.Error("user already exists", http.StatusBadRequest)
			}

			log.Error("failed to register user", logger.Err(err))
			return api.Error("failed to register user", http.StatusBadRequest)
		}

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, dtos.ToRegisterResponse(id))

		return nil
	}
}
