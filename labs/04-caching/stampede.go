package caching

import (
	"context"
	"math/rand"
	"time"

	"golang.org/x/sync/singleflight"
)

// --- VERSI BROKEN (STAMPEDE) ---

type BrokenStampedeService struct {
	cache CacheInterface
	db    *HeavyDB
}

func NewBrokenStampedeService(cache CacheInterface, db *HeavyDB) *BrokenStampedeService {
	return &BrokenStampedeService{cache: cache, db: db}
}

func (s *BrokenStampedeService) GetData(ctx context.Context, key string) (string, error) {
	// 1. Cek cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		return cached, nil
	}

	// 2. Cache miss -> Langsung hajar DB (Tidak ada perlindungan)
	// Jika ada 100 concurrent request, ke-100-nya akan query DB
	data := s.db.FetchHeavyData()

	// 3. Set cache
	_ = s.cache.Set(ctx, key, data, 1*time.Minute)

	return data, nil
}

// --- VERSI PROTECTED (SINGLEFLIGHT) ---

type ProtectedStampedeService struct {
	cache  CacheInterface
	db     *HeavyDB
	flight singleflight.Group
}

func NewProtectedStampedeService(cache CacheInterface, db *HeavyDB) *ProtectedStampedeService {
	return &ProtectedStampedeService{cache: cache, db: db}
}

func (s *ProtectedStampedeService) GetData(ctx context.Context, key string) (string, error) {
	// 1. Cek cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		return cached, nil
	}

	// 2. Cache miss -> Gunakan singleflight
	// Jika ada 100 concurrent request, HANYA SATU yang eksekusi DB.
	// Yang lain menunggu (block) sampai yang pertama selesai, lalu menggunakan hasilnya.
	result, err, _ := s.flight.Do(key, func() (interface{}, error) {
		data := s.db.FetchHeavyData()
		_ = s.cache.Set(ctx, key, data, 1*time.Minute)
		return data, nil
	})

	if err != nil {
		return "", err
	}
	return result.(string), nil
}

// --- TTL JITTER HELPER ---

// TTLWithJitter menambahkan random jitter ke base TTL.
// Tujuannya agar cache keys tidak expire bersamaan (menghindari stampede massal).
func TTLWithJitter(base time.Duration, maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return base
	}
	// Random duration between 0 and maxJitter
	jitter := time.Duration(rand.Int63n(int64(maxJitter)))
	return base + jitter
}