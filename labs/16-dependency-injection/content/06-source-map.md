# Source Map

## Core Concept: Separation of Configuration from Use
Research:
- `research/runs/2026-09-25-dependency-injection/05-report.md` (Finding 1)
- `research/runs/2026-09-25-dependency-injection/02-sources.md` (Source 1 - Martin Fowler, Source 2 - Microsoft)

Implementation:
- `cmd/demo/main.go`
- `internal/di/gateway.go`

## Constructor Injection vs Service Locator Anti-Pattern
Research:
- `research/runs/2026-09-25-dependency-injection/05-report.md` (Finding 3, Finding 4)
- `research/runs/2026-09-25-dependency-injection/02-sources.md` (Source 4 - PSR-11 Meta Document)

Implementation:
- `internal/di/processor.go` (Constructor Injection)
- `internal/di/locator.go` (Service Locator Anti-Pattern)

## What the Tests Prove (Fast, Isolated Unit Testing)
Research:
- `research/runs/2026-09-25-dependency-injection/05-report.md` (Finding 2)
- `engineering/01-design.md`

Implementation & Tests:
- `tests/processor_test.go`
- `internal/di/gateway.go` (RealGateway vs PaymentGateway interface)

## Value Objects Direct Instantiation
Research:
- `research/runs/2026-09-25-dependency-injection/05-report.md` (Finding 5)

Implementation:
- `internal/di/gateway.go` (`Money` struct)
- `internal/di/processor.go`
