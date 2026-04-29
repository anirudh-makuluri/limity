package api

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (s *Server) apiKeyGenerateHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("invalid token: %v", err)})
		return
	}

	userID, err := s.pg.EnsureUserExists(r.Context(), claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to ensure user"})
		return
	}

	key := generateAPIKey()
	id, createdAt, err := s.pg.CreateAPIKey(r.Context(), userID, key)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to create key"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(APIKey{ID: id, UserID: userID, Key: key, CreatedAt: createdAt})
}

func (s *Server) apiKeysListHandler(w http.ResponseWriter, r *http.Request) {
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
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("invalid token: %v", err)})
		return
	}

	userID, err := s.pg.EnsureUserExists(r.Context(), claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to ensure user"})
		return
	}

	keys, err := s.pg.ListAPIKeys(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to fetch keys"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(keys)
}

func (s *Server) apiKeyRevokeHandler(w http.ResponseWriter, r *http.Request) {
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

	keyID := chi.URLParam(r, "id")
	if keyID == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid path"})
		return
	}

	claims, err := extractUserClaimsFromToken(r.Header.Get("Authorization"))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf("invalid token: %v", err)})
		return
	}

	userID, err := s.pg.EnsureUserExists(r.Context(), claims)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to ensure user"})
		return
	}

	revoked, err := s.pg.RevokeAPIKey(r.Context(), keyID, userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "failed to revoke key"})
		return
	}
	if !revoked {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "key not found"})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "key revoked"})
}

func generateAPIKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return "limity_" + hex.EncodeToString(bytes)
}
