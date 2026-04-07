package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/wolf-infra/gatus-api-bridge/internal/gatus"
)

func (s *Server) handleAddEndpoint(w http.ResponseWriter, r *http.Request) {
	var ep gatus.Endpoint
	if err := json.NewDecoder(r.Body).Decode(&ep); err != nil {
		s.logger.Error("Failed to decode JSON payload", slog.Any("error", err))
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	if ep.Interval == "" {
		ep.Interval = "60s"
	}
	if len(ep.Conditions) == 0 {
		ep.Conditions = []string{"[STATUS] == 200"}
	}

	added, err := s.Manager.AddEndpoint(ep)
	if err != nil {
		s.logger.Error("Internal error adding endpoint",
			slog.String("name", ep.Name),
			slog.Any("error", err),
		)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if !added {
		s.logger.Info("Endpoint already exists, skipping",
			slog.String("name", ep.Name),
			slog.String("group", ep.Group),
		)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{"status": "exists"}); err != nil {
			s.logger.Error("Failed to encode response", slog.Any("error", err))
		}
		return
	}

	s.logger.Info("Endpoint created via API",
		slog.String("name", ep.Name),
		slog.String("group", ep.Group),
	)
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"status": "created"}); err != nil {
		s.logger.Error("Failed to encode response", slog.Any("error", err))
	}
}
