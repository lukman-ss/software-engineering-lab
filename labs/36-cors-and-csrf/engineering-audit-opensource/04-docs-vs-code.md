# Docs vs Code

## Mismatch 1: DOC_CODE_MISMATCH / IMPLEMENTATION_OVERCLAIM

Location: README.md vs. filesystem
Claimed Implementation: README documents `internal/cors`, `internal/csrf`, `internal/bank`, `cmd/demo`, and `tests`.
Code: None of these paths exist on disk.
Assessment: WARNING→FAIL
Severity: CRITICAL
Notes: README describes code that does not exist.

## Mismatch 2: RESEARCH_MISMATCH (pending)

Location: README.md claims vs. labs/36-cors-and-csrf/research/
Claim: Implementation matches approved research.
Code: No implementation present to align to research.
Assessment: FAIL
Severity: CRITICAL
Notes: Cannot verify alignment because there is no implementation.
