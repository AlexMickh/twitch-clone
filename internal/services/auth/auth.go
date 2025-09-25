package auth_service

import (
	"context"
	"fmt"

	"github.com/AlexMickh/twitch-clone/internal/consts"
	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	CreateUser(ctx context.Context, login, email, password string) (uuid.UUID, error)
	UserByEmail(ctx context.Context, email string) (entities.User, error)
}

type VerificationSender interface {
	SendVerification(to string, token, login string) error
}

type TokenService interface {
	CreateToken(ctx context.Context, userId uuid.UUID, tokenType string) (string, error)
}

type SessionService interface {
	CreateSession(ctx context.Context, userId uuid.UUID, userAgent string) (uuid.UUID, error)
}

type Service struct {
	userService        UserService
	verificationSender VerificationSender
	tokenService       TokenService
	sessionService     SessionService
}

func New(
	userService UserService,
	verificationSender VerificationSender,
	tokenService TokenService,
	sessionService SessionService,
) *Service {
	return &Service{
		userService:        userService,
		verificationSender: verificationSender,
		tokenService:       tokenService,
		sessionService:     sessionService,
	}
}

func (s *Service) Register(ctx context.Context, req dtos.RegisterRequest) (string, error) {
	const op = "services.auth.Register"

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	id, err := s.userService.CreateUser(ctx, req.Login, req.Email, string(hashPassword))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	token, err := s.tokenService.CreateToken(ctx, id, consts.TokenTypeVerifyEmail)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.verificationSender.SendVerification(req.Email, token, req.Login)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return id.String(), nil
}

func (s *Service) Login(ctx context.Context, req dtos.LoginRequest, userAgent string) (string, error) {
	const op = "services.auth.Login"

	user, err := s.userService.UserByEmail(ctx, req.Email)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, errs.ErrUserNotFound)
	}

	sessionId, err := s.sessionService.CreateSession(ctx, user.ID, userAgent)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return sessionId.String(), nil
}
