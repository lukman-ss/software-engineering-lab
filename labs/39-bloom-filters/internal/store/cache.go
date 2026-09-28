package store

import (
	"sync"
	"sync/atomic"

	"bloomfilters/internal/bloom"
)

// Backend simulates an expensive downstream data source.
type Backend interface {
	Fetch(key string) (string, bool)
}

// Cache is an in-memory cache with an optional Bloom filter admission gate.
// The gate prevents cache-penetration: absent-key queries bypass the cache
// and only reach the backend when the filter cannot rule them out.
type Cache struct {
	mu           sync.RWMutex
	data         map[string]string
	filter       *bloom.SyncFilter // nil when disabled
	backendCalls uint64
	backend      Backend
}

// NewCache creates a cache.
// If epsilon > 0, a SyncFilter admission gate is allocated with capacity n.
func NewCache(n uint, epsilon float64, backend Backend) *Cache {
	c := &Cache{
		data:    make(map[string]string),
		backend: backend,
	}
	if epsilon > 0 && n > 0 {
		c.filter = bloom.NewSync(n, epsilon)
	}
	return c
}

// Add inserts a key-value pair into the cache (and the Bloom filter if enabled).
func (c *Cache) Add(key, value string) {
	c.mu.Lock()
	c.data[key] = value
	c.mu.Unlock()
	if c.filter != nil {
		c.filter.Add([]byte(key))
	}
}

// Get retrieves a value.
// Flow:
//  1. If filter enabled and key is definitely absent → return miss (no backend call).
//  2. Check in-memory cache.
//  3. If not in cache → call backend, store result.
func (c *Cache) Get(key string) (string, bool) {
	if c.filter != nil && !c.filter.Check([]byte(key)) {
		// Definitely not a known key – reject without hitting backend.
		return "", false
	}

	c.mu.RLock()
	v, ok := c.data[key]
	c.mu.RUnlock()
	if ok {
		return v, true
	}

	// Cache miss: call backend.
	atomic.AddUint64(&c.backendCalls, 1)
	v, ok = c.backend.Fetch(key)
	if ok {
		c.mu.Lock()
		c.data[key] = v
		c.mu.Unlock()
		if c.filter != nil {
			c.filter.Add([]byte(key))
		}
	}
	return v, ok
}

// BackendCalls returns the number of times the backend was called.
func (c *Cache) BackendCalls() uint64 {
	return atomic.LoadUint64(&c.backendCalls)
}

// ResetCounters clears backend call counter.
func (c *Cache) ResetCounters() {
	atomic.StoreUint64(&c.backendCalls, 0)
}
