# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- internal/metrics/tracker.go
- internal/slo/evaluator.go
- internal/alerting/engine.go
- cmd/demo/main.go

Tests:
- tests/slo_test.go

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/05-report.md (SLI ratio definition, Error Budget calculation, Multi-window multi-burn-rate alerting, Release Freeze policy, Service criticality tiering)

Main Claims To Verify:
1. SLI calculation accurately computes ratio of good events over total events across sliding windows.
2. Error Budget accurately calculates total budget `(1-SLO)*Total` and remaining budget `TotalBudget - BadEvents`.
3. Release freeze policy (`CanDeploy`) evaluates to `false` when remaining budget <= 0.
4. Multi-window multi-burn-rate alerting triggers when both short and long window burn rates exceed thresholds, avoiding false positives on transient spikes.
5. In-memory metric tracker correctly handles out-of-order events, eviction, and concurrent access without race conditions.
6. Execution outputs in `engineering/03-execution-result.md` and `README.md` reflect actual execution.

Commands To Run:
- `cd labs/24-slo-sli-error-budget && go test ./...`
- `cd labs/24-slo-sli-error-budget && go test -race ./...`
- `cd labs/24-slo-sli-error-budget && go run ./cmd/demo`

Primary Risks:
- Floating point / rounding inaccuracies in error budget and burn rate calculations.
- Incorrect sliding window eviction when inserting out-of-order timestamps.
- Discrepancies between documentation (README / engineering logs) and actual command outputs.
