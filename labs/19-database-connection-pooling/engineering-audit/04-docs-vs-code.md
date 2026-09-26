# Docs vs Code Audit

## Comparison Points

1. **README vs Code**
   - README claims to demonstrate connection pool behavior, failure modes, and connection leaks.
   - Demo and tests successfully demonstrate this. No mismatch.
   - Commands listed in README (`go test -v ./...`, `go test -race -v ./...`) work, though the race detector currently fails.

2. **Engineering Notes vs Code**
   - `01-design.md` states: "Test validates throughput degradation or max connection enforcement".
   - Implementation only validates max connection enforcement (returning `ErrServerOverloaded`). No latency knee/throughput degradation logic exists in `mockdb.go`.
   - Mismatch type: `IMPLEMENTATION_OVERCLAIM` (Minor: the design says "or", but it is a weak demonstration of the overall research claim about throughput).

3. **Research Claims vs Implementation**
   - **Finding 4 (Performance Knee):** Not implemented.
   - **Finding 9 (Connection Lifecycle Best Practices):** Perfectly implemented in `service.go`.
   - **Finding 11 (Pool-locking Deadlock):** Not demonstrated in code. The lab focuses purely on resource sizing and network leaks.
   - **Finding 1 (Connection Exhaustion):** Implemented accurately. Hard limits trigger explicit errors.
   - Mismatch type: `RESEARCH_IMPLEMENTATION_MISMATCH` (The research is broader than the lab, which is acceptable, but the deadlock edge case would have been a strong addition).

## Findings
- **DOC_CODE_MISMATCH**: None.
- **TEST_CLAIM_MISMATCH**: None.
- **RESEARCH_IMPLEMENTATION_MISMATCH**: Implementation focuses heavily on exhaustion and leaks, omitting the performance degradation and deadlock phenomena found in research.
