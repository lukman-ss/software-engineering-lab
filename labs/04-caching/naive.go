package caching

import (
	"context"
	"database/sql"
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

