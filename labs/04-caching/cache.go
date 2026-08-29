// Package caching explores real-world caching patterns and their trade-offs.
// This lab focuses on: Cache Aside, Single Flight, Distributed Lock, Stampede prevention.
package caching

import (
	"context"
	"fmt"
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

// ShouldRefreshEarly analyzes whether a cached item should be refreshed early to prevent stampede.
// It returns whether refresh should happen, the remaining TTL, and whether the entry is expired.
//
// The function uses probabilistic early refresh when remaining TTL falls within the refresh window.
// Refresh window is 20% of the original TTL (configurable via probability).
func ShouldRefreshEarly(now time.Time, expiry time.Time, originalTTL time.Duration, probability float64, randomFloat func() float64) bool {
	if expiry.IsZero() {
		return false
	}

	rem := expiry.Sub(now)
	if rem < 0 {
		// Already expired - handled as cache miss elsewhere
		return false
	}

	// Refresh window is 20% of TTL
	refreshWindow := originalTTL - (originalTTL * 80 / 100) // = 20% of TTL

	if rem <= refreshWindow {
		return randomFloat() < probability
	}
	return false
}

// DashboardCacheKey creates canonical cache key for dashboard statistics.
// Format: cmms:dashboard:v1:tenant:{tenantID}:branch:{branchID}:date:{YYYY-MM-DD}
func DashboardCacheKey(tenantID, branchID int64, businessDate time.Time) string {
	return NewDashboardKey(branchID).WithTenant(tenantID).WithDate(businessDate).Build()
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
