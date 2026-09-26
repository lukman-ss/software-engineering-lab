# Docs vs Code

## README.md claims vs code

### Claim 1
README: Naive read-modify-write causes silent lost update anomaly.
Code/Tests: `NaiveDeduct` reads/stale-writes under sleep. `TestNaiveLostUpdate` expects stock != 50, demo shows 99/100. Matches.

### Claim 2
README: Pessimistic locking (`SELECT ... FOR UPDATE`) prevents conflicts.
Code: `PessimisticDeduct` holds per-row mutex. `TestPessimisticLocking` enforces stock==50. Demo shows final 50. Matches.

### Claim 3
README: Optimistic locking version check in WHERE with retry/conflict detection.
Code: `OptimisticDeduct` checks `curr.Version != p.Version`. Conflict increments `OptimisticFails`. Retry implements exponential backoff. Test `TestOptimisticLockingConflict` observes 1 succeeded, 19 rejected. Matches.

### Claim 4
README: Atomic single-statement (`UPDATE ... SET stock = stock - N WHERE stock >= N`).
Code: `AtomicDeduct` checks-and-decrements under single mutex. Matches.

### Claim 5
README: Demo compares all concurrency approaches.
Demo output printed: naive (lost update), pessimistic (exact), optimistic direct (conflicts), optimistic retry (converged), atomic (exact). Matches.

## Engineering/ notes vs code

### 01-design.md success criteria:
- Lost update demonstrated: ✅ demo/test
- Pessimistic locking maintains exact invariants: ✅ test/demo
- Optimistic conflict detection and retry convergence: ✅ test/demo
- Atomic conditional decrement guarantees consistency: ✅ test/demo
- Race detector passes: ✅ no races

### 02-implementation-notes.md choices:
- In-memory store granular per-row mutexes: ✅ `rowLocks`
- Version checking via atomic lock phase: ✅
- Jittered exponential backoff: ✅
- Zero third-party deps: ✅
- Artificial µs delays in naive: ✅ 100µs

### 03-execution-result.md recorded outputs:
- Tests: recorded naive stock 99, optimistic direct 1/19, retry 20/61, atomic 50.
- Actual run: naive 98, optimistic direct 1/19, retry 20/41, atomic 50.
- Variance: stock fluctuation within expected lost-update range, conflict counter variance due to scheduling.
- No fabricated numbers. Variance due to real concurrency — acceptable.
- Race detector: both recorded and actual show no races.

## DOC_CODE_MISMATCH
None found.

## TEST_CLAIM_MISMATCH
None found. (Tests prove claimed behavior, albeit with nondeterministic thresholds in naive/optimistic paths.)

## RESEARCH_IMPLEMENTATION_MISMATCH
(Not audited per PIPELINE OVERRIDE.)