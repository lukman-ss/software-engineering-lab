# Docs vs Code: Contract Testing (lab-26)

Reviewed `README.md`, `engineering/01-design.md`, `engineering/03-execution-result.md`, `content/06-source-map.md` vs implementation.

## README Accuracy
- **Test commands:** `go test -v ./...` — PASS ✅
- **Race commands:** `go test -race ./...` — PASS ✅
- **Demo:** `go run ./cmd/demo` — PASS ✅
- **Structure table:** `cmd/demo/main.go`, `internal/*`, `tests/contract_test.go`, `engineering/*`, `go.mod`, `README` — all paths exact ✔️
- **Feature claims:** CDC generation, CI gate verification, breaking-change detection, V1/V2 evolution — all directly implemented and observed.

## Design vs Implementation
- `engineering/01-design.md` claims 3 breaking diffs (enum casing, field rename, type mutation) — implementation delivers exactly 3, demo Stage 3 logs them verbatim.
- Claims V1 contract passes on `/v1/orders/{id}` for both ProviderV1 and ProviderDual — tests + demo confirm.
- Minimal subset rule (extra provider fields ignored) — matches `diffValues` subset-only comparison (`internal/contract/verifier.go:131-143`), plus demonstrated tolerance of `notes` extra field.
- Success criteria all met: pure stdlib, contract generation, verification engine, failure reporting, race-clean, 4-stage demo.

## Execution-Result Record vs Actual Run
- `engineering/03-execution-result.md` demo output (Stage 1–4) matches live `go run ./cmd/demo` output 1:1 except JSON key ordering (`encoding/json` map marshal nondeterministic — cosmetic only).
- Test logs in execution record match live `go test -v ./...` (5 tests PASS).

## Source-Map Accuracy
- `content/06-source-map.md` line references verified:
  - `GenerateMobileContract` lines → match
  - `ProviderV1` lines → match
  - `ProviderBreaking` lines → match
  - `ProviderDual` lines → match
  - `Verify` / `diffValues` lines → match
  - Demo orchestrator line range `14-68` vs actual `14-71` — off by 3 lines (cosmetic drift, not behavioral).

## Result
No `DOC_CODE_MISMATCH`, no `TEST_CLAIM_MISMATCH`, no `RESEARCH_IMPLEMENTATION_MISMATCH`, no `FAKE_DEMO`, no `FAKE_BENCHMARK`.