package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{
		"status":  "healthy",
		"service": "gatus-api-bridge",
		"version": "1.0.0",
	}); err != nil {
		s.logger.Error("Failed to encode health response", slog.Any("error", err))
	}
}
