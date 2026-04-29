package postgres

import (
	"context"
	"database/sql"
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

func (s *Store) Ping(ctx context.Context) error {
	if s.db == nil {
		return fmt.Errorf("database not initialized")
	}
	return s.db.PingContext(ctx)
}

func (s *Store) EnsureUserExists(ctx context.Context, claims *api.TokenClaims) (string, error) {
	if s.db == nil {
		return "", fmt.Errorf("database not initialized")
	}
	if claims == nil || claims.Sub == "" {
		return "", fmt.Errorf("invalid user claims")
	}

	userUUID := uuid.NewSHA1(uuid.NameSpaceOID, []byte(claims.Sub)).String()
	email := strings.TrimSpace(claims.Email)
	if email == "" {
		email = claims.Sub + "@unknown.local"
	}

	var id string
	err := s.db.QueryRowContext(
		ctx,
		`INSERT INTO users (id, email, auth0_id)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (auth0_id)
		 DO UPDATE SET email = EXCLUDED.email
		 RETURNING id`,
		userUUID, email, claims.Sub,
	).Scan(&id)
	if err != nil {
		return "", fmt.Errorf("failed to ensure user exists: %w", err)
	}

	return id, nil
}

func (s *Store) CreateAPIKey(ctx context.Context, userID, key string) (string, string, error) {
	var id string
	var createdAt string
	err := s.db.QueryRowContext(
		ctx,
		"INSERT INTO api_keys (user_id, key) VALUES ($1, $2) RETURNING id, created_at",
		userID, key,
	).Scan(&id, &createdAt)
	if err != nil {
		return "", "", fmt.Errorf("failed to create key: %w", err)
	}
	return id, createdAt, nil
}

func (s *Store) ListAPIKeys(ctx context.Context, userID string) ([]api.APIKey, error) {
	rows, err := s.db.QueryContext(
		ctx,
		"SELECT id, user_id, key, created_at, revoked_at FROM api_keys WHERE user_id = $1 ORDER BY created_at DESC",
		userID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch keys: %w", err)
	}
	defer rows.Close()

	keys := make([]api.APIKey, 0)
	for rows.Next() {
		var key api.APIKey
		var revokedAt sql.NullString
		if err := rows.Scan(&key.ID, &key.UserID, &key.Key, &key.CreatedAt, &revokedAt); err != nil {
			return nil, fmt.Errorf("failed to parse keys: %w", err)
		}
		if revokedAt.Valid {
			key.RevokedAt = &revokedAt.String
		}
		keys = append(keys, key)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate keys: %w", err)
	}

	return keys, nil
}

func (s *Store) RevokeAPIKey(ctx context.Context, keyID, userID string) (bool, error) {
	result, err := s.db.ExecContext(
		ctx,
		"UPDATE api_keys SET revoked_at = NOW() WHERE id = $1 AND user_id = $2",
		keyID, userID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to revoke key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to read rows affected: %w", err)
	}

	return rowsAffected > 0, nil
}
