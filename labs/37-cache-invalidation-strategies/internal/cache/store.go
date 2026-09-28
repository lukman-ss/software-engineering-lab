package cache

import (
	"errors"
	"math/rand"
	"sync"
	"time"
)

var ErrCacheMiss = errors.New("cache miss")

type Item struct {
	Value     string
	CreatedAt time.Time
	TTL       time.Duration
	ExpiresAt time.Time
	ReadDelta time.Duration // Δ (compute/fetch duration)
}

type MemoryCache struct {
	mu    sync.RWMutex
	items map[string]Item
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		items: make(map[string]Item),
	}
}

func (c *MemoryCache) Get(key string) (Item, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return Item{}, ErrCacheMiss
	}
	if !item.ExpiresAt.IsZero() && time.Now().After(item.ExpiresAt) {
		return Item{}, ErrCacheMiss
	}
	return item, nil
}

// GetRaw returns item regardless of expiration, needed for SWR and XFetch inspection.
func (c *MemoryCache) GetRaw(key string) (Item, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, ok := c.items[key]
	return item, ok
}

func (c *MemoryCache) Set(key string, value string, ttl time.Duration, readDelta time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()
	var exp time.Time
	if ttl > 0 {
		exp = now.Add(ttl)
	}

	c.items[key] = Item{
		Value:     value,
		CreatedAt: now,
		TTL:       ttl,
		ExpiresAt: exp,
		ReadDelta: readDelta,
	}
}

func (c *MemoryCache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// TTLWithJitter adds random positive jitter to prevent synchronized expiration.
func TTLWithJitter(base time.Duration, maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return base
	}
	// ponytail: standard rand package used for simple jitter generation; add crypto/rand or custom source if security-sensitive
	jitter := time.Duration(rand.Int63n(int64(maxJitter)))
	return base + jitter
}
