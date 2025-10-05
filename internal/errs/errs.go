package errs

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrUserEmailNotVerify = errors.New("user email not verify")
	ErrTokenNotFound      = errors.New("token not found")
	ErrSessionNotFound    = errors.New("session not found")
)
