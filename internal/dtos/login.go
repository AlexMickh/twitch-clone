package dtos

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

func (l LoginRequest) Validate() error {
	const op = "dtos.register.Validate"

	if err := validator.New().Struct(&l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
