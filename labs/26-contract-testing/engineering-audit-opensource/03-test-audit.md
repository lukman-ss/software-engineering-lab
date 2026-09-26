# Test Audit

Test file: `tests/contract_test.go` (package `tests`).

## Coverage matrix

| Scenario | Test | Result |
|---|---|---|
| Contract generation + shape | `TestConsumerContractGeneration` | PASS |
| Provider V1 satisfies consumer contract (happy path) | `TestProviderV1_ContractVerification_Success` | PASS |
| Provider V1 end-to-end via mobile client | (same test, second half) | PASS |
| Breaking provider is detected/blocked (failure path) | `TestProviderBreaking_ContractVerification_Fails` | PASS |
| Breaking schema propagates client error | (same test, second half) | PASS |
| Dual provider V1 path stays compatible | `TestProviderDual_ContractVerification_Success` | PASS |
| Concurrent verification is safe | `TestConcurrentContractVerification` | PASS (race detector clean) |

## Required execution results

```
Command: go test -v -count=1 ./...
Result: ok  labs/26-contract-testing/tests   PASS  (5 tests)

Command: go test -race -count=1 ./...
Result: ok  labs/26-contract-testing/tests   1.556s   (no races)

Command: go build ./...  /  go vet ./...
Result: clean, exit 0
```

## Findings

### Finding T1 — Happy path proven
Location: `tests/contract_test.go:27` (`TestProviderV1_ContractVerification_Success`)
Claim: Provider V1 satisfies the MobileApp contract and the mobile client decodes correctly.
Evidence: Test asserts `result.Passed == true` (no errors) and that the parsed `MobileOrderSummary` has the expected `IN_PROGRESS` status / `Budi Santoso` name / `150000` total. Verified at run time.
Assessment: PASS

### Finding T2 — Failure path proven with 3 distinct breakages
Location: `tests/contract_test.go:50` (`TestProviderBreaking_ContractVerification_Fails`)
Claim: Breaking provider is detected with >= 3 errors and the consumer client errors out.
Evidence: Test asserts `!result.Passed`, `len(result.Errors) >= 3`, and `client.FetchOrder` returns an error. Runtime diff confirms three distinct error classes (missing field, type mismatch, value mismatch) matching the three documented breaking changes.
Assessment: PASS

### Finding T3 — Evolutionary compatibility proven
Location: `tests/contract_test.go:74` (`TestProviderDual_ContractVerification_Success`)
Claim: Dual provider keeps V1 contract compatible while the consumer client succeeds on the V1 route.
Evidence: Test asserts `result.Passed == true` and the client returns expected V1 semantics.
Assessment: PASS

### Finding T4 — Concurrency safety tested
Location: `tests/contract_test.go:96` (`TestConcurrentContractVerification`)
Claim: Verifier is safe under concurrent parallel verification.
Evidence: 20 goroutines call `verifier.Verify` against a shared `httptest` server concurrently; all pass. `-race` reports no data races.
Assessment: PASS

### Finding T5 — No test for the declared ResponseDefinition.Headers schema field
Location: `tests/contract_test.go` (absence)
Claim: None in README or tests explicitly requires response-header validation, but `contract.ResponseDefinition.Headers` is defined and populated in `GenerateMobileContract`.
Evidence: No test asserts response-header behavior; the field is never read by `Verify`.
Assessment: WARNING
Severity: LOW
Notes: Latent gap; the schema declares a field the engine ignores. Not a failed test or fabricated result.

### Finding T6 — V2 provider route is untested
Location: `tests/contract_test.go` (absence)
Claim: README claims "Safe API Evolution: Preserving V1 contract compatibility while exposing V2 schemas."
Evidence: No test exercises `/v2/orders/*`. V2 correctness is only asserted via the demo narrative, not by an automated test.
Assessment: WARNING
Severity: LOW
Notes: V2 is a provider-side schema extension, not a consumer contract obligation of the MobileApp consumer, so strict CDC coverage is arguably out of scope. Flagged because the README lists it as an explicit lab objective.

## Overall test-suite strength
The suite covers happy path, failure path, and concurrency with meaningful assertions (not empty). It is sufficient to prove the core CDC behavior claimed. Two minor coverage gaps (T5, T6) are noted.
