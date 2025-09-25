package entities

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID        uuid.UUID `redis:"-"`
	UserId    uuid.UUID `redis:"-"`
	UserAgent string    `redis:"user_agent"`
	LastSeen  time.Time `redis:"last_seen"`
}
