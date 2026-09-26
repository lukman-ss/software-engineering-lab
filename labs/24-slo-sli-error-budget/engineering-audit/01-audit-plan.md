# Engineering Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
Implementation Files:
- `internal/metrics/tracker.go`
- `internal/slo/evaluator.go`
- `internal/alerting/engine.go`
- `cmd/demo/main.go`

Tests:
- `tests/slo_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `research-revision/03-revision-result.md`
- `research-audit/07-verdict.md`

Main Claims To Verify:
1. SLI Calculation accurately measures good events / total events ratio across a sliding window.
2. Error Budget calculation and exhaustion tracking (`totalErrorBudget - badEvents`), correctly blocking deployments when budget is exhausted (`CanDeploy: false`).
3. Multi-Window Multi-Burn-Rate alerting logic computes burn rates across short/long windows and triggers appropriate alert severities (PAGE / TICKET) without false positives.
4. Concurrency safety of metrics aggregation under concurrent traffic.
5. Zero-traffic handling, edge case correctness, and timestamp eviction logic.
6. Execution behavior matches claims in `README.md` and `engineering/`.

Commands To Run:
```bash
go test -v -count=1 ./...
go test -race -count=1 ./...
go run ./cmd/demo
```

Primary Risks:
- Thread-safety / race conditions during concurrent bucket eviction or writes.
- Edge case division by zero when total requests = 0 in SLI, budget, or burn-rate calculations.
- Timestamp ordering and eviction handling in sliding time buckets.
- Precision/rounding discrepancies in SLI or Error Budget reporting.
- Mismatch between README documentation and actual code implementation.
