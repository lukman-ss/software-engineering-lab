# Source Map

## 1. Introduction: Connection Overhead

Research:
- `research/05-report.md` — Finding 1 (Connection Exhaustion)
- `research/03-evidence.md` — Evidence 1–3 (max_connections semantics, resource consumption)

Implementation:
- `internal/pool/mockdb.go` — `NewMockDriver()` with `connectDelay` simulating TCP/TLS handshake
- `cmd/demo/main.go` — `demoDirectOverhead()` comparing unpooled vs pooled execution time

Tests:
- `tests/pool_test.go` — `TestDirectConnectionOverhead`, `TestTotalCreatedPoolReuse`

## 2. Pool Sizing Formula

Research:
- `research/05-report.md` — Finding 2 (Optimal Pool Size Formula), Finding 4 (PostgreSQL Performance Knee)
- `research/02-sources.md` — Source 2 (PostgreSQL Wiki), Source 3 (HikariCP Wiki)
- `research/03-evidence.md` — Evidence 4 (pool sizing formula), Evidence 7 (SSD considerations)

Implementation:
- `internal/pool/mockdb.go` — `maxConnections` parameter enforces server-side limit
- `cmd/demo/main.go` — `demoOversizedPool()` demonstrating rejection when client pool exceeds server limit

Tests:
- `tests/pool_test.go` — `TestOversizedPoolExhaustsServerConnections`, `TestPoolLockingDeadlock`

## 3. Connection Leaks & Pool Starvation

Research:
- `research/05-report.md` — Finding 3 (Connection Leaks Cause Invisible Failures), Finding 9 (Connection Lifecycle Best Practices)
- `research/02-sources.md` — Source 4 (HikariCP Configuration — `leakDetectionThreshold`)
- `research/03-evidence.md` — Evidence 8 (leakDetectionThreshold default and minimum)

Implementation:
- `internal/pool/service.go` — `ProcessOrderUnsafeLeak()` vs `ProcessOrderSafe()`
- `cmd/demo/main.go` — `demoConnectionLeak()` demonstrating starvation via `ProcessOrderSafe` timeout

Tests:
- `tests/pool_test.go` — `TestConnectionStarvationDueToLeak`, `TestExternalCallErrorPropagation`, `TestUnsafeLeakExecContextFailure`

## 4. Exhaustion via Horizontal Scale

Research:
- `research/05-report.md` — Finding 7 (Pool Sizing Must Account for Global Deployment), Finding 6 (Cloud Providers Recommend External Pooling)
- `research/03-evidence.md` — Evidence 12–13 (Azure reserved connections and pool recommendations), Evidence 14 (AWS RDS formula)

Implementation:
- `cmd/demo/main.go` — `demoOversizedPool()` showing 30 concurrent attempts against 15-server limit

Tests:
- `tests/pool_test.go` — `TestOversizedPoolExhaustsServerConnections`

## 5. Monitoring & Observability

Research:
- `research/05-report.md` — Finding 10 (Monitoring Metrics Beyond CPU/Memory), Finding 11 (Pool-Locking Deadlock Prevention Formula)
- `research/02-sources.md` — Source 10 (Dropwizard Metrics), Source 11 (PostgreSQL Monitoring Stats)
- `research/03-evidence.md` — Evidence 18 (HikariCP metrics), Evidence 22 (PreparedStatement caching anti-pattern)

Implementation:
- `internal/pool/mockdb.go` — `ActiveConnections()`, `TotalCreated()` counters for monitoring

Tests:
- `tests/pool_test.go` — `TestMockConnDoubleClose` (verifying resource tracking), `TestSafeProcessingConcurrently`

## 6. Cloud Provider Recommendations

Research:
- `research/05-report.md` — Finding 6 (Cloud Providers Recommend External Pooling)
- `research/02-sources.md` — Source 7 (Azure Limits), Source 8 (Google Cloud SQL), Source 9 (AWS RDS Connection Limits)
- `research/03-evidence.md` — Evidence 12–13, 14, 15

Implementation: Not directly implemented (conceptual guidance).

## 7. PgBouncer Pool Modes

Research:
- `research/05-report.md` — Finding 8 (PgBouncer Pool Modes Have Different Feature Compatibility)
- `research/02-sources.md` — Source 5 (PgBouncer Configuration), Source 6 (PgBouncer Features)
- `research/03-evidence.md` — Evidence 10–11 (PgBouncer defaults and pooling modes)

Implementation: Not directly implemented (noted as known limitation in engineering docs).

## 8. Oracle 50x Improvement Case Study

Research:
- `research/05-report.md` — Finding 5 (Oracle Real-World Performance Demonstrated 50x Improvement)
- `research/02-sources.md` — Source 3 (HikariCP Wiki)
- `research/03-evidence.md` — Evidence 5 (Oracle video demonstration)

Implementation: Not directly implemented (cited as external case study).

## 9. Implementation Details & Mock Architecture

Engineering Design:
- `engineering/01-design.md` — Architecture overview, components, test strategy
- `engineering/02-implementation-notes.md` — Core design decisions, known limitations, trade-offs
- `engineering/03-execution-result.md` — Build, test, race detector, and demo output

Implementation:
- `internal/pool/mockdb.go` — Full `MockDriver`, `MockConnector`, `mockConn`, `mockStmt`, `mockTx`, `mockRows`
- `internal/pool/service.go` — `OrderService` with `ProcessOrderSafe` and `ProcessOrderUnsafeLeak`
- `tests/pool_test.go` — 10 test functions covering all core scenarios
- `cmd/demo/main.go` — Interactive demo with three demonstration functions

## 10. Research Limitations & Open Questions

Research:
- `research/04-contradictions.md` — Analysis of apparent contradictions (SSD formula, PgBouncer modes, max_connections calculation)
- `research/06-open-questions.md` — Unanswered questions, weak evidence, claims needing deeper research
- `research/05-report.md` — Limitations section

Engineering:
- `engineering-audit/05-gaps.md` — Gap analysis (GAP-001: go.mod version metadata, GAP-002/003: minor test coverage gaps)
