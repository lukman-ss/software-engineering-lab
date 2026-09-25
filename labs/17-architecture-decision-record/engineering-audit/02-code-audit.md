# Code Audit

## Finding 1

Location: `internal/adr/linter.go:21` (Duplicate ID Logic)
Claimed Behavior: Linter identifies duplicate ADR IDs.
Observed Implementation: Duplicate detection exists but silently overwrites `recordMap[r.ID]`. The monotonic numbering loop subsequently executes over `len(recordMap)`, validating `ids` length which hides the missing count constraint if duplicates occurred.
Assessment: WARNING
Severity: LOW
Notes: `recordMap` overwrite causes subsequent topological validation to skip the displaced record.

## Finding 2

Location: `internal/adr/linter.go:42` (Parallel Integrity Check)
Claimed Behavior: Linter validates graph in parallel, safely collecting errors.
Observed Implementation: Uses Goroutines (`go func(rec *Record)`) and appends local errors to a global slice guarded by `sync.Mutex`.
Assessment: PASS
Severity: LOW
Notes: No race conditions. Tested cleanly via `go test -race`.

## Finding 3

Location: `internal/adr/linter.go:50` (Self-Supersession / Cycle Avoidance)
Claimed Behavior: Linter structurally validates decision DAG.
Observed Implementation: It performs bidirectional validation of `StatusSuperseded` and `Supersedes` but lacks explicit checks against self-referential IDs (e.g., `ADR 1` superseding `ADR 1`) or future cycle references (e.g., `ADR 2` superseded by `ADR 1`).
Assessment: WARNING
Severity: MEDIUM
Notes: Temporal directionality constraint is implicitly assumed but not programmatically enforced.

## Finding 4

Location: `internal/adr/parser.go:37` (Status Parsing Logic)
Claimed Behavior: Parser correctly parses valid lifecycle status.
Observed Implementation: `statusRegex` strictly extracts status and optional `by <ID>`. If `Status: Accepted by 2` is supplied, it eagerly populates `SupersededBy` despite the status being `Accepted`. 
Assessment: WARNING
Severity: LOW
Notes: Semantic mismatch in parsing logic, though linter partially ignores it because `if rec.Status == StatusSuperseded` gating applies in validation.

## Finding 5

Location: `internal/adr/parser.go:38` (strings.Title Usage)
Claimed Behavior: Normalizes case for Status matching.
Observed Implementation: Uses `strings.Title` which is officially deprecated since Go 1.18.
Assessment: WARNING
Severity: LOW
Notes: Functionally sound for ASCII keywords (`Proposed`, `Accepted`) but constitutes technical debt.
