# Gap Analysis

## GAP_1: MISSING_IMPLEMENTATION
Location: labs/36-cors-and-csrf/internal/
Claimed: `internal/cors` (spec-compliant CORS middleware), `internal/csrf` (HMAC-SHA256 signed tokens, Fetch Metadata, custom header middleware), `internal/bank` (bank service with cookie-authenticated balance inquiries, vulnerable/protected transfer endpoints)
Observed: No `internal/` directory exists.
Severity: CRITICAL

## GAP_2: MISSING_DEMO
Location: labs/36-cors-and-csrf/cmd/demo
Claimed: Runnable CLI program showcasing attacks against vulnerable vs. protected configurations.
Observed: No `cmd/` directory exists.
Severity: CRITICAL

## GAP_3: MISSING_TESTS
Location: labs/36-cors-and-csrf/tests/
Claimed: Integration test suite verifying cross-origin requests, preflight, race safety, and attack mitigation.
Observed: No `tests/` directory exists.
Severity: CRITICAL

## GAP_4: MISSING_GO_MOD
Location: labs/36-cors-and-csrf/go.mod
Claimed: Go module file for dependency management and build.
Observed: Not present on disk.
Severity: CRITICAL

## GAP_5: DOC_CODE_MISMATCH
Location: README.md vs. filesystem
Claimed: Implementation matches README description.
Observed: None of the described components exist.
Severity: CRITICAL

## GAP_6: RESEARCH_MISMATCH
Location: README claim vs. research/
Claim: Implementation matches approved research.
Observed: No implementation to compare against research.
Severity: CRITICAL