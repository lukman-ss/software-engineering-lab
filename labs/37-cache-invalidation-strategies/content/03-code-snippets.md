## Snippet 1 — Cache-Aside: Read and Update

Source File: `internal/cache/patterns.go:22-47`

Purpose: Demonstrates lazy-load-on-miss with write-then-invalidate ordering.

```go
func (s *CacheAsideService) Get(ctx context.Context, key string) (string, error) {
	item, err := s.cache.Get(key)
	if err == nil {
		return item.Value, nil
	}

	start := time.Now()
	val, err := s.db.Query(ctx, key)
	if err != nil {
		return "", err
	}
	delta := time.Since(start)

	s.cache.Set(key, val, s.ttl, delta)
	return val, nil
}

func (s *CacheAsideService) Update(ctx context.Context, key, val string) error {
	// 1. Store first
	if err := s.db.Write(ctx, key, val); err != nil {
		return fmt.Errorf("db write failed: %w", err)
	}
	// 2. Invalidate cache
	s.cache.Delete(key)
	return nil
}
```

Explanation: Read path: check cache; on miss query DB and populate cache with TTL. Update path: write to DB first, then delete the cache key. Order matters — deleting cache before DB write opens a window where a concurrent reader can re-populate cache with stale data.

---

## Snippet 2 — Write-Through: Synchronous Dual Write

Source File: `internal/cache/patterns.go:78-89`

Purpose: Ensures subsequent reads see the updated value without DB query.

```go
func (s *WriteThroughService) Update(ctx context.Context, key, val string) error {
	start := time.Now()
	// 1. Synchronous write to store
	if err := s.db.Write(ctx, key, val); err != nil {
		return fmt.Errorf("db write failed: %w", err)
	}
	delta := time.Since(start)

	// 2. Synchronous write to cache
	s.cache.Set(key, val, s.ttl, delta)
	return nil
}
```

Explanation: Write both DB and cache in sequence. The `delta` measured is DB write duration and is recorded in `ReadDelta` field (minor semantic mismatch noted in audit Finding 4). Reader immediately hits cache on next access.

---

## Snippet 3 — Write-Behind: Async Flush

Source File: `internal/cache/patterns.go:152-161`

Purpose: Immediate cache update with deferred DB flush via background worker.

```go
func (s *WriteBehindService) Update(key, val string) {
	// Write immediately to cache
	s.cache.Set(key, val, s.ttl, 1*time.Millisecond)
	// Enqueue write to backing store asynchronously
	select {
	case s.writeQueue <- WriteRequest{Key: key, Value: val}:
	default:
		// Queue full (demonstration: drop or handle overflow)
	}
}
```

Explanation: Cache updated immediately; write enqueued to buffered channel. On full queue, the write is silently dropped — an implementation limitation, not a production recommendation (ponytail comment in source). Background `flushWorker` drains queue and calls `db.Write` in loop.

---

## Snippet 4 — SingleFlight Stampede Mitigation

Source File: `internal/cache/stampede.go:56-84`

Purpose: Coalesces concurrent misses on the same key to exactly one DB query.

```go
func (s *SingleFlightService) Get(ctx context.Context, key string) (string, error) {
	item, err := s.cache.Get(key)
	if err == nil {
		return item.Value, nil
	}

	res, err, _ := s.flight.Do(key, func() (interface{}, error) {
		// Double check inside flight execution
		item, err := s.cache.Get(key)
		if err == nil {
			return item.Value, nil
		}

		start := time.Now()
		val, err := s.db.Query(ctx, key)
		if err != nil {
			return "", err
		}
		delta := time.Since(start)

		s.cache.Set(key, val, s.ttl, delta)
		return val, nil
	})

	if err != nil {
		return "", err
	}
	return res.(string), nil
}
```

Explanation: First miss enters `flight.Do`. The closure performs a double-check cache hit before querying DB. All concurrent callers with the same key block on `Do` and receive the same result when the first completes. Verified by test: 20 goroutines → 1 DB query.

---

## Snippet 5 — XFetch Probabilistic Early Expiration Formula

Source File: `internal/cache/stampede.go:125-136`

Purpose: Core mathematical condition for proactive early refresh.

```go
func ShouldRecompute(delta time.Duration, beta float64, ttlRemaining time.Duration, u float64) bool {
	if u <= 0 || u >= 1 {
		return false
	}
	deltaSec := delta.Seconds()
	ttlRemainingSec := ttlRemaining.Seconds()

	// Correct mathematical formula with negative log draw:
	// -delta * beta * log(u)
	expiryCompute := -deltaSec * beta * math.Log(u)
	return expiryCompute > ttlRemainingSec
}
```

Explanation: `u` from Uniform(0,1); `math.Log(u)` is negative in this range, so `-delta*beta*ln(u)` is positive. Condition fires when the computed positive offset exceeds remaining TTL. Guard prevents `u≤0` (log infinity) and `u≥1` (no-op). Test explicitly checks erroneous sign `(delta*beta*ln(u) > remaining)` yields false.

---

## Snippet 6 — Stale-While-Revalidate (SWR)

Source File: `internal/cache/stampede.go:197-224`

Purpose: Serve stale data immediately while triggering async background revalidation.

```go
func (s *SWRService) Get(ctx context.Context, key string) (string, error) {
	item, ok := s.cache.GetRaw(key)
	now := time.Now()

	if ok {
		// 1. Fully fresh
		if item.ExpiresAt.IsZero() || now.Before(item.ExpiresAt) {
			return item.Value, nil
		}

		// 2. Stale but within stale window
		staleUntil := item.ExpiresAt.Add(s.staleDelta)
		if now.Before(staleUntil) {
			s.triggerRevalidate(key)
			return item.Value, nil
		}
	}

	// 3. Completely expired or miss -> synchronous fetch
	start := time.Now()
	val, err := s.db.Query(ctx, key)
	if err != nil {
		return "", err
	}
	delta := time.Since(start)
	s.cache.Set(key, val, s.ttl, delta)
	return val, nil
}
```

Explanation: Three-path logic: fresh → return immediately; stale-within-window → return stale + trigger async revalidate; fully expired/missing → sync fetch. Revalidation deduplicated via `revalidating` map.

---

## Snippet 7 — TTL Jitter

Source File: `internal/cache/store.go:79-86`

Purpose: Add randomized positive offset to prevent synchronized key expiration.

```go
func TTLWithJitter(base time.Duration, maxJitter time.Duration) time.Duration {
	if maxJitter <= 0 {
		return base
	}
	jitter := time.Duration(rand.Int63n(int64(maxJitter)))
	return base + jitter
}
```

Explanation: Returns `base + rand[0, maxJitter)`. Useful when many keys share the same TTL origin (deploy, cron). Does not prevent stampede on a single hot key — only desynchronizes across keys.
