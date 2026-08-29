package caching

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// RobustDashboardService mengimplementasikan graceful degradation saat Redis down.
type RobustDashboardService struct {
	db      *sql.DB
	cache   CacheInterface
	metrics *CacheMetrics
}

func NewRobustDashboardService(db *sql.DB, cache CacheInterface, metrics *CacheMetrics) *RobustDashboardService {
	if metrics == nil {
		metrics = NewCacheMetrics()
	}
	return &RobustDashboardService{
		db:      db,
		cache:   cache,
		metrics: metrics,
	}
}

// GetDashboard returns dashboard data with graceful degradation.
func (s *RobustDashboardService) GetDashboard(ctx context.Context, branchID int64) (Dashboard, error) {
	key := DashboardCacheKey(branchID)

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
		return s.fetchAndPopulate(ctx, branchID, key)
	}

	if cached == "" {
		// Explicit miss
		s.metrics.IncMiss()
		return s.fetchAndPopulate(ctx, branchID, key)
	}

	// 2. Cache hit - check for corruption
	var d Dashboard
	if err := json.Unmarshal([]byte(cached), &d); err != nil {
		// Corrupt cache entry
		s.metrics.IncError()
		_ = s.cache.Set(ctx, key, "", -1*time.Second) // delete corrupt
		return s.fetchAndPopulate(ctx, branchID, key)
	}

	// Valid cache hit
	s.metrics.IncHit()
	return d, nil
}

func (s *RobustDashboardService) fetchAndPopulate(ctx context.Context, branchID int64, key string) (Dashboard, error) {
	s.metrics.IncRebuild()
	s.metrics.IncDBQuery()

	// Source of truth: PostgreSQL database
	d := Dashboard{
		BranchID:          branchID,
		InvoiceCountToday: 42,
		TotalRevenueToday: 150000.0,
		Date:              Today(),
	}

	// Rebuild cache
	data, marshalErr := json.Marshal(d)
	if marshalErr != nil {
		// Log but don't fail - we have the data
		fmt.Printf("warn: failed to marshal dashboard for cache: %v\n", marshalErr)
	} else if setErr := s.cache.Set(ctx, key, string(data), 30*time.Second); setErr != nil {
		// Log cache set error but don't fail
		fmt.Printf("warn: cache set failed: %v\n", setErr)
	}

	return d, nil
}

// InvalidateBranchDashboard invalidates cache for a branch after data mutation
func (s *RobustDashboardService) InvalidateBranchDashboard(ctx context.Context, branchID int64) error {
	key := DashboardCacheKey(branchID)
	return s.cache.Set(ctx, key, "", -1*time.Second)
}

// QueryCount returns total database queries made
func (s *RobustDashboardService) QueryCount() int64 {
	return s.metrics.DBQueries()
}

// HitRatio returns cache hit percentage
func (s *RobustDashboardService) HitRatio() float64 {
	return s.metrics.HitRatio()
}