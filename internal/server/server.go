package server

import (
	"context"
	"errors"
	"fmt"
	"go-books-api/internal/config"
	"go-books-api/pkg/logger"
	"net/http"

	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
}

func NewServer(cfg *config.HTTPConfig, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:           ":" + cfg.Port,
			Handler:        handler,
			ReadTimeout:    cfg.ReadTimeout,
			WriteTimeout:   cfg.WriteTimeout,
			MaxHeaderBytes: cfg.MaxHeaderMegabytes << 20,
		},
	}
}

func (s *Server) Run(ctx context.Context) error {
	log := logger.GetLoggerFromCtx(ctx)
	log.Info(ctx, "Server is running", zap.String("port", s.httpServer.Addr))

	err := s.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("failed to listen and serve: %w", err)
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
