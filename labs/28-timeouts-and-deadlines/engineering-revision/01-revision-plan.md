# Engineering Revision Plan

Target Lab: labs/28-timeouts-and-deadlines
Previous Verdict: APPROVED_WITH_WARNINGS

## Blocking Issues
None.

## Non-Blocking Issues
1. Missing lazy eviction in idempotency store on expired read.
2. Missing test cases for circuit breaker half-open failure trip and zero configuration defaults.
3. Missing test cases for retry jitter bounding and zero configuration defaults.
4. Concurrency assertion strength in idempotency store test.

## Files To Change
- `internal/idempotency/idempotency.go`: Add active deletion of expired records on `Get()`.
- `internal/idempotency/idempotency_test.go`: Add `TestStore_LazyEvictionOnGet` and strengthen concurrent assertions.
- `internal/circuit/circuit_test.go`: Add `TestCircuitBreaker_HalfOpenFailureTripsOpen` and `TestCircuitBreaker_DefaultZeroConfig`.
- `internal/retry/retry_test.go`: Add `TestRetrier_JitterBoundsAndZeroConfig`.

## Tests To Add/Modify
- `TestStore_LazyEvictionOnGet` in `internal/idempotency/idempotency_test.go`
- `TestCircuitBreaker_HalfOpenFailureTripsOpen` in `internal/circuit/circuit_test.go`
- `TestCircuitBreaker_DefaultZeroConfig` in `internal/circuit/circuit_test.go`
- `TestRetrier_JitterBoundsAndZeroConfig` in `internal/retry/retry_test.go`

## Validation Commands
```bash
go test -v ./...
go test -count=1 -race ./...
go run ./cmd/demo
```
