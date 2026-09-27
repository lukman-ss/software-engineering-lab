# Documentation vs Code vs Test vs Demo

## README vs Code

| README Claim | Code Implementation | Match |
|---|---|---|
| "Pessimistic Locking (SELECT ... FOR UPDATE)" | `PessimisticDeduct` uses per-row `sync.Mutex` | PASS (simulated) |
| "Optimistic Locking (Version check in WHERE clause with retry logic)" | `OptimisticDeduct` checks `curr.Version != p.Version`; `DeductOptimisticWithRetry` implements retry | PASS |
| "Atomic Single-Statement Operations (UPDATE ... SET stock = stock - N WHERE stock >= N)" | `AtomicDeduct` checks stock, decrements under lock | PASS (simulated) |
| Directory structure | Matches README (`cmd/`, `internal/inventory/`, `tests/`) | PASS |
| `go test -v ./...` / `go test -race ./...` / `go run ./cmd/demo` | All commands work | PASS |

DOC_CODE_MISMATCH: none found.

## Design vs Implementation

| Design Requirement | Implementation | Match |
|---|---|---|
| "per-row mutex" for pessimistic | `rowLocks map[int]*sync.Mutex` | PASS |
| "version incrementing" | `curr.Version++` on success | PASS |
| "exponential backoff with jitter" | `time.Duration(1<<attempt)*time.Millisecond + rand.Intn(5)` | PASS |
| "zero third-party dependencies" | go.mod has no requires | PASS |
| "artificial micro-delay (100 us)" for naive | `time.Sleep(100 * time.Microsecond)` | PASS |

RESEARCH_IMPLEMENTATION_MISMATCH: none found.

## Engineering Execution Result vs Actual Demo

Recorded (engineering/03-execution-result.md):
- Demo section 3: "Rejected Conflicts: 19", Final Stock 99
- Demo section 4: "Total Attempted Conflicts Retried: 61", "Elapsed Time: 49.859125ms"

Actual run (this audit):
- Section 3: "Rejected Conflicts: 19", Final Stock 99
- Section 4: "Total Attempted Conflicts Retried: 45", "Elapsed Time: 23.075708ms"

TEST_CLAIM_MISMATCH: Non-deterministic values in recorded execution result
(OptimisticFails count and elapsed time differ between recorded and actual run).
This is expected for concurrent code but the recorded values are not labeled as
approximate/non-deterministic.

Assessment: WARNING
Severity: LOW
Notes: The demo produces different numbers across runs for optimistic conflict counts and timing. Engineering record should note these are non-deterministic.

## README Completeness

README does not mention:
- `ErrOptimisticLock`, `ErrInsufficientStock`, `ErrInvalidQuantity` error types
- The 5 demo sections include a 5th atomic operation not enumerated in the README's
  "three primary remediation strategies" (atomic is the 3rd strategy, so this is consistent)

DOC_CODE_MISMATCH: none
