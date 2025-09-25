package user_service

import (
	"context"
	"fmt"

	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/google/uuid"
)

type UserRepository interface {
	SaveUser(ctx context.Context, user entities.User) error
	UserByEmail(ctx context.Context, email string) (entities.User, error)
	ValidateEmail(ctx context.Context, id uuid.UUID) error
}

type TokenService interface {
	Token(ctx context.Context, token string) (entities.Token, error)
	DeleteToken(ctx context.Context, token string) error
}

type Service struct {
	userRepository UserRepository
	tokenService   TokenService
}

func New(userRepository UserRepository, tokenService TokenService) *Service {
	return &Service{
		userRepository: userRepository,
		tokenService:   tokenService,
	}
}

func (s *Service) CreateUser(ctx context.Context, login, email, password string) (uuid.UUID, error) {
	const op = "services.user.CreateUser"

	id := uuid.New()

	user := entities.User{
		ID:              id,
		Login:           login,
		Email:           email,
		Password:        password,
		IsEmailVerified: false,
	}

	err := s.userRepository.SaveUser(ctx, user)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *Service) UserByEmail(ctx context.Context, email string) (entities.User, error) {
	const op = "services.user.UserByEmail"

	user, err := s.userRepository.UserByEmail(ctx, email)
	if err != nil {
		return entities.User{}, fmt.Errorf("%s: %w", op, err)
	}
	if !user.IsEmailVerified {
		return entities.User{}, fmt.Errorf("%s: %w", op, errs.ErrUserEmailNotVerify)
	}

	return user, nil
}

func (s *Service) VerifyEmail(ctx context.Context, req dtos.ValidateEmailRequest) error {
	const op = "services.user.VerifyEmail"

	token, err := s.tokenService.Token(ctx, req.Token)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.userRepository.ValidateEmail(ctx, token.UserId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, errs.ErrUserNotFound)
	}

	err = s.tokenService.DeleteToken(ctx, req.Token)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
