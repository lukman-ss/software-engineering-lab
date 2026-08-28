// Package database provides PostgreSQL connection management with
// a connection pool, health checks, and graceful shutdown.
package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/lib/pq"
)

// Config holds PostgreSQL connection settings.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

func DefaultConfig() Config {
	return Config{
		Host:            "localhost",
		Port:            "5432",
		User:            "postgres",
		Password:        "postgres",
		Database:        "se_lab",
		SSLMode:         "disable",
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
		ConnMaxIdleTime: 1 * time.Minute,
	}
}

func FromEnv() Config {
	c := DefaultConfig()
	if v := os.Getenv("POSTGRES_HOST"); v != "" {
		c.Host = v
	}
	if v := os.Getenv("POSTGRES_PORT"); v != "" {
		c.Port = v
	}
	if v := os.Getenv("POSTGRES_USER"); v != "" {
		c.User = v
	}
	if v := os.Getenv("POSTGRES_PASSWORD"); v != "" {
		c.Password = v
	}
	if v := os.Getenv("POSTGRES_DB"); v != "" {
		c.Database = v
	}
	if v := os.Getenv("POSTGRES_SSLMODE"); v != "" {
		c.SSLMode = v
	}
	if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.MaxOpenConns = n
		}
	}
	if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			c.MaxIdleConns = n
		}
	}
	if v := os.Getenv("DB_CONN_MAX_LIFETIME"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			c.ConnMaxLifetime = d
		}
	}
	if v := os.Getenv("DB_CONN_MAX_IDLE_TIME"); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			c.ConnMaxIdleTime = d
		}
	}
	return c
}

// Connect establishes a connection to PostgreSQL with the configured pool settings
// and verifies connectivity with a ping within a timeout.
func Connect(ctx context.Context, cfg Config) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)
	db.SetConnMaxIdleTime(cfg.ConnMaxIdleTime)

	// Startup ping with context timeout
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}

	return db, nil
}
