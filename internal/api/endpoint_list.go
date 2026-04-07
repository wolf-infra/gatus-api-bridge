package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) handleListEndpoints(w http.ResponseWriter, r *http.Request) {
	groupFilter := r.URL.Query().Get("group")

	endpoints, err := s.Manager.GetEndpoints(groupFilter)
	if err != nil {
		s.logger.Error("Failed to list endpoints", slog.Any("error", err))
		http.Error(w, "Failed to read endpoints", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(endpoints); err != nil {
		s.logger.Error("Failed to encode response", slog.Any("error", err))
	}
}
