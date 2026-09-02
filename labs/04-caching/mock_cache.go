package caching

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// MockCache adalah implementasi in-memory cache untuk testing.
type MockCache struct {
	mu       sync.RWMutex
	data     map[string]string
	expiries map[string]time.Time
	fail     bool // simulate failure
}

func NewMockCache() *MockCache {
	return &MockCache{
		data:     make(map[string]string),
		expiries: make(map[string]time.Time),
	}
}

func NewFailingMockCache() *MockCache {
	return &MockCache{
		data:     make(map[string]string),
		expiries: make(map[string]time.Time),
		fail:     true,
	}
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	if m.fail {
		return "", ErrCacheDown
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.data[key]
	if !ok {
		return "", ErrCacheMiss
	}

	// Check expiry
	if expiry, ok := m.expiries[key]; ok && time.Now().After(expiry) {
		delete(m.data, key)
		delete(m.expiries, key)
		return "", ErrCacheMiss
	}

	return data, nil
}

func (m *MockCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if m.fail {
		return ErrCacheDown
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
	if ttl > 0 {
		m.expiries[key] = time.Now().Add(ttl)
	} else {
		delete(m.expiries, key)
	}

	return nil
}

func (m *MockCache) Delete(ctx context.Context, key string) error {
	if m.fail {
		return ErrCacheDown
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	delete(m.expiries, key)
	return nil
}

// GetWithExpiry is part of CacheInterface
func (m *MockCache) GetWithExpiry(ctx context.Context, key string) (string, time.Time, error) {
	if m.fail {
		return "", time.Time{}, ErrCacheDown
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.data[key]
	if !ok {
		return "", time.Time{}, ErrCacheMiss
	}

	expiry, ok := m.expiries[key]
	if !ok {
		// Key has no expiry
		return data, time.Time{}, nil
	}

	if time.Now().After(expiry) {
		delete(m.data, key)
		delete(m.expiries, key)
		return "", time.Time{}, ErrCacheMiss
	}

	return data, expiry, nil
}

// SetWithExpiry sets data with explicit expiry time (for testing)
func (m *MockCache) SetWithExpiry(key string, data string, actualExpiry time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = data
	m.expiries[key] = actualExpiry
}

// MockCacheWithStats tracks hit/miss for testing
type MockCacheWithStats struct {
	base   *MockCache
	hits   atomic.Int64
	misses atomic.Int64
}

func NewMockCacheWithStats() *MockCacheWithStats {
	return &MockCacheWithStats{
		base: NewMockCache(),
	}
}

func (m *MockCacheWithStats) Get(ctx context.Context, key string) (string, error) {
	data, err := m.base.Get(ctx, key)
	if errors.Is(err, ErrCacheMiss) {
		m.misses.Add(1)
	} else if err == nil {
		m.hits.Add(1)
	}
	return data, err
}

func (m *MockCacheWithStats) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	return m.base.Set(ctx, key, value, ttl)
}

func (m *MockCacheWithStats) Delete(ctx context.Context, key string) error {
	return m.base.Delete(ctx, key)
}

func (m *MockCacheWithStats) GetWithExpiry(ctx context.Context, key string) (string, time.Time, error) {
	data, expiry, err := m.base.GetWithExpiry(ctx, key)
	if errors.Is(err, ErrCacheMiss) {
		m.misses.Add(1)
	} else if err == nil {
		m.hits.Add(1)
	}
	return data, expiry, err
}

func (m *MockCacheWithStats) Hits() int {
	return int(m.hits.Load())
}

func (m *MockCacheWithStats) Misses() int {
	return int(m.misses.Load())
}

// MockRedisClient - implements LockInterface for distributed locking tests.
type MockRedisClient struct {
	data     map[string]string
	expiries map[string]time.Time
	mu       sync.RWMutex
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data:     make(map[string]string),
		expiries: make(map[string]time.Time),
	}
}

// removeExpired cleans up expired keys. Not thread-safe, caller must hold lock.
func (r *MockRedisClient) removeExpired(key string) {
	if expiry, ok := r.expiries[key]; ok && time.Now().After(expiry) {
		delete(r.data, key)
		delete(r.expiries, key)
	}
}

// SetNX atomically sets key only if it doesn't exist.
func (r *MockRedisClient) SetNX(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.removeExpired(key)

	if _, exists := r.data[key]; exists {
		return false, nil // Key already exists
	}

	r.data[key] = value
	if ttl > 0 {
		r.expiries[key] = time.Now().Add(ttl)
	} else {
		delete(r.expiries, key)
	}
	return true, nil
}

// Get retrieves value (for lock token verification).
func (r *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	r.mu.Lock() // Write lock needed since removeExpired can mutate map
	defer r.mu.Unlock()

	r.removeExpired(key)

	val, exists := r.data[key]
	if !exists {
		return "", nil
	}
	return val, nil
}

// CompareAndDel atomically deletes only if value matches.
func (r *MockRedisClient) CompareAndDel(ctx context.Context, key, value string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.removeExpired(key)

	val, exists := r.data[key]
	if !exists || val != value {
		return false, nil // Key missing or value mismatch
	}

	delete(r.data, key)
	delete(r.expiries, key)
	return true, nil
}

// ForceExpire forces a key to expire immediately (for testing).
func (r *MockRedisClient) ForceExpire(key string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[key]; ok {
		r.expiries[key] = time.Now().Add(-1 * time.Minute)
	}
}

// Ensure interfaces are satisfied
var _ CacheInterface = (*MockCache)(nil)
var _ CacheInterface = (*MockCacheWithStats)(nil)
var _ LockInterface = (*MockRedisClient)(nil)
