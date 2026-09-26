# Docs vs Code Audit

## Comparisons

1. **README vs Code**:
   - README claims `internal/adr` provides `models.go`, `parser.go`, and `linter.go`. Matches code exactly.
   - README claims `cmd/demo/main.go` demonstrates Modular Monolith to Microservices. Matches code exactly.
   - README test and run commands (`go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`) are correct and function as described.

2. **Engineering Notes vs Demo Output**:
   - Output listed in `engineering/03-execution-result.md` exactly matches the actual runtime output of `go run ./cmd/demo`.

3. **Mismatches**:
   - None found.
