# Engineering Audit Plan

Target Lab: labs/20-zero-downtime-deployment
Implementation Files:
- internal/db/db.go
- internal/server/server.go
- internal/worker/worker.go
- cmd/demo/main.go
- go.mod
Tests:
- tests/db_test.go
- tests/server_test.go
- tests/worker_test.go
Executable/Demo: cmd/demo/main.go
Approved Research Inputs:
- Zero-Downtime Deployment patterns from research
- Expand and Contract database schema evolution
- Health probe distinction (liveness vs readiness)
- Graceful connection draining mechanisms
- Background worker cooperative termination
Main Claims To Verify:
1. HTTP server readiness check rejects traffic when unready, accepts when ready
2. Graceful HTTP shutdown allows active in-flight requests to complete  
3. PreStop hook waits for configured delay before closing listeners
4. Worker finishes ongoing job upon shutdown signal
5. Expand and Contract database reads/writes work across legacy and modern formats
Commands To Run:
- go build ./...
- go test -v ./...
- go test -race ./...
- go run ./cmd/demo
Primary Risks:
- Race conditions in server request counting during shutdown
- Incorrect ordering of shutdown steps (preStop vs listener close)
- Missing edge cases in Expand/Contract fallback logic
- Deadlock in worker channel operations
- Panic from double-close of worker channel on multiple Stop() calls