# Engineering Audit Plan

Target Lab: labs/27-database-constraints
Implementation Files:
- cmd/demo/main.go
- internal/store/store.go
- internal/store/store_test.go
- internal/engine/engine.go
- internal/model/model.go
- internal/dberr/errors.go
Tests: ./internal/store/...
Executable/Demo: go run ./cmd/demo
Approved Research Inputs: research/ (plan, sources, evidence, report)
Main Claims To Verify:
- Database constraints (NOT NULL, CHECK, UNIQUE, FOREIGN KEY, PARTIAL UNIQUE) are correctly enforced
- Concurrent registration test demonstrates race condition in UnsafeStore vs safety in SafeStore
- Error mapping returns appropriate domain errors
- Demo shows all constraint violations and concurrency behavior
Commands To Run:
- go test ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Implementation only simulates constraints via in-memory maps; not a real DB
- Concurrency safety relies on mutex; correct but may have performance implications not tested
- No actual SQL or database; correctness is limited to simulation