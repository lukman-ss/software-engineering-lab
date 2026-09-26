## Finding 1

Type: RACE_CONDITION
Location: `internal/circuitbreaker/circuit_breaker.go:89-116`
Severity: LOW
Description: `Execute` releases `mu` before invoking `fn`, then re-acquires and calls `onFailureLocked(now)` (recording `openedAt`) if `fn` failed. Between unlock and re-lock, another concurrent `Execute` could open the breaker; the delayed goroutine's failure is then ignored (failures in OPEN are not incremented). This matches expected circuit-breaker semantics (do not count fail-fast rejections as new failures). No data corruption. The `-race` detector found no races during `TestConcurrentAccess`.
Status: ACCEPTED — correct by design; documented as known pattern.

## Finding 2

Type: DOC_CODE_MISMATCH
Location: README.md (Expected Behavior section, lines 114-116)
Severity: MEDIUM
Description: README "Expected Behavior" snippet shows `err=payment failed: status 500 duration=...`, but actual `go run ./cmd/demo` output (and `client.go:42`) produces `err=payment failed: status 500 body internal payment server failure\n duration=...`. README omits the `body` segment. Stale/illustrative example not synchronized with implementation.

## Finding 3

Type: MISSING_TEST
Location: `internal/circuitbreaker/circuit_breaker_test.go`
Severity: LOW
Description: `HalfOpenMaxCalls > 1` (allowing multiple concurrent probes) path is untested; all tests use value `1`. The branch at line 97 `if b.halfOpenIn >= b.cfg.HalfOpenMaxCalls` is only ever exercised with threshold 1.

## Finding 4

Type: MISSING_TEST
Location: `circuit_breaker.go:62-73` (New defaults)
Severity: LOW
Description: No test verifies `New` default fallback behavior (FailureThreshold=3, OpenTimeout=300ms, HalfOpenMaxCalls=1) when callers pass zero-valued Config. Test suite always passes explicit config.

## Finding 5

Type: UNVERIFIED_RESULT
Location: README.md "Expected Behavior"
Severity: LOW
Description: README prints exact microsecond timings (e.g., 366.75µs). Actual runs differ (464µs/164µs/129µs). README line 63 disclaimer says timeouts are illustrative, so this is tolerated, but embedding exact numbers without a "(example, not literal)" note risks misleading readers.

## Finding 6

Type: DOC_CODE_MISMATCH (minor — formatting)
Location: `internal/circuitbreaker/circuit_breaker.go`, `cmd/demo/main.go`
Severity: LOW
Description: `gofmt -l` flags these two files as not gofmt-formatted. Purely cosmetic (import ordering, constant alignment). Does not affect correctness. Not a functional mismatch but flagged for completeness.

## Summary

| Gap | Severity | Blocking? |
|---|---|---|
| Race condition (by-design) | LOW | No |
| README error format stale | MEDIUM | No |
| HalfOpenMaxCalls=1 untested path | LOW | No |
| New() defaults untested | LOW | No |
| README exact timings | LOW | No |
| gofmt formatting | LOW | No |

No HIGH or CRITICAL findings.
