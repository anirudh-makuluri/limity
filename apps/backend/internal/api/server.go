package api

import "context"

type PostgresStore interface {
	Ping(ctx context.Context) error
	EnsureUserWithAPIKey(ctx context.Context, claims *TokenClaims) (*UserProfile, error)
	RotateAPIKey(ctx context.Context, userID string) (string, error)
	InsertRequestEvents(ctx context.Context, events []RequestEvent) error
	CountAPIKeys(ctx context.Context) (int, error)
	ResolveOwnerUserIDByAPIKey(ctx context.Context, apiKey string) (string, error)
}

type RedisStore interface {
	Incr(ctx context.Context, key string) (int, error)
	Expire(ctx context.Context, key string, seconds int) error
}

type Server struct {
	pg             PostgresStore
	redis          RedisStore
	allowedOrigins string
	metrics        *Metrics
	analytics      *AsyncAnalytics
}

func NewServer(pg PostgresStore, redis RedisStore, allowedOrigins string, metrics *Metrics, analytics *AsyncAnalytics) *Server {
	return &Server{pg: pg, redis: redis, allowedOrigins: allowedOrigins, metrics: metrics, analytics: analytics}
}
