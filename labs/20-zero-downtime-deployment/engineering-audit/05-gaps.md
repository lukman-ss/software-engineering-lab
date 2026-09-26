# Gap Analysis

Target Lab: labs/20-zero-downtime-deployment

---

## GAP-01

Type: MISSING_TEST  
Severity: MEDIUM  
Location: internal/worker/worker.go:68-73 + Stop()  
Description: TOCTOU race between `Enqueue` and `Stop`. `Enqueue` checks `stopped.Load()`, passes the check, then `Stop()` calls `close(jobChan)` before the channel send in `Enqueue` executes. Result: send-on-closed-channel panic. No test exercises concurrent `Enqueue` + `Stop`. The race detector did not catch this because no test triggers it.  
Impact: Potential panic in production concurrent usage. Lab-scope usage (sequential enqueue then stop) avoids it.

---

## GAP-02

Type: DOC_CODE_MISMATCH  
Severity: LOW  
Location: engineering/03-execution-result.md  
Description: Execution result records 5 passing tests. Actual test suite at audit time has 14 tests. The engineering revision phase added tests not reflected in the execution result document.  
Impact: Stale documentation. Does not affect correctness.

---

## GAP-03

Type: MISSING_TEST  
Severity: LOW  
Location: tests/worker_test.go  
Description: No test explicitly asserts that `Enqueue` after `Stop` is silently rejected (dropped with log). The behavior exists and is correct, but the contract is unverified by test.  
Impact: Regression risk if the drop behavior changes.

---

## GAP-04

Type: MISSING_TEST  
Severity: LOW  
Location: tests/db_test.go  
Description: No test overwrites a legacy record with `SaveExpand` and re-reads it to verify upgrade path. The expand-in-place migration path is untested.  
Impact: Minor. The dual-write scenario (new version writes to existing legacy ID) is a real production pattern.

---

## GAP-05

Type: MISSING_TEST  
Severity: LOW  
Location: tests/server_test.go  
Description: No test with multiple concurrent in-flight requests during graceful shutdown. Only 1 in-flight request is tested. Multi-request drain correctness relies on `http.Server.Shutdown`'s built-in behavior, which is correct, but not exercised.  
Impact: Low. Go stdlib's Shutdown is well-tested upstream.

---

## GAP-06

Type: MISSING_EDGE_CASE  
Severity: LOW  
Location: internal/worker/worker.go — zero-buffer or zero-concurrency  
Description: No test for `NewWorker(0)` (zero buffer) or `Start(0)` (zero goroutines). Channel send to zero-buffer channel blocks indefinitely if no workers are running. Panic-free but deadlocks.  
Impact: Lab-scope only. Not a real-world concern given explicit usage with buffer=100, concurrency=2.

---

## Summary Table

| ID | Type | Severity | Status |
|----|------|----------|--------|
| GAP-01 | MISSING_TEST (TOCTOU panic risk) | MEDIUM | Open |
| GAP-02 | DOC_CODE_MISMATCH (stale exec result) | LOW | Open |
| GAP-03 | MISSING_TEST (Enqueue rejection) | LOW | Open |
| GAP-04 | MISSING_TEST (legacy overwrite) | LOW | Open |
| GAP-05 | MISSING_TEST (multi-request drain) | LOW | Open |
| GAP-06 | MISSING_EDGE_CASE (zero-buffer/concurrency) | LOW | Open |

No CRITICAL or HIGH gaps. No fabricated results. No fake benchmarks.
