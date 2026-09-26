# Code Audit

## Finding 1

Location: internal/adr/parser.go:12-17
Claimed Behavior: Extract ID, title, status, supersession refs from Markdown.
Observed Implementation: Line scan + title/status/supersedes regex. Case-insensitive status/supersedes prefixes.
Assessment: PASS
Severity: LOW
Notes: Strict `# N. Title` format; documented as limitation in 02-implementation-notes.md.

## Finding 2

Location: internal/adr/parser.go:46
Claimed Behavior: Status parsing.
Observed Implementation: `strings.Title(strings.ToLower(...))` normalizes case.
Assessment: WARNING
Severity: LOW
Notes: `strings.Title` deprecated since Go 1.18; still compiles on go 1.22, vet clean. Prefer `cases.Title` if touched.

## Finding 3

Location: internal/adr/parser.go:26-68
Claimed Behavior: Error propagation on malformed input.
Observed Implementation: Returns errors for missing title, missing status, invalid status. `scanner.Err()` never checked; oversized lines would surface as misleading title-not-found.
Assessment: WARNING
Severity: LOW
Notes: No timeout/recovery relevant; pure function.

## Finding 4

Location: internal/adr/linter.go:15-25
Claimed Behavior: Duplicate ID detection.
Observed Implementation: Map insert with duplicate check, error per duplicate.
Assessment: PASS
Severity: LOW
Notes: Branch implemented but untested (see 03-test-audit.md).

## Finding 5

Location: internal/adr/linter.go:28-39
Claimed Behavior: Monotonic numbering enforcement.
Observed Implementation: Sort IDs, require exactly 1..N contiguous. Breaks after first report.
Assessment: PASS
Severity: LOW
Notes: Matches README claim.

## Finding 6

Location: internal/adr/linter.go:42-79
Claimed Behavior: Concurrent referential-integrity validation.
Observed Implementation: recordMap built before goroutines, read-only inside; errs append under mutex; wg + param capture correct.
Assessment: PASS
Severity: LOW
Notes: `go test -race -count=1 ./...` PASS. Concurrent reads of fully-built map safe.

## Finding 7

Location: internal/adr/linter.go:49-69
Claimed Behavior: Bidirectional supersession validation.
Observed Implementation: Both directions checked: Superseded requires existing SupersededBy with back-link; Supersedes requires target Superseded with forward-link. Dangling refs rejected.
Assessment: PASS
Severity: LOW
Notes: Core invariant proven by tests + demo.

## Finding 8

Location: internal/adr/linter.go:15
Claimed Behavior: Validate slice of records.
Observed Implementation: No nil-element guard; nil `*Record` panics on `r.ID`.
Assessment: WARNING
Severity: LOW
Notes: Caller-controlled input; demo/tests never pass nil. No cleanup/rollback relevant.

## Finding 9

Location: internal/adr/models.go
Claimed Behavior: Statuses Proposed/Accepted/Superseded/Deprecated/Rejected via IsValid.
Observed Implementation: Matches.
Assessment: PASS
Severity: LOW
Notes: No unnecessary complexity; stdlib only.

## Finding 10

Location: cmd/demo/main.go
Claimed Behavior: Parse + lint 3-ADR progression (Monolith -> Microservice, Rejected Event Sourcing).
Observed Implementation: Hardcoded strings parsed via adr.Parse, validated via NewLinter, exit 1 on failure.
Assessment: PASS
Severity: LOW
Notes: Output verified real (see 03-test-audit.md). No benchmark claims.
