package postgres

import (
	"crypto/rand"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/limity/backend/internal/api"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func generateAPIKey() string {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		panic(err)
	}
	return "limity_" + hex.EncodeToString(bytes)
}

func (s *Store) Ping(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("database not initialized")
	}
	return s.db.PingContext(ctx)
}

func (s *Store) EnsureUserWithAPIKey(ctx context.Context, claims *api.TokenClaims) (*api.UserProfile, error) {
	if s.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	if claims == nil || claims.Sub == "" {
		return nil, fmt.Errorf("invalid user claims")
	}

	userUUID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(claims.Sub)).String()
	email := strings.TrimSpace(claims.Email)
	if email == "" {
		email = claims.Sub + "@unknown.local"
	}
	defaultKey := generateAPIKey()

	profile := &api.UserProfile{}
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO users (id, email, external_user_id, api_key)
		 VALUES ($1, $2, $3, $4)
		 ON CONFLICT (external_user_id)
		 DO UPDATE SET email = EXCLUDED.email
		 RETURNING id, external_user_id, email, api_key, created_at`,
		userUUID, email, claims.Sub, defaultKey,
	).Scan(&profile.ID, &profile.ExternalUserID, &profile.Email, &profile.APIKey, &profile.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to ensure user exists: %w", err)
	}

	// Backward-compat for migrated rows where api_key might still be null/empty.
	if strings.TrimSpace(profile.APIKey) == "" {
		newKey := generateAPIKey()
		err = s.db.QueryRowContext(
			ctx,
			`UPDATE users SET api_key = $2 WHERE id = $1 RETURNING api_key`,
			profile.ID, newKey,
		).Scan(&profile.APIKey)
		if err != nil {
			return nil, fmt.Errorf("failed to backfill user api_key: %w", err)
		}
	}
	return profile, nil
}
