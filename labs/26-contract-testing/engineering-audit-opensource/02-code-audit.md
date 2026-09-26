# Code Audit

Target Lab: labs/26-contract-testing

## Finding 1

Location: internal/contract/verifier.go:58-115 (`Verifier.Verify`)
Claimed Behavior: Provider CI gate verifies status + body subset against consumer contract.
Observed Implementation: Builds request from interaction, checks status code, decodes body with `UseNumber`, runs `diffValues` subset check, aggregates errors, sets `Passed=false` on any diff. Extra provider fields ignored.
Assessment: PASS
Severity: LOW
Notes: Core CDC semantics proven. Live run: V1 PASS, breaking FAIL with 3 diffs, dual V1 PASS.

## Finding 2

Location: internal/contract/verifier.go:117-176 (`diffValues`)
Claimed Behavior: Detects enum casing, field rename, int->string mutations.
Observed Implementation: Recursive map walk; missing key -> error; `json.Number` vs `string` -> type mismatch; same-type value compare -> value mismatch. Correctly yields 3 diffs on breaking provider.
Assessment: PASS
Severity: LOW
Notes: Minimal subset rule works: provider `notes`, `customer.id` ignored, V1 still passes.

## Finding 3

Location: internal/contract/verifier.go:117-176 (slice handling)
Claimed Behavior: General JSON comparison (implicit).
Observed Implementation: No `[]interface{}` recursion. Two slices with same type fall to `expected != actual` on interface values holding slices -> runtime panic (uncomparable type).
Assessment: WARNING
Severity: LOW
Notes: Unreachable in current scope (contract has no arrays). Known GAP-01. Scoped correctly per revision plan; no refactor needed unless arrays added.

## Finding 4

Location: internal/contract/verifier.go:58-115 (headers)
Claimed Behavior: Engineering 01-design/components claim verification of "headers".
Observed Implementation: Request headers set; response headers (`Content-Type: application/json` in contract) never checked.
Assessment: WARNING
Severity: LOW
Notes: Header contract unenforced. No test covers it. Out of lab's core claims (status/body diffs).

## Finding 5

Location: internal/contract/verifier.go:51-55 (`NewVerifier`)
Claimed Behavior: Executes contracts against running HTTP server.
Observed Implementation: `&http.Client{}` with no `Timeout`. Hang on unresponsive provider blocks gate forever.
Assessment: WARNING
Severity: LOW
Notes: Safe in lab scope (httptest servers). Production gate would need timeout. No timeout test.

## Finding 6

Location: internal/contract/verifier.go:58-115 + tests/contract_test.go:96-115 (concurrency)
Claimed Behavior: Concurrent verification safe.
Observed Implementation: No shared mutable state; `http.Client.Do` concurrent-safe; `diffValues` pure. 20-goroutine shared-verifier run passes under `-race` (1.184s).
Assessment: PASS
Severity: LOW
Notes: Race detector clean on live run.

## Finding 7

Location: internal/provider/server.go (ProviderV1/Breaking/Dual)
Claimed Behavior: V1 passes, breaking fails, dual preserves V1 while exposing V2.
Observed Implementation: Prefix routing `/v1/orders/`, `/v2/orders/`; correct DTOs per model; method guard + `NotFound` fallback. Dual V1 path byte-identical semantics to ProviderV1 (minus `Notes` field, which subset rule ignores).
Assessment: PASS
Severity: LOW
Notes: V2 endpoint (`/v2/orders/`) exists but no test verifies it (see Finding 9).

## Finding 8

Location: internal/consumer/client.go:34-79 (`FetchOrder`)
Claimed Behavior: Mobile client parses V1, fails on breaking schema.
Observed Implementation: Strict unmarshal; rejects non-200, missing `customer.name`, unknown status (`in_progress` rejected by allowlist). Test asserts failure on breaking provider.
Assessment: PASS
Severity: LOW
Notes: Allowlist (`IN_PROGRESS`/`COMPLETED`) brittle but matches contract scope. Error propagation typed with `%w`/context.

## Finding 9

Location: engineering/01-design.md:18 vs code/tests
Claimed Behavior: "V1 contract passes against /v1, while V2 contract passes against /v2".
Observed Implementation: No V2 contract exists; `GenerateMobileContract` only builds V1; no test hits `/v2/orders/`. Dual test only re-verifies V1 path.
Assessment: WARNING
Severity: LOW
Notes: V2 handler code present and simple; evolutionary claim partially unproven. Downgraded to LOW because core claim (V1 preserved) is proven.

## Finding 10

Location: cmd/demo/main.go
Claimed Behavior: 4-stage lifecycle demo with CI gate messaging, `os.Exit(1)` on unexpected outcome.
Observed Implementation: Real `httptest` servers, real `Verify` calls, correct branch logic. Live output matches engineering/03-execution-result.md (same 3 diffs; order varies by map iteration).
Assessment: PASS
Severity: LOW
Notes: Demo output real, no fabrication. Diff order nondeterministic (map range) but tests assert count only, so stable.
