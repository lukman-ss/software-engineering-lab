# Engineering Audit Plan

Target Lab: labs/19-database-connection-pooling
Implementation Files:
- internal/pool/mockdb.go
- internal/pool/service.go
- cmd/demo/main.go
Tests:
- tests/pool_test.go
Executable/Demo:
- cmd/demo/main.go
Approved Research Inputs:
- research/05-report.md
- engineering/01-design.md
- engineering/02-implementation-notes.md
- engineering/03-execution-result.md
Main Claims To Verify:
1. Connection reuse avoids expensive connection establishment latency.
2. Sizing the connection pool larger than backend max_connections causes connection errors/rejections.
3. Holding database connections during non-database operations (leaks) causes pool starvation and timeouts.
4. Concurrent requests with safely bounded connection usage complete successfully.
Commands To Run:
- `go test -v ./...`
- `go test -race -v ./...`
- `go run ./cmd/demo`
Primary Risks:
- Race conditions around MockDriver state tracking (mix of atomics and mutex).
- Flaky tests due to tight timing and time.Sleep durations.
- Overclaiming performance knee/throughput degradation simulation when only simple rejection is implemented.
