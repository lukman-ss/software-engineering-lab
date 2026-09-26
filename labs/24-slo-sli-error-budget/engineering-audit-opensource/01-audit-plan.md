# Engineering Audit Plan

Target Lab: labs/24-slo-sli-error-budget
Implementation Files:
- internal/metrics/tracker.go
- internal/slo/evaluator.go
- internal/alerting/engine.go
Tests:
- tests/slo_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/01-plan.md
- research/02-sources.md
- research/03-evidence.md
- research/04-contradictions.md
- research/05-report.md
- research/06-open-questions.md
Main Claims To Verify:
1. Accurate SLI evaluation (ratio of good events to total events)
2. Error Budget calculation: remaining = (1-target)*total - bad; triggers CanDeploy=false when budget exhausted
3. Multi-window burn-rate alerting: short and long windows evaluated concurrently; alert triggers when both exceed burnRateFactor
4. Thread-safe concurrent access to metrics tracker (no race conditions)
5. Correct handling of zero traffic (SLI=1.0, CanDeploy=true)
6. Correct eviction of stale events outside window
7. Demo correctly illustrates baseline, incident, alerting, and multi-SLO comparison
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Concurrency bugs in WindowTracker (mutex locking correctness)
- Off-by-one errors in budget calculations
- Alert logic requiring both short and long windows to exceed threshold (may miss fast spikes in only one window)
- Time truncation/bucket alignment edge cases
- Integer overflow in high-throughput scenarios (not applicable here due to int64)