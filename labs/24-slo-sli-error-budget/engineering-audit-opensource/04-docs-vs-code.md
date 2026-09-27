# Docs vs Code — labs/24-slo-sli-error-budget

Compared: README.md, engineering/01-design.md, engineering/02-implementation-notes.md, engineering/03-execution-result.md vs tracker.go, evaluator.go, engine.go, demo/main.go, tests/slo_test.go. Research excluded per override.

## D1 — Execution result matches reality
Docs: engineering/03-execution-result.md records build PASS, 6/6 tests, race ok, full demo transcript.
Code/reality: `go vet` clean; `go test -count=1 -v` 6/6 PASS; `-race` clean; `go run ./cmd/demo` output numerically identical (1000/0 → 1100/1090/10, -8.90, 9.09x TICKET-only). No FAKE_DEMO / FAKE_BENCHMARK.
Type: none. Status: VERIFIED.

## D2 — README matches code
Docs: README structure list, test/race/demo commands.
Code: Paths and commands all valid; ran as written.
Type: none. Status: VERIFIED.

## D3 — Dead field: slo.Config.LatencyThreshold
Docs: Design implies evaluator enforces latency threshold.
Code: evaluator.go never reads LatencyThreshold; latency judged only by caller's isGood closure.
Type: DOC_CODE_MISMATCH. Severity: LOW. Fix docs or wire field.

## D4 — Per-rule window fields unimplemented
Docs/code: BurnRateRule carries LongWindow/ShortWindow/BudgetConsumedPct; design says per-rule multi-window engine.
Code: engine.Check ignores all three; all rules share two construction-time trackers.
Type: IMPLEMENTATION_OVERCLAIM (scoped to alerting generality, core AND-gating works). Severity: MEDIUM.

## D5 — Unproven claims in design success criteria
Docs: "100% test coverage", demo "recovery" phase, "1-second = 1-hour scale", "histogram latency buckets", component named SLOEvaluator.
Code/tests: No coverage report (edge branches uncovered); demo has no recovery phase; no time-scaling logic; plain bucket counts; type is Evaluator.
Type: TEST_CLAIM_MISMATCH / DOC_CODE_MISMATCH. Severity: LOW (core math proven; wording overclaims).

## D6 — Demo Phase 4 differentiation unproven
Docs/demo comment: Reports wider tolerance differentiates criticality.
Code: Identical 10% error rate vs both 99.9% and 95% SLOs freezes both (Reports -5.00).
Type: IMPLEMENTATION_OVERCLAIM (scenario choice, math correct). Severity: LOW.
