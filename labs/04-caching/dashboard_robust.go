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
	key := NewTenantDashboardKey(tenantID, branchID, businessDate).Build()

	// 1. Check cache with latency measurement
	startGet := time.Now()
	cached, err := s.cache.Get(ctx, key)
	s.metrics.IncCacheGetOp()
	s.metrics.RecordCacheGetLatency(time.Since(startGet))

	if err != nil {
		switch {
		case errors.Is(err, ErrCacheMiss):
			s.metrics.IncMiss()
			return s.fetchAndPopulate(ctx, tenantID, branchID, businessDate, key, false)
		default:
			s.metrics.IncError()
			s.metrics.IncDBFallback()
			return s.fetchAndPopulate(ctx, tenantID, branchID, businessDate, key, true)
		}
	}

	if cached == "" {
		s.metrics.IncMiss()
		return s.fetchAndPopulate(ctx, tenantID, branchID, businessDate, key, false)
	}

	// 2. Cache hit - check for corruption
	var d Dashboard
	if err := json.Unmarshal([]byte(cached), &d); err != nil {
		s.metrics.IncError()
		s.metrics.IncCacheInvalidateOp()
		if delErr := s.cache.Delete(ctx, key); delErr != nil && !errors.Is(delErr, ErrCacheMiss) {
			s.metrics.IncCacheInvalidationError()
			fmt.Printf("warn: failed to delete corrupt cache key %s: %v\n", key, delErr)
		}
		return s.fetchAndPopulate(ctx, tenantID, branchID, businessDate, key, false)
	}

	// Valid cache hit
	s.metrics.IncHit()
	return d, nil
}

func (s *RobustDashboardService) fetchAndPopulate(ctx context.Context, tenantID, branchID int64, businessDate time.Time, key string, isFallback bool) (Dashboard, error) {
	s.metrics.IncRebuildAttempt()
	startRebuild := time.Now()

	s.metrics.IncDBQuery()
	startDB := time.Now()
	d, err := s.repo.GetDashboard(ctx, tenantID, branchID, businessDate)
	if isFallback {
		s.metrics.RecordDBFallbackLatency(time.Since(startDB))
	}
	if err != nil {
		return Dashboard{}, fmt.Errorf("repo get: %w", err)
	}

	// Rebuild cache
	data, marshalErr := json.Marshal(d)
	if marshalErr != nil {
		fmt.Printf("warn: failed to marshal dashboard for cache: %v\n", marshalErr)
		s.metrics.IncError()
	} else {
		jitteredTTL := TTLWithJitter(30*time.Second, 10*time.Second)
		s.metrics.IncCacheSetOp()
		startSet := time.Now()
		setErr := s.cache.Set(ctx, key, string(data), jitteredTTL)
		s.metrics.RecordCacheSetLatency(time.Since(startSet))
		if setErr != nil {
			s.metrics.IncError()
			s.metrics.IncCacheSetError()
			fmt.Printf("warn: cache set failed: %v\n", setErr)
		} else {
			s.metrics.IncRebuildSuccess()
		}
	}

	s.metrics.RecordRebuildLatency(time.Since(startRebuild))
	return d, nil
}

// InvalidateDashboard invalidates cache for a branch and specific business date.
func (s *RobustDashboardService) InvalidateDashboard(ctx context.Context, tenantID, branchID int64, businessDate time.Time) error {
	key := NewTenantDashboardKey(tenantID, branchID, businessDate).Build()
	s.metrics.IncCacheInvalidateOp()
	err := s.cache.Delete(ctx, key)
	if err != nil && !errors.Is(err, ErrCacheMiss) {
		s.metrics.IncError()
		s.metrics.IncCacheInvalidationError()
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
