package caching

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
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
//
// Metric classification:
//   - hit: cache returned valid value
//   - miss: cache returned ErrCacheMiss (key absent / expired)
//   - error: cache returned backend error (Redis down, network) OR corrupt JSON
//   - db_fallback: only recorded when cache returned a backend error
func (s *CacheAsideService) GetProduct(ctx context.Context, id string) (Product, error) {
	key := CacheKey("product", id, 1)

	// 1. Check cache with latency measurement
	startGet := time.Now()
	cached, err := s.cache.Get(ctx, key)
	s.metrics.IncCacheGetOp()
	s.metrics.RecordCacheGetLatency(time.Since(startGet))

	switch {
	case err == nil && cached != "":
		// Cache HIT - return data
		var p Product
		if err := json.Unmarshal([]byte(cached), &p); err == nil {
			s.metrics.IncHit()
			return p, nil
		}
		// Cache CORRUPT - unmarshal gagal
		s.metrics.IncError()
		s.metrics.IncCacheInvalidateOp()
		_ = s.cache.Delete(ctx, key) // attempt cleanup, but don't fail read
	case errors.Is(err, ErrCacheMiss):
		// Normal cache miss - bukan error
		s.metrics.IncMiss()
	case err != nil:
		// Cache backend error (Redis down, network, timeout)
		s.metrics.IncError()
		s.metrics.IncCacheOperationError()
		s.metrics.IncDBFallback()
	default:
		// Empty value without error - treat as miss per empty-value policy
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
		// Marshal error = cache side error, bukan business error.
		// DB read sukses, business operation tetap success.
		s.metrics.IncError()
		return p, nil
	}

	jitteredTTL := TTLWithJitter(s.ttl, s.jitterTTL)
	s.metrics.IncCacheSetOp()
	startSet := time.Now()
	err = s.cache.Set(ctx, key, string(data), jitteredTTL)
	s.metrics.RecordCacheSetLatency(time.Since(startSet))
	if err != nil {
		// Cache SET failed - DB success, data tersedia
		// Record cache-side error, bukan fail business read
		s.metrics.IncError()
		s.metrics.IncCacheSetError()
	}

	return p, nil
}

// GetMetrics mengembalikan metrics untuk observability.
func (s *CacheAsideService) GetMetrics() *CacheMetrics {
	return s.metrics
}
