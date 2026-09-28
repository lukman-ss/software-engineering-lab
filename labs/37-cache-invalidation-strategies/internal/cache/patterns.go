package cache

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// --- CACHE-ASIDE PATTERN ---

type CacheAsideService struct {
	cache *MemoryCache
	db    *MockDB
	ttl   time.Duration
}

func NewCacheAsideService(cache *MemoryCache, db *MockDB, ttl time.Duration) *CacheAsideService {
	return &CacheAsideService{cache: cache, db: db, ttl: ttl}
}

func (s *CacheAsideService) Get(ctx context.Context, key string) (string, error) {
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

func (s *CacheAsideService) Update(ctx context.Context, key, val string) error {
	// 1. Store first
	if err := s.db.Write(ctx, key, val); err != nil {
		return fmt.Errorf("db write failed: %w", err)
	}
	// 2. Invalidate cache
	s.cache.Delete(key)
	return nil
}

// --- WRITE-THROUGH PATTERN ---

type WriteThroughService struct {
	cache *MemoryCache
	db    *MockDB
	ttl   time.Duration
}

func NewWriteThroughService(cache *MemoryCache, db *MockDB, ttl time.Duration) *WriteThroughService {
	return &WriteThroughService{cache: cache, db: db, ttl: ttl}
}

func (s *WriteThroughService) Get(ctx context.Context, key string) (string, error) {
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

func (s *WriteThroughService) Update(ctx context.Context, key, val string) error {
	start := time.Now()
	// 1. Synchronous write to store
	if err := s.db.Write(ctx, key, val); err != nil {
		return fmt.Errorf("db write failed: %w", err)
	}
	delta := time.Since(start)

	// 2. Synchronous write to cache
	s.cache.Set(key, val, s.ttl, delta)
	return nil
}

// --- WRITE-BEHIND (WRITE-BACK) PATTERN ---

type WriteRequest struct {
	Key   string
	Value string
}

type WriteBehindService struct {
	cache      *MemoryCache
	db         *MockDB
	ttl        time.Duration
	writeQueue chan WriteRequest
	wg         sync.WaitGroup
	quit       chan struct{}
}

func NewWriteBehindService(cache *MemoryCache, db *MockDB, ttl time.Duration, bufferSize int) *WriteBehindService {
	s := &WriteBehindService{
		cache:      cache,
		db:         db,
		ttl:        ttl,
		writeQueue: make(chan WriteRequest, bufferSize),
		quit:       make(chan struct{}),
	}
	s.wg.Add(1)
	go s.flushWorker()
	return s
}

func (s *WriteBehindService) flushWorker() {
	defer s.wg.Done()
	for {
		select {
		case req := <-s.writeQueue:
			_ = s.db.Write(context.Background(), req.Key, req.Value)
		case <-s.quit:
			// ponytail: background flush drains remaining queue on graceful shutdown
			for len(s.writeQueue) > 0 {
				req := <-s.writeQueue
				_ = s.db.Write(context.Background(), req.Key, req.Value)
			}
			return
		}
	}
}

func (s *WriteBehindService) Get(ctx context.Context, key string) (string, error) {
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

func (s *WriteBehindService) Update(key, val string) {
	// Write immediately to cache
	s.cache.Set(key, val, s.ttl, 1*time.Millisecond)
	// Enqueue write to backing store asynchronously
	select {
	case s.writeQueue <- WriteRequest{Key: key, Value: val}:
	default:
		// Queue full (demonstration: drop or handle overflow)
	}
}

func (s *WriteBehindService) Close() {
	close(s.quit)
	s.wg.Wait()
}
