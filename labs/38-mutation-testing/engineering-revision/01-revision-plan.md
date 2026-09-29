# Engineering Revision Plan

Target Lab: labs/38-mutation-testing
Previous Verdict: APPROVED (with non-blocking warnings and missing test gaps)

## Blocking Issues
None.

## Non-Blocking Issues
1. Unused `mu sync.Mutex` in `engine.Runner` struct (`internal/engine/runner.go:16`).
2. Missing test coverage for degenerate/negative inputs (`TotalAmount: 0`, `ItemCount: 0`) and below-boundary values (`499.99`, `99.99`).

## Files To Change
- `internal/engine/runner.go`
- `internal/service/discount_strong_test.go`

## Tests To Add/Modify
- `internal/service/discount_strong_test.go`: Added test cases for zero amounts, zero items, below-boundary amounts (`499.99`, `99.99`), and zero-item coupon flags.

## Validation Commands
```bash
go test -v ./...
go test -race ./...
go run ./cmd/demo
```
