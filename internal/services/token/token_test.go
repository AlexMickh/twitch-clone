package token_service

import (
	"context"
	"testing"

	"github.com/AlexMickh/twitch-clone/internal/consts"
	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateToken(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx       context.Context
		userId    uuid.UUID
		tokenType string
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
				tokenType: consts.TokenTypeVerifyEmail,
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
				tokenType: consts.TokenTypeVerifyEmail,
			},
			wantMockErr: errs.ErrTokenNotFound,
			wantErr:     errs.ErrTokenNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().SaveToken(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("entities.Token"),
			).Return(tt.wantMockErr).Once()

			s := &Service{
				repository: tt.fields.repository,
			}
			_, err := s.CreateToken(tt.args.ctx, tt.args.userId, tt.args.tokenType)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestService_Token(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx   context.Context
		token string
	}

	m := NewMockRepository(t)

	token := entities.Token{
		Token:  uuid.NewString(),
		UserId: uuid.New(),
		Type:   consts.TokenTypeVerifyEmail,
	}

	tests := []struct {
		name        string
		fields      fields
		args        args
		want        entities.Token
		wantMockErr error
		wantErr     error
	}{
		{
			name: "good case",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx:   context.Background(),
				token: token.Token,
			},
			want:        token,
			wantMockErr: nil,
			wantErr:     nil,
		},
		{
			name: "repository error case",
			fields: fields{
				repository: m,
			},
			args: args{
				ctx:   context.Background(),
				token: "invalid token",
			},
			want:        entities.Token{},
			wantMockErr: errs.ErrTokenNotFound,
			wantErr:     errs.ErrTokenNotFound,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().Token(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
			).Return(tt.want, tt.wantMockErr).Once()

			s := &Service{
				repository: tt.fields.repository,
			}
			got, err := s.Token(tt.args.ctx, tt.args.token)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, got, tt.want)
		})
	}
}

func TestService_DeleteToken(t *testing.T) {
	type fields struct {
		repository Repository
	}
	type args struct {
		ctx   context.Context
		token string
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
				ctx:   context.Background(),
				token: uuid.NewString(),
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
				ctx:   context.Background(),
				token: uuid.NewString(),
			},
			wantMockErr: errs.ErrTokenNotFound,
			wantErr:     errs.ErrTokenNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			m.EXPECT().DeleteToken(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
			).Return(tt.wantMockErr).Once()

			s := &Service{
				repository: tt.fields.repository,
			}
			err := s.DeleteToken(tt.args.ctx, tt.args.token)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
