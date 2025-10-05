package session_repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
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

func (r *Repository) SessionById(ctx context.Context, sessionId string) (entities.Session, error) {
	const op = "repository.redis.session.SessionById"

	var (
		err    error
		cursor uint64
		keys   []string
	)
	//nolint
	keys, cursor, err = r.rdb.Scan(ctx, cursor, sessionId+":*", 1).Result()
	if err != nil {
		return entities.Session{}, fmt.Errorf("%s: %w", op, err)
	}
	if len(keys) == 0 {
		return entities.Session{}, fmt.Errorf("%s: %w", op, errs.ErrSessionNotFound)
	}

	var session entities.Session
	err = r.rdb.HGetAll(ctx, keys[0]).Scan(&session)
	if err != nil {
		return entities.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	userId := strings.Split(keys[0], ":")[1]

	id, err := uuid.Parse(sessionId)
	if err != nil {
		return entities.Session{}, fmt.Errorf("%s: %w", op, err)
	}
	uId, err := uuid.Parse(userId)
	if err != nil {
		return entities.Session{}, fmt.Errorf("%s: %w", op, err)
	}

	session.ID = id
	session.UserId = uId

	return session, nil
}

func genKey(id uuid.UUID, userId uuid.UUID) string {
	return id.String() + ":" + userId.String()
}
