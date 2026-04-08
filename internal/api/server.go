package api

import (
	"log/slog"
	"net/http"

	"github.com/wolf-infra/gatus-api-bridge/internal/gatus"
)

type Server struct {
	Manager *gatus.Manager
	logger  *slog.Logger
	version string
	apiKey  string // <-- Store it here!
}

func NewServer(manager *gatus.Manager, logger *slog.Logger, version, apiKey string) *Server {
	return &Server{
		Manager: manager,
		logger:  logger,
		version: version,
		apiKey:  apiKey,
	}
}

// authMiddleware protects our endpoints checking the Bearer token
func (s *Server) authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if s.apiKey != "" {
			authHeader := r.Header.Get("Authorization")
			expectedHeader := "Bearer " + s.apiKey

			if authHeader != expectedHeader {
				s.logger.Warn("Unauthorized access attempt", slog.String("ip", r.RemoteAddr))
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
		}
		next(w, r)
	}
}

func (s *Server) Mount() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/v1/endpoints", s.authMiddleware(s.handleEndpoints))

	return mux
}

func (s *Server) handleEndpoints(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.handleListEndpoints(w, r)
	case http.MethodPost:
		s.handleAddEndpoint(w, r)
	case http.MethodDelete:
		s.handleDeleteEndpoint(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
