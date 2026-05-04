package postgres

import (
	"crypto/rand"
	"context"
	"database/sql"
	"encoding/hex"
	"fmt"
	"errors"
	"strconv"
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
	bytes := make([]byte, 16)
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

func (s *Store) RotateAPIKey(ctx context.Context, userID string) (string, error) {
	if s.db == nil {
		return "", fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(userID) == "" {
		return "", fmt.Errorf("invalid user id")
	}

	newKey := generateAPIKey()
	var key string
	err := s.db.QueryRowContext(
		ctx,
		`UPDATE users SET api_key = $2 WHERE id = $1 RETURNING api_key`,
		userID, newKey,
	).Scan(&key)
	if err != nil {
		return "", fmt.Errorf("failed to rotate api key: %w", err)
	}
	return key, nil
}

func (s *Store) CountAPIKeys(ctx context.Context) (int, error) {
	if s.db == nil {
		return 0, fmt.Errorf("database not initialized")
	}
	var count int
	if err := s.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE api_key IS NOT NULL AND api_key <> ''`).Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count api keys: %w", err)
	}
	return count, nil
}

func (s *Store) InsertRequestEvents(ctx context.Context, events []api.RequestEvent) error {
	if s.db == nil {
		return fmt.Errorf("database not initialized")
	}
	if len(events) == 0 {
		return nil
	}

	var b strings.Builder
	args := make([]any, 0, len(events)*9)
	b.WriteString(`INSERT INTO request_events (timestamp, method, route, path, status_code, duration_ms, client_ip, user_agent, owner_user_id) VALUES `)

	for i, e := range events {
		if i > 0 {
			b.WriteString(",")
		}
		offset := i*9 + 1
		b.WriteString("(")
		b.WriteString("$" + strconv.Itoa(offset))
		b.WriteString(",$" + strconv.Itoa(offset+1))
		b.WriteString(",$" + strconv.Itoa(offset+2))
		b.WriteString(",$" + strconv.Itoa(offset+3))
		b.WriteString(",$" + strconv.Itoa(offset+4))
		b.WriteString(",$" + strconv.Itoa(offset+5))
		b.WriteString(",$" + strconv.Itoa(offset+6))
		b.WriteString(",$" + strconv.Itoa(offset+7))
		b.WriteString(",$" + strconv.Itoa(offset+8))
		b.WriteString(")")

		args = append(args,
			e.Timestamp,
			e.Method,
			e.Route,
			e.Path,
			e.StatusCode,
			e.DurationMs,
			e.ClientIP,
			e.UserAgent,
			e.OwnerUserID,
		)
	}

	if _, err := s.db.ExecContext(ctx, b.String(), args...); err != nil {
		return fmt.Errorf("failed to insert request events: %w", err)
	}
	return nil
}

func (s *Store) ResolveOwnerUserIDByAPIKey(ctx context.Context, apiKey string) (string, error) {
	if s.db == nil {
		return "", fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(apiKey) == "" {
		return "", nil
	}

	var ownerUserID string
	err := s.db.QueryRowContext(ctx, `SELECT external_user_id FROM users WHERE api_key = $1`, apiKey).Scan(&ownerUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", fmt.Errorf("failed to resolve owner user id by api key: %w", err)
	}
	return ownerUserID, nil
}
