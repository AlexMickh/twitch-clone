package token_repository

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/AlexMickh/twitch-clone/internal/consts"
	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/AlexMickh/twitch-clone/pkg/clients/mongodb"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestRepository_SaveToken(t *testing.T) {
	isSkip(t)
	type fields struct {
		coll *mongo.Collection
	}
	type args struct {
		ctx   context.Context
		token entities.Token
	}

	client, coll := initRepository(t)
	defer func() {
		_ = client.Disconnect(t.Context())
	}()

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
				ctx: nil,
				token: entities.Token{
					Token:  uuid.NewString(),
					UserId: uuid.New(),
					Type:   consts.TokenTypeVerifyEmail,
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				coll: tt.fields.coll,
			}
			err := r.SaveToken(tt.args.ctx, tt.args.token)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestRepository_Token(t *testing.T) {
	isSkip(t)
	type fields struct {
		coll *mongo.Collection
	}
	type args struct {
		ctx   context.Context
		token string
	}

	client, coll := initRepository(t)
	defer func() {
		_ = client.Disconnect(t.Context())
	}()

	token := entities.Token{
		Token:  uuid.NewString(),
		UserId: uuid.New(),
		Type:   consts.TokenTypeVerifyEmail,
	}

	_, err := coll.InsertOne(t.Context(), token)
	require.NoError(t, err)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.Token
		wantErr error
	}{
		{
			name: "good case",
			fields: fields{
				coll: coll,
			},
			args: args{
				ctx:   t.Context(),
				token: token.Token,
			},
			want:    token,
			wantErr: nil,
		},
		{
			name: "not found case",
			fields: fields{
				coll: coll,
			},
			args: args{
				ctx:   t.Context(),
				token: "not existing token",
			},
			want:    entities.Token{},
			wantErr: errs.ErrTokenNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				coll: tt.fields.coll,
			}
			got, err := r.Token(tt.args.ctx, tt.args.token)
			require.ErrorIs(t, err, tt.wantErr)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Repository.Token() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRepository_DeleteToken(t *testing.T) {
	isSkip(t)
	type fields struct {
		coll *mongo.Collection
	}
	type args struct {
		ctx   context.Context
		token string
	}

	client, coll := initRepository(t)
	defer func() {
		_ = client.Disconnect(t.Context())
	}()

	token := entities.Token{
		Token:  uuid.NewString(),
		UserId: uuid.New(),
		Type:   consts.TokenTypeVerifyEmail,
	}

	_, err := coll.InsertOne(t.Context(), token)
	require.NoError(t, err)

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
				ctx:   t.Context(),
				token: token.Token,
			},
			wantErr: nil,
		},
		{
			name: "not found case",
			fields: fields{
				coll: coll,
			},
			args: args{
				ctx:   t.Context(),
				token: "not existing token",
			},
			wantErr: errs.ErrTokenNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				coll: tt.fields.coll,
			}
			err := r.DeleteToken(tt.args.ctx, tt.args.token)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func isSkip(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" {
		t.Skip("skiping in ci")
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

	coll := client.Database("tests").Collection("tokens")

	_, err = coll.Indexes().CreateOne(
		t.Context(),
		mongo.IndexModel{
			Keys:    bson.D{{Key: "token", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	)
	require.NoError(t, err, fmt.Sprintf("failed to create index: %v", err))

	return client, coll
}
