# Docs vs Code Audit — labs/24-slo-sli-error-budget

## Comparison Source -> Target

### README.md -> code
- metrics tracker description ("Sliding-window time-bucketed event tracker") -> internal/metrics/tracker.go : accurate. PASS
- slo evaluator ("SLI ratios, remaining Error Budget, release freeze policy") -> internal/slo/evaluator.go : accurate. PASS
- alerting ("Multi-window burn rate alert calculator") -> internal/alerting/engine.go : accurate. PASS
- cmd/demo runnable -> present and executes. PASS
- test commands (`go test ./...`, `go test -race ./...`) -> present, pass. PASS

### engineering/01-design.md -> code
- "SLOEvaluator: Calculates SLI, remaining Error Budget, and current Burn Rate." -> internal/slo/evaluator.go : no burn-rate calculation exists (only alerting engine computes burn rate). DOC_CODE_MISMATCH, LOW
- "internal/metrics: ... histogram latency buckets & success counts" -> internal/metrics/tracker.go : no latency histogram buckets; Duration stored on Event, evaluated only by caller's isGood closure. DOC_CODE_MISMATCH, LOW
- "100% test coverage on core math and sliding window calculations" -> tests/slo_test.go : 6 tests pass, but core math not exhaustively covered (see test audit). Implementation overclaim, MEDIUM

### engineering/03-execution-result.md -> actual demo
- Recorded demo output lists only PHASE 1-3. Real `go run ./cmd/demo` (run during audit) prints PHASE 4 (Endpoint Criticality Comparison). DOC_CODE_MISMATCH, LOW (stale result doc)
- Recorded result claims build/test/race PASS : matches actual run. PASS

### cmd/demo/main.go inline narrative -> computed values
- Printed comment "// ... burn rate = 0.10 / 0.001 = 100x" but aggregated window (baseline + incident) yields 9.09x burn, so only the 6.0x (TICKET) rule fires. DOC_CODE_MISMATCH in narrative comment, MEDIUM (misleading, not a logic bug; alert engine behaves consistently)

## Verdict on docs
README matches implementation. Design/execution docs contain aspirational/stale claims. No claims invented to paper over failures.
