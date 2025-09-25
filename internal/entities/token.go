package entities

import "github.com/google/uuid"

type Token struct {
	Token  string    `bson:"token"`
	UserId uuid.UUID `bson:"user_id"`
	Type   string    `bson:"type"`
}
