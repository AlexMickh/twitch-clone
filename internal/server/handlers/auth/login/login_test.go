package login

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexMickh/twitch-clone/internal/config"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/AlexMickh/twitch-clone/pkg/api"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestLogin_New(t *testing.T) {
	cases := []struct {
		name           string
		email          string
		password       string
		respStatus     int
		respMessage    string
		wantLoginError error
	}{
		{
			name:           "good case",
			email:          "test@test.com",
			password:       "qwerty",
			respStatus:     http.StatusCreated,
			respMessage:    "",
			wantLoginError: nil,
		},
		{
			name:           "invalid request case",
			email:          "test@test.com",
			password:       `"`,
			respStatus:     http.StatusBadRequest,
			respMessage:    "failed to decode body",
			wantLoginError: nil,
		},
		{
			name:           "invalid email case",
			email:          "test",
			password:       "qwerty",
			respStatus:     http.StatusBadRequest,
			respMessage:    "failed to validate body",
			wantLoginError: nil,
		},
		{
			name:           "invalid password case",
			email:          "test@test.com",
			password:       "1",
			respStatus:     http.StatusBadRequest,
			respMessage:    "failed to validate body",
			wantLoginError: nil,
		},
		{
			name:           "user not found case",
			email:          "test@test.com",
			password:       "qwerty",
			respStatus:     http.StatusNotFound,
			respMessage:    errs.ErrUserNotFound.Error(),
			wantLoginError: errs.ErrUserNotFound,
		},
		{
			name:           "email not verify case",
			email:          "test@test.com",
			password:       "qwerty",
			respStatus:     http.StatusForbidden,
			respMessage:    errs.ErrUserEmailNotVerify.Error(),
			wantLoginError: errs.ErrUserEmailNotVerify,
		},
		{
			name:           "login error case",
			email:          "test@test.com",
			password:       "qwerty",
			respStatus:     http.StatusInternalServerError,
			respMessage:    "failed to login user",
			wantLoginError: errors.New("some error"),
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mLogin := NewMockLoginer(t)

			mLogin.EXPECT().Login(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("dtos.LoginRequest"),
				mock.AnythingOfType("string"),
			).Return("some id", tt.wantLoginError).Maybe()

			handler := api.ErrorWrapper(New(mLogin, config.SessionConfig{
				Name:     "session",
				HttpOnly: true,
				Secure:   false,
				MaxAge:   3600,
			}))

			input := fmt.Sprintf(`{"email": "%s", "password": "%s"}`, tt.email, tt.password)

			req, err := http.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.respStatus, rr.Code)

			if tt.respStatus >= 400 {
				var resp api.ErrorResponse
				err = json.NewDecoder(rr.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, tt.respMessage, resp.Error)
			}
		})
	}
}
