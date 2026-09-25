# Gap Analysis

Target Lab: labs/18-deadlock

## Gaps

### 1. MISSING_EDGE_CASE
- **Location:** `internal/transfer/transfer.go`
- **Details:** `TransferOrdered` and `TransferNaive` do not validate if `from == to` (self-transfer). In `TransferOrdered`, attempting a self-transfer results in attempting to lock the same channel twice, causing a self-inflicted deadlock.
- **Severity:** LOW. This is a concurrency-focused lab rather than a banking business logic lab.

### 2. MISSING_TEST
- **Location:** `tests/transfer_test.go`
- **Details:** Test suite focuses entirely on concurrency primitives (deadlocks, recovery, durations). Standard unit testing for basic assertions (e.g., negative balances, insufficient funds, exact numeric invariants under normal operation) is largely omitted.
- **Severity:** LOW. Concurrency coverage is excellent, which is the primary objective of the lab.

## Impact
None of the identified gaps compromise the core claims of the research or the implementation's ability to demonstrate the intended concepts.
