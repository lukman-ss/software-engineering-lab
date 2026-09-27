# Engineering Audit Plan

Target Lab: `labs/26-contract-testing`
Implementation Files:
- `internal/consumer/client.go`
- `internal/contract/verifier.go`
- `internal/model/order.go`
- `internal/provider/server.go`
- `cmd/demo/main.go`

Tests:
- `tests/contract_test.go`

Executable/Demo:
- `cmd/demo/main.go`

Approved Research Inputs:
- `research/05-report.md`
- `engineering/01-design.md`
- `engineering/02-implementation-notes.md`

Main Claims To Verify:
1. Implementation matches Consumer-Driven Contract (CDC) testing and CI deployment gate principles.
2. Code compiles cleanly without external runtime/C-binding dependencies.
3. Tests run and pass 100% including race detection (`go test -race ./...`).
4. Contract verifier correctly flags provider breaking changes (enum mismatch, field rename, primitive type change).
5. Dual provider maintains backward compatibility for V1 while supporting V2.
6. Demo output is authentic, reproducible, and reflects actual code execution.
7. README accurately describes project structure, test commands, and execution instructions.

Commands To Run:
```bash
cd labs/26-contract-testing
go test -v ./...
go test -race ./...
go run ./cmd/demo
```

Primary Risks:
- Type mismatch detection when JSON numbers vs strings or floats are decoded in standard Go map representations.
- Concurrency race conditions in HTTP client or server verification loops.
- Inconsistency between README commands/structure and actual repository layout.
