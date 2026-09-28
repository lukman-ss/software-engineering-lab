# Engineering Audit Verdict

Target Lab: labs/36-cors-and-csrf
Audit Date: Mon Sep 28 2026

## Summary

Code Files Reviewed: 0 (no .go files found)
Tests Reviewed: 0 (no test files found)
Commands Executed: None (lab not runnable due to missing go.mod and source)
Failures: Compilation impossible; test suite absent; demo absent.
Warnings: None (all issues are blocking).

## Quality Gates

Compilation: FAIL (no Go source files, no go.mod)
Tests: FAIL (no test files)
Race Detector: NOT_APPLICABLE (no code to analyze)
Demo: NOT_APPLICABLE (no cmd/demo)
Research Alignment: FAIL (no implementation to align with research)
Documentation Accuracy: FAIL (README describes components that do not exist)

## Blocking Issues
1. Missing implementation: no internal/cors, internal/csrf, internal/bank directories or .go files.
2. Missing tests: no tests/ directory.
3. Missing demo: no cmd/demo.
4. Missing go.mod: essential for Go toolchain.
5. DOC_CODE_MISMATCH: README claims existence of code that is absent.
6. RESEARCH_MISMATCH: Cannot verify research alignment without implementation.

## Non-Blocking Issues
None (all identified issues prevent any verification of claimed behavior).

## Required Revisions
1. Implement the components described in README: internal/cors (CORS middleware), internal/csrf (anti-CSRF with HMAC-SHA256 tokens, Fetch Metadata, custom headers), internal/bank (bank service with cookie auth and transfer endpoints), cmd/demo (CLI demo showing attacks vs protections).
2. Add a proper go.mod file.
3. Write a comprehensive test suite in tests/ that validates happy path, failure cases, edge cases, concurrency, and attack mitigations.
4. Ensure demo runs and shows the claimed behaviors.
5. Update README if any implementation deviates from described design (but preferably implement as described).
6. Align implementation with approved research in labs/36-cors-and-csrf/research/ (review research to ensure fidelity).

## Final Status

REJECTED