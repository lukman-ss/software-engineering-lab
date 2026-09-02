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
	db        *sql.DB
	cache     CacheInterface
	metrics   *CacheMetrics
	ttl       time.Duration
	jitterTTL time.Duration
}

// NewCacheAsideService membuat service dengan TTL 5 menit dan jitter 15 detik.
// metrics diwajibkan untuk observability yang benar.
func NewCacheAsideService(db *sql.DB, cache CacheInterface, metrics *CacheMetrics) *CacheAsideService {
	if metrics == nil {
		metrics = NewCacheMetrics()
	}
	return &CacheAsideService{
		db:        db,
		cache:     cache,
		metrics:   metrics,
		ttl:       5 * time.Minute,
		jitterTTL: 15 * time.Second,
	}
}

// GetProduct menggunakan cache-aside pattern untuk membaca produk.
// API menerima productID langsung, bukan cache key (menghindari coupling yang rapuh).
// Key dibuat menggunakan canonical key builder di dalam service.
func (s *CacheAsideService) GetProduct(ctx context.Context, id string) (Product, error) {
	key := CacheKey("product", id, 1)

	// 1. Check cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		// Cache HIT - return data
		var p Product
		if err := json.Unmarshal([]byte(cached), &p); err == nil {
			s.metrics.IncHit()
			return p, nil
		}
		// Cache CORRUPT - unmarshal gagal
		s.metrics.IncError()
		_ = s.cache.Delete(ctx, key) // attempt cleanup, but don't fail read
		// Note: In production, log structured error for observability
	} else if err != nil {
		// Cache ERROR (Redis down, network, timeout)
		// Bukan cache miss - error teknis
		s.metrics.IncError()
		s.metrics.IncDBFallback()
		// TODO: Use structured logging instead of fmt.Printf
		// fmt.Printf("cache error on GET %s: %v\n", key, err)
	} else {
		// Cache MISS - key tidak ada
		s.metrics.IncMiss()
	}

	// 2. Cache miss / error → query DB
	s.metrics.IncDBQuery()
	var p Product
	row := s.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = $1", id)
	err = row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return Product{}, fmt.Errorf("query DB: %w", err)
	}

	// 3. Populate cache (best-effort)
	data, err := json.Marshal(p)
	if err != nil {
		// Marshal error = data tidak valid untuk cache, bukan error bisnis
		// Return data saja, log error
		// fmt.Printf("warn: marshal product failed for cache: %v\n", err)
		return p, nil
	}

	jitteredTTL := TTLWithJitter(s.ttl, s.jitterTTL)
	if err := s.cache.Set(ctx, key, string(data), jitteredTTL); err != nil {
		// Cache SET failed - DB success, data tersedia
		// Log error, bukan fail business read
		s.metrics.IncError()
		// fmt.Printf("warn: cache set failed for %s: %v\n", key, err)
	}

	return p, nil
}

// GetMetrics mengembalikan metrics untuk observability.
func (s *CacheAsideService) GetMetrics() *CacheMetrics {
	return s.metrics
}