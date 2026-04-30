package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

type supabaseUserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func extractBearerToken(authHeader string) (string, error) {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return "", fmt.Errorf("missing Authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)
	if token == authHeader || token == "" {
		return "", fmt.Errorf("invalid Authorization header format (expected 'Bearer <token>')")
	}

	return token, nil
}

func extractUserClaimsFromToken(authHeader string) (*TokenClaims, error) {
	token, err := extractBearerToken(authHeader)
	if err != nil {
		return nil, err
	}

	supabaseURL := strings.TrimSpace(os.Getenv("SUPABASE_URL"))
	supabaseAnonKey := strings.TrimSpace(os.Getenv("SUPABASE_ANON_KEY"))
	if supabaseURL == "" || supabaseAnonKey == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_ANON_KEY are required")
	}

	req, err := http.NewRequest(http.MethodGet, strings.TrimRight(supabaseURL, "/")+"/auth/v1/user", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("apikey", supabaseAnonKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to verify token with supabase: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("supabase token verification failed with status %d", resp.StatusCode)
	}

	var user supabaseUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode supabase user response: %w", err)
	}
	if user.ID == "" {
		return nil, fmt.Errorf("supabase user response missing id")
	}

	return &TokenClaims{
		Sub:   user.ID,
		Email: user.Email,
	}, nil
}
