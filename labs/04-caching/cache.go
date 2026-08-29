// Package caching explores real-world caching patterns and their trade-offs.
// This lab focuses on: Cache Aside, Single Flight, Distributed Lock, Stampede prevention.
package caching

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Product adalah entitas domain untuk contoh caching.
type Product struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// CacheInterface mendefinisikan contract untuk cache yang bisa diganti.
type CacheInterface interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key, value string, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
	GetWithExpiry(ctx context.Context, key string) (string, time.Time, error)
}

// LockInterface mendefinisikan operation untuk distributed locking.
// SetNX = Set if Not eXists, dengan TTL.
type LockInterface interface {
	SetNX(ctx context.Context, key, value string, ttl time.Duration) (bool, error)
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, key string) (bool, error)
	// CompareAndDel atomically deletes key only if value matches (for safe release).
	CompareAndDel(ctx context.Context, key, value string) (bool, error)
}

// CacheKey membuat cache key dengan format: entity:id:vVersion
func CacheKey(entity, id string, version int) string {
	return fmt.Sprintf("%s:%s:v%d", entity, id, version)
}

// ShouldRefreshEarly menggunakan probabilitas untuk mencegah stampede.
// Jika data sudah 80% TTL berlalu, random 50% request akan refresh.
func ShouldRefreshEarly(ctx context.Context, cache CacheInterface, key string) bool {
	data, expiry, err := cache.GetWithExpiry(ctx, key)
	if err != nil || data == "" {
		return false
	}

	ttl := 5 * time.Minute
	age := time.Since(expiry)

	// Jika sudah 80% TTL berlalu, refresh randomly
	if age > ttl*8/10 {
		return rand.Float64() < 0.5
	}
	return false
}

// DashboardCacheKey membuat cache key untuk dashboard statistik.
// Format: cmms:dashboard:v1:branch:{branchID}:date:{YYYY-MM-DD}
func DashboardCacheKey(branchID int64) string {
	return fmt.Sprintf("cmms:dashboard:v1:branch:%d:date:%s", branchID, Today())
}

// extractID mengambil segment ID dari cache key format "entity:id[:...]".
func extractID(key string) string {
	parts := splitKey(key)
	if len(parts) >= 2 {
		return parts[1]
	}
	return key
}

// splitKey memisahkan key berdasarkan karakter ':'.
func splitKey(key string) []string {
	var parts []string
	current := ""
	for _, c := range key {
		if c == ':' {
			parts = append(parts, current)
			current = ""
		} else {
			current += string(c)
		}
	}
	return append(parts, current)
}
