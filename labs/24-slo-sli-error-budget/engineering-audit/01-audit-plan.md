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
1. Sliding window time-bucketed metric tracking accurate for good/total request count.
2. SLI calculation matches good/total events ratio.
3. Error budget total, consumed, and remaining calculations conform to SRE principles.
4. Release freeze policy (`CanDeploy`) correctly toggles when budget is <= 0.
5. Multi-window multi-burn-rate alerting correctly triggers when both short and long window burn rates exceed rule factor.
6. Code handles out-of-order event timestamps gracefully without corruption.
7. Concurrency operations on WindowTracker are data-race free under high parallel load.

Commands To Run:
```bash
cd labs/24-slo-sli-error-budget
go test ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Thread-safety / race conditions on metric slice mutations.
- Stale bucket eviction errors when timestamps arrive out-of-order.
- Floating-point precision / rounding issues in SLI or Error Budget evaluation.
- False positive alerts if multi-window conditions are not properly evaluated.
