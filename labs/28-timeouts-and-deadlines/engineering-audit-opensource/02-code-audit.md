# Code Audit

## Finding 1
Location: `internal/deadline/deadline.go:3-28`
Claimed Behavior: `ExecuteWithBudget` runs `fn` within the smaller of parent deadline or `budget`, propagating parent cancellation.
Observed Implementation: Uses `context.WithTimeout(ctx, budget)` then selects on `childCtx.Done()` vs `fn` completion.
Assessment: PASS
Severity: LOW
Notes: Parent deadline propagation works since child inherits parent. However the goroutine executing `fn` is never signalled to stop when the child deadline fires — the func is expected to watch `childCtx.Done()` itself. Tests rely on this. Goroutine leaks if `fn` ignores the context (mitigated because test functions and demo respect it). No `sync.WaitGroup`/cleanup; acceptable for bounded short calls but worth flagging as latent leak risk.

## Finding 2
Location: `internal/retry/retry.go:35-48` (`CalculateBackoff`)
Claimed Behavior: Exponential backoff with full jitter (AWS architecture blog pattern).
Observed Implementation: `multiplier := 1 << uint(attempt-1)`; `temp = baseBackoff * multiplier` clamped to `MaxBackoff`; `sleep = rand.Float64() * temp` in `[0,temp]` = full jitter.
Assessment: PASS
Severity: LOW
Notes: Full jitter formula correct (`random[0, temp]`). Edge: `attempt` overflow on `1 << uint(attempt-1)` for attempt ≥ 64 but capped by clamp; unreachable for realistic configs. `attempt <= 0` returns 0 — fine because loop starts at 1.

## Finding 3
Location: `internal/retry/retry.go:50-77` (`Do`)
Claimed Behavior: Retries up to `MaxAttempts`, checking context before each attempt and during backoff sleep.
Observed Implementation: `select` on `ctx.Done()` pre-attempt and pre-backoff-sleep; returns joined `ErrMaxRetriesExceeded` + last error on exhaustion.
Assessment: PASS
Severity: LOW
Notes: Last-attempt path skips backoff (correct). On `MaxAttempts`-th failure returns joined error; tests assert `errors.Is(err, ErrMaxRetriesExceeded)` — true because `errors.Join` preserves it.

## Finding 4
Location: `internal/circuit/circuit.go:64-139`
Claimed Behavior: State machine CLOSED→OPEN→HALF_OPEN→CLOSED/OPEN with failure/success counters and cooldown.
Observed Implementation: `checkCooldown` flips `OPEN→HALF_OPEN` after `Cooldown`; `HALF_OPEN` flips to `CLOSED` after `SuccessThreshold` successes or back to `OPEN` on any failure; `CLOSED` trips to `OPEN` at `FailureThreshold`.
Assessment: PASS (with caveats)
Severity: MEDIUM
Notes: `State()` takes write lock and mutates — works but surprising naming. More substantive: `RecordSuccess` in `CLOSED` resets `failures=0`, discarding accumulated-but-sub-threshold failures. Behaviorally fine for the documented semantics but means partial failure counts are lost. Not a correctness bug for current tests. Test verifies exactly the 2-failure→OPEN→cooldown→HALF_OPEN→2-success→CLOSED transition.

## Finding 5
Location: `internal/idempotency/idempotency.go:29-51`
Claimed Behavior: Get/Set deduplication with TTL eviction.
Observed Implementation: Lazy TTL eviction in `Get`; coarse `RWMutex`; `Set` overwrites `CreatedAt`.
Assessment: PASS
Severity: LOW
Notes: TTL checked in `Get` only (no background sweeper). Acceptable for demo. Concurrent test exercises 100 goroutines without races reported by `-race`.

## Finding 6
Location: `cmd/demo/main.go:23-30` (Demo 1)
Claimed Behavior: Parent context (50ms) is tighter than budget (100ms); parent deadline wins.
Observed Implementation: `ExecuteWithBudget(parentCtx{50ms}, 100ms, fn{80ms sleep})`.
Assessment: PASS — `childCtx` derived from parent with smaller budget, returns `context.DeadlineExceeded`.
Severity: LOW

## Finding 7
Location: `cmd/demo/main.go:35-49` (Demo 2)
Claimed Behavior: Retry recovers after transient failures (3 attempts succeed).
Observed Implementation: Returns `nil`, 3 attempts. Matches execution result.
Assessment: PASS
Severity: LOW

## Finding 8
Location: `cmd/demo/main.go:52-73` (Demo 3)
Claimed Behavior: Circuit opens after 2 failures, blocks while OPEN, moves HALF_OPEN, closes on success.
Observed Implementation: Matches test; output reflects states correctly.
Assessment: PASS
Severity: LOW

## Finding 9
Location: `cmd/demo/main.go:76-94` (Demo 4)
Claimed Behavior: First payment charged; retried request deduplicated.
Observed Implementation: `processPayment` checks store, returns dedup string on hit. Output matches.
Assessment: PASS
Severity: LOW
