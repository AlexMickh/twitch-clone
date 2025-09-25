package user_repository

import (
	"context"
	"fmt"

	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	coll *mongo.Collection
}

func New(ctx context.Context, client *mongo.Client, db string, collection string) (*Repository, error) {
	const op = "repository.mongo.user.New"

	coll := client.Database(db).Collection(collection)

	_, err := coll.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Repository{
		coll: coll,
	}, nil
}

func (r *Repository) SaveUser(ctx context.Context, user entities.User) error {
	const op = "repository.mongo.user.SaveUser"

	_, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		if writeErr, ok := err.(mongo.WriteException); ok {
			for _, e := range writeErr.WriteErrors {
				if e.Code == 11000 {
					return fmt.Errorf("%s: %w", op, errs.ErrUserAlreadyExists)
				}
			}
		}
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) UserByEmail(ctx context.Context, email string) (entities.User, error) {
	const op = "repository.mongo.user.UserByEmail"

	filter := bson.D{{Key: "email", Value: email}}
	result := r.coll.FindOne(ctx, filter)
	if result.Err() != nil {
		return entities.User{}, fmt.Errorf("%s: %w", op, errs.ErrUserNotFound)
	}

	var user entities.User
	if err := result.Decode(&user); err != nil {
		return entities.User{}, fmt.Errorf("%s: %w", op, errs.ErrUserNotFound)
	}

	return user, nil
}

func (r *Repository) ValidateEmail(ctx context.Context, id uuid.UUID) error {
	const op = "repository.mongo.user.ValidateEmail"

	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "is_email_verified", Value: true},
		}},
	}
	result, err := r.coll.UpdateByID(ctx, id, update)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if result.ModifiedCount == 0 {
		return fmt.Errorf("%s: %w", op, errs.ErrUserNotFound)
	}

	return nil
}
