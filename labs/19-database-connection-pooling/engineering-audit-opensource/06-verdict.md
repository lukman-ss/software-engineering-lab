# Engineering Audit Verdict

Target Lab:
  /Users/tthi/Documents/LUKMAN/software-engineering-lab/labs/19-database-connection-pooling

Audit Date:
  2026-09-26

## Summary

Code Files Reviewed:
  - go.mod
  - internal/pool/mockdb.go
  - internal/pool/service.go
  - cmd/demo/main.go
  - tests/pool_test.go
  - README.md
  - engineering/01-design.md, 02-implementation-notes.md, 03-execution-result.md

Tests Reviewed:
  tests/pool_test.go (4 tests)

Commands Executed:
  - go version                                  → go1.26.7 darwin/arm64
  - go build ./...                              → exit 0, no output (Success)
  - go vet ./...                                → exit 0, no output (Success)
  - go test -v ./...                            → 4 PASS, 0 FAIL (exit 0)
  - go test -race -v ./...                      → 4 PASS, 0 DATA RACES (exit 0)
  - go run ./cmd/demo                           → exit 0, output reproduced (see below)

Demo re-execution output (verbatim):
  --- Database Connection Pooling Demo ---

  1. Direct Connection Overhead Penalty
  Unpooled (5 requests): 54.242792ms
  Pooled (5 requests): 2.833µs

  2. Oversized Pool Exhausting Server Connections
  Client attempted: 30, Succeeded: 15, Server Rejected: 15

  3. Connection Leak Starving Pool
  Starting 2 unsafe orders (holding connection during slow external IO)
  Attempting 3rd order with safe flow and short timeout...
  Order 3 Failed: context deadline exceeded

Failures:
  None — all quality gates executed and passed.

Warnings:
  1. LOW: MockDriver max-connection cap is enforced with a check-then-act (TOCTOU)
     gap (atomic increment outside the mutex). Race detector is clean; tests tolerate it.
  2. LOW: ErrAcquireTimeout (service.go) is declared but unused/untested. The
     timeout behavior itself is tested and works.
  3. LOW: No test pins mockRows/mockTx rollback paths or post-error connection reuse.

## Quality Gates

Compilation:            PASS  (go build ./... exit 0; go vet ./... exit 0)
Tests:                 PASS  (4/4 PASS, go test -v ./...)
Race Detector:          PASS  (4/4 PASS, 0 races, go test -race -v ./...)
Demo:                  PASS  (go run ./cmd/demo exit 0, output matches claimed structure)
Research Alignment:     PASS  (implementation realizes the documented safe-vs-unsafe
                              narrative; demo reproduces claims)
Documentation Accuracy: PASS  (README lists exact components and reproducible commands;
                              README claims verified by re-execution)

## Blocking Issues
None.

## Non-Blocking Issues
1. LOW — MockDriver cap TOCTOU (see above). Move atomic.AddInt32 inside d.mu
   if exact server-saturation counts must be enforced.
2. LOW — ErrAcquireTimeout sentinel unused. Either return it on acquire timeout
   or remove it.
3. LOW — Minor test-coverage gaps (rollback path, post-error reuse, exact cap).

## Required Revisions
None for approval. Optional improvements (non-blocking):
1. Enforce MockDriver cap atomically (move increment under d.mu).
2. Wire ErrAcquireTimeout or delete it.
3. Add a test asserting connection reuse after an error path.

## Final Status

APPROVED_WITH_WARNINGS

The lab compiles, all 4 tests pass, the race detector reports zero data races, and the
demo reproduces the documented output exactly (integer-valued behavioral claims match
verbatim; only wall-clock sub-millisecond timings differ, which is expected). The
implementation faithfully realizes the claimed "safe vs unsafe connection use"
behavior including pool exhaustion and leak-induced starvation. Documentation matches
code and the verified execution results. Remaining warnings are low-severity pedagogical
completeness gaps in a teaching mock and do not undermine the claimed behavior.
