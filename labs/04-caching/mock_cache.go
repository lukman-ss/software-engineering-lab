package caching

import (
	"context"
	"errors"
	"fmt"
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

func NewMockCacheWithTTL(ttl time.Duration) *MockCache {
	mc := NewMockCache()
	// TTL is implicit in expiry times
	_ = ttl
	return mc
}

func NewFailingMockCache() *MockCache {
	return &MockCache{data: make(map[string]string), fail: true}
}

func (m *MockCache) Get(ctx context.Context, key string) (string, error) {
	if m.fail {
		return "", errors.New("cache connection failed")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.data[key]
	if !ok {
		return "", errors.New("cache miss")
	}

	// Check expiry
	if expiry, ok := m.expiries[key]; ok && time.Now().After(expiry) {
		delete(m.data, key)
		delete(m.expiries, key)
		return "", errors.New("cache expired")
	}

	return data, nil
}

func (m *MockCache) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	if m.fail {
		return errors.New("cache connection failed")
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

func (m *MockCache) GetWithExpiry(ctx context.Context, key string) (string, time.Time, error) {
	if m.fail {
		return "", time.Time{}, errors.New("cache connection failed")
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.data[key]
	if !ok {
		return "", time.Time{}, errors.New("cache miss")
	}

	expiry, ok := m.expiries[key]
	if !ok {
		return data, time.Time{}, nil
	}

	if time.Now().After(expiry) {
		delete(m.data, key)
		delete(m.expiries, key)
		return "", time.Time{}, errors.New("cache expired")
	}

	return data, expiry, nil
}

// SetWithExpiry sets data with explicit expiry time (for testing)
func (m *MockCache) SetWithExpiry(key, data string, ttl, actualExpiry time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = data
	m.expiries[key] = actualExpiry
}

// Statistik untuk testing/debugging
type CacheStats struct {
	Hits   int
	Misses int
}

func (m *MockCache) Stats() CacheStats {
	// Not implemented for simplicity
	return CacheStats{}
}

// MockCacheWithStats tracks hit/miss for testing
type MockCacheWithStats struct {
	*MockCache
	mu       sync.RWMutex
	hits     int
	misses   int
	expireAt time.Time
}

func NewMockCacheWithStats() *MockCacheWithStats {
	return &MockCacheWithStats{
		MockCache: NewMockCache(),
	}
}

func (m *MockCacheWithStats) Get(ctx context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	data, ok := m.data[key]
	if !ok {
		m.misses++
		return "", errors.New("cache miss")
	}

	// Check expiry
	if !m.expireAt.IsZero() {
		if time.Now().After(m.expireAt) {
			m.misses++
			return "", errors.New("cache expired")
		}
	}

	m.hits++
	return data, nil
}

func (m *MockCacheWithStats) Set(ctx context.Context, key, value string, ttl time.Duration) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = value
	m.expireAt = time.Now().Add(ttl)
	return nil
}

func (m *MockCacheWithStats) GetWithExpiry(ctx context.Context, key string) (string, time.Time, error) {
	result, err := m.Get(ctx, key)
	if err != nil {
		return "", time.Time{}, err
	}
	return result, m.expireAt, nil
}

func (m *MockCacheWithStats) SetWithExpiry(key, data string, ttl, actualExpiry time.Time) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.data[key] = data
	m.expiries[key] = actualExpiry
}

// MockRedisClient - a simplified interface for Redis-like operations
// Used for testing distributed lock scenarios
type MockRedisClient struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

func (r *MockRedisClient) SetNX(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[key]; exists {
		return false, nil // Key already exists
	}

	r.data[key] = value
	return true, nil
}

func (r *MockRedisClient) Get(ctx context.Context, key string) (string, bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	val, exists := r.data[key]
	return val, exists, nil
}

func (r *MockRedisClient) Del(ctx context.Context, key string) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[key]; exists {
		delete(r.data, key)
		return 1, nil
	}
	return 0, nil
}

// HeavyDB adalah simulasi DB untuk testing stampede
type HeavyDB struct {
	rebuildCount atomic.Int64
}

func NewHeavyDB() *HeavyDB {
	return &HeavyDB{}
}

func (db *HeavyDB) FetchHeavyData() string {
	db.rebuildCount.Add(1)
	time.Sleep(50 * time.Millisecond) // Simulate expensive query
	return "heavy_data_result"
}

func (db *HeavyDB) RebuildCount() int64 {
	return db.rebuildCount.Load()
}

func (db *HeavyDB) Reset() {
	db.rebuildCount.Store(0)
}

// Ensure interfaces are satisfied
var _ CacheInterface = (*MockCache)(nil)
var _ CacheInterface = (*MockCacheWithStats)(nil)