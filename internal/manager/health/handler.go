package health

import (
	"context"
	"net/http"
	"time"

	"github.com/sirajul777/genieacs-platform/internal/shared/response"
	"github.com/sirajul777/genieacs-platform/internal/shared/version"
)

type DatabasePinger interface {
	Ping(ctx context.Context) error
}

type Handler struct{ database DatabasePinger }

func NewHandler(database DatabasePinger) *Handler { return &Handler{database: database} }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := http.StatusOK
	databaseStatus := "ok"
	if h.database != nil {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()
		if err := h.database.Ping(ctx); err != nil {
			status = http.StatusServiceUnavailable
			databaseStatus = "unavailable"
		}
	}
	response.JSON(w, status, map[string]any{
		"status":   "ok",
		"service":  "manager",
		"version":  version.Version,
		"database": databaseStatus,
	})
}
