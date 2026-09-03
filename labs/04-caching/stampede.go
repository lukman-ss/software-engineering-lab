package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/sync/singleflight"
)

// TTLWithJitter adds positive random jitter to prevent synchronized expiration (stampede).
// Returns base + random[0, maxJitter). Never reduces TTL below base.
// Upper bound is exclusive: maxJitter is not included.
func TTLWithJitter(base time.Duration, maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return base
	}
	jitter := time.Duration(rand.Int63n(int64(maxJitter)))
	return base + jitter
}

var fixedNow = time.Now()

// CounterRepository counts calls for stampede demonstration.
// Provides deterministic blocking/coordinating capabilities for tests.
type CounterRepository struct {
	callCount atomic.Int64
	mu        sync.Mutex
	blockCh   chan struct{}
	enteredCh chan struct{} // buffered, signals that a goroutine has entered GetDashboard
}

// Block sets up blocking behavior: all GetDashboard calls will block until Unblock is called.
// Returns immediately. Signals via enteredCh when a call enters GetDashboard.
func (r *CounterRepository) Block() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.blockCh = make(chan struct{})
	r.enteredCh = make(chan struct{}, 100)
}

// WaitUntilEntered blocks until a goroutine signals it has entered GetDashboard.
// Call after Block() - returns when first goroutine enters the repository.
func (r *CounterRepository) WaitUntilEntered() {
	r.mu.Lock()
	ch := r.enteredCh
	r.mu.Unlock()
	if ch != nil {
		<-ch
	}
}

// Unblock releases all blocked GetDashboard calls.
func (r *CounterRepository) Unblock() {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.blockCh != nil {
		close(r.blockCh)
		r.blockCh = nil
	}
}

func NewCounterRepository() *CounterRepository {
	return &CounterRepository{}
}

// queryDelay is the artificial delay before each DB query returns.
// This ensures a realistic time-of-check-to-time-of-use window for cache stampede demonstrations.
const queryDelay = 5 * time.Millisecond

func (r *CounterRepository) GetDashboard(ctx context.Context, tenantID, branchID int64, businessDate time.Time) (Dashboard, error) {
	// Signal entry into repository (for deterministic test coordination)
	r.mu.Lock()
	if r.enteredCh != nil {
		select {
		case r.enteredCh <- struct{}{}:
		default:
		}
	}
	// Check for blocking (thread-safe)
	ch := r.blockCh
	r.mu.Unlock()

	if ch != nil {
		select {
		case <-ch:
		case <-ctx.Done():
			return Dashboard{}, ctx.Err()
		}
	}

	// Simulate DB query latency so concurrent cache misses accumulate
	// before any single goroutine finishes and writes to cache.
	select {
	case <-ctx.Done():
		return Dashboard{}, ctx.Err()
	case <-time.After(queryDelay):
	}

	newCount := r.callCount.Add(1)
	return Dashboard{InvoiceCountToday: int(newCount)}, nil
}

func (r *CounterRepository) CallCount() int64 {
	return r.callCount.Load()
}

// --- VERSI BROKEN (STAMPEDE) ---

// BrokenStampedeService demonstrates cache stampede vulnerability.
// Multiple concurrent cache misses result in parallel DB queries.
// INTENTIONALLY BROKEN FOR LAB DEMONSTRATION.
// Do not use as production code.
type BrokenStampedeService struct {
	cache CacheInterface
	repo  *CounterRepository
	ttl   time.Duration
}

func NewBrokenStampedeService(cache CacheInterface, repo *CounterRepository) *BrokenStampedeService {
	return &BrokenStampedeService{cache: cache, repo: repo, ttl: 5 * time.Minute}
}

func (s *BrokenStampedeService) GetData(ctx context.Context, branchID int64) (Dashboard, error) {
	key := fmt.Sprintf("dash:%d", branchID)

	// 1. Check cache
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var d Dashboard
		if unmarshalErr := json.Unmarshal([]byte(cached), &d); unmarshalErr == nil {
			return d, nil
		}
		// If unmarshal fails, we still fall through to DB
	}

	// 2. Cache miss -> parallel DB queries (NO protection)
	// In real system: 100 concurrent requests → 100 DB queries
	d, err := s.repo.GetDashboard(ctx, 1, branchID, fixedNow.UTC())
	if err != nil {
		return Dashboard{}, err
	}

	// 3. Set cache
	data, err := json.Marshal(d)
	if err != nil {
		return Dashboard{}, fmt.Errorf("marshal: %w", err)
	}
	if err := s.cache.Set(ctx, key, string(data), s.ttl); err != nil {
		// Log but don't fail
		fmt.Printf("warn: cache set failed: %v\n", err)
	}

	return d, nil
}

// --- VERSI PROTECTED (SINGLEFLIGHT) ---

// ProtectedStampedeService protects against cache stampede using singleflight.
// Concurrent requests for same key share a single DB query result.
type ProtectedStampedeService struct {
	cache  CacheInterface
	repo   *CounterRepository
	flight singleflight.Group
	ttl    time.Duration
}

func NewProtectedStampedeService(cache CacheInterface, repo *CounterRepository) *ProtectedStampedeService {
	return &ProtectedStampedeService{cache: cache, repo: repo, ttl: 5 * time.Minute}
}

func (s *ProtectedStampedeService) GetData(ctx context.Context, branchID int64) (Dashboard, error) {
	key := fmt.Sprintf("dash:%d", branchID)

	// 1. First cache check (caller's context)
	cached, err := s.cache.Get(ctx, key)
	if err == nil && cached != "" {
		var d Dashboard
		if unmarshalErr := json.Unmarshal([]byte(cached), &d); unmarshalErr == nil {
			return d, nil
		}
		// If unmarshal fails, fall through to singleflight
	}

	// 2. Cache miss -> use singleflight with DoChan for context-awareness
	// DoChan returns a channel that delivers the result, allowing
	// callers to also select on ctx.Done() for proper cancellation.
	ch := s.flight.DoChan(key, func() (interface{}, error) {
		// Create bounded context for shared rebuild FIRST:
		// - created before any cache/DB operations to ensure all use same bounded context
		// - context.WithoutCancel strips the leader's cancellation
		// - Individual callers can still cancel their own wait via ctx.Done()
		// - Shared rebuild and cache operations continue even if leader cancels
		// - Timeout added so rebuild doesn't run forever
		rebuildCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 30*time.Second)
		defer cancel()

		// Double-check cache inside singleflight gate using rebuildCtx
		// This handles the window where another goroutine may have populated cache
		// between our first check and acquiring the singleflight slot
		cached, err := s.cache.Get(rebuildCtx, key)
		if err == nil && cached != "" {
			var d Dashboard
			if unmarshalErr := json.Unmarshal([]byte(cached), &d); unmarshalErr == nil {
				return d, nil
			}
		}

		// Only one goroutine reaches here - it fetches from DB using rebuildCtx
		d, err := s.repo.GetDashboard(rebuildCtx, 1, branchID, fixedNow.UTC())
		if err != nil {
			return Dashboard{}, err
		}

		// Populate cache using rebuildCtx (not leader's ctx)
		// This ensures cache SET succeeds even if leader cancelled
		data, err := json.Marshal(d)
		if err != nil {
			return Dashboard{}, fmt.Errorf("marshal: %w", err)
		}
		if setErr := s.cache.Set(rebuildCtx, key, string(data), s.ttl); setErr != nil {
			fmt.Printf("warn: cache set failed: %v\n", setErr)
		}
		return d, nil
	})

	// 3. Wait with context awareness - caller can cancel waiting independently
	select {
	case <-ctx.Done():
		return Dashboard{}, ctx.Err()
	case result := <-ch:
		if result.Err != nil {
			return Dashboard{}, result.Err
		}
		return result.Val.(Dashboard), nil
	}
}
