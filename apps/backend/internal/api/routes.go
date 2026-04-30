package api

import "github.com/go-chi/chi/v5"

func (s *Server) RegisterRoutes(r chi.Router) {
	r.HandleFunc("/check", s.corsMiddleware(s.checkHandler))
	r.HandleFunc("/health", s.corsMiddleware(s.healthHandler))

	r.HandleFunc("/api/me", s.corsMiddleware(s.meHandler))
	r.HandleFunc("/api/me/refresh-key", s.corsMiddleware(s.refreshAPIKeyHandler))
}
