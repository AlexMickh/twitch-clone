package session_service

import (
	"context"
	"testing"
	"time"

	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateSession(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx       context.Context
		userId    uuid.UUID
		userAgent string
	}

	m := NewMockRepository(t)

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMockErr error
		wantErr     error
	}{
		{
			name: "good case",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx:       context.Background(),
				userId:    uuid.New(),
				userAgent: "firefox",
			},
			wantMockErr: nil,
			wantErr:     nil,
		},
		{
			name: "repository error case",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx:       context.Background(),
				userId:    uuid.New(),
				userAgent: "firefox",
			},
			wantMockErr: errs.ErrSessionNotFound,
			wantErr:     errs.ErrSessionNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().SaveSession(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("entities.Session"),
			).Return(tt.wantMockErr).Once()

			s := &Service{
				repository: tt.fields.repository,
			}
			_, err := s.CreateSession(tt.args.ctx, tt.args.userId, tt.args.userAgent)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestService_SessionById(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx       context.Context
		sessionId string
	}

	m := NewMockRepository(t)

	sessionId := uuid.NewString()

	tests := []struct {
		name        string
		fields      fields
		args        args
		want        entities.Session
		wantMockErr error
		wantErr     error
	}{
		{
			name: "good case",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx:       context.Background(),
				sessionId: sessionId,
			},
			want: entities.Session{
				ID:        uuid.MustParse(sessionId),
				UserId:    uuid.New(),
				UserAgent: "firefox",
				LastSeen:  time.Now(),
			},
			wantMockErr: nil,
			wantErr:     nil,
		},
		{
			name: "repository error case",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx:       context.Background(),
				sessionId: "invalid session id",
			},
			want:        entities.Session{},
			wantMockErr: errs.ErrSessionNotFound,
			wantErr:     errs.ErrSessionNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().SessionById(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
			).Return(tt.want, tt.wantMockErr).Once()

			s := &Service{
				repository: tt.fields.repository,
			}
			got, err := s.SessionById(tt.args.ctx, tt.args.sessionId)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestService_ValidateSession(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx       context.Context
		sessionId string
	}

	m := NewMockRepository(t)

	tests := []struct {
		name        string
		fields      fields
		args        args
		wantMockErr error
		wantErr     error
	}{
		{
			name: "good case",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx:       context.Background(),
				sessionId: uuid.NewString(),
			},
			wantMockErr: nil,
			wantErr:     nil,
		},
		{
			name: "repository error case",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx:       context.Background(),
				sessionId: "invalid session id",
			},
			wantMockErr: errs.ErrSessionNotFound,
			wantErr:     errs.ErrSessionNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().SessionById(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
			).Return(entities.Session{}, tt.wantMockErr).Once()

			s := &Service{
				repository: tt.fields.repository,
			}
			_, err := s.ValidateSession(tt.args.ctx, tt.args.sessionId)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
