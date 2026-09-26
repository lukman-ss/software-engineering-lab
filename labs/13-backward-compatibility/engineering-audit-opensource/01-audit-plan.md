# Engineering Audit Plan

Target Lab: labs/13-backward-compatibility
Implementation Files:
- internal/compat/model.go
- internal/compat/store.go
- internal/compat/flags.go
- internal/compat/metrics.go
- internal/compat/backfill.go
- internal/compat/service.go
- internal/compat/handler.go
- cmd/demo/main.go
Tests:
- internal/compat/service_test.go
- tests/migration_test.go
- tests/concurrency_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: 
- engineering/01-design.md (approved design)
Main Claims To Verify:
1. Expand-Migrate-Contract pattern implementation (dual-write, backfill, fallback reads, contract drop)
2. Backward compatibility (legacy clients work throughout)
3. Forward compatibility (modern clients work)
4. Safe rollback capability (no data loss during rollback from dual-write)
5. Data consistency (drift detection and reconciliation)
6. Observability and deprecation headers
7. Concurrency safety (no races under -race)
8. Demo matches claimed lifecycle
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
- go vet ./...
Primary Risks:
- Incorrect dual-write implementation causing data divergence
- Backfill not being idempotent/resumable
- Contract enforcement not properly guarded (dropping legacy while traffic exists)
- Fallback read not working causing data starvation
- Concurrency bugs in shared state (maps, counters)
- Demo not matching actual implementation