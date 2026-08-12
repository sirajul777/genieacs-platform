package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirajul777/genieacs-platform/internal/agent/config"
	"github.com/sirajul777/genieacs-platform/internal/agent/health"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(cfg config.ServerConfig, logger *zap.Logger) *Server {
	router := chi.NewRouter()
	router.Get("/health", health.NewHandler().ServeHTTP)
	return &Server{
		httpServer: &http.Server{Addr: cfg.Address, Handler: router, ReadHeaderTimeout: 5 * time.Second},
		logger:     logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("agent server listening", zap.String("address", s.httpServer.Addr))
		if err := s.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
			return
		}
		errCh <- nil
	}()
	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return s.httpServer.Shutdown(shutdownCtx)
	case err := <-errCh:
		return err
	}
}
