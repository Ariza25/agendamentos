package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open establishes a connection pool to Postgres with production-tuned pool defaults
// and configurable limits via environment variables.
func Open(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}

	maxOpen := resolveIntEnv("DB_MAX_OPEN_CONNS", 25)
	maxIdle := resolveIntEnv("DB_MAX_IDLE_CONNS", 10)
	maxLifetime := resolveDurationEnv("DB_CONN_MAX_LIFETIME", 30*time.Minute)
	maxIdleTime := resolveDurationEnv("DB_CONN_MAX_IDLE_TIME", 5*time.Minute)

	db.SetMaxOpenConns(maxOpen)
	db.SetMaxIdleConns(maxIdle)
	db.SetConnMaxLifetime(maxLifetime)
	db.SetConnMaxIdleTime(maxIdleTime)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return db, nil
}

func resolveIntEnv(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	val, err := strconv.Atoi(raw)
	if err != nil || val <= 0 {
		return fallback
	}
	return val
}

func resolveDurationEnv(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	dur, err := time.ParseDuration(raw)
	if err != nil || dur <= 0 {
		return fallback
	}
	return dur
}
