package api

import (
	"log/slog"
	"net/http"

	"github.com/wolf-infra/gatus-api-bridge/internal/gatus"
)

type Server struct {
	Manager *gatus.Manager
	logger  *slog.Logger
}

func NewServer(manager *gatus.Manager, logger *slog.Logger) *Server {
	return &Server{
		Manager: manager,
		logger:  logger,
	}
}

func (s *Server) Mount() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/v1/endpoints", s.handleEndpoints)
	return mux
}

func (s *Server) handleEndpoints(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleListEndpoints(w, r)
	case http.MethodPost:
		s.handleAddEndpoint(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
