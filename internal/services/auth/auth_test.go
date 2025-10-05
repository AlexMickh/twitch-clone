package auth_service

import (
	"context"
	"testing"

	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/errs"

	// auth_service_mocks "github.com/AlexMickh/twitch-clone/internal/services/auth/mocks"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func TestService_Register(t *testing.T) {
	type args struct {
		ctx context.Context
		req dtos.RegisterRequest
	}
	tests := []struct {
		name                string
		args                args
		wantUserErr         error
		wantVerificationErr error
		wantTokenErr        error
		wantErr             error
	}{
		{
			name: "good case",
			args: args{
				ctx: context.Background(),
				req: dtos.RegisterRequest{
					Login:    "some login",
					Email:    "test@test.com",
					Password: "test",
				},
			},
			wantUserErr:         nil,
			wantVerificationErr: nil,
			wantTokenErr:        nil,
			wantErr:             nil,
		},
		{
			name: "invalid password case",
			args: args{
				ctx: context.Background(),
				req: dtos.RegisterRequest{
					Login:    "some login",
					Email:    "test@test.com",
					Password: "edfdhgjkwsedgfcghjkwedgfjhgdgfhjwedgfhjasdgcfhgwsdgsdadgwedgweghgdfashgsa",
				},
			},
			wantUserErr:         nil,
			wantVerificationErr: nil,
			wantTokenErr:        nil,
			wantErr:             bcrypt.ErrPasswordTooLong,
		},
		{
			name: "user error case",
			args: args{
				ctx: context.Background(),
				req: dtos.RegisterRequest{
					Login:    "some login",
					Email:    "test@test.com",
					Password: "test",
				},
			},
			wantUserErr:         errs.ErrUserAlreadyExists,
			wantVerificationErr: nil,
			wantTokenErr:        nil,
			wantErr:             errs.ErrUserAlreadyExists,
		},
		{
			name: "token error case",
			args: args{
				ctx: context.Background(),
				req: dtos.RegisterRequest{
					Login:    "some login",
					Email:    "test@test.com",
					Password: "test",
				},
			},
			wantUserErr:         nil,
			wantVerificationErr: nil,
			wantTokenErr:        errs.ErrTokenNotFound,
			wantErr:             errs.ErrTokenNotFound,
		},
		{
			name: "send error case",
			args: args{
				ctx: context.Background(),
				req: dtos.RegisterRequest{
					Login:    "some login",
					Email:    "test@test.com",
					Password: "test",
				},
			},
			wantUserErr:         nil,
			wantVerificationErr: errs.ErrTokenNotFound,
			wantTokenErr:        nil,
			wantErr:             errs.ErrTokenNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mUserService := NewMockUserService(t)
			mVerificationSender := NewMockVerificationSender(t)
			mTokenService := NewMockTokenService(t)

			mUserService.EXPECT().CreateUser(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
			).Return(uuid.New(), tt.wantUserErr).Maybe()

			mTokenService.EXPECT().CreateToken(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("uuid.UUID"),
				mock.AnythingOfType("string"),
			).Return("token", tt.wantTokenErr).Maybe()

			mVerificationSender.EXPECT().SendVerification(
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
				mock.AnythingOfType("string"),
			).Return(tt.wantVerificationErr).Maybe()

			s := &Service{
				userService:        mUserService,
				verificationSender: mVerificationSender,
				tokenService:       mTokenService,
			}
			_, err := s.Register(tt.args.ctx, tt.args.req)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestService_Login(t *testing.T) {
	type args struct {
		ctx       context.Context
		req       dtos.LoginRequest
		userAgent string
	}

	password := "test"
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	tests := []struct {
		name           string
		args           args
		wantUserErr    error
		wantSessionErr error
		wantErr        error
	}{
		{
			name: "good case",
			args: args{
				ctx: context.Background(),
				req: dtos.LoginRequest{
					Email:    "test@test.com",
					Password: password,
				},
				userAgent: "firefox",
			},
			wantUserErr:    nil,
			wantSessionErr: nil,
			wantErr:        nil,
		},
		{
			name: "user error case",
			args: args{
				ctx: context.Background(),
				req: dtos.LoginRequest{
					Email:    "test@test.com",
					Password: password,
				},
				userAgent: "firefox",
			},
			wantUserErr:    errs.ErrUserEmailNotVerify,
			wantSessionErr: nil,
			wantErr:        errs.ErrUserEmailNotVerify,
		},
		{
			name: "invalid password case",
			args: args{
				ctx: context.Background(),
				req: dtos.LoginRequest{
					Email:    "test@test.com",
					Password: "invalid",
				},
				userAgent: "firefox",
			},
			wantUserErr:    nil,
			wantSessionErr: nil,
			wantErr:        errs.ErrUserNotFound,
		},
		{
			name: "session error case",
			args: args{
				ctx: context.Background(),
				req: dtos.LoginRequest{
					Email:    "test@test.com",
					Password: password,
				},
				userAgent: "firefox",
			},
			wantUserErr:    nil,
			wantSessionErr: errs.ErrSessionNotFound,
			wantErr:        errs.ErrSessionNotFound,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mUserService := NewMockUserService(t)
			mSessionService := NewMockSessionService(t)

			mUserService.EXPECT().UserByEmail(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("string"),
			).Return(entities.User{
				Password: string(hash),
			}, tt.wantUserErr).Once()

			mSessionService.EXPECT().CreateSession(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("uuid.UUID"),
				mock.AnythingOfType("string"),
			).Return(uuid.New(), tt.wantSessionErr).Maybe()

			s := &Service{
				userService:    mUserService,
				sessionService: mSessionService,
			}
			_, err := s.Login(tt.args.ctx, tt.args.req, tt.args.userAgent)
			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
