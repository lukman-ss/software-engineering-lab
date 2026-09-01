package caching

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// WriteThroughService menggunakan pola Write Through:
// Data ditulis secara sinkron ke Database dan Cache pada transaksi yang sama.
// Request -> Update DB -> Update Cache -> Return
type WriteThroughService struct {
	db        *sql.DB
	cache     CacheInterface
	ttl       time.Duration
	jitterTTL time.Duration
}

func NewWriteThroughService(db *sql.DB, cache CacheInterface) *WriteThroughService {
	return &WriteThroughService{
		db:        db,
		cache:     cache,
		ttl:       5 * time.Minute,
		jitterTTL: 15 * time.Second,
	}
}

// UpdateProduct melakukan Write Through.
// 1. Update Database
// 2. Update Cache (hanya jika DB sukses)
func (s *WriteThroughService) UpdateProduct(ctx context.Context, p Product) error {
	// 1. Update ke Source of Truth (Database)
	startDB := time.Now()
	res, err := s.db.ExecContext(ctx, "UPDATE products SET name = $1, price = $2 WHERE id = $3", p.Name, p.Price, p.ID)
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
	_ = startDB

	// 2. Write-through ke Cache
	key := CacheKey("product", p.ID, 1)
	data, err := json.Marshal(p)
	if err != nil {
		// PENTING: DB sukses, tapi persiapan cache gagal.
		// Lebih aman untuk INVALIDATE (delete) dibanding membiarkan stale data.
		_ = s.cache.Delete(ctx, key)
		return fmt.Errorf("marshal product (DB updated successfully): %w", err)
	}

	jitteredTTL := TTLWithJitter(s.ttl, s.jitterTTL)
	if err := s.cache.Set(ctx, key, string(data), jitteredTTL); err != nil {
		// PENTING: DB sukses, tapi update cache gagal (Redis down).
		// Kita harus log error ini. Request tetap dianggap sukses karena DB (source of truth) sudah berubah.
		// Pada arsitektur yang sangat ketat, kita mungkin mengharuskan retry/outbox.
		fmt.Printf("warn: write-through cache set failed for %s: %v\n", key, err)
		// Kita menghapus key (invalidate) sebagai safety fallback, karena set() gagal
		// dan kita tidak ingin meninggalkan stale data.
		_ = s.cache.Delete(ctx, key)
	}

	return nil
}
