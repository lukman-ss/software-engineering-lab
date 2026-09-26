# Key Takeaways

1. **Circuit Breaker does not heal a broken dependency** — it prevents the caller from self-destruction by failing fast and preserving system resources (threads, memory, connections).

2. **Three-state machine** — CLOSED (normal), OPEN (fail-fast), HALF_OPEN (probe) — enables automated recovery without manual intervention.

3. **Fail-fast executes in nanoseconds** — once OPEN, requests return `ErrCircuitOpen` in microseconds vs millisecond HTTP timeouts, demonstrated in demo Scenario 2.

4. **Downstream calls stop when OPEN** — demo verified: only 3 downstream calls made before circuit tripped, subsequent requests never reach network.

5. **Probe mechanism prevents thundering herd** — `HalfOpenMaxCalls` limits concurrent probes during recovery; excess requests fail fast.

6. **Consecutive successes required** — HALF_OPEN requires `HalfOpenMaxCalls` successful probes before transitioning to CLOSED.

6. **Thread-safe by design** — all state mutations protected by `sync.Mutex`; zero data races verified under `go test -race ./...`.

7. **Panic safety** — panics during execution don't corrupt internal state; failure counters updated before re-panicking.

8. **Default configuration is sensible** — FailureThreshold=3, OpenTimeout=5s, HalfOpenMaxCalls=1 work out of the box.

9. **Lab timeouts are illustrative** — 100ms HTTP timeout and 300ms cooldown are for fast testing; production must tune to actual SLA and recovery profiles.

10. **Implementation uses consecutive failure counting only** — not sliding window or error rate; production may need more sophisticated metrics.