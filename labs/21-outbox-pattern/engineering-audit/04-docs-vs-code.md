# Docs vs Code Comparison

Target Lab: labs/21-outbox-pattern

## 1. README vs Code & Execution

- `README.md` lists files accurately (`internal/outbox/db.go`, `broker.go`, `service.go`, `relay.go`, `consumer.go`).
- Running instructions `go test ./...`, `go test -race ./...`, and `go run ./cmd/demo` function as specified.
- README accurate alignment: PASS.

## 2. Engineering Notes vs Code

- `engineering/01-design.md` mentions:
  - "SQLite Database (- orders table, - outbox_events table)" in architecture diagram and components.
  - Actual implementation in `internal/outbox/db.go` is a pure Go thread-safe in-memory transactional database (`DB` struct with `Tx` staging).
  - Note: `01-design.md` section "Implementation Decisions" acknowledges using an in-memory transactional database driver to avoid external CGO dependencies.
  - Severity: LOW (Minor wording discrepancy in architecture diagram text vs actual in-memory implementation).

- `engineering/02-implementation-notes.md` matches `internal/outbox/*` files and demo output accurately.

- `engineering/03-execution-result.md` accurately records actual output from running `go build`, `go test`, `go test -race`, and `go run ./cmd/demo`.

## Discrepancies Found

- **DOC_CODE_MISMATCH (Minor)**: `engineering/01-design.md` diagram mentions "SQLite Database", whereas actual code implements a custom thread-safe in-memory transactional database (`internal/outbox/db.go`). Implementation notes (`02-implementation-notes.md`) explicitly clarify this design choice.
