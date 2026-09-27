# Engineering Audit Plan

Target Lab: labs/21-outbox-pattern
Implementation Files: 
- internal/outbox/model.go
- internal/outbox/db.go
- internal/outbox/broker.go
- internal/outbox/service.go
- internal/outbox/relay.go
- internal/outbox/consumer.go
- cmd/demo/main.go

Tests: tests/outbox_test.go

Executable/Demo: cmd/demo/main.go

Approved Research Inputs: 
- research/05-report.md (Research Report: APPROVED)
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md

Main Claims To Verify:
1. Atomic persistence of business data and outbox message in single transaction
2. Transaction rollback discards both business data and outbox message
3. Polling relay fetches pending outbox messages, publishes to broker, marks as processed
4. Consumer enforces idempotency via message ID tracking
5. Dual-write pattern (naive) leads to inconsistency when broker fails after DB commit
6. Relay retries after broker failure (messages remain PENDING)
7. Thread-safety under concurrent writes (no race conditions)
8. Processed outbox messages can be purged

Commands To Run:
- go build ./...
- go test ./...
- go test -race ./...
- go run ./cmd/demo

Primary Risks:
- Inaccurate representation of production outbox schema (missing aggregate IDs)
- Relay uses polling not CDC/log-tailing (higher latency)
- No exponential backoff or retry limits in relay
- No outbox cleanup automation (manual purge only)
- Test concurrency uses same IDs, limiting effectiveness of concurrent test