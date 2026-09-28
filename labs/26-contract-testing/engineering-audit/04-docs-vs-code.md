# Docs vs Code Audit

Target Lab: labs/26-contract-testing

## Comparison Summary

| Document Element | Claimed in README / Engineering Notes | Observed in Implementation | Status |
|---|---|---|---|
| Project Structure | Lists `cmd/demo/main.go`, `internal/...`, `tests/contract_test.go` | Tree structure matches exact path layout | MATCH |
| Test Command | `go test -v ./...` & `go test -race ./...` | All commands run successfully and pass | MATCH |
| Demo Command | `go run ./cmd/demo` | Demo output prints 4 stages matching described workflow | MATCH |
| Contract Schema | `id`, `status`, `customer.name`, `total` | Exact fields in `GenerateMobileContract` | MATCH |
| Breaking Changes | Casing, type mutation int->str, field rename | Exact 3 violations emitted by `ProviderBreaking` | MATCH |

## Discrepancy Findings

No mismatches found between `README.md`, `engineering/` design documents, and actual Go code implementation.
