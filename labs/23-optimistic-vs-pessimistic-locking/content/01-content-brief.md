# Content Brief

Topic:
Optimistic vs Pessimistic Locking & Atomic Updates — preventing the "lost update" anomaly in concurrent read-modify-write operations.

Target Reader:
Software developers, database engineers, and system architects familiar with basic transaction concepts who want to understand when and how to apply concurrency-control strategies on resources accessed concurrently by many clients.

Problem:
Two concurrent processes read the same value, compute a new value independently, then write back sequentially — the second transaction silently overwrites the first. This "lost update" corrupts data without raising an error, and occurs even under default isolation levels (READ COMMITTED) on PostgreSQL and Oracle. A bare transaction does not prevent it.

Core Mental Model:
Locking is a concurrency strategy. Pessimistic = prevent conflict by blocking (SELECT FOR UPDATE). Optimistic = detect conflict at commit time (version guard + 0-rows-affected). Atomic = eliminate the read-modify-write gap with a single conditional UPDATE statement. Each strategy is a trade-off between consistency, concurrency, and operational complexity. Choose by conflict frequency and data criticality.

Approved Research Status:
APPROVED (research-audit/07-verdict.md)

Approved Engineering Status:
APPROVED (engineering-audit/06-verdict.md)

Main Concepts:
1. Lost update — definition, mechanism, and reproduction under READ COMMITTED
2. Pessimistic locking — SELECT ... FOR UPDATE, row-level TX lock, held until commit/rollback
3. Optimistic locking — version/timestamp guard in WHERE clause, affected_rows = 0 as conflict signal, requires retry or 409 Conflict
4. Atomic single-statement — UPDATE ... SET stock = stock - N WHERE stock >= N, statement-level atomicity
5. Isolation level limitation — lock/version guard/atomic statement is required; isolation level alone does not prevent lost update
6. Anti-patterns — transaction alone is insufficient, holding locks during network I/O, ignoring 0-rows-affected, naive read-modify-write, unnecessary distributed locks (Redis)
7. Deadlock — auto-detected by PostgreSQL/Oracle; one transaction aborted; mitigated by consistent lock ordering and short transactions
8. MVCC — foundation enabling read consistency without blocking readers on PostgreSQL and Oracle

Verified Behaviors:
- 50 goroutines call DeductNaive(1, 1) on initial stock 100 → final stock 99 (lost update manifests; not 50)
- 50 goroutines call DeductPessimistic(2, 1) on initial stock 100 → final stock 50 (invariant maintained)
- Pessimistic insufficient-stock boundary: Seed(2), DeductPessimistic(20, 2) succeeds, then DeductPessimistic(20, 1) returns ErrInsufficientStock
- 20 goroutines call DeductOptimisticDirect(3, 1) on initial stock 100 → 1 succeeded, 19 conflicts rejected (affected-rows = 0 equivalent)
- Invariant holds for optimistic: successCount + finalStock == 100
- 20 goroutines with retry (maxRetries=10, jittered exponential backoff) → all 20 converge, final stock 80
- 50 goroutines call DeductAtomic(5, 1) on initial stock 100 → final stock 50 (lockless, statement-atomic)
- `go test -race ./...` passes with zero warnings
- `go run ./cmd/demo` runs, terminates, produces real output (not hardcoded)

Available Case Studies:
1. Demo CLI (`go run ./cmd/demo`): five side-by-side scenarios — Naive, Pessimistic, Optimistic Direct, Optimistic With Retry, Atomic — all starting from stock = 100, output recorded in engineering/03-execution-result.md
2. Test suite `tests/locking_test.go`: 6 automated tests (TestNaiveLostUpdate, TestPessimisticLocking, TestPessimisticLockingInsufficientStock, TestOptimisticLockingConflict, TestOptimisticLockingWithRetry, TestAtomicConditionalUpdate)

Warnings:
- Lab is an in-memory simulation using sync.Mutex — not a real RDBMS. Row locks are modeled, not enforced by a database engine.
- Artificial micro-delays (100µs in NaiveDeduct, 50µs in OptimisticDeduct) are inserted to make lost update/conflict reproducible.
- Deadlock detection is not simulated (single-resource lock ordering).
- MySQL docs were NOT directly accessible (HTTP 403); MySQL-related claims rely on SQL standard and secondary sources. Confidence: MEDIUM.
- Quantitative performance benchmarks are NOT included (documented as research gap). Demo elapsed time (~50ms) is illustrative, not a benchmark.
- Concurrency distribution among goroutines is nondeterministic; final stock values in naive/optimistic scenarios may vary between runs. The anomaly and invariance hold; exact figures do not.
- Lab omits error-path tests per engineering audit: ErrInvalidQuantity, ErrNotFound, Optimistic retry exhaustion (maxRetries reached), and concurrent overdraft under atomic/pessimistic are implemented but not individually tested.
- Optimistic Version integer overflow is not tested; research recommends 64-bit counter or timestamp for high-frequency systems.
- Demo elapsed-time figures are illustrative, not reproducible benchmarks (not run with fixed seed).
