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
		if s.metrics != nil {
			s.metrics.checkTotal.WithLabelValues("method_not_allowed").Inc()
		}
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
		return
	}

	if s.redis == nil {
		if s.metrics != nil {
			s.metrics.checkTotal.WithLabelValues("error").Inc()
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "redis store not initialized"})
		return
	}

	var req CheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		if s.metrics != nil {
			s.metrics.checkTotal.WithLabelValues("bad_request").Inc()
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid request"})
		return
	}

	if req.Key == "" || req.Limit == 0 || req.Window == 0 {
		if s.metrics != nil {
			s.metrics.checkTotal.WithLabelValues("bad_request").Inc()
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "missing required fields"})
		return
	}
	// Attribution key for analytics should come from Bearer API key used by SDK hosted mode.
	if bearerKey, err := extractBearerToken(r.Header.Get("Authorization")); err == nil && bearerKey != "" {
		setAPIKeyFromCheckRequest(r.Context(), bearerKey)
	} else if len(req.Key) > 7 && req.Key[:7] == "limity_" {
		// Backward-compat for manual testing where API key may be passed as body key.
		setAPIKeyFromCheckRequest(r.Context(), req.Key)
	}

	allowed, remaining, reset, err := s.rateLimitCheck(r.Context(), req.Key, req.Limit, req.Window)
	if err != nil {
		if s.metrics != nil {
			s.metrics.checkTotal.WithLabelValues("error").Inc()
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}
	if s.metrics != nil {
		if allowed {
			s.metrics.checkTotal.WithLabelValues("allowed").Inc()
		} else {
			s.metrics.checkTotal.WithLabelValues("blocked").Inc()
		}
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
		if s.metrics != nil {
			s.metrics.redisErrorsTotal.Inc()
		}
		return false, 0, 0, err
	}

	if count == 1 {
		if err := s.redis.Expire(ctx, redisKey, window); err != nil {
			if s.metrics != nil {
				s.metrics.redisErrorsTotal.Inc()
			}
			return false, 0, 0, err
		}
	}

	allowed := count <= limit
	remaining := max(0, limit-count)

	return allowed, remaining, reset, nil
}
