package app

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"

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

	r := chi.NewRouter()
	srv := api.NewServer(pgStore, redisStore, allowedOrigins)
	srv.RegisterRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on :%s\n", port)
	return http.ListenAndServe(":"+port, r)
}
