package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func (s *Server) checkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	if s.redis == nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "redis store not initialized"})
		return
	}

	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	if req.Key == "" || req.Limit == 0 || req.Window == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing required fields"})
		return
	}

	allowed, remaining, reset, err := s.rateLimitCheck(r.Context(), req.Key, req.Limit, req.Window)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CheckResponse{Allowed: allowed, Remaining: remaining, Reset: reset})
}

func (s *Server) rateLimitCheck(ctx context.Context, key string, limit int, window int) (bool, int, int64, error) {
	now := time.Now().Unix()
	windowStart := now - (now % int64(window))
	reset := windowStart + int64(window)

	redisKey := fmt.Sprintf("ratelimit:%s:%d", key, windowStart)

	count, err := s.redis.Incr(ctx, redisKey)
	if err != nil {
		return false, 0, 0, err
	}

	if count == 1 {
		if err := s.redis.Expire(ctx, redisKey, window); err != nil {
			return false, 0, 0, err
		}
	}

	allowed := count <= limit
	remaining := max(0, limit-count)

	return allowed, remaining, reset, nil
}
