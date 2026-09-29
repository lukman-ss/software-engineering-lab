# Gap Analysis

Target Lab: labs/37-cache-invalidation-strategies

## GAP-01: Write-Behind Flush Error Ignored Without Test or Metric

Gap Type: UNHANDLED_ERROR
Severity: MEDIUM
Location: `internal/cache/patterns.go:125, 130`
Description:
In `WriteBehindService.flushWorker`, DB write errors are explicitly discarded with `_ = s.db.Write(...)`. If the backing store is unreachable or returns an error, the write is silently lost with no retry, error callback, dead-letter queue, or metric. While documented as a demonstration in implementation notes, no test exercises this failure path.
Impact: Silent data loss on DB failure.
Remediation: Document the error-handling limitation clearly in code comments or add an error callback/channel to `WriteBehindService` so failures can be observed and tested.

---

## GAP-02: Missing Test for SWR Hard-Miss Path

Gap Type: MISSING_TEST
Severity: LOW
Location: `internal/cache/stampede.go:215-224`
Description:
`SWRService.Get` has three branches:
1. Cache hit (fresh)
2. Stale within stale window (serve stale + async revalidate)
3. Hard miss / expired beyond stale window (synchronous fetch)
Branch 3 is executed on the initial fetch (before key exists), but the condition where an existing key has aged *past* `ExpiresAt.Add(s.staleDelta)` is never explicitly tested.
Impact: Edge case branch coverage incomplete for key expiration exceeding stale window.
Remediation: Add a test case where `time.Sleep` exceeds `ttl + staleDelta` and verify synchronous query occurs and blocks.

---

## GAP-03: Missing Test for XFetch DB Failure Stale Fallback

Gap Type: MISSING_TEST
Severity: LOW
Location: `internal/cache/stampede.go:158-163`
Description:
`XFetchService.Get` contains defensive fallback logic:
```go
if err != nil {
    if ok {
        return item.Value, nil
    }
    return "", err
}
```
If DB recompute fails during proactive early refresh, it returns the stale cached value rather than returning the error to the caller. This behavior is intentional and robust, but no test covers it.
Impact: Unverified fallback behavior during database outage on proactive recompute.
Remediation: Add a unit test injecting a failing DB query into `XFetchService` on a populated key and asserting that stale value is returned with `nil` error.

---

## GAP-04: Design Document Lists Non-Existent `jitter.go`

Gap Type: DOC_CODE_MISMATCH
Severity: LOW
Location: `engineering/01-design.md:46`
Description:
The design document specifies:
```text
internal/cache:
  ...
  - jitter.go: TTL jitter calculation.
```
In reality, `TTLWithJitter` was placed in `internal/cache/store.go`. There is no `jitter.go` file. The README correctly omits `jitter.go`.
Impact: Minor confusion for developers reading the design document.
Remediation: Update `engineering/01-design.md` line 46 to reflect that jitter logic is in `store.go`.

---

## GAP-05: Missing Test for SWR Concurrent Revalidation Deduplication

Gap Type: MISSING_TEST
Severity: LOW
Location: `internal/cache/stampede.go:226-234`
Description:
The SWR implementation includes a deduplication guard `revalidating map[string]bool` to ensure only one background revalidation goroutine runs per key. While the unit test checks that async revalidation happens, it does not issue concurrent reads during the stale window to verify that `revalCount` remains 1 and duplicate revalidations are suppressed.
Impact: Deduplication logic correctness under concurrency is unverified by tests.
Remediation: Add a concurrency test launching multiple goroutines during the stale window and asserting `revalCount` is 1.

---

## Summary

| Gap ID | Gap Type | Severity | Description |
|---|---|---|---|
| GAP-01 | UNHANDLED_ERROR | MEDIUM | DB write errors in WriteBehind flush worker are discarded silently without test or metric |
| GAP-02 | MISSING_TEST | LOW | SWR hard-miss path (past staleDelta) not explicitly tested |
| GAP-03 | MISSING_TEST | LOW | XFetch fallback to stale value on DB error not tested |
| GAP-04 | DOC_CODE_MISMATCH | LOW | `engineering/01-design.md` lists `jitter.go` which lives in `store.go` |
| GAP-05 | MISSING_TEST | LOW | SWR in-flight revalidation deduplication not verified under concurrent reads |
