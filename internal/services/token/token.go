package token_service

import (
	"context"
	"fmt"

	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/google/uuid"
)

type Repository interface {
	SaveToken(ctx context.Context, token entities.Token) error
	Token(ctx context.Context, token string) (entities.Token, error)
	DeleteToken(ctx context.Context, token string) error
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateToken(ctx context.Context, userId uuid.UUID, tokenType string) (string, error) {
	const op = "services.token.CreateToken"

	token := entities.Token{
		Token:  uuid.NewString(),
		UserId: userId,
		Type:   tokenType,
	}

	err := s.repository.SaveToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token.Token, nil
}

func (s *Service) Token(ctx context.Context, token string) (entities.Token, error) {
	const op = "services.token.Token"

	tokenEntity, err := s.repository.Token(ctx, token)
	if err != nil {
		return entities.Token{}, fmt.Errorf("%s: %w", op, err)
	}

	return tokenEntity, nil
}

func (s *Service) DeleteToken(ctx context.Context, token string) error {
	const op = "services.token.DeleteToken"

	err := s.repository.DeleteToken(ctx, token)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
