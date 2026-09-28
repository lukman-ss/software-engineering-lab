package cache

import (
	"context"
	"math"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

// --- NAIVE STAMPEDE SERVICE ---

type NaiveStampedeService struct {
	cache *MemoryCache
	db    *MockDB
	ttl   time.Duration
}

func NewNaiveStampedeService(cache *MemoryCache, db *MockDB, ttl time.Duration) *NaiveStampedeService {
	return &NaiveStampedeService{cache: cache, db: db, ttl: ttl}
}

func (s *NaiveStampedeService) Get(ctx context.Context, key string) (string, error) {
	item, err := s.cache.Get(key)
	if err == nil {
		return item.Value, nil
	}

	start := time.Now()
	val, err := s.db.Query(ctx, key)
	if err != nil {
		return "", err
	}
	delta := time.Since(start)

	s.cache.Set(key, val, s.ttl, delta)
	return val, nil
}

// --- SINGLEFLIGHT STAMPEDE SERVICE ---

type SingleFlightService struct {
	cache  *MemoryCache
	db     *MockDB
	ttl    time.Duration
	flight singleflight.Group
}

func NewSingleFlightService(cache *MemoryCache, db *MockDB, ttl time.Duration) *SingleFlightService {
	return &SingleFlightService{cache: cache, db: db, ttl: ttl}
}

func (s *SingleFlightService) Get(ctx context.Context, key string) (string, error) {
	item, err := s.cache.Get(key)
	if err == nil {
		return item.Value, nil
	}

	res, err, _ := s.flight.Do(key, func() (interface{}, error) {
		// Double check inside flight execution
		item, err := s.cache.Get(key)
		if err == nil {
			return item.Value, nil
		}

		start := time.Now()
		val, err := s.db.Query(ctx, key)
		if err != nil {
			return "", err
		}
		delta := time.Since(start)

		s.cache.Set(key, val, s.ttl, delta)
		return val, nil
	})

	if err != nil {
		return "", err
	}
	return res.(string), nil
}

// --- XFETCH (PROBABILISTIC EARLY EXPIRATION) SERVICE ---

type XFetchService struct {
	cache *MemoryCache
	db    *MockDB
	ttl   time.Duration
	beta  float64
	rand  *rand.Rand
	mu    sync.Mutex
	// Custom random supplier for deterministic tests
	randFunc func() float64
}

func NewXFetchService(cache *MemoryCache, db *MockDB, ttl time.Duration, beta float64) *XFetchService {
	return &XFetchService{
		cache: cache,
		db:    db,
		ttl:   ttl,
		beta:  beta,
		rand:  rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (s *XFetchService) SetRandFunc(f func() float64) {
	s.randFunc = f
}

func (s *XFetchService) getRand() float64 {
	if s.randFunc != nil {
		return s.randFunc()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.rand.Float64()
}

// ShouldRecompute calculates XFetch condition:
// -Δ * β * ln(U) > TTL_remaining
// where Δ is read duration, β > 0, U ~ Uniform(0,1).
func ShouldRecompute(delta time.Duration, beta float64, ttlRemaining time.Duration, u float64) bool {
	if u <= 0 || u >= 1 {
		return false
	}
	deltaSec := delta.Seconds()
	ttlRemainingSec := ttlRemaining.Seconds()

	// Correct mathematical formula with negative log draw:
	// -delta * beta * log(u)
	expiryCompute := -deltaSec * beta * math.Log(u)
	return expiryCompute > ttlRemainingSec
}

func (s *XFetchService) Get(ctx context.Context, key string) (string, error) {
	item, ok := s.cache.GetRaw(key)
	now := time.Now()

	u := s.getRand()

	var needRecompute bool
	if !ok || (!item.ExpiresAt.IsZero() && now.After(item.ExpiresAt)) {
		needRecompute = true
	} else {
		remaining := item.ExpiresAt.Sub(now)
		needRecompute = ShouldRecompute(item.ReadDelta, s.beta, remaining, u)
	}

	if !needRecompute {
		return item.Value, nil
	}

	start := time.Now()
	val, err := s.db.Query(ctx, key)
	if err != nil {
		// If recompute fails but raw item exists, return stale value as fallback
		if ok {
			return item.Value, nil
		}
		return "", err
	}
	delta := time.Since(start)

	s.cache.Set(key, val, s.ttl, delta)
	return val, nil
}

// --- STALE-WHILE-REVALIDATE (SWR) SERVICE ---

type SWRService struct {
	cache      *MemoryCache
	db         *MockDB
	ttl        time.Duration
	staleDelta time.Duration
	revalidating map[string]bool
	mu         sync.Mutex
	revalCount atomic.Int64
}

func NewSWRService(cache *MemoryCache, db *MockDB, ttl time.Duration, staleDelta time.Duration) *SWRService {
	return &SWRService{
		cache:        cache,
		db:           db,
		ttl:          ttl,
		staleDelta:   staleDelta,
		revalidating: make(map[string]bool),
	}
}

func (s *SWRService) RevalidateCount() int64 {
	return s.revalCount.Load()
}

func (s *SWRService) Get(ctx context.Context, key string) (string, error) {
	item, ok := s.cache.GetRaw(key)
	now := time.Now()

	if ok {
		// 1. Fully fresh
		if item.ExpiresAt.IsZero() || now.Before(item.ExpiresAt) {
			return item.Value, nil
		}

		// 2. Stale but within stale window
		staleUntil := item.ExpiresAt.Add(s.staleDelta)
		if now.Before(staleUntil) {
			s.triggerRevalidate(key)
			return item.Value, nil
		}
	}

	// 3. Completely expired or miss -> synchronous fetch
	start := time.Now()
	val, err := s.db.Query(ctx, key)
	if err != nil {
		return "", err
	}
	delta := time.Since(start)
	s.cache.Set(key, val, s.ttl, delta)
	return val, nil
}

func (s *SWRService) triggerRevalidate(key string) {
	s.mu.Lock()
	if s.revalidating[key] {
		s.mu.Unlock()
		return
	}
	s.revalidating[key] = true
	s.mu.Unlock()

	s.revalCount.Add(1)

	go func() {
		defer func() {
			s.mu.Lock()
			delete(s.revalidating, key)
			s.mu.Unlock()
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		start := time.Now()
		val, err := s.db.Query(ctx, key)
		if err == nil {
			delta := time.Since(start)
			s.cache.Set(key, val, s.ttl, delta)
		}
	}()
}
