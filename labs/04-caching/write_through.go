package caching

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// WriteThroughService menggunakan pola Write Through.
//
// PENTING: Database dan Redis adalah sistem TERPISAH. Tidak ada atomic commit lintar-keduanya.
// Redis tidak berada dalam SQL transaction sehingga tidak ada "transaction commit" berarti
// otomatis sinkron dengan Redis.
//
// TTL tetap berfungsi sebagai safety net saat proses crash antara DB commit dan cache update.
type WriteThroughService struct {
	db        *sql.DB
	cache     CacheInterface
	ttl       time.Duration
	jitterTTL time.Duration
}

// NewWriteThroughService membuat service dengan TTL default 5 menit.
func NewWriteThroughService(db *sql.DB, cache CacheInterface) *WriteThroughService {
	return &WriteThroughService{
		db:        db,
		cache:     cache,
		ttl:       5 * time.Minute,
		jitterTTL: 15 * time.Second,
	}
}

// UpdateProduct melakukan Write Through:
// 1. Update Database (Source of Truth)
// 2. Best-effort update cache
//
// CATATAN: Jika cache update gagal, business operation tetap SUCCESS.
// - Log metric error
// - Best-effort invalidate stale key
// - Stale cache masih mungkin bertahan sampai TTL (safety net)
//
// Flow yang benar:
// validate input -> DB write -> DB success -> best-effort cache update -> return success
func (s *WriteThroughService) UpdateProduct(ctx context.Context, p Product) error {
	// 0. Validasi input
	if p.ID == "" {
		return fmt.Errorf("product ID tidak boleh kosong")
	}

	// 1. Serialize dulu sebelum DB mutation (detect error lebih awal)
	data, err := json.Marshal(p)
	if err != nil {
		// Serialization error = data tidak valid untuk cache
		// Return error - ini bukan cache failure, data tidak akan pernah sengaja disimpan
		return fmt.Errorf("data tidak valid untuk cache: %w", err)
	}

	// 2. Update ke Source of Truth (Database)
	res, err := s.db.ExecContext(ctx,
		"UPDATE products SET name = $1, price = $2 WHERE id = $3",
		p.Name, p.Price, p.ID)
	if err != nil {
		return fmt.Errorf("update DB: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("check rows affected: %w", err)
	}
	if rows == 0 {
		return fmt.Errorf("product not found")
	}

	// 3. Best-effort cache update (setelah DB commit!)
	key := CacheKey("product", p.ID, 1)
	jitteredTTL := TTLWithJitter(s.ttl, s.jitterTTL)
	if err := s.cache.Set(ctx, key, string(data), jitteredTTL); err != nil {
		// Cache set gagal - DB sudah sukses, business operation tetap success
		// Log error, bukan return error

		// Best-effort: delete stale key sebagai safety fallback
		// INI BUKAN guaranteed - stale cache masih mungkin bertahan sampai TTL
		if delErr := s.cache.Delete(ctx, key); delErr != nil {
			// stale cache persists until TTL
		}
	}

	return nil
}

// GetProduct menggunakan cache-aside pattern untuk read.
//
// Cache-Aside: Application mengontrol cache interaction, bukan cache layer.
// - GET cache
// - MISS -> GET database
// - SET cache
//
// Bukan Read-Through (di mana cache layer sendiri meload dari backing store).
func (s *WriteThroughService) GetProduct(ctx context.Context, key string) (Product, error) {
	// 1. Check cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var p Product
		if err := json.Unmarshal([]byte(cached), &p); err == nil {
			return p, nil // cache hit
		}
		// Unmarshal gagal → corrupt cache, delete
		_ = s.cache.Delete(ctx, key)
	}

	// 2. Cache miss → query DB
	var p Product
	row := s.db.QueryRowContext(ctx, "SELECT id, name, price FROM products WHERE id = $1", extractID(key))
	err = row.Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return Product{}, fmt.Errorf("query DB: %w", err)
	}

	// 3. Populate cache (best-effort)
	data, err := json.Marshal(p)
	if err != nil {
		return p, fmt.Errorf("marshal product: %w", err) // return data but log error
	}
	_ = s.cache.Set(ctx, key, string(data), s.ttl) // log but don't fail

	return p, nil
}