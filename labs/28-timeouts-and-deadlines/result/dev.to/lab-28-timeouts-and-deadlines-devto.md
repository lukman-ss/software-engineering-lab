---
title: "Timeouts and Deadlines: Preventing Slow Dependencies from Bringing Down the System"
published: false
description: "Learn how context deadlines, exponential backoff with full jitter, circuit breakers, and idempotency keys protect distributed systems from cascading failures."
tags: go, architecture, distributed, resilience
---

# Timeouts and Deadlines: Preventing Slow Dependencies from Bringing Down the System

## 1. Introduction and Little's Law

Slow dependencies kill systems more reliably than dependencies that fail fast. A fast failure releases the goroutine, releases the connection, releases the slot in the worker pool. A slow dependency holds all three. Under load, held resources accumulate. Throughput collapses while CPU looks idle. This is the standard shape of a cascading failure.

Little's Law makes the mechanism quantitative: `L = lambda * W`. `L` is the average number of requests in flight. `lambda` is arrival rate. `W` is average time in system. When a downstream dependency shifts latency from 30ms to 30s, `W` grows by 1000x. For constant `lambda`, `L` grows by 1000x. Every one of those in-flight requests holds a goroutine, a socket, possibly a database connection. The process hits file descriptor limits, pool exhaustion, GC pressure. Independent endpoints sharing the same process starve even though their own dependencies remain healthy. That cross-endpoint contamination is the reason timeouts belong at every boundary.

Lab 28 implements four defenses against this dynamic in Go: context deadline propagation in `internal/deadline`, exponential backoff with full jitter in `internal/retry`, a circuit breaker state machine in `internal/circuit`, and idempotency key deduplication in `internal/idempotency`. Integration tests in `tests/integration_test.go` verify the components operate correctly in combination. The demo in `cmd/demo/main.go` exercises all four paths end to end.

## 2. The 4 Mental Models of Timeouts

Four models structure every timeout decision in this lab.

**Model 1: Indeterminate status.** A timeout reports that no response arrived within budget. It never reports that the operation failed. The downstream operation may have completed, may be running, may never have started. Any retry policy built on timeouts operates under this ambiguity. Section 7 covers the deduplication mechanism that makes retries safe despite it.

**Model 2: Hierarchical inheritance.** Deadlines propagate down the call tree. A parent context with a 50ms deadline constrains every child derived from it, regardless of the larger budget a child requests. `TestExecuteWithBudget_ParentTimeoutInherited` proves this behavior: a parent with 20ms forces failure on a child requesting 500ms. Timeouts compose by taking the minimum along the path.

**Model 3: Jittered retry versus retry storm.** A fixed backoff synchronizes clients. When a dependency recovers after an outage, every waiting client retries at the same instant. The synchronized spike re-overloads the dependency. This thundering herd is called a retry storm. Full jitter randomizes each sleep over a uniform distribution, spreading load across time. Section 5 shows the exact math from the AWS Architecture Blog.

**Model 4: Storage-level timeouts.** Application timeouts cap waiting inside Go code. They do not cap work inside PostgreSQL. A query cancelled by `context.WithTimeout` on the client side can continue consuming database CPU until it finishes unless the database enforces its own `statement_timeout`. Section 9 covers the three PostgreSQL settings that close this gap.

## 3. Architecture and 4 Core Components

The lab is a standard Go module with four internal packages, one integration test package, and one demo binary:

```
labs/28-timeouts-and-deadlines/
  internal/deadline/deadline.go
  internal/retry/retry.go
  internal/circuit/circuit.go
  internal/idempotency/idempotency.go
  tests/integration_test.go
  cmd/demo/main.go
```

| Component | File | Responsibility |
|---|---|---|
| Deadline budget | `internal/deadline/deadline.go` | Execute a worker under a derived `context.WithTimeout`, return `context.DeadlineExceeded` on expiry |
| Jittered retrier | `internal/retry/retry.go` | Retry a fallible function up to `MaxAttempts` with full-jitter backoff capped at `MaxBackoff`, abort early on context cancellation |
| Circuit breaker | `internal/circuit/circuit.go` | Track consecutive failures under `sync.RWMutex`, trip CLOSED to OPEN at threshold, admit a probe after cooldown via HALF_OPEN |
| Idempotency store | `internal/idempotency/idempotency.go` | Map idempotency keys to stored responses with TTL, evict expired entries lazily on `Get` under write lock |

Each component has zero third-party dependencies. Standard library only: `context`, `time`, `sync`, `math/rand/v2`, `errors`. Defaults are defensive: `NewRetrier` fills zero `Config` with 3 attempts, 10ms base, 100ms max. `NewBreaker` fills zero `Config` with threshold 3, success threshold 2, 100ms cooldown. `NewStore` fills non-positive TTL with 1 minute. `TestRetrier_JitterBoundsAndZeroConfig` and `TestCircuitBreaker_DefaultZeroConfig` lock this behavior.

## 4. Context Deadline Propagation in Go

`ExecuteWithBudget` is the full implementation:

```go
func ExecuteWithBudget(ctx context.Context, budget time.Duration, fn WorkerFunc) error {
    childCtx, cancel := context.WithTimeout(ctx, budget)
    defer cancel()

    done := make(chan error, 1)
    go func() {
        done <- fn(childCtx)
    }()

    select {
    case <-childCtx.Done():
        return childCtx.Err()
    case err := <-done:
        return err
    }
}
```

Three details carry the design.

First, `context.WithTimeout(ctx, budget)` derives from the incoming context. When the parent already carries an earlier deadline, the child inherits the earlier one. Go resolves nested timeouts by earliest expiry. The test with a 20ms parent and 500ms child budget asserts `context.DeadlineExceeded`, confirming inheritance works without additional code.

Second, `done := make(chan error, 1)` uses a buffered channel with capacity 1. The worker goroutine sends exactly once. When the `select` takes the timeout branch, no receiver remains on `done`. With an unbuffered channel, the worker would block on send forever, leaking the goroutine. The buffer of 1 lets the late worker deposit its result and exit. This single integer is the goroutine leak fix.

Third, `defer cancel()` releases the context timer. `context.WithTimeout` allocates a timer that lives until expiry or cancellation. Omitting `cancel` keeps the timer alive until the full budget elapses, even after the worker returns. Under high request rates those live timers accumulate. The deferred call frees them on return along either `select` branch.

Callers pass `childCtx` into the worker so blocking operations inside select on `<-childCtx.Done()`. The demo shows the pattern: an 80ms simulated workload under a 50ms parent deadline returns `context deadline exceeded` because the inner select observes cancellation.

## 5. Exponential Backoff with Full Jitter

`CalculateBackoff` implements the AWS Architecture Blog formula for full jitter:

```go
func (r *Retrier) CalculateBackoff(attempt int) time.Duration {
    if attempt <= 0 {
        return 0
    }
    multiplier := 1 << uint(attempt-1)
    temp := float64(r.cfg.BaseBackoff) * float64(multiplier)
    maxVal := float64(r.cfg.MaxBackoff)
    if temp > maxVal {
        temp = maxVal
    }
    // Full Jitter: random duration in [0, temp]
    sleep := rand.Float64() * temp
    return time.Duration(sleep)
}
```

The computation has three stages. Exponential growth: `BaseBackoff * 2^(attempt-1)` via bit shift. Capping: `min(temp, MaxBackoff)` bounds growth so attempt 10 with 10ms base stays at the configured ceiling instead of reaching 5 seconds. Randomization: multiply by a uniform sample in `[0, 1)` from `math/rand/v2`, yielding a sleep uniformly distributed over `[0, temp]`.

Concrete sequence with base 10ms and max 50ms: attempt 1 draws from `[0, 10ms]`, attempt 2 from `[0, 20ms]`, attempt 3 from `[0, 40ms]`, attempt 4 from `[0, 50ms]` after capping, and every later attempt stays in `[0, 50ms]`. Expected sleep per attempt equals half the capped exponential. Variance is maximal by construction, which is the property that desynchronizes clients.

`Do` drives the loop with two context checks per iteration:

```go
for attempt := 1; attempt <= r.cfg.MaxAttempts; attempt++ {
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }
    err := fn(ctx)
    if err == nil {
        return nil
    }
    lastErr = err
    if attempt == r.cfg.MaxAttempts {
        break
    }
    backoff := r.CalculateBackoff(attempt)
    select {
    case <-ctx.Done():
        return ctx.Err()
    case <-time.After(backoff):
    }
}
return errors.Join(ErrMaxRetriesExceeded, lastErr)
```

The pre-attempt check aborts before invoking a doomed call. The sleep check aborts during backoff. `TestRetrier_ContextCanceled` sets a 20ms context against 50ms base backoff and asserts the return joins context expiry rather than hanging through all 5 attempts. Terminal failure joins `ErrMaxRetriesExceeded` with the last error so callers can match either with `errors.Is`. `TestRetrier_JitterBoundsAndZeroConfig` samples attempts 1 through 5 and asserts every value stays within `[0, MaxBackoff]`.

## 6. Circuit Breaker State Machine

The breaker has three states: `CLOSED`, `OPEN`, `HALF_OPEN`. `CLOSED` passes traffic and counts failures. `OPEN` rejects traffic immediately with `ErrCircuitOpen`. `HALF_OPEN` admits a bounded probe after cooldown; success closes the circuit, failure re-opens it.

```go
type Breaker struct {
    mu           sync.RWMutex
    cfg          Config
    state        State
    failures     int
    successes    int
    lastStateChg time.Time
}
```

All state transitions hold `mu` as a write lock. `State()` and `Allow()` call `checkCooldown()` under that lock before reading, so a goroutine observing the breaker after cooldown sees `HALF_OPEN` without a separate timer goroutine. Lazy evaluation keeps the implementation free of background goroutines and their shutdown logic.

Transition rules in code:

- `RecordFailure` in `CLOSED` increments `failures`. At `failures >= FailureThreshold`, state becomes `OPEN` with counters reset and `lastStateChg` stamped.
- `RecordFailure` in `HALF_OPEN` returns directly to `OPEN` on the first failure. One failed probe proves the dependency remains unhealthy.
- `RecordSuccess` in `HALF_OPEN` increments `successes`. At `successes >= SuccessThreshold`, state becomes `CLOSED`.
- `RecordSuccess` in `CLOSED` resets `failures` to zero. A single success clears the consecutive failure count.
- `Allow` returns `ErrCircuitOpen` while `OPEN` and nil otherwise.

`Execute` composes the guard with accounting:

```go
func (b *Breaker) Execute(fn func() error) error {
    if err := b.Allow(); err != nil {
        return err
    }
    err := fn()
    if err != nil {
        b.RecordFailure()
        return err
    }
    b.RecordSuccess()
    return nil
}
```

The value for retry loops: wrapping each attempt in `cb.Execute` converts repeated downstream failure into fast local rejection. `TestCircuitBreaker_StateTransitions` walks the full cycle with threshold 2, cooldown 50ms: two failures trip `OPEN`, an immediate call gets `ErrCircuitOpen`, sleep 60ms surfaces `HALF_OPEN`, two successes restore `CLOSED`. `TestCircuitBreaker_HalfOpenFailureTripsOpen` confirms one probe failure re-opens immediately.

## 7. Idempotency Key Deduplication

Timeouts create ambiguity. The client gives up waiting while the server may still process the request. A naive retry of `POST /charges` executes the charge twice. Idempotency keys resolve the ambiguity by making the second execution return the stored first result.

The lab store is a TTL map:

```go
type Record struct {
    Response  string
    CreatedAt time.Time
}

type Store struct {
    mu      sync.RWMutex
    records map[string]Record
    ttl     time.Duration
}
```

`Get` performs lazy eviction under a write lock:

```go
func (s *Store) Get(key string) (string, bool) {
    s.mu.Lock()
    defer s.mu.Unlock()

    rec, ok := s.records[key]
    if !ok {
        return "", false
    }
    if time.Since(rec.CreatedAt) > s.ttl {
        delete(s.records, key)
        return "", false
    }
    return rec.Response, true
}
```

The write lock is required because expiry deletes from the map. A read lock would need upgrading, which Go `sync.RWMutex` does not support. `TestStore_LazyEvictionOnGet` sets a key with 20ms TTL, sleeps 30ms, asserts `Get` returns missing, then inspects the map directly and asserts the entry was physically removed. No background sweeper goroutine exists. Memory stays bounded by live keys plus one stale entry per key until its next read.

The usage pattern from the demo:

```go
processPayment := func(key string, amount int) (string, error) {
    if res, ok := store.Get(key); ok {
        return fmt.Sprintf("%s (DEDUPLICATED)", res), nil
    }
    result := fmt.Sprintf("Charged $%d successfully", amount)
    store.Set(key, result)
    return result, nil
}
```

Check before execute, store after execute. Production systems persist this mapping in a database with a unique constraint on the key column, following the Stripe API model. The lab keeps it in memory to isolate the deduplication logic. `TestStore_ConcurrentAccess` runs 50 goroutine pairs of `Set` and `Get` against one key, verifying race safety. The race detector run in section 10 confirms zero data races.

## 8. Integration: Retry Plus Breaker Plus Idempotency

The two integration tests verify composition.

`TestIntegration_RetryWithCircuitBreaker` configures a breaker with failure threshold 2 and cooldown 50ms, plus a retrier with 4 attempts and 5ms base backoff. Every attempt executes `cb.Execute` around a function returning `backend down`. The first two attempts record failures and trip the breaker to `OPEN`. Attempts 3 and 4 receive `ErrCircuitOpen` without touching the backend. The test asserts a non-nil terminal error and `cb.State() == StateOpen`. Total backend load is 2 calls instead of 4. Under production thresholds the savings scale: with threshold 5 and 10 retry attempts, the breaker sheds the majority of doomed calls.

`TestIntegration_IdempotentRetry` configures a store with 1 minute TTL and a retrier with 3 attempts. The `process` closure checks `store.Get(idKey)` first, increments `actualExecutions` only on a miss, stores `order-created`, returns it. The first direct call executes once. The retry loop then invokes `process` up to 3 times. Every retry hits the stored response. The test asserts `actualExecutions == 1`. The counter is the critical assertion: retries happened, side effects did not repeat.

Correct composition order for production handlers: check idempotency key, check breaker admission, execute under context deadline, record breaker outcome, store idempotent result on success, retry transient failures with jittered backoff while the context remains live.

## 9. Database Timeouts: The Missing Layer

Go-level timeouts bound waiting. They leave server-side work unbounded. When `context.WithTimeout` fires on an in-flight PostgreSQL query, the Go `database/sql` driver closes the connection or returns it to the pool, while PostgreSQL continues executing the query until completion unless a server-side limit intervenes. Under pile-up conditions those orphaned queries saturate database CPU and hold row locks, blocking healthy traffic at the storage layer.

Three PostgreSQL settings close the gap. `statement_timeout` aborts any statement running longer than the limit, measured server side from execution start. `lock_timeout` aborts any statement waiting longer than the limit to acquire a table or row lock, preventing queue pile-up behind a held lock. `idle_in_transaction_session_timeout` aborts sessions that hold a transaction open while idle, releasing locks held by application code that forgot to commit or roll back after a client-side timeout.

Set all three per role or per database for defense in depth. Representative values for an OLTP API: `statement_timeout = '5s'`, `lock_timeout = '2s'`, `idle_in_transaction_session_timeout = '30s'`. Long-running analytics or migrations override them per session with `SET LOCAL`. The lab omits a live database to keep the module dependency-free, and this section documents the boundary explicitly so adopters do not mistake client-side budgets for complete protection.

## 10. Test Verification and Results

Full suite with the race detector passes. The run `go test -race -v ./...` exercises five packages: `internal/circuit` with 3 tests covering full cycle, half-open failure, zero config defaults; `internal/deadline` with 3 tests covering success, timeout, parent inheritance; `internal/idempotency` with 3 tests covering get/set/expiry, 50-goroutine concurrent access, lazy eviction map removal; `internal/retry` with 5 tests covering first-try success, retry until success, max attempts exhaustion, context cancellation, jitter bounds; `tests` with 2 integration tests covering breaker-gated retry and deduplicated retry. Every test reports PASS. The race detector reports zero warnings.

Demo output from `go run cmd/demo/main.go` matches the four components:

```
=== Timeouts & Deadlines Laboratory Demo ===

--- Demo 1: Context Deadline & Budget Propagation ---
Deadline propagation result: context deadline exceeded

--- Demo 2: Exponential Backoff with Full Jitter ---
  Attempt #1 executed
  Attempt #2 executed
  Attempt #3 executed
Retry execution result: <nil> (total attempts: 3)

--- Demo 3: Circuit Breaker State Transitions ---
Initial state: CLOSED
State after 2 failures: OPEN
Execution attempt while OPEN: circuit breaker is open
State after cooldown: HALF_OPEN
Execution attempt in HALF_OPEN (success): <nil> -> New State: CLOSED

--- Demo 4: Idempotent Request Retry Protection ---
First request execution: Charged $100 successfully
Retried request execution: Charged $100 successfully (DEDUPLICATED)

=== Demo Completed Successfully ===
```

Demo 1 confirms parent deadline wins over the larger child budget. Demo 2 confirms two transient failures followed by success on attempt 3. Demo 3 walks CLOSED to OPEN to HALF_OPEN to CLOSED with exact threshold counts. Demo 4 confirms the second payment with the same key returns the stored result tagged deduplicated.

## 11. Production Considerations and Lab Boundaries

Three boundaries separate this lab from production deployment.

In-memory idempotency state works for a single process. Multiple replicas need a shared store such as Redis with `SET key value NX EX ttl` or a database table with a unique key column and atomic insert-or-select. Key scope needs a namespace per operation type. TTL needs alignment with the business reconciliation window: too short allows duplicate charges after expiry, too long retains PII beyond necessity.

Single-process `context.Context` propagates deadlines within one binary. Distributed calls need wire propagation: gRPC deadlines over HTTP/2 headers, W3C `traceparent` plus timeout headers over HTTP, per-hop budget subtraction at each service. Without wire propagation, an edge 500ms budget becomes five sequential 500ms downstream calls totaling 2.5s.

Static budgets in this lab use fixed durations. Production systems under variable load benefit from adaptive budgets derived from latency histograms, such as the p99 over a sliding window plus headroom, clamped between a floor and a ceiling. The Google SRE Book documents deadline budgets per RPC with hedged requests as the mature form of this pattern.

Retry safety requires classifying errors before retrying. Retry timeouts and 503s. Never retry 400s. Retry POST only with an idempotency key attached. Cap attempts, cap backoff, jitter every sleep, abort on context cancellation. The lab retrier implements the mechanics. Error classification belongs in the caller.

## 12. Commit Hash and Repo Link

Lab source at commit `bdb285d35a71c9a309e83a0e4497b9dbd842ceb8`.

Repository: https://github.com/lukman-ss/software-engineering-lab/tree/main/labs/28-timeouts-and-deadlines

Key references: Google SRE Book chapter on handling overload, AWS Architecture Blog post on exponential backoff and jitter, gRPC deadline propagation documentation, PostgreSQL 18 documentation for `statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout`, Stripe API idempotency key model.
