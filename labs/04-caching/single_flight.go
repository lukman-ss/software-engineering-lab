package caching

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"golang.org/x/sync/singleflight"
)

// SingleFlightService menggunakan golang.org/x/sync/singleflight
// untuk mencegah duplicate DB queries untuk key yang sama.
type SingleFlightService struct {
	db     *sql.DB
	cache  CacheInterface
	ttl    time.Duration
	flight singleflight.Group
}

func NewSingleFlightService(db *sql.DB, cache CacheInterface) *SingleFlightService {
	return &SingleFlightService{
		db:    db,
		cache: cache,
		ttl:   5 * time.Minute,
	}
}

// GetProduct menggunakan singleflight dengan double-check pattern.
// Menerima productID secara eksplisit untuk domain-driven design.
func (s *SingleFlightService) GetProduct(ctx context.Context, productID string) (Product, error) {
	key := CacheKey("product", productID, 1)

	// 1. Check cache dulu
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var p Product
		if err := json.Unmarshal([]byte(cached), &p); err == nil {
			return p, nil
		}
		// Unmarshal gagal - corrupt cache, delete and proceed to DB
		_ = s.cache.Delete(ctx, key)
	}

	// 2. Single flight - satu DB query untuk semua request yang sama
	result, err, _ := s.flight.Do(key, func() (interface{}, error) {
		start := time.Now()
		var p Product

		// Double-check cache setelah mendapat singleflight slot
		cached2, err2 := s.cache.Get(ctx, key)
		if err2 == nil && cached2 != "" {
			var p2 Product
			if err3 := json.Unmarshal([]byte(cached2), &p2); err3 == nil {
				_ = start // silence unused variable warning before return
				return p2, nil
			}
		}

		row := s.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = $1", productID)
		err := row.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			if err == sql.ErrNoRows {
				return Product{}, fmt.Errorf("product not found: %w", err)
			}
			return Product{}, fmt.Errorf("query DB: %w", err)
		}
		_ = start // silence unused variable warning

		// 3. Populate cache
		data, marshalErr := json.Marshal(p)
		if marshalErr != nil {
			return p, fmt.Errorf("marshal product: %w", marshalErr)
		}
		if setErr := s.cache.Set(ctx, key, string(data), s.ttl); setErr != nil {
			// Log but don't fail - DB read succeeded
			fmt.Printf("warn: cache set failed: %v\n", setErr)
		}

		return p, nil
	})

	if err != nil {
		return Product{}, err
	}

	return result.(Product), nil
}