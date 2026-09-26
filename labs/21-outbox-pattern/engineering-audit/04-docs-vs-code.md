# Docs vs Code Audit

## 1. README.md vs Code
- **README Claims**:
  - Components: `db.go` (in-memory transactional DB), `broker.go` (mock message broker), `service.go` (dual-write vs atomic outbox writes), `relay.go` (asynchronous polling worker), `consumer.go` (idempotent subscriber).
  - Test command: `go test ./...` and `go test -race ./...`.
  - Demo command: `go run ./cmd/demo`.
- **Reality**: All listed files exist, exports and responsibilities match exactly. Commands run without errors.
- **Assessment**: PASS

## 2. engineering/01-design.md vs Code
- **Design Claims**:
  - Mentions SQLite DB in Section 4. Architecture (`SQLite Database`, `orders`, `outbox_events`).
  - Mentions cleanup worker in Section 2 & 3.
- **Reality**:
  - Section 7 of design doc clarifies implementation decision: in-memory mock transactional DB is chosen for zero external CGO dependencies and portable testing.
  - Cleanup worker is listed in design doc but omitted in `engineering/02-implementation-notes.md` as non-core scope.
- **Assessment**: WARNING (Minor discrepancy regarding cleanup worker in early design doc, clarified in implementation notes).

## 3. engineering/03-execution-result.md vs Execution
- **Recorded Results**:
  - `go build ./...`: PASS (exit code 0).
  - `go test -count=1 ./...`: PASS.
  - `go test -count=1 -race ./...`: PASS.
  - `go run ./cmd/demo`: Output exactly matches Scenario 1, Scenario 2, Scenario 3 outputs.
- **Reality**: All executed outputs reproduced 1:1 against current codebase.
- **Assessment**: PASS
