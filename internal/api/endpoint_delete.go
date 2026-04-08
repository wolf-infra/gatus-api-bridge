package api

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func (s *Server) handleDeleteEndpoint(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	group := r.URL.Query().Get("group")

	if name == "" {
		http.Error(w, "name parameter is required", http.StatusBadRequest)
		return
	}
	if group == "" {
		group = "Infrastructure" // Fallback to our default
	}

	deleted, err := s.Manager.DeleteEndpoint(name, group)
	if err != nil {
		s.logger.Error("Internal error deleting endpoint", slog.Any("error", err))
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !deleted {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "not found"}); err != nil {
			s.logger.Error("Failed to encode response", slog.Any("error", err))
		}
		return
	}

	s.logger.Info("Endpoint deleted via API", slog.String("name", name))
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "deleted"}); err != nil {
		s.logger.Error("Failed to encode response", slog.Any("error", err))
	}
}
