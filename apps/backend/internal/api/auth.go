package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

func extractUserClaimsFromToken(authHeader string) (*TokenClaims, error) {
	authHeader = strings.TrimSpace(authHeader)
	if authHeader == "" {
		return nil, fmt.Errorf("missing Authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	token = strings.TrimSpace(token)
	if token == authHeader || token == "" {
		return nil, fmt.Errorf("invalid Authorization header format (expected 'Bearer <token>')")
	}

	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid JWT format: expected 3 parts, got %d", len(parts))
	}

	payload := parts[1]
	padding := 4 - len(payload)%4
	if padding != 4 {
		payload += strings.Repeat("=", padding)
	}

	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token payload: %v", err)
	}

	var claims TokenClaims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, fmt.Errorf("failed to parse token claims: %v", err)
	}

	if claims.Sub == "" {
		return nil, fmt.Errorf("no sub claim in token")
	}

	return &claims, nil
}
