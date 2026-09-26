# Docs vs Code Audit

## Consistency Checklist

1. **README File Structure vs Codebase**:
   - `README.md` lists `cmd/demo`, `internal/server`, `internal/loadtest`, `tests`, `engineering/`. All exist and match structure.
   - Status: MATCH

2. **Run Instructions vs Implementation**:
   - `README.md` specifies `go run ./cmd/demo`. Verified working and produces expected benchmark output.
   - `README.md` specifies `go test -v ./...` and `go test -race ./...`. Verified working and passing.
   - Status: MATCH

3. **Research Claims vs Implementation**:
   - Research outlines load testing principles (Smoke vs Stress), tail latency degradation under queueing, and percentiles (P50, P95, P99).
   - Implementation in `cmd/demo/main.go` and `internal/loadtest` directly implements these concepts.
   - Status: MATCH

4. **Demo Execution Output Verification**:
   - Running `go run ./cmd/demo` produces real Smoke (2 VUs) vs Stress (50 VUs) results demonstrating queuing latency jump (Smoke P95 ~21ms vs Stress P95 ~794ms).
   - Status: MATCH (No fake results detected).

## Discrepancies
None detected.
