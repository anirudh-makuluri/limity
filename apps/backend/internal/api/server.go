package api

import "context"

type PostgresStore interface {
	Ping(ctx context.Context) error
	EnsureUserWithAPIKey(ctx context.Context, claims *TokenClaims) (*UserProfile, error)
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
