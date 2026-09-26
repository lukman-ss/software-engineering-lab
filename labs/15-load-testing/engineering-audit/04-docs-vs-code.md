# Docs vs Code Audit

Target Lab: labs/15-load-testing

## Claim vs Implementation

### Claim 1: "Average response time conceals tail latency spikes; percentiles (P95, P99) are necessary to uncover degradation."
- Source: `engineering/01-design.md`, `research/05-report.md`
- Implementation: In `cmd/demo/main.go`, Smoke Test produces Average: ~21ms, P95: ~21ms. Stress Test produces Average: ~200ms, P95: ~211ms.
- Assessment: `DOC_CODE_MISMATCH` / `IMPLEMENTATION_OVERCLAIM`. Average does not conceal the degradation, it tracks P95 almost 1:1 because of constant simulated query time under a closed workload model.

### Claim 2: "Semaphore pattern (buffered channel) used in the server handler to simulate database connection pool bottlenecks, accurately mimicking tail latency growth when saturated."
- Source: `engineering/02-implementation-notes.md`
- Implementation: Present in `internal/server/server.go`. Matches the documentation.
- Assessment: PASS.

### Claim 3: "Demonstration clearly contrasts Smoke test metrics against Stress test metrics."
- Source: `engineering/01-design.md`
- Implementation: The demo shows a clear 10x degradation in response time from smoke to stress, demonstrating resource saturation clearly, but fails to show tail skewness.
- Assessment: PASS (for resource saturation), WARNING (for tail skewness).
