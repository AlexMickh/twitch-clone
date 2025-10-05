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
	SessionById(ctx context.Context, sessionId string) (entities.Session, error)
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
	const op = "services.session.CreateSession"

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

func (s *Service) SessionById(ctx context.Context, sessionId string) (entities.Session, error) {
	const op = "services.session.SessionById"

	session, err := s.repository.SessionById(ctx, sessionId)
	if err != nil {
		return entities.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	return session, nil
}

func (s *Service) ValidateSession(ctx context.Context, sessionId string) (uuid.UUID, error) {
	const op = "services.session.ValidateSession"

	session, err := s.SessionById(ctx, sessionId)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("%s: %w", op, err)
	}

	return session.UserId, nil
}
