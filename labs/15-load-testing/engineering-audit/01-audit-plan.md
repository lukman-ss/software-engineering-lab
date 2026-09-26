# Engineering Audit Plan

Target Lab: labs/15-load-testing
Implementation Files: cmd/demo/main.go, internal/server/server.go, internal/loadtest/runner.go, internal/loadtest/metrics.go
Tests: tests/loadtest_test.go, internal/loadtest/metrics_test.go
Executable/Demo: cmd/demo
Approved Research Inputs: research/
Main Claims To Verify: Load testing runner accurately tracks latencies, handles concurrency safely, computes percentiles correctly. Server adequately mocks a constrained resource. Smoke vs Stress metrics show expected non-linear latency queuing behavior.
Commands To Run: `go test -v ./...`, `go test -race -v ./...`, `go run ./cmd/demo`
Primary Risks: Race conditions in load generator, incorrect percentile math, test flakiness.
