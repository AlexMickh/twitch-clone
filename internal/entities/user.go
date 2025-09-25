package entities

import "github.com/google/uuid"

type User struct {
	ID              uuid.UUID `bson:"_id"`
	Login           string    `bson:"login"`
	Email           string    `bson:"email"`
	Password        string    `bson:"password"`
	IsEmailVerified bool      `bson:"is_email_verified"`
}
