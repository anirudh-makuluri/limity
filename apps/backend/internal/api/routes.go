package api

import "github.com/go-chi/chi/v5"

func (s *Server) RegisterRoutes(r chi.Router) {
	r.HandleFunc("/check", s.corsMiddleware(s.checkHandler))
	r.HandleFunc("/health", s.corsMiddleware(s.healthHandler))

	r.HandleFunc("/api/keys/generate", s.corsMiddleware(s.apiKeyGenerateHandler))
	r.HandleFunc("/api/keys", s.corsMiddleware(s.apiKeysListHandler))
	r.HandleFunc("/api/keys/{id}/revoke", s.corsMiddleware(s.apiKeyRevokeHandler))
}
