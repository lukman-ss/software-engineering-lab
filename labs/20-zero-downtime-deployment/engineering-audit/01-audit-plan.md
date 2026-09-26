# Engineering Audit Plan

Target Lab: labs/20-zero-downtime-deployment
Implementation Files: 
- internal/server/server.go
- internal/worker/worker.go
- internal/db/db.go
- cmd/demo/main.go
Tests: 
- tests/server_test.go
- tests/worker_test.go
- tests/db_test.go
Executable/Demo: go run ./cmd/demo
Approved Research Inputs: research/runs/2026-09-25-zero-downtime-deployment/05-report.md (assumed base)
Main Claims To Verify:
- DB: "Expand and Contract" pattern backward compatibility
- Server: PreStop delay, liveness/readiness probes, graceful shutdown of HTTP requests
- Worker: Graceful shutdown of queue, completes active jobs, stops pulling new jobs
Commands To Run:
- go test -v ./...
- go test -race -v ./...
- go run ./cmd/demo
Primary Risks:
- Worker channel closure drops enqueued jobs ungracefully
- Race conditions in Worker active job completion vs queue consumption
- Incomplete test coverage for worker's buffered channel during shutdown
