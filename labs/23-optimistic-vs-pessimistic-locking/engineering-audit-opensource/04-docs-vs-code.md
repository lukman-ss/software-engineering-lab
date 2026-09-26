# Docs vs Code Audit

## README.md

### Claim: "demonstrates how concurrent processes cause the silent 'lost update' anomaly"
- Code: `NaiveDeduct` (store.go:65-89) intentionally splits read-modify-write across two lock acquisitions with a 100µs delay.
- Tests: `TestNaiveLostUpdate` asserts stock != 50 after 50 concurrent deductions.
- Actual execution: stock = 98-99 (50 attempted, only 1 unit deducted).
- Assessment: **MATCH**. The claim is proven.

### Claim: "three primary remediation strategies: Pessimistic Locking, Optimistic Locking, Atomic Single-Statement Operations"
- Code: `PessimisticDeduct`, `OptimisticDeduct`/`DeductOptimisticWithRetry`, `AtomicDeduct` all implemented in service.go.
- Tests: `TestPessimisticLocking`, `TestOptimisticLockingConflict`, `TestOptimisticLockingWithRetry`, `TestAtomicConditionalUpdate`.
- Assessment: **MATCH**. All three strategies are implemented and tested.

### Claim: Run commands `go test -v ./...`, `go test -race ./...`, `go run ./cmd/demo`
- Actual execution: all three commands succeed. `go test -race` passes with zero warnings.
- Assessment: **MATCH**. Commands are accurate and reproducible.

### Claim: Directory structure listing
- README lists: cmd/demo/main.go, internal/inventory/{model,service,store}.go, tests/locking_test.go, engineering/{01-design,02-implementation-notes,03-execution-result}.md, go.mod, README.md
- Actual: all listed files exist. Additional directories (research/, research-audit/, engineering-revision/) exist but are not part of the README structure claim.
- Assessment: **MATCH** for listed structure.

---

## Engineering Docs (01-design.md, 02-implementation-notes.md, 03-execution-result.md)

### Claim: "Race detector passes with zero race warnings" (01-design.md:26)
- Actual: `go test -race ./...` passed with zero warnings.
- Assessment: **MATCH**.

### Claim: "In-memory thread-safe datastore (`Store`) simulating SQL storage engine semantics (row-level locks, versioned rows, atomic updates)" (01-design.md:29)
- Code: `Store` struct (store.go:9-19) has `sync.Mutex`, `map[int]*sync.Mutex` for row locks, `map[int]*Product` with `Version` field, and `AtomicDeduct` simulating single-statement updates.
- Assessment: **MATCH**.

### Claim: "Jittered exponential backoff implemented in optimistic retry loop" (01-design.md:53)
- Code: `DeductOptimisticWithRetry` (service.go:41): `time.Duration(1<<attempt)*time.Millisecond + time.Duration(rand.Intn(5))*time.Millisecond`
- Assessment: **MATCH**.

### Claim: "Visual CLI runner showing all 4 scenarios" (01-design.md:33)
- Actual: `cmd/demo/main.go` has 5 printed sections: [1] Naive, [2] Pessimistic, [3] Optimistic Direct, [4] Optimistic With Retry, [5] Atomic.
- Assessment: **MISMATCH**. The design doc says "4 scenarios" but the demo has 5 sections. The discrepancy arises because optimistic locking is split into two demo scenarios (direct + retry). This is a minor documentation inaccuracy.

### Claim: "Automated test validates optimistic conflict detection and retry convergence with backoff" (01-design.md:23)
- Tests: `TestOptimisticLockingConflict` validates conflict detection. `TestOptimisticLockingWithRetry` validates retry convergence (but see NOTE below).
- Assessment: **PARTIAL MATCH**. Conflict detection is proven. Retry convergence is demonstrated in practice but the test does not assert convergence (see 03-test-audit.md Finding 5).

### Claim: Execution result output (03-execution-result.md:57-90)
- Actual execution output differs in specific numbers due to timing variance:
  - Naive stock: doc shows 99, actual shows 98-99 (both demonstrate lost update correctly)
  - Optimistic conflicts: doc shows 46, actual shows 46-61 (timing-dependent)
  - Elapsed time: doc shows 49.859125ms, actual shows ~46ms (timing-dependent)
- Assessment: **MATCH (with expected variance)**. The patterns and outcomes are consistent. Concurrent programs exhibit inherent timing variance.

### Claim: "Artificial Micro-delays in Naive Method: Inserted a small time.Sleep (100 µs)" (02-implementation-notes.md:23)
- Code: store.go:79: `time.Sleep(100 * time.Microsecond)` in `NaiveDeduct`.
- Assessment: **MATCH**.

### Claim: "Artificial computation delay" (02-implementation-notes.md, store.go:132-133)
- Code: store.go:133: `time.Sleep(50 * time.Microsecond)` in `OptimisticDeduct`.
- Assessment: **MATCH**.

### Claim: "Zero Third-Party Dependencies: Pure Go standard library" (02-implementation-notes.md:22)
- go.mod: only `module` and `go 1.22` directives. No `require` block.
- Assessment: **MATCH**.

---

## Research Claims vs Implementation

### Research Finding 1: "Lost update anomaly under default isolation" (research/05-report.md:21-33)
- Implementation: `NaiveDeduct` reproduces the anomaly with unsynchronized read-modify-write across goroutines.
- Assessment: **MATCH**.

### Research Finding 2: "Pessimistic locking via SELECT ... FOR UPDATE blocks concurrent writers" (research/05-report.md:35-46)
- Implementation: `PessimisticDeduct` uses per-row `sync.Mutex` to serialize access.
- Assessment: **MATCH** (in-memory simulation of the concept).

### Research Finding 3: "Pessimistic locking trade-offs: concurrency loss, deadlock risk, hold-time sensitivity" (research/05-report.md:48-58)
- Implementation: Deadlock simulation is explicitly omitted (02-implementation-notes.md:27: "Deadlock detection simulation is omitted as lock ordering is single-resource"). No network call holding is demonstrated.
- Assessment: **SCOPE LIMITATION** (transparently documented, not a mismatch).

### Research Finding 4: "Optimistic locking detects conflicts via version guard and affected-rows check" (research/05-report.md:60-73)
- Implementation: `OptimisticDeduct` compares `curr.Version != p.Version`, returns `ErrOptimisticLock` on mismatch, increments version on success.
- Assessment: **MATCH**.

### Research Finding 5: "Selection criteria: pessimistic prevents, optimistic detects" (research/05-report.md:75-87)
- Implementation: No selection logic or guidance implemented. The lab demonstrates individual strategies, not a selection framework.
- Assessment: **NOT APPLICABLE** (lab demonstrates strategies, not selection criteria).

### Research Finding 6: "Atomic single-statement updates eliminate read-modify-write window" (research/05-report.md:89-100)
- Implementation: `AtomicDeduct` performs check-and-decrement atomically under `s.mu`.
- Assessment: **MATCH**.

### Research Finding 7-8: Isolation levels differ across databases (research/05-report.md:102-128)
- Implementation: In-memory simulation; not applicable.
- Assessment: **NOT APPLICABLE** (lab is an in-process simulation, not a real DB).

### Research Finding 9: "Common anti-patterns" (research/05-report.md:129-143)
- Implementation: The naive `DeductNaive` demonstrates anti-pattern #5 (naive read-modify-write). Anti-patterns #1-4 are not explicitly demonstrated in code.
- Assessment: **PARTIAL MATCH** (anti-pattern #5 demonstrated).

---

## Summary

| Category | Matches | Mismatches | Notes |
|---|---|---|---|
| README claims | 6 | 0 | All commands and descriptions accurate |
| Design doc claims | 5 | 1 | "4 scenarios" should be 5 (doc:01-design.md:33) |
| Execution result | 1 | 0 | Timing variance expected for concurrent code |
| Research-Findings→Implementation | 4 MATCH | 2 PARTIAL | Findings 5, 7-8, 9 not in scope for implementation lab |
