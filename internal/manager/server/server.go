package server

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirajul777/genieacs-platform/internal/manager/config"
	"github.com/sirajul777/genieacs-platform/internal/manager/health"
	managerhttp "github.com/sirajul777/genieacs-platform/internal/manager/http"
	"github.com/sirajul777/genieacs-platform/internal/manager/repository/postgres"
	agentusecase "github.com/sirajul777/genieacs-platform/internal/manager/usecase/agent"
	"go.uber.org/zap"
)

type Server struct {
	httpServer *http.Server
	logger     *zap.Logger
}

func New(cfg config.ServerConfig, logger *zap.Logger, database *pgxpool.Pool) *Server {
	router := chi.NewRouter()
	router.Get("/health", health.NewHandler(database).ServeHTTP)
	agentHandler := managerhttp.NewAgentHandler(agentusecase.NewUseCase(postgres.NewAgentRepository(database), postgres.NewHeartbeatRepository(database)))
	router.Post("/api/v1/agents/register", agentHandler.Register)
	router.Post("/api/v1/agents/heartbeat", agentHandler.Heartbeat)
	return &Server{
		httpServer: &http.Server{Addr: cfg.Address, Handler: router, ReadHeaderTimeout: 5 * time.Second},
		logger:     logger,
	}
}

func (s *Server) Start(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		s.logger.Info("manager server listening", zap.String("address", s.httpServer.Addr))
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
