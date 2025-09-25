package dtos

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ValidateEmailRequest struct {
	Token string `validate:"required,uuid4"`
}

func (v ValidateEmailRequest) Validate() error {
	const op = "dtos.register.Validate"

	if err := validator.New().Struct(&v); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
