## Revision 1

Audit Issue: Gap 1
Severity: LOW
Files Changed: 
- `internal/adr/linter.go`
- `tests/linter_test.go`
Action: Added explicit rejection of self-supersession in the linter and added a corresponding test case.
Verification: Ran tests to confirm the new error is generated when `Supersedes` or `SupersededBy` equals the ADR's own ID.
Status: RESOLVED