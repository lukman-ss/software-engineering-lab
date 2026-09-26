# Content Audit Plan

Target Lab: labs/20-zero-downtime-deployment
Audit Date: 2026-09-26
Scope: Content only (01-content-brief.md through 06-source-map.md)

## Files Reviewed

- `content/01-content-brief.md`
- `content/02-master-draft.md`
- `content/03-code-snippets.md`
- `content/04-diagrams.md`
- `content/05-key-takeaways.md`
- `content/06-source-map.md`

## Reference Materials (Cross-Check)

- Research: `research/runs/2026-09-26-zero-downtime-deployment/05-report.md`, `03-evidence.md`, `04-contradictions.md`, `06-open-questions.md`
- Research Audit: `research-audit/07-verdict.md`, `03-claim-audit.md`, `04-contradictions.md`, `06-gaps.md`
- Engineering: `engineering/01-design.md`, `engineering/02-implementation-notes.md`
- Engineering Audit: `engineering-audit/06-verdict.md`, `02-code-audit.md`, `03-test-audit.md`, `04-docs-vs-code.md`, `05-gaps.md`
- Implementation: `internal/server/server.go`, `internal/worker/worker.go`, `internal/db/db.go`, `cmd/demo/main.go`
- Tests: `tests/server_test.go`, `tests/worker_test.go`, `tests/db_test.go`

## Audit Criteria

1. Accuracy against implementation (code-level verification)
2. Accuracy against research claims (research-level verification)
3. Clarity, formatting, and readability
4. Absence of hallucinated facts
5. Absence of unjustified platform-specific bias
6. Completeness of coverage against content brief
