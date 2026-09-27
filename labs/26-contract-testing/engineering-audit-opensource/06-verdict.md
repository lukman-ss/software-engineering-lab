# Engineering Audit Verdict

Target Lab: `labs/26-contract-testing`
Audit Date: 2026-09-27

## Summary

Code Files Reviewed:
- `internal/model/order.go`
- `internal/provider/server.go`
- `internal/consumer/client.go`
- `internal/contract/verifier.go`
- `cmd/demo/main.go`
- `tests/contract_test.go`
- `go.mod`, `README.md`

Tests Reviewed: 5 tests in `tests/contract_test.go`

Commands Executed (verification run by auditor, not taken on faith):
- `go build ./...` → SUCCESS (exit 0)
- `go test -v ./...` → 5/5 PASS, exit 0
- `go test -race ./...` → `ok labs/26-contract-testing/tests 1.425s`, exit 0
- `go test -race -count=1 ./...` → `ok ... 1.142s` (fresh run, race clean), exit 0
- `go run ./cmd/demo` → exit 0, reproduced all 4 stages with real output

Failures: 0 (none observed during verification)

Warnings:
- Research recommends Pact/Go; implementation uses documented pure-Go stdlib custom verifier (intentional, `ponytail:`-noted).
- Documentation mismatches (stage count, test name, "Interactive CLI" label, stale diff ordering in 03-execution-result.md).
- V2 endpoint `/v2/orders` verified only by demo print, not by an automated test (MEDIUM gap).

## Quality Gates

| Gate | Result | Evidence |
|---|---|---|
| Compilation | PASS | `go build ./...` exit 0 |
| Tests | PASS | 5/5 pass `go test -v ./...` |
| Race Detector | PASS | `go test -race ./...` clean (fresh rerun exit 0) |
| Demo | PASS | `go run ./cmd/demo` exit 0; stages reproduced |
| Research Alignment | WARNING | Deviation from Pact recommendation documented via `ponytail:` (custom verifier preserves CDC semantics + CI-gate blocking). See GAP-5. |
| Documentation Accuracy | WARNING | 6 doc-vs-code mismatches recorded in 04-docs-vs-code.md (stage count, test-name, interactive label, loose test count, Pact deviation, demo diff-ordering). |

## Blocking Issues
(none)

No HIGH/CRITICAL issues. All breaking-change detection paths verified against real execution — the three documented breaking mutations (enum casing, missing `customer.name`, int→string `total`) are each correctly emitted. Demo claims are reproducible (not fabricated): the demo's actual output was captured and matches the described behavior (V1 PASS, Breaking BLOCKED with 3 errors, Dual PASS).

## Non-Blocking Issues
1. DOC_CODE_MISMATCH (LOW): Design doc describes 3 demo stages; actual demo has 4. — See 04-docs-vs-code.md Mismatch 1.
2. DOC_CODE_MISMATCH (LOW): Test named `TestConcurrentContractVerification` in code but `TestConcurrentVerification` in design doc. — Mismatch 2.
3. DOC_CODE_MISMATCH (LOW): Demo labeled "interactive" but is non-interactive. — Mismatch 3.
4. TEST_CLAIM_MISMATCH (LOW): Test asserts `len(errors) >= 3` not exactly 3 named diffs. — Mismatch 4.
5. RESEARCH_IMPLEMENTATION_MISMATCH (MEDIUM): Upstream research verdict recommends active CDC tool (Pact Go/JS); implementation intentionally substitutes a pure-Go stdlib verifier. Well-documented deviation (`ponytail:`), and all researched semantics are fulfilled, so not a failure. — Mismatch 5.
6. DOC_DEMO_MISMATCH (LOW): 03-execution-result.md lists breaking diffs in a fixed order that does not match actual (map iteration is randomized); same 3 diffs are always produced. — Mismatch 6.
7. MISSING_TEST (MEDIUM): V2 `/v2/orders` route not asserted by an automated test (only by demo print). — GAP 1 in 05-gaps.md.

## Required Revisions
(Optional hardening — none block approval)
1. Add automated test asserting V2 endpoint contract compliance (`TestProviderDual_V2ContractVerification_Success`).
2. Tighten `TestProviderBreaking_*` to assert each of the three breaking-diff categories, and/or assert exactly 3.
3. Align design doc stage count and test name to match implementation (4 stages; `TestConcurrentContractVerification`).
4. Either correct "Interactive CLI" label or make CLI prompt-based.
5. (Optional) Note in README that the demo is non-interactive.

## Final Status

APPROVED_WITH_WARNINGS

Reasoning: Code compiles. All 5 tests pass. Race detector clean (re-run fresh, exit 0). Demo reproduces every claimed stage and every claimed breaking diff (verified by execution, not trusted from docs). Core behavior — consumer contract generation, V1 pass, breaking-provider CI-block with 3 diffs, dual-provider V1 compatibility — is proven. Warnings remain: documented deviation from research's tool recommendation (custom verifier in place of Pact) and a set of LOW documentation mismatches. The MEDIUM gap (V2 route lacks unit test), while worth hardening, is compensated by the demo proving the V2 route returns the expected schema, so it does not block approval. Lab is trustworthy for the Technical Writer subject to the noted doc accuracy caveats.
