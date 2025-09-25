package session_repository

import (
	"context"
	"fmt"
	"time"

	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type Repository struct {
	rdb    *redis.Client
	expire time.Duration
}

func New(rdb *redis.Client, expire time.Duration) *Repository {
	return &Repository{
		rdb:    rdb,
		expire: expire,
	}
}

func (r *Repository) SaveSession(ctx context.Context, session entities.Session) error {
	const op = "repository.redis.session.SaveSession"

	key := genKey(session.ID, session.UserId)
	pipeline := r.rdb.Pipeline()

	err := pipeline.HSet(ctx, key, session).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	err = pipeline.Expire(ctx, key, r.expire).Err()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = pipeline.Exec(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func genKey(id uuid.UUID, userId uuid.UUID) string {
	return id.String() + ":" + userId.String()
}
