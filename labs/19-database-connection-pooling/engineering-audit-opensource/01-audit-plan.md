# Engineering Audit Plan

Target Lab: labs/19-database-connection-pooling
Implementation Files:
- internal/pool/mockdb.go
- internal/pool/service.go
Tests:
- tests/pool_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs: (Out of scope for this audit per pipeline override)
Main Claims To Verify:
1. Direct connection creation penalty (overhead) is demonstrated and measurable.
2. Connection exhaustion occurs when client pool exceeds server max_connections.
3. Connection leaks (holding DB connection during external I/O) cause pool starvation.
4. Safe processing (external I/O outside DB transaction) does not leak connections and allows concurrent operations.
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race -v ./...
- go run ./cmd/demo
Primary Risks:
- Implementation may not accurately simulate real-world database connection pooling behavior.
- Tests may not cover edge cases or failure modes adequately.
- Demo output may be misleading or not representative of actual behavior.
- Concurrency safety issues in the mock driver or service.