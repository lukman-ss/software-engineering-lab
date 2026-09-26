# Docs vs Code Audit

Target Lab: `labs/26-contract-testing`

## Comparison Matrix

| Claim / Section | README / Engineering Notes Claim | Code / Test Implementation | Alignment |
|---|---|---|---|
| Project Structure | Lists `cmd/demo`, `internal/consumer`, `internal/contract`, `internal/model`, `internal/provider`, `tests/` | Exact match with directory contents | PASS |
| Consumer-Driven Contract | Mobile consumer specifies minimal required subset (id, status, customer.name, total) | `internal/consumer/client.go:82-110` defines interaction | PASS |
| CI Verification Gate | Provider validates against consumer contract before deployment | `internal/contract/verifier.go:58-115` executes validation | PASS |
| Breaking Change Detection | Detects enum casing, field rename/omission, primitive type mutation | `internal/provider/server.go:61-72` exercises all three; detected by test and demo | PASS |
| Safe API Evolution | Dual versioning allows V1 compatibility alongside V2 evolution | `internal/provider/server.go:83-130` and `tests/contract_test.go:74-94` | PASS |
| Run Commands | `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo` | All commands run without error and pass | PASS |

## Discrepancies
None detected. Documentation accurately reflects codebase structure, functionality, and execution outputs.
