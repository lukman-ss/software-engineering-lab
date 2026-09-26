# Docs vs Code Audit

Target Lab: labs/28-timeouts-and-deadlines

## Comparison Matrix

| Claimed Feature / Item | Documented Location | Code Location | Status |
|------------------------|---------------------|---------------|--------|
| Context Deadline Propagation | README.md / engineering/01-design.md | internal/deadline/deadline.go | MATCH |
| Exponential Backoff with Full Jitter | README.md / engineering/01-design.md | internal/retry/retry.go | MATCH |
| 3-State Circuit Breaker | README.md / engineering/01-design.md | internal/circuit/circuit.go | MATCH |
| Idempotency Key Deduplication | README.md / engineering/01-design.md | internal/idempotency/idempotency.go | MATCH |
| Runnable Demo Script | README.md:21-23 / engineering/03-execution-result.md | cmd/demo/main.go | MATCH |
| Test Execution Commands | README.md:14-18 | go test ./... / go test -race ./... | MATCH |

## Discrepancies Found

- None. All four components described in `README.md` and `engineering/01-design.md` match the source code structure, behavior, and demo output verbatim.
- No fake benchmarks or fictitious outputs detected. Demo output in `engineering/03-execution-result.md` exactly reproduces when running `go run ./cmd/demo`.
