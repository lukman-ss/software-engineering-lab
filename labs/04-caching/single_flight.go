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

func (s *SingleFlightService) GetProduct(ctx context.Context, key string) (Product, error) {
	// 1. Check cache dulu
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var p Product
		if err := json.Unmarshal([]byte(cached), &p); err == nil {
			return p, nil
		}
	}

	// 2. Single flight - satu DB query untuk semua request yang sama
	result, err, _ := s.flight.Do(key, func() (interface{}, error) {
		start := time.Now()
		var p Product
		row := s.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = $1", extractID(key))
		err := row.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return Product{}, fmt.Errorf("query DB: %w", err)
		}
		_ = start

		// 3. Populate cache
		data, _ := json.Marshal(p)
		_ = s.cache.Set(ctx, key, string(data), s.ttl)

		return p, nil
	})

	if err != nil {
		return Product{}, err
	}

	return result.(Product), nil
}