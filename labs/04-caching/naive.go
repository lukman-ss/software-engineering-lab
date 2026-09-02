package caching

import (
	"context"
	"database/sql"
	"fmt"
)

// NaiveService tidak menggunakan cache apapun.
// Setiap request selalu query DB langsung — baseline untuk perbandingan latency.
type NaiveService struct {
	db *sql.DB
}

func NewNaiveService(db *sql.DB) *NaiveService {
	return &NaiveService{db: db}
}

// GetProduct langsung query DB tanpa cache.
func (s *NaiveService) GetProduct(ctx context.Context, productID string) (Product, error) {
	var p Product
	row := s.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = $1", productID)
	if err := row.Scan(&p.ID, &p.Name, &p.Price); err != nil {
		return Product{}, fmt.Errorf("query DB: %w", err)
	}
	return p, nil
}
