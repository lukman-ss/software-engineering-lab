package ratelimit

import "sync"

// Registry maintains per-tenant rate limiters to avoid shared IP / CGNAT degradation (RFC 6598).
type Registry struct {
	mu          sync.RWMutex
	buckets     map[string]*TokenBucket
	capacity    float64
	refillRate  float64
}

func NewRegistry(defaultCapacity, defaultRefillRate float64) *Registry {
	return &Registry{
		buckets:     make(map[string]*TokenBucket),
		capacity:    defaultCapacity,
		refillRate:  defaultRefillRate,
	}
}

func (r *Registry) Get(tenantKey string) *TokenBucket {
	r.mu.RLock()
	tb, exists := r.buckets[tenantKey]
	r.mu.RUnlock()
	if exists {
		return tb
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	if tb, exists = r.buckets[tenantKey]; exists {
		return tb
	}
	tb = NewTokenBucket(r.capacity, r.refillRate)
	r.buckets[tenantKey] = tb
	return tb
}
