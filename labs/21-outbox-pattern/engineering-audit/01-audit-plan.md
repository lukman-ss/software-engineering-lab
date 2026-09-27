# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern
Implementation Files:
- internal/outbox/model.go
- internal/outbox/db.go
- internal/outbox/broker.go
- internal/outbox/service.go
- internal/outbox/relay.go
- internal/outbox/consumer.go
Tests:
- tests/outbox_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md
- engineering/01-design.md
- engineering/02-implementation-notes.md
Main Claims To Verify:
1. Atomicity of business record and outbox record commit in single transaction.
2. Rollback discards both business and outbox records.
3. Dual-write naive pattern leaves system inconsistent on broker failure.
4. Polling relay asynchronously queries pending records and publishes to broker, marking records processed.
5. Idempotent consumer discards duplicate deliveries safely.
6. Processed outbox records can be purged.
7. Concurrency safety under Go race detector.
Commands To Run:
- go test -v -count=1 ./...
- go test -race -v -count=1 ./...
- go run ./cmd/demo
Primary Risks:
- In-memory mock DB mimicking transaction semantics rather than real SQL engine (modernc.org/sqlite noted in 01-design.md).
- Relay error handling and retry loop behavior under persistent broker failures.
- Concurrency test asserting race safety rather than throughput or ordering consistency.
