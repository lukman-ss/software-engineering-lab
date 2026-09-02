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
	// CompareAndDel atomically deletes key only if value matches (for safe release).
	CompareAndDel(ctx context.Context, key, value string) (bool, error)
}

// CacheKey membuat cache key dengan format: entity:id:vVersion
func CacheKey(entity, id string, version int) string {
	return fmt.Sprintf("%s:%s:v%d", entity, id, version)
}

// ShouldRefreshEarly analyzes whether a cached item should be refreshed early to prevent stampede.
// Returns true if refresh should happen, false otherwise.
//
// The function uses probabilistic early refresh when remaining TTL falls within the refresh window.
// Refresh window is 20% of the original TTL (configurable via probability).
//
// Guard conditions:
// - originalTTL <= 0 → false (invalid TTL)
// - probability <= 0 → false (no refresh chance)
// - probability >= 1 → true if within refresh window (always refresh)
// - expiry zero → false (unknown expiry)
// - expired → false (handled as cache miss)
func ShouldRefreshEarly(now time.Time, expiry time.Time, originalTTL time.Duration, probability float64, randomFloat func() float64) bool {
	// Guard: invalid TTL
	if originalTTL <= 0 {
		return false
	}

	// Guard: zero expiry means unknown/undefined expiry
	if expiry.IsZero() {
		return false
	}

	rem := expiry.Sub(now)

	// Guard: already expired - handled as cache miss elsewhere
	if rem < 0 {
		return false
	}

	// Guard: negative or zero probability
	if probability <= 0 {
		return false
	}

	// Guard: probability >= 1 means always refresh within window
	if probability >= 1 {
		if rem <= originalTTL/5 { // 20% of TTL
			return true
		}
		return false
	}

	// Refresh window is 20% of TTL
	refreshWindow := originalTTL / 5 // = 20% of TTL

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
