package api

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/wolf-infra/gatus-api-bridge/internal/gatus"
)

type Server struct {
	Manager *gatus.Manager
	logger  *slog.Logger
	version string
}

func NewServer(manager *gatus.Manager, logger *slog.Logger, version string) *Server {
	return &Server{
		Manager: manager,
		logger:  logger,
		version: version,
	}
}

// authMiddleware protects our endpoints with a simple API Key
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		expectedKey := os.Getenv("API_KEY")
		// If API_KEY is set, enforce it
		if expectedKey != "" && r.Header.Get("X-API-Key") != expectedKey {
			s.logger.Warn("Unauthorized access attempt", slog.String("ip", r.RemoteAddr))
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func (s *Server) Mount() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check remains public for Docker/Traefik
	mux.HandleFunc("/health", s.handleHealth)

	// Wrap the API endpoints in our auth middleware
	mux.HandleFunc("/api/v1/endpoints", s.authMiddleware(s.handleEndpoints))

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
