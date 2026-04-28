package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type CheckRequest struct {
	Key    string `json:"key"`
	Limit  int    `json:"limit"`
	Window int    `json:"window"`
}

type CheckResponse struct {
	Allowed   bool   `json:"allowed"`
	Remaining int    `json:"remaining"`
	Reset     int64  `json:"reset"`
}

type RedisResponse struct {
	Result interface{} `json:"result"`
	Error  *string     `json:"error"`
}

func main() {
	http.HandleFunc("/check", checkHandler)
	http.HandleFunc("/health", healthHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on :%s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func checkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(map[string]string{"error": "method not allowed"})
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

	// Rate limit check
	allowed, remaining, reset, err := rateLimitCheck(req.Key, req.Limit, req.Window)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(CheckResponse{
		Allowed:   allowed,
		Remaining: remaining,
		Reset:     reset,
	})
}

func rateLimitCheck(key string, limit int, window int) (bool, int, int64, error) {
	now := time.Now().Unix()
	windowStart := now - (now % int64(window))
	reset := windowStart + int64(window)

	redisKey := fmt.Sprintf("ratelimit:%s:%d", key, windowStart)

	// Increment the counter
	count, err := redisIncr(redisKey)
	if err != nil {
		return false, 0, 0, err
	}

	// Set expiry on first increment
	if count == 1 {
		if err := redisExpire(redisKey, window); err != nil {
			return false, 0, 0, err
		}
	}

	// Check if limit exceeded
	allowed := count <= limit
	remaining := max(0, limit-count)

	return allowed, remaining, reset, nil
}

func redisIncr(key string) (int, error) {
	redisURL := os.Getenv("UPSTASH_REDIS_REST_URL")
	redisToken := os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if redisURL == "" || redisToken == "" {
		return 0, fmt.Errorf("missing redis configuration")
	}

	url := fmt.Sprintf("%s/incr/%s", redisURL, key)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", redisToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var redisResp RedisResponse
	if err := json.Unmarshal(body, &redisResp); err != nil {
		return 0, err
	}

	if redisResp.Error != nil {
		return 0, fmt.Errorf("redis error: %s", *redisResp.Error)
	}

	// Redis returns the count as a float64
	count, ok := redisResp.Result.(float64)
	if !ok {
		return 0, fmt.Errorf("unexpected redis response type")
	}

	return int(count), nil
}

func redisExpire(key string, seconds int) error {
	redisURL := os.Getenv("UPSTASH_REDIS_REST_URL")
	redisToken := os.Getenv("UPSTASH_REDIS_REST_TOKEN")

	if redisURL == "" || redisToken == "" {
		return fmt.Errorf("missing redis configuration")
	}

	url := fmt.Sprintf("%s/expire/%s/%d", redisURL, key, seconds)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", redisToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var redisResp RedisResponse
	if err := json.Unmarshal(body, &redisResp); err != nil {
		return err
	}

	if redisResp.Error != nil {
		return fmt.Errorf("redis error: %s", *redisResp.Error)
	}

	return nil
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
