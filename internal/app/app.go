package app

import (
	"context"
	"go-books-api/internal/config"
	"go-books-api/internal/database/postgres"
	"go-books-api/internal/delivery/http/handler"
	"go-books-api/internal/delivery/http/routes"
	"go-books-api/internal/repository"
	"go-books-api/internal/server"
	"go-books-api/internal/service"
	"go-books-api/pkg/logger"
	"go-books-api/pkg/migrator"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const shutdownTimeout = 5 * time.Second

func Run() {
	logs, err := logger.New()
	if err != nil {
		panic(err)
	}

	ctx := logger.SetToCtx(context.Background(), logs)

	cfg, err := config.New()
	if err != nil {
		logs.Fatal(ctx, "failed to create config", zap.Error(err))
	}

	db, err := postgres.New(&cfg.Postgres)
	if err != nil {
		logs.Fatal(ctx, "failed to create database connection", zap.Error(err))
	}

	err = migrator.Start(cfg)
	if err != nil {
		logs.Fatal(ctx, "failed to start migratoins", zap.Error(err))
	}

	repo := repository.New(db)
	service := service.New(repo)

	app := gin.Default()

	routes.RegistrateRoutes(app, handler.NewHandler(service, logs))

	srv := server.NewServer(&cfg.HTTP, app)

	go func() {
		if err := srv.Run(ctx); err != nil {
			logs.Fatal(ctx, "failed running the server", zap.Error(err))
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	<-c

	ctx, shutdown := context.WithTimeout(ctx, shutdownTimeout)
	defer shutdown()

	if err := srv.Stop(ctx); err != nil {
		logs.Error(ctx, "failed shutting down the server", zap.Error(err))
	}

	if err := db.Close(); err != nil {
		logs.Error(ctx, "failed to close database connection", zap.Error(err))
	}

	logs.Info(ctx, "Server gracefully stopped")
}
