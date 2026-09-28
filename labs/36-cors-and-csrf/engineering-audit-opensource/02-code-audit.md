# Code Audit

## Finding 1

Location: labs/36-cors-and-csrf/ (entire implementation tree)
Claimed Behavior: README describes `internal/cors`, `internal/csrf`, `internal/bank`, `cmd/demo`, and `tests` directories containing spec-compliant CORS middleware, anti-CSRF mechanisms, a bank service, and an integration test suite.
Observed Implementation: None of these directories or files exist. The lab directory contains only `engineering-audit-opensource/` and `research/`. No `.go` files, no `cmd/`, no `internal/`, no `tests/`, no `go.mod` (the `go.mod` referenced in the directory listing is not present on disk).
Assessment: FAIL
Severity: CRITICAL
Notes: The entire implementation is missing. Nothing to audit for correctness, state transitions, failure handling, timeout behavior, recovery, concurrency, cleanup, or error propagation.

## Finding 2

Location: labs/36-cors-and-csrf/go.mod
Claimed Behavior: Go module file expected for a Go lab.
Observed Implementation: Not present on disk.
Assessment: FAIL
Severity: CRITICAL
Notes: Without `go.mod`, `go test ./...` and `go run ./cmd/demo` cannot execute.

## Finding 3

Location: labs/36-cors-and-csrf/cmd/demo
Claimed Behavior: Runnable CLI program showcasing attacks against vulnerable vs. protected configurations.
Observed Implementation: Directory does not exist.
Assessment: FAIL
Severity: CRITICAL
Notes: No demo executable or source exists.

## Finding 4

Location: labs/36-cors-and-csrf/tests
Claimed Behavior: Integration test suite verifying cross-origin requests, preflight, race safety, and attack mitigation.
Observed Implementation: Directory does not exist.
Assessment: FAIL
Severity: CRITICAL
Notes: No tests exist to verify claimed behavior.