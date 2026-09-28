# Test Audit

Target Lab: labs/26-contract-testing
Date: 2026-09-28

## Commands Executed (fresh, uncached)

- `go test -v -count=1 ./...` → PASS (5/5)
- `go test -race -count=1 ./...` → PASS (no data races)
- `go vet ./...` → clean (no output)
- `go build ./...` → clean
- `go run ./cmd/demo` → PASS, exit 0, all 4 stages

## Coverage Matrix

| Area | Covered | Test |
|---|---|---|
| Contract generation (names, 1 interaction, path) | Yes | TestConsumerContractGeneration |
| Happy path verifier (V1 pass) | Yes | TestProviderV1_ContractVerification_Success |
| Happy path client e2e (V1 values) | Yes | TestProviderV1 (FetchOrder asserts name/status/total) |
| Failure path verifier (breaking fails, ≥3 errors) | Yes | TestProviderBreaking_ContractVerification_Fails |
| Failure path client (breaking → non-nil err) | Yes | TestProviderBreaking (FetchOrder err asserted) |
| Dual V1 compatibility (verifier + client) | Yes | TestProviderDual_ContractVerification_Success |
| Concurrency (20 goroutines, shared Verifier) | Yes | TestConcurrentContractVerification |
| Status-code mismatch (e.g. 404/405 vs 200) | No | MISSING_TEST (LOW) |
| Non-JSON body → verification error | No | MISSING_TEST (LOW) |
| Dual `/v2` route shape assertion | No | MISSING_TEST (LOW) |
| Response-header enforcement | No | Not implemented (see 02-code-audit Finding 2) |
| Timeout / cancellation | No | Not implemented (see Finding 3) |

## Strengths

- Failure test asserts both verifier rejection AND real client breakage — proves the breaking change is consumer-visible, not just a diff-engine artifact.
- Dual test proves V1 compatibility survives alongside V2 routing (no accidental V2-on-V1 serve).
- Concurrent test shares one `Verifier` across 20 goroutines; `-race` clean → `Verify` is stateless/goroutine-safe (fresh `VerificationResult` per call, no shared mutation).
- `t.Errorf` from goroutines is safe (`testing.T` is goroutine-safe for Errorf); `wg.Wait` before return — correct pattern.

## Weaknesses (all LOW, non-blocking)

1. Breaking test asserts `len(errors) >= 3`, not which fields — a verifier emitting 3 wrong diffs could pass. Mitigated: live demo shows the 3 diffs are exactly status/customer.name/total.
2. No negative test for wrong-status or malformed-body providers; verifier branches for those (verifier.go:81-102) are unproven by tests but simple.
3. `/v2` handler output never asserted — routing bug serving V1 payload on `/v2` would go undetected (V1-path tests would still pass).

## Assessment

Suite is small (5 tests, 1 file) but precisely targeted: every core claim (generate → V1 pass → breaking fail → dual pass → concurrency-safe) has a passing, race-clean test backed by live reproduction. Passing suite is not vacuous — the breaking test would fail if the verifier missed any of the 3 deltas (it requires ≥3 errors AND client failure).

Result: PASS with 3 LOW MISSING_TEST notes.
