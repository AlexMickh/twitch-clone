package user_repository

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/AlexMickh/twitch-clone/pkg/clients/mongodb"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestRepository_SaveUser(t *testing.T) {
	type fields struct {
		coll *mongo.Collection
	}
	type args struct {
		ctx  context.Context
		user entities.User
	}

	client, coll := initRepository(t)
	defer client.Disconnect(t.Context())

	existingUser := entities.User{
		ID:       uuid.New(),
		Login:    gofakeit.FirstName(),
		Email:    gofakeit.Email(),
		Password: "some password",
	}

	_, err := coll.InsertOne(t.Context(), existingUser)
	require.NoError(t, err, fmt.Sprintf("failed to save user: %v", err))

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "good case",
			fields: fields{
				coll: coll,
			},
			args: args{
				ctx: context.TODO(),
				user: entities.User{
					ID:       uuid.New(),
					Login:    gofakeit.FirstName(),
					Email:    gofakeit.Email(),
					Password: "some password",
				},
			},
			wantErr: nil,
		},
		{
			name: "user already exists case",
			fields: fields{
				coll: coll,
			},
			args: args{
				ctx:  context.TODO(),
				user: existingUser,
			},
			wantErr: errs.ErrUserAlreadyExists,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				coll: tt.fields.coll,
			}
			err := r.SaveUser(tt.args.ctx, tt.args.user)
			require.ErrorIs(t, err, tt.wantErr, fmt.Sprintf("Repository.SaveUser() error = %v, wantErr %v", err, tt.wantErr))
		})
	}
}

func initRepository(t *testing.T) (*mongo.Client, *mongo.Collection) {
	t.Helper()

	connString := fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/?authSource=admin",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)

	client, err := mongo.Connect(options.Client().ApplyURI(connString).SetRegistry(mongodb.UUIDRegistry))
	require.NoError(t, err, fmt.Sprintf("failed to connect to db: %v", err))

	coll := client.Database("tests").Collection("users")

	_, err = coll.Indexes().CreateOne(
		t.Context(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	require.NoError(t, err, fmt.Sprintf("failed to create index: %v", err))

	return client, coll
}
