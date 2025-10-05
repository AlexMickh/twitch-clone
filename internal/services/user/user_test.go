package user_service

import (
	"context"
	"testing"

	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestService_CreateUser(t *testing.T) {
	type fields struct {
		userRepository UserRepository
		tokenService   TokenService
	}
	type args struct {
		ctx      context.Context
		login    string
		email    string
		password string
	}

	mr := NewMockUserRepository(t)
	ms := NewMockTokenService(t)

	tests := []struct {
		name              string
		fields            fields
		args              args
		wantRepositoryErr error
		wantErr           error
	}{
		{
			name: "good case",
			fields: fields{
				userRepository: mr,
				tokenService:   ms,
			},
			args: args{
				ctx:      context.Background(),
				login:    "some login",
				email:    "example@test.com",
				password: "qwerty123",
			},
			wantRepositoryErr: nil,
			wantErr:           nil,
		},
		{
			name: "repository error case",
			fields: fields{
				userRepository: mr,
				tokenService:   ms,
			},
			args: args{
				ctx:      context.Background(),
				login:    "some login",
				email:    "example@test.com",
				password: "qwerty123",
			},
			wantRepositoryErr: errs.ErrUserAlreadyExists,
			wantErr:           errs.ErrUserAlreadyExists,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mr.EXPECT().SaveUser(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("entities.User"),
			).Return(tt.wantRepositoryErr).Once()

			s := &Service{
				userRepository: tt.fields.userRepository,
				tokenService:   tt.fields.tokenService,
			}
			_, err := s.CreateUser(tt.args.ctx, tt.args.login, tt.args.email, tt.args.password)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestService_UserByEmail(t *testing.T) {
	type fields struct {
		userRepository UserRepository
		tokenService   TokenService
	}
	type args struct {
		ctx   context.Context
		email string
	}

	mr := NewMockUserRepository(t)
	ms := NewMockTokenService(t)

	email := "example@mail.com"
	userId := uuid.New()

	tests := []struct {
		name              string
		fields            fields
		args              args
		wantRepository    entities.User
		want              entities.User
		wantRepositoryErr error
		wantErr           error
	}{
		{
			name: "good case",
			fields: fields{
				userRepository: mr,
				tokenService:   ms,
			},
			args: args{
				ctx:   context.Background(),
				email: email,
			},
			wantRepository: entities.User{
				ID:              userId,
				Login:           "login",
				Email:           email,
				Password:        "1234",
				IsEmailVerified: true,
			},
			want: entities.User{
				ID:              userId,
				Login:           "login",
				Email:           email,
				Password:        "1234",
				IsEmailVerified: true,
			},
			wantRepositoryErr: nil,
			wantErr:           nil,
		},
		{
			name: "email not verify case",
			fields: fields{
				userRepository: mr,
				tokenService:   ms,
			},
			args: args{
				ctx:   context.Background(),
				email: email,
			},
			wantRepository: entities.User{
				ID:              uuid.New(),
				Login:           "login",
				Email:           email,
				Password:        "1234",
				IsEmailVerified: false,
			},
			want:              entities.User{},
			wantRepositoryErr: nil,
			wantErr:           errs.ErrUserEmailNotVerify,
		},
		{
			name: "repository error case",
			fields: fields{
				userRepository: mr,
				tokenService:   ms,
			},
			args: args{
				ctx:   context.Background(),
				email: email,
			},
			wantRepository:    entities.User{},
			want:              entities.User{},
			wantRepositoryErr: errs.ErrUserNotFound,
			wantErr:           errs.ErrUserNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mr.EXPECT().UserByEmail(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
			).Return(tt.wantRepository, tt.wantRepositoryErr).Once()

			s := &Service{
				userRepository: tt.fields.userRepository,
				tokenService:   tt.fields.tokenService,
			}
			got, err := s.UserByEmail(tt.args.ctx, tt.args.email)
			require.ErrorIs(t, err, tt.wantErr)
			require.Equal(t, tt.want, got)
		})
	}
}

func TestService_VerifyEmail(t *testing.T) {
	type args struct {
		ctx context.Context
		req dtos.ValidateEmailRequest
	}

	tests := []struct {
		name                 string
		args                 args
		wantRepositoryErr    error
		wantServiceGetErr    error
		wantServiceDeleteErr error
		wantErr              error
	}{
		{
			name: "good case",
			args: args{
				ctx: context.Background(),
				req: dtos.ValidateEmailRequest{
					Token: uuid.NewString(),
				},
			},
			wantRepositoryErr:    nil,
			wantServiceGetErr:    nil,
			wantServiceDeleteErr: nil,
			wantErr:              nil,
		},
		{
			name: "service get error case",
			args: args{
				ctx: context.Background(),
				req: dtos.ValidateEmailRequest{
					Token: uuid.NewString(),
				},
			},
			wantRepositoryErr:    nil,
			wantServiceGetErr:    errs.ErrTokenNotFound,
			wantServiceDeleteErr: nil,
			wantErr:              errs.ErrTokenNotFound,
		},
		{
			name: "repository error case",
			args: args{
				ctx: context.Background(),
				req: dtos.ValidateEmailRequest{
					Token: uuid.NewString(),
				},
			},
			wantRepositoryErr:    errs.ErrUserNotFound,
			wantServiceGetErr:    nil,
			wantServiceDeleteErr: nil,
			wantErr:              errs.ErrUserNotFound,
		},
		{
			name: "service delete case",
			args: args{
				ctx: context.Background(),
				req: dtos.ValidateEmailRequest{
					Token: uuid.NewString(),
				},
			},
			wantRepositoryErr:    nil,
			wantServiceGetErr:    nil,
			wantServiceDeleteErr: errs.ErrTokenNotFound,
			wantErr:              errs.ErrTokenNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			mr := NewMockUserRepository(t)
			ms := NewMockTokenService(t)

			ms.EXPECT().Token(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
			).Return(entities.Token{}, tt.wantServiceGetErr).Once()

			mr.EXPECT().ValidateEmail(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("uuid.UUID"),
			).Return(tt.wantRepositoryErr).Maybe()

			ms.EXPECT().DeleteToken(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
			).Return(tt.wantServiceDeleteErr).Maybe()

			s := &Service{
				userRepository: mr,
				tokenService:   ms,
			}
			err := s.VerifyEmail(tt.args.ctx, tt.args.req)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
