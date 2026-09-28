# Docs vs Code

Target Lab: labs/26-contract-testing
Scope: README, engineering notes, code, tests, demo. Research/content excluded per pipeline override.

## README vs Code
- README structure matches repo layout (cmd/demo, internal/*, tests, engineering, go.mod). PASS.
- README commands (`go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`) all executed live and succeed. PASS.
- README claims (CDC, CI gate, 3 breaking mutations, V1+V2 dual) each proven by code + tests. PASS.

## Engineering Notes vs Code
- 03-execution-result.md transcript reproduces live demo modulo diff ordering (map iteration). Real output, no fabrication. PASS with note below.
- 03-execution-result.md lists 5 tests; repo now has 7 (2 added by revision). Stale count, results still accurate. See GAP-4.
- 01-design.md "exit code non-zero" for breaking path overstates demo (demo exits 0 on expected block; exits 1 only on inverted logic). See GAP-1.
- 01-design.md architecture references `contracts/mobile_order_v1.json` file; implementation is in-memory (JSON marshal for display only). See GAP-2.

## Mismatch Register
| ID | Type | Location | Severity |
|----|------|----------|----------|
| D1 | DOC_CODE_MISMATCH | design exit-code claim vs demo/main.go:44-54 | LOW |
| D2 | DOC_CODE_MISMATCH | design `contracts/*.json` path vs in-memory contract | LOW |
| D3 | DOC_CODE_MISMATCH | 03-execution-result.md test count (5 vs 7) | LOW |

No TEST_CLAIM_MISMATCH. No FAKE_DEMO. No FAKE_BENCHMARK.
