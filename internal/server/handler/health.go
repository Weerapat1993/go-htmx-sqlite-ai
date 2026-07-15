package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Weerapat1993/go-htmx-sqlite-ai/internal/version"
)

// Health returns a health check response, including DB connectivity when a database is wired.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	status := "ok"
	statusCode := http.StatusOK

	if h.database != nil {
		if err := h.database.DB().PingContext(r.Context()); err != nil {
			h.logger.Error("health check: database ping failed", "error", err)
			status = "degraded"
			statusCode = http.StatusServiceUnavailable
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(healthResponse{Status: status, Version: version.Value}); err != nil {
		h.logger.Error("failed to encode health response", "error", err)
	}
}

type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}
