# Engineering Audit Plan

Target Lab: `labs/26-contract-testing` (Consumer-Driven Contract Testing in Go)

Implementation Files:
- `internal/model/order.go` — V1, Breaking, V2 DTOs
- `internal/provider/server.go` — ProviderV1 / ProviderBreaking / ProviderDual HTTP handlers
- `internal/consumer/client.go` — MobileOrderClient + GenerateMobileContract
- `internal/contract/verifier.go` — CDC verification engine (diffValues, Verify)
- `cmd/demo/main.go` — CLI demo (Stages 1–4)
- `go.mod` — module `labs/26-contract-testing`

Tests:
- `tests/contract_test.go` — 5 tests: contract generation, V1 pass, breaking fail, dual pass, concurrent

Executable/Demo:
- `go run ./cmd/demo`

Approved Research Inputs (audited):
- `research/01-plan.md` — CDC methodology, breaking-change detection, CI gate semantics
- `research-audit/07-verdict.md` — APPROVED, notes: use active tools (Pact Go/JS) preferred
- `engineering/01-design.md` — success criteria, 7 gates incl. race detector + 3-stage demo
- `engineering/02-implementation-notes.md` — pure Go stdlib decision (ponytail: documented deviation)
- `engineering/03-execution-result.md` — claimed build/test/race/demo output

NOTE (pipeline override): research content NOT audited here. Implementation + tests + demo + docs vs code only.

Main Claims To Verify:
1. Build compiles (`go build ./...`) → SUCCESS claim matches code.
2. All unit/contract tests pass (`go test -v ./...`) → 5 tests listed, all PASS.
3. Race detector clean (`go test -race ./...`) → `ok ... tests 1.425s` claim.
4. Demo reproduces 4-stage lifecycle with CI gates (Stage 2 V1 PASS, Stage 3 Breaking BLOCKED w/ 3 diffs, Stage 4 Dual PASS).
5. Breaking provider yields exactly 3 breaking diffs (enum casing, missing field, type mismatch).
6. V1 and Dual providers pass verification against consumer contract.
7. Concurrency safe: 20 goroutines × Verify, race clean.
8. Minimal subset rule: consumer declares {id, status, customer.name, total}; extra fields ignored.
9. V2 safe-evolution path: `/v2/orders` adds currency field while `/v1` stays contract-compatible.

Commands To Run:
- `cd labs/26-contract-testing && go build ./...`
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Demo stage numbering/name divergence from design doc (design says 3 stages; demo has 4).
- Design doc lists test name `TestConcurrentVerification`; actual test named `TestConcurrentContractVerification`.
- Research verdict recommends active tools (Pact); implementation intentionally substitutes pure-Go stdlib verifier (ponytail:). Auditor must confirm this is documented intent, not a research mismatch.
- Demo output error ordering vs documented (Stage 4 in 03-execution-result lists `total` diff before `status`; actual code emits `status` first).
