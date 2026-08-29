package caching

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// CacheAsideService menggunakan pola Cache Aside:
// - Check cache dulu
// - Cache miss → query DB
// - Populate cache setelah DB hit
type CacheAsideService struct {
	db          *sql.DB
	cache       CacheInterface
	ttl         time.Duration
	jitterTTL   time.Duration // max jitter to add to TTL
}

func NewCacheAsideService(db *sql.DB, cache CacheInterface) *CacheAsideService {
	return &CacheAsideService{
		db:        db,
		cache:     cache,
		ttl:       5 * time.Minute,
		jitterTTL: 15 * time.Second, // Jitter to prevent stampede on expiration
	}
}

func (s *CacheAsideService) GetProduct(ctx context.Context, key string) (Product, error) {
	// 1. Check cache dulu
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var p Product
		if err := json.Unmarshal([]byte(cached), &p); err == nil {
			return p, nil // CACHE HIT
		}
		// Jika unmarshal gagal, hapus cache (stale/corrupt data)
		_ = s.cache.Delete(ctx, key)
	}

	// 2. Cache miss → query DB
	startDB := time.Now()
	var p Product
	row := s.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = $1", extractID(key))
	err = row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return Product{}, fmt.Errorf("query DB: %w", err)
	}
	_ = startDB

	// 3. Populate cache with jitter to prevent synchronized expiration
	data, err := json.Marshal(p)
	if err != nil {
		return Product{}, fmt.Errorf("marshal product: %w", err)
	}
	// Use jittered TTL to disperse expiration times across cache entries
	jitteredTTL := TTLWithJitter(s.ttl, s.jitterTTL)
	if err := s.cache.Set(ctx, key, string(data), jitteredTTL); err != nil {
		// Log but don't fail - we have the data from DB
		fmt.Printf("warn: cache set failed: %v\n", err)
	}

	return p, nil
}