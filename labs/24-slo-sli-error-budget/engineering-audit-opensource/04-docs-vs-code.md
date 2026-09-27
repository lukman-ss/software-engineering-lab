# Documentation vs Code Audit — labs/24-slo-sli-error-budget

## Claim Comparison

| Claim (source) | Code reality | Status |
|---|---|---|
| SLI = good/total ratio | internal/slo/evaluator.go:46-47 | MATCH |
| Error Budget = (1-SLO)*total | internal/slo/evaluator.go:49-52 | MATCH |
| Budget consumed on bad events | internal/slo/evaluator.go:51 | MATCH |
| CanDeploy=false when budget exhausted | internal/slo/evaluator.go:54-57 | MATCH |
| Multi-window burn-rate alerting (short&&long) | internal/alerting/engine.go:73 | MATCH |
| WindowTracker holds dead-letter retention | internal/metrics/tracker.go (rolling buckets only) | MISMATCH |
| 100% test coverage on core math | tests/slo_test.go (6 tests; no coverage tooling run) | UNVERIFIED |
| Go race detector passes | verified: ok 1.109s | MATCH |
| Demo shows eviction, budget depletion, alert, freeze | cmd/demo + recorded output | MATCH |

## Doc_Code_Mismatch

engineering/01-design.md §Architecture lists `internal/metrics` as "histogram latency buckets & success counts" — code stores only good/bad/total counts, no latency histogram. Low severity; latency modeled via isGood closure.

engineering/03-execution-result.md Phase 4 caption reads "Reports has wider 5% error tolerance" yet both Payment and Reports report CanDeploy=false; Reports SLI=90% < 95% target, budget -5.00 → correctly frozen. Claim wording misleading but behavior correct.

BurnRateRule.LongWindow / ShortWindow documented as part of rule definition but unused in code (implementation-specific: windows are engine-level).

Config.LatencyThreshold listed in design as latency threshold driver; evaluator never references it (closer defined by tracker closure).

## Test_Claim_Mismatch

"100% test coverage on core math" not verified by `go test -cover`; no coverage report present. Tests are strong but coverage not proven.

## Research_Implementation_Mismatch

N/A — research audit skipped per pipeline override. Cross-checked only against design/implementation docs.

## Summary

README and engineering notes align with code on core formulas. Several documentation claims overstate scope (latency histograms, per-rule windows, coverage %) and one caption is misleading but not incorrect. No fabricated results; demo output reproduced verbatim by fresh execution.
