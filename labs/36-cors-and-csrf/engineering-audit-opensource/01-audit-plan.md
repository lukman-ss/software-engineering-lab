# Engineering Audit Plan

Target Lab: labs/36-cors-and-csrf
Workspace Root: /Users/tthi/Documents/LUKMAN/software-engineering-lab
Audit Scope: Implementation and tests only (no audit of research/content).

Implementation Files:
- NONE. The lab directory contains only an empty `research/` directory.
- `find labs/36-cors-and-csrf -type f` returns no source files, tests, README, or demo.

Tests:
- NONE. No test files present.

Executable/Demo:
- NONE. No `cmd/`, `main.go`, or runnable artifact present.

Approved Research Inputs:
- `research/` directory exists but is empty (no files). No implementation to align against.

Main Claims To Verify:
- CORS allowlist/denylist behavior.
- CSRF token validation on state-changing requests.
- SameSite cookie enforcement.
- Origin/A referer validation.

Commands To Run:
- `go test ./...` in the lab (expect failure: no Go module / no source).
- `go run ./cmd/demo` in the lab (expect failure: no entry point).
- `go build ./...` in the lab (expect failure).

Primary Risks:
- MISSING_IMPLEMENTATION: no engineering artifacts exist to audit.
- FAKE_DEMO / UNVERIFIED_RESULT: cannot be demonstrated because there is no code.
- RESEARCH_IMPLEMENTATION_MISMATCH: not applicable; no implementation to compare.

Notes:
- This audit is limited to implementation/tests by instruction.
- The lab appears to be an empty scaffold with only an empty research directory.
- The verdict will be REJECTED at the implementation level (cannot pass quality gates) but recorded as missing scaffold.
