# Docs vs Code Audit

## Consistency Matrix

| Topic | README / Engineering Docs | Implementation / Tests | Status |
|---|---|---|---|
| Query Count (Naive) | 1 query for authors + N queries for posts (N=3, total=4) | `TestGetAuthorsWithPostsNPlusOne` expects count == 4 | MATCH |
| Query Count (Eager) | 1 query for authors + 1 batched query for posts (total=2) | `TestGetAuthorsWithPostsEager` expects count == 2 | MATCH |
| Implementation Path | `cmd/demo/main.go`, `internal/blog/*` | Exact path and files present | MATCH |
| Test Commands | `go test -v ./...`, `go test -race ./...` | Both commands pass cleanly | MATCH |
| Demo Command | `go run ./cmd/demo` | Runs and produces matching text output | MATCH |

## Discrepancies Found

None. Documentation strictly reflects code structure, test commands, and demo execution output.
