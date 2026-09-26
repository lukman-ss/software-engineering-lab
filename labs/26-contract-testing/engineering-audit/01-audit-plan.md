# Engineering Audit Plan

Target Lab: labs/26-contract-testing
Implementation Files:
- `internal/consumer/client.go`
- `internal/contract/verifier.go`
- `internal/model/order.go`
- `internal/provider/server.go`

Tests:
- `tests/contract_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

Main Claims To Verify:
1. Consumer generates minimal required schema/interaction specification contract.
2. Custom lightweight Go contract verifier detects non-breaking and breaking API changes.
3. Breaking changes (casing change `IN_PROGRESS` -> `in_progress`, field rename `name` -> `full_name`, type mutation `int64` -> `string`) are caught and fail contract verification.
4. Dual provider (V1 & V2) deployment maintains backwards compatibility for V1 consumer while enabling API evolution.
5. All code compiles, tests pass, race detector passes, and interactive demo runs correctly.
6. Documentation (`README.md`) accurately reflects codebase paths, usage, and commands.

Commands To Run:
- `go test -v ./...`
- `go test -race ./...`
- `go run ./cmd/demo`

Primary Risks:
- Partial diff engine in contract verifier missing slice or unhandled type recursive checks.
- False pass on breaking provider schema changes.
- Unhandled HTTP transport errors during verifier execution causing panic or false result.
