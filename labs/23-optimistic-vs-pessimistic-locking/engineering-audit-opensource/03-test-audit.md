# Test Audit

## Finding 1
Location: tests/locking_test.go:10-38 (`TestNaiveLostUpdate`)
Coverage: happy path / lost-update demonstration.
Assessment: WARNING
Severity: MEDIUM
Notes: Assertion only checks `p.Stock == 50`. With 50 goroutines deducting 1 from stock 100, any final value other than 50 passes — including 99 (49 silent overwrites). Assertion is too permissive: it cannot distinguish a true lost update from other concurrency artifacts. Stronger check would verify `p.Stock > 50` or reconcile against successful-return count. Does NOT verify that all 50 calls returned nil (success) while state was corrupted.

## Finding 2
Location: tests/locking_test.go:40-67 (`TestPessimisticLocking`)
Coverage: happy path — 50 concurrent deductions.
Assessment: PASS
Severity: LOW
Notes: Asserts exact final stock 50. Strong invariant.

## Finding 3
Location: tests/locking_test.go:69-83 (`TestPessimisticLockingInsufficientStock`)
Coverage: failure path — insufficient stock.
Assessment: PASS
Severity: LOW
Notes: Good edge-case test. Deducts exactly available amount then tries one more.

## Finding 4
Location: tests/locking_test.go:85-121 (`TestOptimisticLockingConflict`)
Coverage: conflict detection, no retry.
Assessment: PASS
Severity: LOW
Notes: Verifies invariant `successCount + p.Stock == 100` and asserts at least one conflict. Well-designed.

## Finding 5
Location: tests/locking_test.go:123-145 (`TestOptimisticLockingWithRetry`)
Coverage: retry convergence.
Assessment: WARNING
Severity: MEDIUM
Notes: Asserts `p.Stock == 100 - store.Optimistically` — uses the implementation's own counter as source of truth. If `Optimistically` counter were incremented on a non-applied update, the test would still pass. Better assertion: independent invariant check (success + stock == initial).

## Finding 6
Location: tests/locking_test.go:147-171 (`TestAtomicConditionalUpdate`)
Coverage: atomic decrement happy path — 50 concurrent.
Assessment: PASS
Severity: LOW
Notes: Strong assertion: exact stock 50.

## Missing Tests
Location: tests/locking_test.go
Coverage Gap: `ErrInvalidQuantity` (qty <= 0) never tested.
Assessment: MISSING_TEST
Severity: MEDIUM
Notes: All four deduct methods guard against `qty <= 0` but no test exercises this path.

## Missing Tests
Location: tests/locking_test.go
Coverage Gap: `ErrNotFound` (deduct non-seeded product) never tested.
Assessment: MISSING_TEST
Severity: MEDIUM
Notes: `Get`, `NaiveDeduct`, `PessimisticDeduct`, `OptimisticDeduct`, `AtomicDeduct` all return `ErrNotFound` but path untested.

## Missing Tests
Location: tests/locking_test.go
Coverage Gap: Atomic decrement with insufficient stock.
Assessment: MISSING_TEST
Severity: LOW
Notes: Atomic method guards stock but no test covers the insufficient path.

## Race Detection
Command: `go test -race ./...`
Result: ok — zero race warnings.
Assessment: PASS
