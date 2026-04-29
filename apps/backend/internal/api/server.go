package api

import "context"

type PostgresStore interface {
	Ping(ctx context.Context) error
	EnsureUserExists(ctx context.Context, claims *TokenClaims) (string, error)
	CreateAPIKey(ctx context.Context, userID, key string) (string, string, error)
	ListAPIKeys(ctx context.Context, userID string) ([]APIKey, error)
	RevokeAPIKey(ctx context.Context, keyID, userID string) (bool, error)
}

type RedisStore interface {
	Incr(ctx context.Context, key string) (int, error)
	Expire(ctx context.Context, key string, seconds int) error
}

type Server struct {
	pg             PostgresStore
	redis          RedisStore
	allowedOrigins string
}

func NewServer(pg PostgresStore, redis RedisStore, allowedOrigins string) *Server {
	return &Server{pg: pg, redis: redis, allowedOrigins: allowedOrigins}
}
