package api

import (
	"encoding/json"
	"net/http"
)

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if s.pg == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "database not initialized"})
		return
	}

	if err := s.pg.Ping(r.Context()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"status": "error", "message": "database connection failed"})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
