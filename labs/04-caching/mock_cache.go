package caching

import (
	"context"
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
		return "", ErrCacheMiss
	}

	// Check expiry
	if !m.expireAt.IsZero() {
		if time.Now().After(m.expireAt) {
			m.misses++
			return "", ErrCacheMiss
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

func (m *MockCacheWithStats) Delete(ctx context.Context, key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.data, key)
	delete(m.expiries, key)
	m.expireAt = time.Time{}
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

// MockRedisClient - implements LockInterface for distributed locking tests.
// Provides atomic SET NX and compare-and-delete (via ReleaseLock in distributed_lock.go).
type MockRedisClient struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMockRedisClient() *MockRedisClient {
	return &MockRedisClient{
		data: make(map[string]string),
	}
}

// SetNX atomically sets key only if it doesn't exist.
// Returns true if key was set, false if key already existed.
func (r *MockRedisClient) SetNX(ctx context.Context, key, value string, ttl time.Duration) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[key]; exists {
		return false, nil // Key already exists
	}

	r.data[key] = value
	return true, nil
}

// Get retrieves value (for lock token verification).
func (r *MockRedisClient) Get(ctx context.Context, key string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	val, exists := r.data[key]
	if !exists {
		return "", nil
	}
	return val, nil
}

// Del deletes a key by value verification.
// In Redis, this would be a Lua script: if GET == ARGV then DEL.
// For simplicity, this returns whether key existed and was deleted.
func (r *MockRedisClient) Del(ctx context.Context, key string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.data[key]; exists {
		delete(r.data, key)
		return true, nil
	}
	return false, nil
}

// CompareAndDel atomically deletes only if value matches.
// Returns true if deleted, false if value didn't match or key missing.
func (r *MockRedisClient) CompareAndDel(ctx context.Context, key, value string) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	val, exists := r.data[key]
	if !exists || val != value {
		return false, nil // Key missing or value mismatch
	}

	delete(r.data, key)
	return true, nil
}

// HeavyDB is for stampede testing

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
var _ LockInterface = (*MockRedisClient)(nil)
