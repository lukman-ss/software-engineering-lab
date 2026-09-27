# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Audit Output: labs/26-contract-testing/engineering-audit-opensource/
Pipeline Override: implementation + tests only. No research/content audit. No code modification.

## Implementation Files

- internal/contract/verifier.go (CDC engine, Verify + diffValues)
- internal/model/order.go (V1, Breaking, V2 DTOs)
- internal/provider/server.go (ProviderV1, ProviderBreaking, ProviderDual)
- internal/consumer/client.go (MobileOrderClient.FetchOrder, GenerateMobileContract)
- cmd/demo/main.go (4-stage lifecycle demo)
- go.mod (module labs/26-contract-testing, go 1.22, stdlib only)

## Tests

- tests/contract_test.go (5 tests: generation, V1 pass, breaking fail, dual pass, concurrent)

## Executable/Demo

- go run ./cmd/demo (httptest-backed 4-stage demo, no external deps)

## Approved Research Inputs

Not audited per pipeline override. Design reference only: engineering/01-design.md.

## Main Claims To Verify

1. Consumer generates minimal contract (id, status, customer.name, total).
2. Provider V1 passes verification, exit 0, gate ALLOWED.
3. Breaking provider (enum casing, field rename, int->string) fails with exact diffs, gate BLOCKED.
4. Dual provider preserves V1 compat while exposing V2.
5. Concurrency-safe verification under -race.
6. README commands and structure match code.
7. Demo output real, no fake benchmark.

## Commands To Run

- go build ./...
- go test -v ./...
- go test -race -count=1 ./...
- go vet ./...
- go run ./cmd/demo

## Primary Risks

- Verifier claims header comparison but may not implement it.
- V2 endpoint may be untested (dual test may cover V1 path only).
- Verifier http.Client may lack timeout.
- Negative paths (bad status, bad JSON) may be untested.
