package caching

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// NaiveService tidak menggunakan cache apapun.
// Setiap request selalu query DB - latensi tinggi.
type NaiveService struct {
	db *sql.DB
}

func NewNaiveService(db *sql.DB) *NaiveService {
	return &NaiveService{db: db}
}

func (s *NaiveService) GetProduct(ctx context.Context, key string) (Product, error) {
	start := time.Now()

	// Extract ID dari cache key
	id := extractID(key)

	// Selalu query DB
	var p Product
	row := s.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = $1", id)
	err := row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return Product{}, fmt.Errorf("query DB: %w", err)
	}

	elapsed := time.Since(start)
	_ = elapsed

	return p, nil
}

func extractID(key string) string {
	// Key format: entity:id
	parts := splitKey(key)
	if len(parts) >= 2 {
		return parts[1]
	}
	return key
}

func splitKey(key string) []string {
	var parts []string
	current := ""
	for _, c := range key {
		if c == ':' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	parts = append(parts, current)
	return parts
}