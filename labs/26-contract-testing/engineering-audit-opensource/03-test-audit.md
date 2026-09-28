# Test Audit

Target Lab: labs/26-contract-testing

Coverage Matrix vs engineering/01-design.md Test Strategy:

| Test Name | Claimed | Actual | Status |
|---|---|---|---|
| TestConsumerContractGeneration | Asserts contract shape | asserts consumer/provider names + interaction path | PASS |
| TestProviderV1_ContractVerification_Success | V1 satisfies contract | Verify pass + end-to-end client parse of all fields | PASS |
| TestProviderBreaking_ContractVerification_Fails | Breaking fails w/ diffs | Verify fail + ≥3 diff errors + client errors | PASS |
| TestProviderDual_ContractVerification_Success | Dual maintains V1 | Verify pass + client parse V1 fields | PASS |
| TestConcurrentContractVerification | concurrent safe | 20 goroutines share verifier+httptest server | PASS |

## Finding 1

Location: tests/contract_test.go
Claimed Behavior: Tests cover happy path, failure path, edge cases, transitions, recovery, rollback, concurrency.
Observed Implementation: Happy (V1), failure (Breaking), concurrency covered. No negative-path assertions on client: client contract-violation branches (`customer.name is missing`, `unknown status %q`) are unreachable under breaking provider because `total` int vs string type-mismatch fails `json.Unmarshal` first — client error path tested only generically via `err == nil` check.
Assessment: WARNING
Severity: MEDIUM
Notes: No test isolates client enum-missing or missing-field (non-type) violations. Recovery/rollback: no provider state to recover; out of scope. Edge cases: empty order ID, 405 method, 404 path — untested.

## Finding 2

Location: tests/contract_test.go:104-109 (`TestConcurrentContractVerification`)
Claimed Behavior: `go test -race ./...` clean for concurrent contract verification.
Observed Implementation: Spawns 20 goroutines sharing one `*contract.Verifier` against a shared `*httptest.Server`. `Verify` mutates only the local `result`; per-request `http.Client` default is thread-safe. Race detector passed (exit 0).
Assessment: PASS
Severity: LOW
Notes: Concurrency claim proven.

## Finding 3

Location: tests/contract_test.go — missing suite
Claimed Behavior: Dual provider V2 path verified.
Observed Implementation: `TestProviderDual` only hits `/v1/orders/ORD-123`. V2 endpoint (`/v2/orders/...`, served by ProviderDual) is implemented but never asserted. No test that `/v2/orders/ORD-123` returns V2 DTO shape.
Assessment: WARNING
Severity: MEDIUM
Notes: Safe-API-evolution success criterion is partially tested (V1 compat proven; V2 endpoint not).

## Finding 4

Location: tests/contract_test.go
Claimed Behavior: Tests verify headers in interactions.
Observed Implementation: Contract `Response.Headers` (`Content-Type`) set but `Verify` ignores them. No test asserts response Content-Type. Request header setting exercised but response headers never validated.
Assessment: WARNING
Severity: MEDIUM
Notes: Design claims header comparison (01-design.md Components §4) but neither code nor tests implement it.
