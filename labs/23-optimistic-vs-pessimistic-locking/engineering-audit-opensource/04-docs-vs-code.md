# Documentation vs Code Audit

## DOC_CODE_MISMATCH Analysis

### Claim from README.md (line 34-36):
> ### How to Run
> 
> ### Run Unit & Concurrency Tests
> ```bash
> go test -v ./...
> ```
> 
> ### Run Race Detector
> ```bash
> go test -race ./...
> ```
> 
> ### Run Demonstration
> ```bash
> go run ./cmd/demo
> ```

### Actual Code/Execution:
- `go test -v ./...` → PASS (all 6 tests pass)
- `go test -race ./...` → PASS (race detector passes, 1.125s)
- `go run ./cmd/demo` → PASS (demo runs successfully)

**Assessment**: MATCH
**Severity**: NONE

---

### Claim from README.md (lines 4-7):
> 1. **Pessimistic Locking** (`SELECT ... FOR UPDATE`)
> 2. **Optimistic Locking** (Version check in `WHERE` clause with retry logic)
> 3. **Atomic Single-Statement Operations** (`UPDATE ... SET stock = stock - N WHERE stock >= N`)

### Actual Code:
- `PessimisticDeduct` (store.go:93-116): Acquires per-product row mutex (`sync.Mutex`), then performs read-modify-write under engine mutex (`s.mu`). This simulates row-level locking.
- `OptimisticDeduct` (store.go:120-152): Reads version, sleeps, checks version match under engine mutex, then writes stock and increments version. This simulates `WHERE id = ? AND version = ?`.
- `AtomicDeduct` (store.go:154-173): Performs read-check-write under engine mutex (`s.mu`). This simulates `WHERE stock >= N` but uses a mutex, not a single atomic SQL statement.

**Assessment**: MOSTLY MATCH (with simulation notes)
**Severity**: LOW (for Atomic being mutex-guarded, not truly lockless)
**Notes**: 
- The README claims "Atomic Single-Statement Operation" and demo output line 116 says "Lockless Single Statement". 
- Actual implementation uses `s.mu` (engine-level mutex), so concurrent atomic operations on different products are serialized (real SQL would only lock the specific row). 
- This is a simulation limitation documented in the engineering notes but not emphasized in README.

---

### Claim from engineering/01-design.md (lines 6-10):
> Prove how concurrent transactions executing naive read-modify-write patterns result in silent lost updates under default isolation levels, and demonstrate the three approved remediation patterns:
> 1. Pessimistic Locking (`SELECT ... FOR UPDATE` row exclusivity)
> 2. Optimistic Locking (version guard in WHERE clause with retry / conflict detection)
> 3. Atomic Single-Statement Operations (`UPDATE ... SET stock = stock - N WHERE stock >= N`)

### Actual Demo Output (/tmp/opencode/demo_output.txt and engineering/03-execution-result.md):
- [1] Naive Read-Modify-Write: Shows stock = 99 (expected 50) → "LOST UPDATE DETECTED!" ✓
- [2] Pessimistic Locking: Shows stock = 50 → "SUCCESS - Fully Synchronized" ✓
- [3] Optimistic Locking Direct: 1 successful, 19 conflicts, stock = 99 → "SUCCESS - State Guarded, Zero Corruption" ✓
- [4] Optimistic Locking With Exponential Backoff Retry: 20 successful, stock = 80 → "SUCCESS - All retries eventually converged" ✓
- [5] Atomic Single-Statement Operation: Shows stock = 50 → "SUCCESS - Lockless Single Statement" (see note above)

**Assessment**: MATCH (core behavior verified)
**Severity**: NONE

---

### Claim from engineering/01-design.md (lines 21-25):
> ## Success Criteria
> - Automated test demonstrates lost update under unsynchronized concurrent read-modify-write.
> - Automated test validates pessimistic locking maintains exact inventory invariants.
> - Automated test validates optimistic conflict detection and retry convergence with backoff.
> - Automated test validates atomic conditional decrement guarantees consistency.
> - Race detector passes with zero race warnings.

### Actual Test Results:
- Lost update demonstrated: `TestNaiveLostUpdate` shows stock = 99 (not 50) ✓
- Pessimistic locking invariant: `TestPessimisticLocking` shows stock = 50 ✓
- Optimistic conflict detection: `TestOptimisticLockingConflict` shows conflictCount > 0 and successCount+stock=100 ✓
- Optimistic retry convergence: `TestOptimisticLockingWithRetry` shows all 20 succeed and stock=80 ✓
- Atomic invariant: `TestAtomicConditionalUpdate` shows stock = 50 ✓
- Race detector: PASS with zero warnings ✓

**Assessment**: EXACT MATCH
**Severity**: NONE

---

### Claim from engineering/02-implementation-notes.md (line 11):
> - `internal/inventory/store.go`: Storage engine simulating row-level locking, versioning, and atomic updates.

### Actual Code:
- Row-level locking: `rowLocks map[int]*sync.Mutex` + `GetRowLock` ✓
- Versioning: `Product.Version` field + version increment in `OptimisticDeduct` ✓
- Atomic updates: `AtomicDeduct` method (though mutex-guarded) ✓

**Assessment**: MATCH
**Severity**: NONE

---

### Claim from engineering/02-implementation-notes.md (lines 21-23):
> ## Implementation-Specific Choices
> - **Zero Third-Party Dependencies**: Pure Go standard library (`sync`, `sync/atomic`, `time`, `math/rand`, `testing`) to guarantee maximum portability and instantaneous test execution.

### Actual Dependencies:
- go.mod: `go 1.22` (standard library only)
- Imports: `sync`, `sync/atomic`, `time`, `math/rand`, `testing` (in test file) ✓

**Assessment**: EXACT MATCH
**Severity**: NONE

---

### Claim from engineering/02-implementation-notes.md (lines 34-39):
> ## What Is Demonstrated
> - The lost update concurrency anomaly when concurrent processes execute naive read-modify-write.
> - Total consistency guarantee with row-level pessimistic locking.
> - Detection and safe rejection of concurrent conflicts using optimistic version checking.
> - Convergence of optimistic retries using exponential backoff.
> - Lockless concurrency protection using single-statement atomic operations.

### Actual Demo Behavior:
- Lost update anomaly: Demonstrated (stock = 99 instead of 50) ✓
- Pessimistic locking consistency: Exact stock = 50 after 50 deductions ✓
- Optimistic conflict detection: Conflicts seen (19/20 in demo), no corruption (invariant holds) ✓
- Optimistic retry convergence: All 20 requests eventually succeed ✓
- Atomic operations: Exact stock = 50 after 50 deductions (simulated, not truly lockless) ⚠️

**Assessment**: MOSTLY MATCH (atomic claim slightly overstated)
**Severity**: LOW
**Notes**: The atomic operations are mutex-protected, not lockless. The demo's "Lockless Single Statement" label is aspirational (describing the SQL being simulated) rather than literal (describing the Go implementation).

---

### Claim from engineering/02-implementation-notes.md (lines 23-24):
> - **Artificial Micro-delays in Naive Method**: Inserted a small `time.Sleep` (100 µs) between read and write in `NaiveDeduct` to reliably replicate the application-level processing window that triggers lost updates in production.

### Actual Code:
- `NaiveDeduct` line 79: `time.Sleep(100 * time.Microsecond)` ✓

**Assessment**: EXACT MATCH
**Severity**: NONE

---

### DOC_CODE_MISMATCH Summary

| Claim | Source | Match? | Severity | Notes |
|---|---|---|---|---|
| Test commands | README.md | EXACT | NONE | All commands work as documented |
| Three remediation patterns | README.md / Design | MOSTLY | LOW | Atomic is mutex-guarded, not truly lockless/single-statement |
| Success criteria | Design.md | EXACT | NONE | All verified by tests and race detector |
| Implementation details | Implementation Notes | EXACT | NONE | Matches file list and standard library claim |
| What is demonstrated | Implementation Notes | MOSTLY | LOW | Atomic operations simulated, not literally lockless |
| Artificial sleep in naive method | Implementation Notes | EXACT | NONE | Matches code line 79 |
| Demo output consistency | engineering/03-execution-result.md | MOSTLY | LOW | Conflict count varies (timing-dependent), but core behavior matches |

**Overall Documentation Accuracy**: PASS with minor warnings
**Primary Issue**: Slight overstatement in describing atomic operations as "lockless" or "single-statement" when the implementation uses a mutex. This is a common and acceptable simulation technique, but the wording could be more precise.

No **DOC_CODE_MISMATCH** findings reach MEDIUM or HIGH severity. The documentation accurately reflects the code's behavior and limitations.