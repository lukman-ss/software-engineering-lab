# Docs vs Code

Target Lab: labs/19-database-connection-pooling

Actual Outputs:
- Demo: Unpooled 55.689ms vs Pooled 9µs; oversized 15/15; starvation "context deadline exceeded"
- Tests: All PASS, race clean

## Comparison Matrix

| Doc Claim Location | Corresponding Evidence | Match |
|--------------------|------------------------|-------|
| README: mockdb.go simulates max_connections and latency | Confirmed in mockdb.go Open | YES |
| README: service.go safe vs unsafe operations | Confirmed in service.go safe/unsafe | YES |
| README: tests cover direct overhead, pool exhaustion, connection leaks, concurrent | All present in tests/pool_test.go | YES |
| README: demo demonstrates overhead, rejection, starvation | demo runs, shows all | YES |
| Engineering: expected behavior (oversized pools → exhaustion, external I/O holds → starvation) | Tests deadlock, exhaustedConn, starvation | YES |

## Discrepancies
None.

## Findings
DOC vs CODE: MATCH
TEST vs CLAIM: MATCH
RESEARCH vs IMPLEMENTATION: OUT_OF_SCOPE (per pipeline)
DEMO: REAL (observed above); no fabricated benchmarks
