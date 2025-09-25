package session_service

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/google/uuid"
)

type Repository interface {
	SaveSession(ctx context.Context, session entities.Session) error
}

type Service struct {
	repository Repository
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) CreateSession(ctx context.Context, userId uuid.UUID, userAgent string) (uuid.UUID, error) {
	const op = "services.user.CreateSession"

	session := entities.Session{
		ID:        uuid.New(),
		UserId:    userId,
		UserAgent: userAgent,
		LastSeen:  time.Now(),
	}

	err := s.repository.SaveSession(ctx, session)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return session.ID, nil
}
