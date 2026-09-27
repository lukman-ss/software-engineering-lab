# Test Audit

Target Lab: `labs/26-contract-testing`

Test file reviewed: `tests/contract_test.go` (5 tests, 115 lines)

## Coverage Matrix vs Engineering Design (`engineering/01-design.md` lines 77-81)

| Design Test | Status |
|---|---|
| `TestConsumerContractGeneration` | Present, PASS |
| `TestProviderV1_ContractVerification_Success` | Present, PASS |
| `TestProviderBreaking_ContractVerification_Fails` | Present, PASS |
| `TestProviderDual_ContractVerification_Success` | Present, PASS |
| `TestConcurrentVerification` | Actual name: `TestConcurrentContractVerification`, PASS |

Design doc (line 81 + line 91 success criterion #6) calls the concurrency test `TestConcurrentVerification`. The implemented name is `TestConcurrentContractVerification`. Same coverage, name mismatch only.

## Happiness Path

Location: `TestProviderV1_ContractVerification_Success`
Assessment: PASS
Notes: Spins up `httptest.Server(provider.NewProviderV1())`, generates contract, verifies, asserts `Passed==true`. Additionally exercises `MobileOrderClient.FetchOrder` and asserts parsed values (CustomerName=Budi Santoso, Status=IN_PROGRESS, Total=150000). End-to-end coverage of V1 path.

## Failure / Breaking Path

Location: `TestProviderBreaking_ContractVerification_Fails`
Assessment: PASS
Notes: Against Breaking provider, asserts `Passed==false` and `len(result.Errors) >= 3`. Exactly the three documented breaking mutations are detected (enum casing, missing field, type mismatch). Also asserts the runtime `MobileOrderClient` returns an error. Good negative-case coverage.

## Evolution Path (Dual V1+V2)

Location: `TestProviderDual_ContractVerification_Success`
Assessment: PASS
Notes: Verifies V1 contract still passes against `/v1/orders` on dual provider, and that `FetchOrder` parses correctly. V2 endpoint (`/v2/orders`) is exercised by the demo (Stage 4) but NOT by an automated test (see GAP below).

## Concurrency Safety

Location: `TestConcurrentContractVerification`
Assessment: PASS
Notes: 20 goroutines share one `Verifier` + one `httptest.Server`, each calling `verifier.Verify`. Passes. Race detector (`go test -race ./...`) run separately and PASS (see below).

### Execution: `go test -v ./...`
Result: PASS — all 5 tests PASS, exit 0.
```
=== RUN   TestConsumerContractGeneration          --- PASS
=== RUN   TestProviderV1_ContractVerification_Success  --- PASS
=== RUN   TestProviderBreaking_ContractVerification_Fails --- PASS
=== RUN   TestProviderDual_ContractVerification_Success --- PASS
=== RUN   TestConcurrentContractVerification        --- PASS (0.01s)
PASS  ok  labs/26-contract-testing/tests  (cached)
```

### Execution: `go test -race ./...`
Result: PASS — exit 0.
```
ok  	labs/26-contract-testing/tests  1.425s
```
(Ran fresh, no cache: `go test -race -count=1 ./...` also used; race detector adds no warnings.)

## Gap Analysis (tests vs claims)

### GAP 1 (MISSING_TEST): No automated assertion that V2 endpoint `/v2/orders` is contract-compliant
- Claim (design §4): "V2 contract passes against /v2/orders/{id}" and "introducing V2 DTO while maintaining V1 compatibility."
- The `TestProviderDual_*` test only asserts the V1 route on Dual provider. The V2 route is only demonstrated by the demo CLI (Stage 4 prints PASSED but does not assert). No test constructs a V2-consumer contract and verifies `/v2/orders`.
- Severity: MEDIUM (V2 safe-evolution claim not covered by unit test; only runtime demo print)

### GAP 2 (MISSING_EDGE_CASE): No assertion of *which exact* 3 breaking diffs are produced
- `TestProviderBreaking_*` only asserts `len >= 3`, not the specific diffs (enum / missing field / type). The demo prints them, but no test pins the three named categories. A future regression could surface 3 different errors and still pass.
- Severity: LOW

### GAP 3 (MISSING_EDGE_CASE): ProviderStates in contract are not driven
- Contract interaction declares `provider_state: "Order ORD-123 exists and is IN_PROGRESS"`, but no provider-state setup machinery honors it. Documented as a known limitation (`implementation-notes.md` line 32) — acceptable, but means the demo is a happy-path-only setup, not true CDC provider-state replay.
- Severity: LOW

### GAP 4 (MISSING_EDGE_CASE): No negative test for non-GET / 404 / unexpected status path
- `ProviderV1.ServeHTTP` returns 405 for non-GET and 404 for unknown paths. No test asserts these paths.
- Severity: LOW

## Verdict on Test Suite Strength

Suite covers: happy path (V1 pass), negative path (breaking fail with ≥3 errors), success-evolution (dual V1 pass), concurrency (race clean). Strength adequate for the lab's scope. Weaknesses are the missing V2-route automated assertion (MEDIUM) and lack of precise breaking-diff pinning (LOW). A passing suite does NOT over-claim: the test asserting `>= 3` is loose but the three changes are real (verified by demo output).

Overall: PASS with caveats. Test claim mismatch is limited to naming (`TestConcurrentVerification` vs `TestConcurrentContractVerification`), not to coverage.
