package server

import (
	"context"
	"fmt"
	"net/http"

	_ "github.com/AlexMickh/twitch-clone/docs"
	"github.com/AlexMickh/twitch-clone/internal/config"
	"github.com/AlexMickh/twitch-clone/internal/dtos"
	"github.com/AlexMickh/twitch-clone/internal/entities"
	"github.com/AlexMickh/twitch-clone/internal/server/handlers/auth/login"
	"github.com/AlexMickh/twitch-clone/internal/server/handlers/auth/register"
	"github.com/AlexMickh/twitch-clone/internal/server/handlers/session/current_session"
	"github.com/AlexMickh/twitch-clone/internal/server/handlers/user/verify_email"
	"github.com/AlexMickh/twitch-clone/internal/server/middlewares"
	"github.com/AlexMickh/twitch-clone/pkg/api"
	"github.com/AlexMickh/twitch-clone/pkg/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type Server struct {
	srv *http.Server
}

type AuthService interface {
	Register(ctx context.Context, req dtos.RegisterRequest) (string, error)
	Login(ctx context.Context, req dtos.LoginRequest, userAgent string) (string, error)
}

type UserService interface {
	VerifyEmail(ctx context.Context, req dtos.ValidateEmailRequest) error
}

type SessionService interface {
	SessionById(ctx context.Context, sessionId string) (entities.Session, error)
	ValidateSession(ctx context.Context, sessionId string) (uuid.UUID, error)
}

// @title						Your API
// @version					1.0
// @description				Your API description
// @securityDefinitions.apikey	SessionAuth
// @in							cookie
// @name						session_id
func New(
	ctx context.Context,
	cfg config.ServerConfig,
	authService AuthService,
	userService UserService,
	sessionService SessionService,
) *Server {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(logger.ChiMiddleware(ctx))
	r.Use(middleware.Recoverer)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	}).Handler)

	// validator := validator.New(validator.WithRequiredStructEnabled())

	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://%s/swagger/doc.json", cfg.Addr)), //The url pointing to API definition
	))

	r.Get("/health-check", api.ErrorWrapper(func(w http.ResponseWriter, r *http.Request) error {
		w.WriteHeader(200)
		return nil
	}))

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", api.ErrorWrapper(register.New(authService)))
		r.Post("/login", api.ErrorWrapper(login.New(authService, cfg.Session)))
	})

	r.Route("/user", func(r chi.Router) {
		r.Get("/verify-email/{token}", api.ErrorWrapper(verify_email.New(userService)))
	})

	r.Route("/session", func(r chi.Router) {
		r.Use(middlewares.Auth(cfg.Session, sessionService))
		r.Get("/current", api.ErrorWrapper(current_session.New(sessionService, cfg.Session)))
	})

	return &Server{
		srv: &http.Server{
			Addr:         cfg.Addr,
			Handler:      r,
			ReadTimeout:  cfg.Timeout,
			WriteTimeout: cfg.Timeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
	}
}

func (s *Server) Run() error {
	const op = "server.Run"

	if err := s.srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Server) GracefulStop(ctx context.Context) error {
	const op = "server.GracefulStop"

	if err := s.srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
