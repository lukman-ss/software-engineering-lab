# Docs vs Code Audit

## README.md vs Implementation

### Claim 1: Naive read-modify-write causes lost update
README states naive patterns cause silent lost update anomaly.
Code: NaiveDeduct implements this. Test and demo confirm.
Status: PASS.

### Claim 2: Pessimistic Locking (SELECT ... FOR UPDATE)
README references SQL `SELECT ... FOR UPDATE`.
Code: PessimisticDeduct uses per-row `sync.Mutex` (documented in design as in-memory simulation). Demo labels output as "SELECT ... FOR UPDATE".
Status: WARNING (DOC_CODE_MISMATCH): README implies real SQL but implementation simulates with in-memory mutex. Design doc acknowledges this is a deliberate scoping decision, so acceptable for lab. Not misrepresentation if reader understands simulation.

### Claim 3: Optimistic Locking (version check with retry)
README references version check in WHERE clause with retry logic.
Code: OptimisticDeduct checks version, returns ErrOptimisticLock; service provides retry with backoff. Test confirms conflict and retry convergence.
Status: PASS.

### Claim 4: Atomic Single-Statement Operations
README references `UPDATE ... SET stock = stock - N WHERE stock >= N`.
Code: AtomicDeduct checks `stock >= qty` and decrements under one mutex hold. Comment documents the SQL mapping.
Status: PASS (simulation noted in design doc).

## Engineering Notes vs Code

### 01-design.md Expected Behavior
- Naive: lost updates. Code/test/demo confirm. PASS.
- Pessimistic: per-row lock serialization. Code confirms. PASS.
- Optimistic: version guard + safe rejection + retry convergence. Code/test/demo confirm. PASS.
- Atomic: database statement-level serialization. Code simulates via mutex. PASS (scoped).

### Implementation Decisions
- In-memory store using per-row mutexes instead of external DB. Documented. Acceptable for dependency-free lab. PASS.
- Jittered exponential backoff in retry loop. Code implements. PASS.

## Test Strategy vs Reality

### Test Strategy: "20-50 parallel goroutines per scenario"
- Tests use 20-50 goroutines. PASS.

### Test Strategy: "InitialStock - Deductions == FinalStock"
- Asserted in pessimistic, atomic, optimistic-retry, optimistic-conflict tests. PASS.

### Test Strategy: "Race detector validation"
- Ran `go test -race`. Zero warnings. PASS.

## Demo vs Code

### Demo claims
- 5 scenarios with specific output.
- Ran `go run ./cmd/demo` and recorded actual output.
- All output values match execution (no hardcoded prints). PASS.

## Content/Research alignment (reference only)
Research claims: three remediation patterns. Implementation covers all three plus naive baseline. Research implementation alignment: PASS.

## Mismatches Summary
- DOC_CODE_MISMATCH (LOW): README uses SQL syntax terminology (`SELECT ... FOR UPDATE`, `UPDATE ... SET`) while implementation is in-memory mutex simulation. Design doc transparently documents this as a scoping decision. Acceptable for lab but reader should be aware.
- No TEST_CLAIM_MISMATCH found.
- No RESEARCH_IMPLEMENTATION_MISMATCH found.
