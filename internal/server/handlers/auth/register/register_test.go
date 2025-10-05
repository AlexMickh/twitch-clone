package register

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/errs"
	"github.com/AlexMickh/twitch-clone/pkg/api"
	"github.com/google/uuid"
	mock "github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRegister_New(t *testing.T) {
	successID := uuid.New().String()

	cases := []struct {
		name               string
		req                dtos.RegisterRequest
		respStatus         int
		respMessage        string
		wantRegisterError  error
		wantRegisterReturn string
	}{
		{
			name: "good case",
			req: dtos.RegisterRequest{
				Login:    "test",
				Email:    "test@test.com",
				Password: "test",
			},
			respStatus:         http.StatusCreated,
			respMessage:        successID,
			wantRegisterError:  nil,
			wantRegisterReturn: successID,
		},
		{
			name: "invalid request case",
			req: dtos.RegisterRequest{
				Login:    `"`,
				Email:    "test@test.com",
				Password: "test",
			},
			respStatus:         http.StatusBadRequest,
			respMessage:        "failed to decode body",
			wantRegisterError:  nil,
			wantRegisterReturn: successID,
		},
		{
			name: "invalid login case",
			req: dtos.RegisterRequest{
				Login:    "t",
				Email:    "test@test.com",
				Password: "test",
			},
			respStatus:         http.StatusBadRequest,
			respMessage:        "failed to validate body",
			wantRegisterError:  nil,
			wantRegisterReturn: successID,
		},
		{
			name: "invalid email case",
			req: dtos.RegisterRequest{
				Login:    "test",
				Email:    "test",
				Password: "test",
			},
			respStatus:         http.StatusBadRequest,
			respMessage:        "failed to validate body",
			wantRegisterError:  nil,
			wantRegisterReturn: successID,
		},
		{
			name: "invalid password case",
			req: dtos.RegisterRequest{
				Login:    "test",
				Email:    "test@test.com",
				Password: "t",
			},
			respStatus:         http.StatusBadRequest,
			respMessage:        "failed to validate body",
			wantRegisterError:  nil,
			wantRegisterReturn: successID,
		},
		{
			name: "user already exists case",
			req: dtos.RegisterRequest{
				Login:    "test",
				Email:    "test@test.com",
				Password: "test",
			},
			respStatus:         http.StatusBadRequest,
			respMessage:        errs.ErrUserAlreadyExists.Error(),
			wantRegisterError:  errs.ErrUserAlreadyExists,
			wantRegisterReturn: "",
		},
		{
			name: "register error case",
			req: dtos.RegisterRequest{
				Login:    "test",
				Email:    "test@test.com",
				Password: "test",
			},
			respStatus:         http.StatusInternalServerError,
			respMessage:        "failed to register user",
			wantRegisterError:  errors.New("some error"),
			wantRegisterReturn: "",
		},
	}

	for _, tt := range cases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mRegister := NewMockRegisterer(t)

			mRegister.EXPECT().Register(
				mock.AnythingOfType("context.backgroundCtx"),
				mock.AnythingOfType("dtos.RegisterRequest"),
			).Return(tt.wantRegisterReturn, tt.wantRegisterError).Maybe()

			handler := api.ErrorWrapper(New(mRegister))

			input := fmt.Sprintf(`{"email": "%s", "login": "%s", "password": "%s"}`, tt.req.Email, tt.req.Login, tt.req.Password)

			req, err := http.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, tt.respStatus, rr.Code)

			if tt.respStatus == http.StatusCreated {
				var resp dtos.RegisterResponse
				err = json.NewDecoder(rr.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, tt.respMessage, resp.ID)
			} else {
				var resp api.ErrorResponse
				err = json.NewDecoder(rr.Body).Decode(&resp)
				require.NoError(t, err)

				require.Equal(t, tt.respMessage, resp.Error)
			}
		})
	}
}
