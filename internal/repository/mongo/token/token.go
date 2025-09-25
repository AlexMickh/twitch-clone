package token_repository

import (
	"context"
	"fmt"

	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Repository struct {
	coll *mongo.Collection
}

func New(ctx context.Context, client *mongo.Client, db string, collection string) (*Repository, error) {
	const op = "repository.mongo.token.New"

	coll := client.Database(db).Collection(collection)

	_, err := coll.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.D{{Key: "token", Value: 1}},
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

func (r *Repository) SaveToken(ctx context.Context, token entities.Token) error {
	const op = "repository.mongo.token.New"

	_, err := r.coll.InsertOne(ctx, token)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Repository) Token(ctx context.Context, token string) (entities.Token, error) {
	const op = "repository.mongo.token.Token"

	filter := bson.D{{Key: "token", Value: token}}
	result := r.coll.FindOne(ctx, filter)
	if result.Err() != nil {
		return entities.Token{}, fmt.Errorf("%s: %w", op, errs.ErrTokenNotFound)
	}

	var tokenEntity entities.Token
	err := result.Decode(&tokenEntity)
	if err != nil {
		return entities.Token{}, fmt.Errorf("%s: %w", op, errs.ErrTokenNotFound)
	}

	return tokenEntity, nil
}

func (r *Repository) DeleteToken(ctx context.Context, token string) error {
	const op = "repository.mongo.token.DeleteToken"

	filter := bson.D{{Key: "token", Value: token}}
	result, err := r.coll.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("%s: %w", op, errs.ErrTokenNotFound)
	}

	return nil
}
