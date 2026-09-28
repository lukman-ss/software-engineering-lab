# Gap Analysis

## 1. Write-Behind Queue Overflow Drops Silently

Type: MISSING_EDGE_CASE
Severity: LOW
Location: internal/cache/patterns.go:156-160
Details: On channel buffer overflow, `Update` silently drops writes using the `default` branch of a non-blocking select. No error is returned and no test verifies this overflow behavior. The design decision is documented in implementation notes with a `ponytail:` comment.
Impact: Low — educational lab explicitly documents this as a known limitation.

## 2. No Test for SWR Synchronous Fallback on Fully Expired Item

Type: MISSING_TEST
Severity: LOW
Location: tests/cache_test.go, internal/cache/stampede.go:215-223
Details: The SWR service has 3 branches: (1) fully fresh, (2) stale-within-window (async revalidate), (3) fully expired / cold miss (synchronous fetch). Tests cover branch 2 only. Branch 3 (item present and past stale window, or entirely cold) is exercised only in the initial load path, not in an explicit test asserting the synchronous fallback.
Impact: Low — branch 3 was verified to execute correctly through demo and initial SWR test, but no isolated unit test explicitly covers it.

## 3. XFetch Uses Global `math/rand` for Production Seeding

Type: MISSING_EDGE_CASE
Severity: LOW
Location: internal/cache/store.go:84, internal/cache/stampede.go:6,105
Details: `TTLWithJitter` uses global `rand.Int63n` (unseeded implicit global source). `XFetchService` uses a locally seeded `rand.Rand` at creation time. For an educational lab this is acceptable, noted via `ponytail:` comment. Not a race or correctness bug.
Impact: Low — lab use only.

## 4. SWR Test Uses time.Sleep (Timing Sensitive)

Type: MISSING_EDGE_CASE
Severity: LOW
Location: tests/cache_test.go:202-220
Details: `TestStaleWhileRevalidate` relies on `time.Sleep` to wait for background async revalidation. This may flake under extreme system load (noted in implementation notes). No synchronization primitive (e.g. channel, `WaitGroup`) is used to confirm background goroutine completion.
Impact: Low — observed stable across multiple runs; documented in implementation notes.

## No CRITICAL or HIGH gaps found.
