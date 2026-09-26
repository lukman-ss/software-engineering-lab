# Test Audit

## Coverage Evaluation

### Happy Path Coverage
- `TestParse_Valid`: Proves basic parsing of title, status, and metadata.
- `TestParse_Superseded` & `TestParse_Supersedes`: Proves relationship extraction.
- `TestLinter_ValidSequence`: Proves end-to-end validation of a 3-decision sequence.

### Failure Path Coverage
- `TestParse_Invalid`: Proves rejection of missing title, missing status, and invalid lifecycle statuses.
- `TestLinter_BrokenReferences`:
  - `superseded_points_to_non-existent`: Validates dangling reference detection.
  - `supersedes_points_to_non-existent`: Validates dangling predecessor detection.
  - `mismatched_supersession_link`: Validates one-sided supersession rejection.
  - `non-monotonic_numbering`: Validates gaps in sequence detection.

### Edge Cases and Gaps
- Duplicate ADR ID: Handled in `linter.go:21-23`, but not exercised in `linter_test.go`.
- Large batch concurrency stress test: Validation runs concurrently across all records, but test suite only exercises up to 3 records.

## Execution Verification

### 1. `go test -v ./...`
```text
?   	labs/17-architecture-decision-record/cmd/demo	[no test files]
?   	labs/17-architecture-decision-record/internal/adr	[no test files]
=== RUN   TestLinter_ValidSequence
--- PASS: TestLinter_ValidSequence (0.00s)
=== RUN   TestLinter_BrokenReferences
=== RUN   TestLinter_BrokenReferences/superseded_points_to_non-existent
=== RUN   TestLinter_BrokenReferences/supersedes_points_to_non-existent
=== RUN   TestLinter_BrokenReferences/mismatched_supersession_link
=== RUN   TestLinter_BrokenReferences/non-monotonic_numbering
--- PASS: TestLinter_BrokenReferences (0.00s)
    --- PASS: TestLinter_BrokenReferences/superseded_points_to_non-existent (0.00s)
    --- PASS: TestLinter_BrokenReferences/supersedes_points_to_non-existent (0.00s)
    --- PASS: TestLinter_BrokenReferences/mismatched_supersession_link (0.00s)
    --- PASS: TestLinter_BrokenReferences/non-monotonic_numbering (0.00s)
=== RUN   TestParse_Valid
--- PASS: TestParse_Valid (0.00s)
=== RUN   TestParse_Superseded
--- PASS: TestParse_Superseded (0.00s)
=== RUN   TestParse_Supersedes
--- PASS: TestParse_Supersedes (0.00s)
=== RUN   TestParse_Invalid
=== RUN   TestParse_Invalid/missing_title
=== RUN   TestParse_Invalid/missing_status
=== RUN   TestParse_Invalid/invalid_status
--- PASS: TestParse_Invalid (0.00s)
    --- PASS: TestParse_Invalid/missing_title (0.00s)
    --- PASS: TestParse_Invalid/missing_status (0.00s)
    --- PASS: TestParse_Invalid/invalid_status (0.00s)
PASS
ok  	labs/17-architecture-decision-record/tests	0.269s
```

### 2. `go test -race ./...`
```text
?   	labs/17-architecture-decision-record/cmd/demo	[no test files]
?   	labs/17-architecture-decision-record/internal/adr	[no test files]
PASS
ok  	labs/17-architecture-decision-record/tests	0.162s
```
Result: PASS without data races.

### 3. `go run ./cmd/demo`
```text
=== Architecture Decision Record (ADR) Lab ===
Successfully parsed 3 ADRs:
  - [0001] Use Modular Monolith for Core SaaS ERP        | Status: Superseded -> Superseded by ADR 2
  - [0002] Extract Notification Service to Microservice  | Status: Accepted   -> Supersedes ADR 1
  - [0003] Reject Event Sourcing for Order Management    | Status: Rejected  

Running ADR Integrity Linter...
Integrity check passed! Decision lineage and lifecycle invariants are intact.
```
Result: PASS with expected output.
