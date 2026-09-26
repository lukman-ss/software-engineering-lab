# Gap Analysis

Allowed gap types enumerated in auditor guidance.

| Gap ID | Type                  | Location                    | Description                                                                 | Severity |
|--------|-----------------------|-----------------------------|-----------------------------------------------------------------------------|----------|
| G1     | MISSING_TEST          | internal/server/server.go   | 10% overload spike branch (`activeReq > MaxDB` → 25x duration) untested.   | MEDIUM   |
| G2     | MISSING_TEST          | internal/loadtest/runner.go | Zero/negative Duration, empty URL, empty Method, VUs<=0 defaulting untested. | LOW      |
| G3     | UNHANDLED_ERROR       | cmd/demo/main.go:85         | JSON encode error discarded (benign but unlogged).                          | LOW      |
| G4     | MISSING_EDGE_CASE     | tests/loadtest_test.go      | ContentType header propagation, exact body/method not verified in tests.    | LOW      |
| G5     | DOC_CODE_MISMATCH     | README.md                   | Does not mention `go vet`, `go build`, existence of engineering-* dirs.     | LOW      |
| G6     | IMPLEMENTATION_OVERCLAIM | engineering/02-implementation-notes.md:20 | Claims exact sorting suitable for test scale; correct but no quantification of scale boundary. | LOW      |
| G7     | WEAK_ASSERTION        | tests/loadtest_test.go:52-56 | `stressRes.P95Latency <= stressRes.AvgLatency` not guaranteed (P95 can be ≤ Avg in some distributions). Flaky if tail not pronounced. | MEDIUM   |

Gap summary:
- MEDIUM: G1, G7
- LOW: G2, G3, G4, G5, G6

No HIGH/CRITICAL gaps: no fake benchmark, no race condition, no broken core behavior, no doc-code contradiction on primary claims, no missing happy-path test.