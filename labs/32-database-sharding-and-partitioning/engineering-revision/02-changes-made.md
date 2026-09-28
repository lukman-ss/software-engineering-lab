## Revision 1

Audit Issue: `ExtractTimeFromUUIDv7` utility function lacks explicit test assertion in `TestIDGenerators`.
Severity: LOW
Files Changed:
- `internal/idgen/idgen.go`
- `tests/sharding_test.go`
Action: Refactored `ExtractTimeFromUUIDv7` byte parsing logic and added explicit assertion in `TestIDGenerators` verifying timestamp extraction from generated UUIDv7.
Verification: `go test -v -count=1 ./tests/...`, `go test -race -v -count=1 ./tests/...`
Status: RESOLVED
