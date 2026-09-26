# 10 Final Research Summary

## Synthesis Points

### 1. Circuit Breakers isolate network blast radius
A failing downstream dependency can take down callers through thread exhaustion, connection pool depletion, and memory pressure as blocked requests wait for timeouts. The circuit breaker decouples the caller's resources from the downstream's state: after threshold failures, it fails fast in microseconds, freeing caller capacity while giving the dependency time to recover.

### 2. Caller fails in microseconds instead of multi-second HTTP timeout hangs
Azure Circuit Breaker: "This pattern quickly rejects a request for an operation that's likely to fail, rather than waiting for the operation to time out." Google SRE Workbook Dressy case study: synchronous retries on overload caused 20x traffic spikes; fail-fast with CB + jitter prevents replay storms.

### 3. Asynchronous non-critical flows can be decoupled via queues
Queue-Based Load Leveling pattern (Microsoft Azure) enables accepting work into durable queue even when downstream not available; worker processes queue asynchronously. Enables: degraded/fallback behavior when dependency unavailable; idempotent reprocessing; prevents 4xx client errors from tripping CB.

## Architecture Takeaways for Lab
- Use concurrent-failure counter (gobreaker default 5 consecutive) → deterministic for tests.
- HalfOpen probes = 1 (minimize thundering herd).
- OpenTimeout small enough for demo (500ms) vs production (timeout + jitter backoff window).
- HTTP client with explicit timeout; timeout not CB responsibility.
- CB returns sentinel error; callers can detect `ErrCircuitOpen` for metrics/logic.
- No retries inside protected call; CB gates retries externally via fail-fast.

## Confidence Levels
- Three-state machine: HIGH (verified by Azure, Fowler, Hystrix)
- Threshold models: MEDIUM (consecutive count vs rolling window — implementation variance, not contradiction)
- Retry + jitter: HIGH (AWS Architecture Blog, Google SRE)
- Bulkhead vs CB distinction: HIGH (Azure docs explicit)

## NOT VERIFIED Claims
- Adaptive hysteresis (AI-powered thresholds) — mentioned in Azure but not studied; production systems use static for predictability.
- Netflix Hystrix specific parameter tuning — referenced for concepts; lab uses simpler model.
- Exact production threshold values — NOT VERIFIED; sources say "tuned per SLO, no universal".

## Next Steps for Implementation Lab
1. Build `circuitbreaker/` with mutex-protected state machine (CLOSED→OPEN→HALF_OPEN→CLOSED).
2. Build payments/ fake HTTP server with HEALTHY/SLOW/DOWN states (controlled by test/demo flags).
3. Build checkout/ service invoking Payment Client wrapped by CB.
4. Implement cmd/demo/main.go running Scenarios 1-4 with real timing counters.
5. Write tests per quality gate (state transitions, fail-fast, downstream skip).

## Sources Summary
All evidence gathered from Tier 1 sources:
- Microsoft Azure Architecture Center (Circuit Breaker, Retry, Bulkhead, Queue-Based Load Leveling)
- Netflix Hystrix Wiki (How it Works, Configuration, Operations)
- Google SRE (Handling Overload, Managing Load)
- AWS Architecture Blog (Exponential Backoff and Jitter)
- Primary Go libraries (sony/gobreaker, cep21/circuit) as implementation references.