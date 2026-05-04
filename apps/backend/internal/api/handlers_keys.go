package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (s *Server) meHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	if s.pg == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database not initialized"})
		return
	}

	claims, err := extractUserClaimsFromToken(r.Header.Get("Authorization"))
	if err != nil {
		if s.metrics != nil {
			s.metrics.authFailuresTotal.Inc()
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("invalid token: %v", err)})
		return
	}
	setOwnerUserIDFromVerifiedAuth(r.Context(), claims.Sub)

	userProfile, err := s.pg.EnsureUserWithAPIKey(r.Context(), claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to ensure user"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userProfile)
}

func (s *Server) refreshAPIKeyHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	if s.pg == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "database not initialized"})
		return
	}

	claims, err := extractUserClaimsFromToken(r.Header.Get("Authorization"))
	if err != nil {
		if s.metrics != nil {
			s.metrics.authFailuresTotal.Inc()
		}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("invalid token: %v", err)})
		return
	}
	setOwnerUserIDFromVerifiedAuth(r.Context(), claims.Sub)

	userProfile, err := s.pg.EnsureUserWithAPIKey(r.Context(), claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to ensure user"})
		return
	}

	newKey, err := s.pg.RotateAPIKey(r.Context(), userProfile.ID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to refresh api key"})
		return
	}

	userProfile.APIKey = newKey
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userProfile)
}
