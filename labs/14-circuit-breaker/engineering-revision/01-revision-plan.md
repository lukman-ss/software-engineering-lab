# Engineering Revision Plan

Target Lab: labs/14-circuit-breaker
Previous Verdict: APPROVED

## Blocking Issues
None.

## Non-Blocking Issues
1. Missing edge-case test: verify that a successful request in the CLOSED state resets the consecutive failures counter.
2. Missing edge-case test: verify that when the circuit is in the HALF_OPEN state, requests exceeding `HalfOpenMaxCalls` are rejected immediately with `ErrCircuitOpen`.

## Files To Change
- `internal/circuitbreaker/circuit_breaker_test.go`

## Tests To Add/Modify
- Add `TestSuccessInClosedResetsFailures` to assert that `failures` count resets upon success before hitting the failure threshold.
- Add `TestHalfOpenThrottlesExcessCalls` to assert that `HalfOpenMaxCalls` limit blocks excess probes with `ErrCircuitOpen`.

## Validation Commands
- `cd labs/14-circuit-breaker`
- `go test ./...`
- `go test -race ./...`
- `go run ./cmd/demo`
