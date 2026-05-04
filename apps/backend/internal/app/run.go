package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/limity/backend/internal/api"
	postgresstore "github.com/limity/backend/internal/store/postgres"
	redisstore "github.com/limity/backend/internal/store/redis"
)

func Run() error {
	_ = godotenv.Load()

	var db *sql.DB
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		fmt.Println("Warning: DATABASE_URL not set. API key endpoints will not work.")
	} else {
		conn, err := sql.Open("postgres", dbURL)
		if err != nil {
			fmt.Printf("Warning: Failed to connect to database: %v\n", err)
		} else {
			if err := conn.Ping(); err != nil {
				fmt.Printf("Warning: Database ping failed: %v\n", err)
			}
			db = conn
			defer db.Close()
		}
	}

	allowedOrigins := os.Getenv("ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000,http://localhost:5173"
		fmt.Printf("ALLOWED_ORIGINS not set, using defaults: %s\n", allowedOrigins)
	}

	pgStore := postgresstore.NewStore(db)
	redisStore := redisstore.NewStore(
		os.Getenv("UPSTASH_REDIS_REST_URL"),
		os.Getenv("UPSTASH_REDIS_REST_TOKEN"),
	)
	metrics := api.NewMetrics()

	var analytics *api.AsyncAnalytics
	analyticsEnabled := envBool("ANALYTICS_ENABLED", true)
	if analyticsEnabled && db != nil {
		analytics = api.NewAsyncAnalytics(
			pgStore,
			metrics,
			envInt("ANALYTICS_QUEUE_SIZE", 10000),
			envInt("ANALYTICS_BATCH_SIZE", 200),
			time.Duration(envInt("ANALYTICS_FLUSH_INTERVAL_MS", 1000))*time.Millisecond,
			time.Duration(envInt("ANALYTICS_FLUSH_TIMEOUT_MS", 15000))*time.Millisecond,
		)
		defer analytics.Close()
	} else if analyticsEnabled && db == nil {
		fmt.Println("Analytics disabled: DATABASE_URL not configured")
	}

	r := chi.NewRouter()
	srv := api.NewServer(pgStore, redisStore, allowedOrigins, metrics, analytics)
	r.Use(srv.ObservabilityMiddleware)
	srv.RegisterObservabilityRoutes(r)
	srv.RegisterRoutes(r)

	go srv.StartAPIKeyGaugeUpdater(time.Duration(envInt("API_KEY_GAUGE_REFRESH_SEC", 60)) * time.Second)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on :%s\n", port)
	return http.ListenAndServe(":"+port, r)
}

func envInt(name string, fallback int) int {
	v := os.Getenv(name)
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}

func envBool(name string, fallback bool) bool {
	v := os.Getenv(name)
	if v == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return parsed
}
