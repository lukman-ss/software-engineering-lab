# Engineering Audit Plan

Target Lab: labs/20-zero-downtime-deployment
Implementation Files:
- internal/db/db.go
- internal/server/server.go
- internal/worker/worker.go
- cmd/demo/main.go
Tests:
- tests/db_test.go
- tests/server_test.go
- tests/worker_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs: research/ and research-audit/ directories
Main Claims To Verify:
1. Zero-downtime deployment via coordination of traffic routing (readiness/liveness probes)
2. Graceful termination with preStop delay handling
3. Backward-compatible database schema changes (expand/contract pattern)
4. Cooperative worker shutdown (finishing active jobs)
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in concurrent access to shared state
- Inadequate testing of edge cases (e.g., shutdown during active requests)
- Mismatch between claimed behavior and actual implementation (e.g., preStop not actually delaying)
- Fake demo output (pre-recorded logs)
