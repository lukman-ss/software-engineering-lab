# Lab 24: SLO, SLI, and Error Budget Implementation

This lab implements Service Level Indicators (SLI), Service Level Objectives (SLO), Error Budget tracking, and Multi-Window Multi-Burn-Rate alerting in Go based on Google SRE principles.

## Structure

- `internal/metrics`: Sliding-window time-bucketed event tracker for recording requests and measuring good vs. total events.
- `internal/slo`: Evaluator calculating SLI ratios, remaining Error Budget, and release freeze policy enforcement.
- `internal/alerting`: Multi-window burn-rate alert calculator evaluating fast and slow budget burn rates against SLO thresholds.
- `cmd/demo`: Executable demonstration illustrating baseline SLO tracking, error budget depletion during an incident, and burn rate alert triggering.
- `tests/`: Unit and concurrency tests ensuring thread-safety and mathematical correctness.

## Running Tests

```bash
go test ./...
go test -race ./...
```

## Running Demo

```bash
go run ./cmd/demo
```
