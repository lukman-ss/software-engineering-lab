# Engineering Audit Plan

Target Lab: labs/19-database-connection-pooling

Implementation Files:
- internal/pool/mockdb.go
- internal/pool/service.go
- cmd/demo/main.go

Tests:
- tests/pool_test.go (10 test functions)

Executable/Demo:
- cmd/demo/main.go

Approved Research Inputs:
- research/01-plan.md through research/06-open-questions.md
- engineering/01-design.md (Research Status: APPROVED)

Main Claims To Verify:
1. Direct connection creation has measurable latency overhead vs pooled reuse
2. Oversized pools exhaust server max_connections and cause rejections
3. Holding a DB connection during external I/O causes pool starvation
4. Safe pattern (external call before acquiring connection) avoids starvation
5. Mock driver correctly enforces connection limits and cleanup

Commands To Run:
```
go build ./...
go test -v ./...
go test -race -v ./...
go run ./cmd/demo
```

Primary Risks:
- Timing-dependent tests (TestDirectConnectionOverhead, TestTotalCreatedPoolReuse) may be flaky on loaded machines
- TestOversizedPoolExhaustsServerConnections relies on race between 20 goroutines and driver-side lock
- MockDriver.Open uses both mutex and atomic — mixed sync primitives need verification
- errCount in TestOversizedPoolExhaustsServerConnections uses atomic via sync.Mutex inconsistently
