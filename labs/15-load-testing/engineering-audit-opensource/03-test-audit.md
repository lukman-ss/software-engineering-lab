# Test Audit

## Coverage matrix

| Package        | Tests                                  | Purpose                                    | Result |
|----------------|----------------------------------------|--------------------------------------------|--------|
| loadtest (unit)| TestCalculateMetrics                   | Min/Max/ Avg, P50, P95, P99 exactness      | PASS   |
|                | TestCalculateMetrics_Empty             | nil latencies → zero Result                | PASS   |
|                | TestCalculateMetrics_Invariants        | P50<=P90<=P95<=P99<=Max ordering           | PASS   |
| tests (integ)  | TestLoadTest_SmokeVsStress             | 10x pool limit causes P95 tail spike       | PASS   |
|                | TestLoadTest_ErrorCount                | 500 responses → all errors, 0 success     | PASS   |
|                | TestServer_MethodNotAllowed            | GET /booking → 405                         | PASS   |
|                | TestLoadTest_DialError                 | refused conn → all errors                 | PASS   |
|                | TestServer_ContextCanceled             | pre-canceled ctx → no 201                  | PASS   |

Verified results (count=1, real run):
- `go test -v -count=1 ./...` → 8/8 PASS
- `go test -race -count=1 ./...` → PASS, 0 races
- `go vet ./...` → clean
- `go build ./...` → SUCCESS

## Covered scenarios

- Happy path: smoke load produces 186 requests, 0 errors (live). Unit metrics exact match.
- Failure path: HTTP 500 (all errors), connection refused (all errors) — both assert `ErrorCount == TotalRequests`.
- Negative case: method-not-allowed returns 405.
- Edge case: empty latencies slice returns zero Result without panic.
- Transition / recovery: canceled context aborts handler (no 201 written, no leak).
- Invariant: latency percentiles strictly ordered by value.

## Weakness / gaps

1. MISSING_TEST — Server overload spike branch (`activeReq > MaxDB` → 10% chance of 25x duration) has NO test forcing > MaxDB concurrency; the slow path is undocumented and unverified.
2. MISSING_TEST — Zero/negative `Duration`, empty/blank `URL`, empty `Method`, `VUs<=0` defaulting: not exercised.
3. RACE_DESIGN — Smoke test asserts `P95 > Avg` and `stress P95 > smoke P95`. With only 1 VU / 500ms these are timing-sensitive; on a contended host the stress tail could collapse below smoke, flipping the comparison. Passes now; latent flakiness.
4. UNVERIFIED_RESULT — `engineering/03-execution-result.md` "sample run" figures not byte-for-byte reproducible (host-dependent). Demo note acknowledges; acceptable, but the doc should avoid implying fixed values.
5. MISSING_EDGE_CASE — No test for `ContentType` header propagation; no test that load generator sends the configured `Body`/`Method` (errors/dial tests use GET).

Assessment: Test suite proves stated behavior but has medium coverage gaps in server slow-path and config-validation edge cases. No fake results; live run reproduced documented output shape.
