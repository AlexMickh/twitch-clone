package app

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/AlexMickh/twitch-clone/internal/config"
	"github.com/AlexMickh/twitch-clone/internal/lib/email"
	token_repository "github.com/AlexMickh/twitch-clone/internal/repository/mongo/token"
	user_repository "github.com/AlexMickh/twitch-clone/internal/repository/mongo/user"
	session_repository "github.com/AlexMickh/twitch-clone/internal/repository/redis/session"
	"github.com/AlexMickh/twitch-clone/internal/server"
	auth_service "github.com/AlexMickh/twitch-clone/internal/services/auth"
	session_service "github.com/AlexMickh/twitch-clone/internal/services/session"
	token_service "github.com/AlexMickh/twitch-clone/internal/services/token"
	user_service "github.com/AlexMickh/twitch-clone/internal/services/user"
	"github.com/AlexMickh/twitch-clone/pkg/clients/mongodb"
	redis_client "github.com/AlexMickh/twitch-clone/pkg/clients/redis"
	"github.com/AlexMickh/twitch-clone/pkg/logger"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type App struct {
	cfg  *config.Config
	db   *mongo.Client
	cash *redis.Client
	srv  *server.Server
}

func New(ctx context.Context, cfg *config.Config) *App {
	const op = "app.New"
	log := logger.FromCtx(ctx).With(slog.String("op", op))

	log.Info("initing mongo")
	db, err := mongodb.New(
		ctx,
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
	)
	if err != nil {
		log.Error("failed to init mongo", logger.Err(err))
		os.Exit(1)
	}

	userRepository, err := user_repository.New(ctx, db, cfg.DB.Database, cfg.DB.Collections["users"])
	if err != nil {
		log.Error("failed to init mongo", logger.Err(err))
		os.Exit(1)
	}

	tokenRepository, err := token_repository.New(ctx, db, cfg.DB.Database, cfg.DB.Collections["tokens"])
	if err != nil {
		log.Error("failed to init mongo", logger.Err(err))
		os.Exit(1)
	}

	log.Info("initing redis")
	cash, err := redis_client.New(
		ctx,
		fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		cfg.Redis.User,
		cfg.Redis.Password,
		cfg.Redis.DB,
	)
	if err != nil {
		log.Error("failed to init redis", logger.Err(err))
		os.Exit(1)
	}
	sessionRepository := session_repository.New(cash, cfg.Redis.Expiration)

	mailService := email.New(cfg.Mail)

	log.Info("initing service layer")
	tokenService := token_service.New(tokenRepository)
	userService := user_service.New(userRepository, tokenService)
	sessionService := session_service.New(sessionRepository)
	authService := auth_service.New(userService, mailService, tokenService, sessionService)

	log.Info("initing server")
	srv := server.New(ctx, cfg.Server, authService, userService, sessionService)

	return &App{
		cfg:  cfg,
		db:   db,
		srv:  srv,
		cash: cash,
	}
}

func (a *App) Run(ctx context.Context) {
	const op = "app.Run"
	log := logger.FromCtx(ctx).With(slog.String("op", op))

	log.Info("server started", slog.String("addr", a.cfg.Server.Addr))

	go func() {
		if err := a.srv.Run(); err != nil {
			log.Error("failed to start server", logger.Err(err))
			os.Exit(1)
		}
	}()
}

func (a *App) Close(ctx context.Context) {
	_ = a.srv.GracefulStop(ctx)
	_ = a.db.Disconnect(ctx)
	_ = a.cash.Close()
}
