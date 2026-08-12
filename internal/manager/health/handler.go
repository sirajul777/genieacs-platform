package health

import (
	"net/http"

	"github.com/sirajul777/genieacs-platform/internal/shared/response"
	"github.com/sirajul777/genieacs-platform/internal/shared/version"
)

type Handler struct{}

func NewHandler() *Handler { return &Handler{} }

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "manager",
		"version": version.Version,
	})
}
