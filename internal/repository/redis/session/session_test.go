package session_repository

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
)

func TestRepository_SaveSession(t *testing.T) {
	isSkip(t)
	type fields struct {
		rdb    *redis.Client
		expire time.Duration
	}
	type args struct {
		ctx     context.Context
		session entities.Session
	}

	rdb := initRepository(t)
	defer func() {
		_ = rdb.Close()
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
				rdb:    rdb,
				expire: 10 * time.Minute,
			},
			args: args{
				ctx: t.Context(),
				session: entities.Session{
					ID:        uuid.New(),
					UserId:    uuid.New(),
					UserAgent: "chrome",
					LastSeen:  time.Now(),
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				rdb:    tt.fields.rdb,
				expire: tt.fields.expire,
			}
			err := r.SaveSession(tt.args.ctx, tt.args.session)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

// func TestRepository_SessionById(t *testing.T) {
// 	type fields struct {
// 		rdb    *redis.Client
// 		expire time.Duration
// 	}
// 	type args struct {
// 		ctx       context.Context
// 		sessionId string
// 	}

// 	rdb := initRepository(t)
// 	defer rdb.Close()

// 	session := entities.Session{
// 		ID:        uuid.New(),
// 		UserId:    uuid.New(),
// 		UserAgent: "firefox",
// 		LastSeen:  time.Now(),
// 	}

// 	err := rdb.HSet(t.Context(), genKey(session.ID, session.UserId), session).Err()
// 	require.NoError(t, err)

// 	err = rdb.Expire(t.Context(), genKey(session.ID, session.UserId), 5*time.Minute).Err()
// 	require.NoError(t, err)

// 	time.Sleep(2 * time.Second)

// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		want    entities.Session
// 		wantErr error
// 	}{
// 		{
// 			name: "good case",
// 			fields: fields{
// 				rdb:    rdb,
// 				expire: 5 * time.Minute,
// 			},
// 			args: args{
// 				ctx:       t.Context(),
// 				sessionId: session.ID.String(),
// 			},
// 			want:    session,
// 			wantErr: nil,
// 		},
// 		{
// 			name: "not found case",
// 			fields: fields{
// 				rdb:    rdb,
// 				expire: 5 * time.Minute,
// 			},
// 			args: args{
// 				ctx:       t.Context(),
// 				sessionId: "not existing id",
// 			},
// 			want:    entities.Session{},
// 			wantErr: errs.ErrSessionNotFound,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			r := &Repository{
// 				rdb:    tt.fields.rdb,
// 				expire: tt.fields.expire,
// 			}
// 			got, err := r.SessionById(tt.args.ctx, tt.args.sessionId)
// 			require.ErrorIs(t, err, tt.wantErr)
// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("Repository.SessionById() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }

func isSkip(t *testing.T) {
	t.Helper()
	if os.Getenv("CI") != "" {
		t.Skip("skiping in ci")
	}
}

func initRepository(t *testing.T) *redis.Client {
	t.Helper()

	db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	require.NoError(t, err)

	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       db,
	})

	err = rdb.Ping(t.Context()).Err()
	require.NoError(t, err)

	return rdb
}
