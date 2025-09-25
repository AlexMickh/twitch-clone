package dtos

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type RegisterRequest struct {
	Login    string `json:"login" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}

func (r RegisterRequest) Validate() error {
	const op = "dtos.register.Validate"

	if err := validator.New().Struct(&r); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func ToRegisterResponse(id string) RegisterResponse {
	return RegisterResponse{
		ID: id,
	}
}
