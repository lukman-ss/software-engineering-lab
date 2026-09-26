# Gap Analysis

## MISSING_TEST
**Description:** Network dial errors (e.g., closed ports, unroutable URLs) increment error counts in `runner.go`, but there is no specific unit or integration test simulating an unreachable network dial. `TestLoadTest_ErrorCount` only tests HTTP 500 status code errors.
**Severity:** LOW
**Impact:** Minimal. The runner's logic handles it correctly, but explicit test coverage for the network-failure edge case is missing.

## UNHANDLED_ERROR (Warning)
**Description:** In `server.go`, the handler unconditionally waits to acquire the connection pool semaphore (`s.semaphore <- struct{}{}`). It does not listen to `r.Context().Done()` while queuing. If a client disconnects, the server still eventually consumes a DB query slot for 20ms.
**Severity:** LOW
**Impact:** Minor in this lab environment since client timeouts (5s) are much higher than max queue wait times (~210ms). However, this represents an unhandled request cancellation propagation in the server mock.

## IMPLEMENTATION_OVERCLAIM (Warning)
**Description:** `runner.go` closes the response body via `_ = resp.Body.Close()` without draining it via `io.Copy(io.Discard, resp.Body)`. In high-throughput load testers, failing to drain the body can sever TCP connections instead of returning them to `http.Transport`'s idle pool.
**Severity:** LOW
**Impact:** The mock server writes tiny payloads that are immediately buffered by the OS networking stack, so the close operation succeeds without noticeable connection churn. For a production load testing tool, this would be a defect.
