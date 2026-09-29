# Engineering Gap Analysis

Allowed Gap Types:
- `MISSING_TEST`
- `BROKEN_IMPLEMENTATION`
- `DOC_CODE_MISMATCH`
- `RACE_CONDITION`
- `UNHANDLED_ERROR`
- `MISSING_EDGE_CASE`
- `IMPLEMENTATION_OVERCLAIM`
- `RESEARCH_MISMATCH`
- `FAKE_DEMO`
- `FAKE_BENCHMARK`
- `UNVERIFIED_RESULT`

---

## Gap 1

- **Type**: `MISSING_TEST`
- **Location**: `tests/cache_test.go`
- **Severity**: LOW
- **Description**: Error paths (such as `MockDB` returning an error or `ErrNotFound` during Cache-Aside / Write-Through operations) are not exercised in the unit test suite.
- **Impact**: Code handles errors properly (`if err != nil { return "", err }`), but there is no regression test verifying error propagation.
- **Action Required**: Add negative test cases simulating database failures.

---

## Gap 2

- **Type**: `MISSING_TEST`
- **Location**: `tests/cache_test.go`
- **Severity**: LOW
- **Description**: `XFetchService.Get` end-to-end integration is demonstrated in `cmd/demo/main.go`, but `tests/cache_test.go` tests only the calculation function `ShouldRecompute`.
- **Impact**: Formula correctness is proven; service integration logic is demonstrated at runtime in demo, but lacks an automated unit test with mock time.
- **Action Required**: Add an integration unit test for `XFetchService.Get` using `SetRandFunc`.

---

## Gap 3

- **Type**: `MISSING_EDGE_CASE`
- **Location**: `internal/cache/patterns.go:158`
- **Severity**: LOW
- **Description**: `WriteBehindService.Update` silently drops writes when `writeQueue` is full.
- **Impact**: In high-load scenarios exceeding the buffer size, writes to DB will be dropped without error notification. Explicitly labeled in code as intentional demo design (`ponytail:`).
- **Action Required**: Acceptable for educational lab; document as non-production tradeoff.

---

## Gap 4

- **Type**: `MISSING_EDGE_CASE`
- **Location**: `internal/cache/stampede.go:180`
- **Severity**: LOW
- **Description**: `SWRService` does not have a graceful shutdown (`Close()`) or context cancellation method to wait for in-flight background revalidations.
- **Impact**: Background goroutines triggered during SWR revalidation cannot be cleanly drained or awaited by callers.
- **Action Required**: Acceptable for lab scope; add `Close()` with wait group if productionized.
