package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/sync/singleflight"
)

// TTLWithJitter adds random jitter to prevent synchronized expiration (stampede).
func TTLWithJitter(base time.Duration, maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return base
	}
	jitter := time.Duration(rand.Int63n(int64(maxJitter)))
	return base + jitter
}

// CounterRepository counts calls for stampede demonstration.
type CounterRepository struct {
	callCount int64
}

func NewCounterRepository() *CounterRepository {
	return &CounterRepository{}
}

func (r *CounterRepository) GetDashboard(ctx context.Context, branchID int64, businessDate time.Time) (Dashboard, error) {
	r.callCount++
	return Dashboard{InvoiceCountToday: int(r.callCount)}, nil
}

func (r *CounterRepository) CallCount() int64 {
	return r.callCount
}

// --- VERSI BROKEN (STAMPEDE) ---

type BrokenStampedeService struct {
	cache   CacheInterface
	repo    *CounterRepository
	keyFunc func(branchID int64) string
}

func NewBrokenStampedeService(cache CacheInterface, repo *CounterRepository) *BrokenStampedeService {
	return &BrokenStampedeService{cache: cache, repo: repo}
}

func (s *BrokenStampedeService) GetData(ctx context.Context, branchID int64) (Dashboard, error) {
	key := fmt.Sprintf("dash:%d", branchID)

	// 1. Cek cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var d Dashboard
		_ = json.Unmarshal([]byte(cached), &d)
		return d, nil
	}

	// 2. Cache miss -> Langsung hajar DB (Tidak ada perlindungan)
	// Jika ada 100 concurrent request, ke-100-nya akan query DB
	d, _ := s.repo.GetDashboard(ctx, branchID, time.Now())

	// 3. Set cache
	data, _ := json.Marshal(d)
	_ = s.cache.Set(ctx, key, string(data), 1*time.Minute)

	return d, nil
}

// --- VERSI PROTECTED (SINGLEFLIGHT) ---

type ProtectedStampedeService struct {
	cache   CacheInterface
	repo    *CounterRepository
	flight  singleflight.Group
	keyFunc func(branchID int64) string
}

func NewProtectedStampedeService(cache CacheInterface, repo *CounterRepository) *ProtectedStampedeService {
	return &ProtectedStampedeService{cache: cache, repo: repo}
}

func (s *ProtectedStampedeService) GetData(ctx context.Context, branchID int64) (Dashboard, error) {
	key := fmt.Sprintf("dash:%d", branchID)

	// 1. Cek cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var d Dashboard
		_ = json.Unmarshal([]byte(cached), &d)
		return d, nil
	}

	// 2. Cache miss -> Gunakan singleflight
	// Jika ada 100 concurrent request, HANYA SATU yang eksekusi DB.
	result, err, _ := s.flight.Do(key, func() (interface{}, error) {
		d, err := s.repo.GetDashboard(ctx, branchID, time.Now())
		if err != nil {
			return Dashboard{}, err
		}
		data, _ := json.Marshal(d)
		_ = s.cache.Set(ctx, key, string(data), 1*time.Minute)
		return d, nil
	})

	if err != nil {
		return Dashboard{}, err
	}
	return result.(Dashboard), nil
}