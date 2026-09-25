# Implementation Notes

## Files Added
- `cmd/demo/main.go`
- `internal/circuitbreaker/circuit_breaker.go`
- `internal/circuitbreaker/circuit_breaker_test.go`
- `internal/checkout/service.go`
- `internal/payment/client.go`
- `internal/payment/fake_server.go`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`
- `engineering/03-execution-result.md`

## Core Design Decisions
- `CircuitBreaker` manages states via internal `sync.Mutex` protection.
- State checks automatically calculate elapsed time on demand when checked or executed (`checkStateTransitionLocked`).
- Testability achieved via clock injection (`cb.now`).

## Implementation-Specific Choices
- Kept configuration minimal: `FailureThreshold`, `OpenTimeout`, `HalfOpenMaxCalls`.
- Demo uses 100ms HTTP timeout and 300ms circuit cooldown for quick execution.

## Known Limitations
- Implements simple consecutive failure counting rather than a sliding time-window or ring-buffer error rate calculation.
- Concurrency during `HALF_OPEN` permits up to `HalfOpenMaxCalls` probes; excess requests fail fast.

## Trade-offs
- Mutex contention on high throughput: Chosen over complex atomic ring-buffers to keep the implementation clear, correct, and self-contained.

## What Is Demonstrated
- Cascading delay mitigation.
- Instant fail-fast behavior during circuit OPEN.
- Probe testing in HALF_OPEN.
- Automatic recovery to CLOSED when downstream heals.
- Circuit re-tripping to OPEN when downstream remains degraded.

## What Is Not Demonstrated
- Dynamic failure rate percentage thresholds.
- Sliding window metrics.
- External distributed state synchronization (e.g., Redis-backed breakers).
