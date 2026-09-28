# Docs vs Code Audit

## Consistency Checks

### 1. README vs Implementation
- README lists components: `internal/fault`, `internal/circuitbreaker`, `internal/monitor`, `internal/experiment`, `cmd/demo`, `tests`.
- Actual directory structure exactly matches the listed layout.
- README describes 4 core principles:
  1. Steady-state health metric monitoring -> implemented in `internal/monitor/monitor.go`.
  2. Injected downstream service latency and errors -> implemented in `internal/fault/injector.go`.
  3. Resilience validation via Circuit Breakers and Fallback mechanisms -> implemented in `internal/circuitbreaker/circuitbreaker.go`.
  4. Blast radius control with automated experiment abort upon metric degradation -> implemented in `internal/experiment/runner.go`.

### 2. Engineering Notes vs Execution Results
- `engineering/03-execution-result.md` captures actual output matching demo and test execution byte-for-byte.
- Demo logs match live runtime execution of `go run ./cmd/demo`.

### 3. Discrepancies / Mismatches
- Minor mismatch in `engineering/01-design.md`: Architecture section references `pkg/fault`, `pkg/monitor`, `pkg/circuitbreaker`, `pkg/experiment` instead of `internal/fault`, etc. However, `engineering/02-implementation-notes.md` and `README.md` correctly reference `internal/` packages.
- No functional or test-claim discrepancies found.
