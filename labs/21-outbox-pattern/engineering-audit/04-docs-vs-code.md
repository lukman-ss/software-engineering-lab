# Docs vs Code Audit

## Documentation Sources

- `README.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

## Comparison

### Architecture & Components
- README lists `db.go`, `broker.go`, `service.go`, `relay.go`, `consumer.go`. Matches codebase exactly.
- README describes in-memory transactional database simulating `BeginTx`, `Commit`, and `Rollback`. Matches `internal/outbox/db.go`.
- `engineering/01-design.md` proposed SQLite (`modernc.org/sqlite` or mock DB driver), while implementation chose the in-memory mock DB driver to maintain zero-dependency pure Go. `engineering/02-implementation-notes.md` accurately documents this decision.

### Execution Commands
- README lists `go test ./...`, `go test -race ./...`, `go run ./cmd/demo`.
- All commands execute cleanly with matching output.

### Discrepancy Findings
- DOC_CODE_MISMATCH: None between README and codebase. Minor divergence between initial design plan (SQLite) and final implementation (in-memory DB), acknowledged in implementation notes.
