package caching

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// RobustDashboardService mengimplementasikan graceful degradation saat Redis down.
type RobustDashboardService struct {
	repo    DashboardRepository
	cache   CacheInterface
	metrics *CacheMetrics
}

func NewRobustDashboardService(repo DashboardRepository, cache CacheInterface, metrics *CacheMetrics) *RobustDashboardService {
	if metrics == nil {
		metrics = NewCacheMetrics()
	}
	return &RobustDashboardService{
		repo:    repo,
		cache:   cache,
		metrics: metrics,
	}
}

// GetDashboard returns dashboard data with graceful degradation.
func (s *RobustDashboardService) GetDashboard(ctx context.Context, branchID int64) (Dashboard, error) {
	return s.GetDashboardWithTenant(ctx, 1, branchID, time.Now().UTC())
}

// GetDashboardWithTenant returns dashboard data with tenant isolation.
// Primary entry point for multi-tenant deployments.
func (s *RobustDashboardService) GetDashboardWithTenant(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error) {
	key := NewDashboardKey(branchID).WithTenant(tenantID).WithDate(businessDate).Build()

	// 1. Check cache
	cached, err := s.cache.Get(ctx, key)
	if err != nil {
		// Categorize error
		switch {
		case errors.Is(err, ErrCacheMiss):
			// Cache miss - expected, go to DB
			s.metrics.IncMiss()
		default:
			// Cache down/network error - fallback
			s.metrics.IncError()
		}
		return s.fetchAndPopulate(ctx, tenantID, branchID, businessDate, key)
	}

	if cached == "" {
		// Explicit miss
		s.metrics.IncMiss()
		return s.fetchAndPopulate(ctx, tenantID, branchID, businessDate, key)
	}

	// 2. Cache hit - check for corruption
	var d Dashboard
	if err := json.Unmarshal([]byte(cached), &d); err != nil {
		// Corrupt cache entry
		s.metrics.IncError()
		if delErr := s.cache.Delete(ctx, key); delErr != nil {
			fmt.Printf("warn: failed to delete corrupt cache key %s: %v\n", key, delErr)
		}
		return s.fetchAndPopulate(ctx, tenantID, branchID, businessDate, key)
	}

	// Valid cache hit
	s.metrics.IncHit()
	return d, nil
}

func (s *RobustDashboardService) fetchAndPopulate(ctx context.Context, tenantID, branchID int64, businessDate time.Time, key string) (Dashboard, error) {
	s.metrics.IncRebuildAttempt()

	// Query repository for real data - use tenant-aware method for proper isolation
	s.metrics.IncDBQuery() // Ensure DBQuery metric increments exactly when repo is called
	d, err := s.repo.GetDashboard(ctx, tenantID, branchID, businessDate)
	if err != nil {
		return Dashboard{}, fmt.Errorf("repo get: %w", err)
	}

	// Rebuild cache
	data, marshalErr := json.Marshal(d)
	if marshalErr != nil {
		// Log but don't fail - we have the data
		fmt.Printf("warn: failed to marshal dashboard for cache: %v\n", marshalErr)
	} else {
		// Add jitter to prevent synchronized cache expiration across branches
		jitteredTTL := TTLWithJitter(30*time.Second, 10*time.Second)
		if setErr := s.cache.Set(ctx, key, string(data), jitteredTTL); setErr != nil {
			// Log cache set error but don't fail the request
			s.metrics.IncError()
			fmt.Printf("warn: cache set failed: %v\n", setErr)
		} else {
			s.metrics.IncRebuildSuccess()
		}
	}

	return d, nil
}

// InvalidateDashboard invalidates cache for a branch and specific business date.
//
// LOW-LEVEL INVALIDATE METHOD:
// This method can return error if cache DELETE fails (e.g., Redis down, network error).
// Callers should be aware that this is a cache side-effect failure, NOT a transaction failure.
//
// BUSINESS MUTATION FLOW:
// DB mutation → COMMIT success → cache invalidation attempt → if invalidation fails, record error but business write succeeded → TTL serves as safety net for stale cache
//
// The error return allows monitoring/debugging, but business operations must handle
// cache errors as non-blocking side-effects, not transaction rollbacks.
func (s *RobustDashboardService) InvalidateDashboard(ctx context.Context, tenantID, branchID int64, businessDate time.Time) error {
	key := NewDashboardKey(branchID).WithTenant(tenantID).WithDate(businessDate).Build()
	err := s.cache.Delete(ctx, key)
	if err != nil && !errors.Is(err, ErrCacheMiss) {
		s.metrics.IncError()
		// Log and return error - caller should handle as cache side-effect
		// Business operation already succeeded if we're here
		fmt.Printf("warn: failed to invalidate cache key %s: %v\n", key, err)
		return fmt.Errorf("invalidate cache: %w", err)
	}
	return nil
}

// InvalidateCurrentDashboard invalidates cache for today using the injected clock.
func (s *RobustDashboardService) InvalidateCurrentDashboard(ctx context.Context, tenantID, branchID int64) error {
	return s.InvalidateDashboard(ctx, tenantID, branchID, defaultClock.Now().UTC())
}

// QueryCount returns total database queries made
func (s *RobustDashboardService) QueryCount() int64 {
	return s.metrics.DBQueries()
}

// HitRatio returns cache hit percentage
func (s *RobustDashboardService) HitRatio() float64 {
	return s.metrics.HitRatio()
}
