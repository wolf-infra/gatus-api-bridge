package api

import (
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) handleListEndpoints(w http.ResponseWriter, r *http.Request) {
	groupFilter := r.URL.Query().Get("group")

	endpoints, err := s.Manager.GetEndpoints(groupFilter)
	if err != nil {
		log.Printf("Failed to list endpoints: %v", err)
		http.Error(w, "Failed to read endpoints", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(endpoints)
}
