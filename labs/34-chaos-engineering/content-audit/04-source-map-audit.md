# Source Map Audit

Source: `content/06-source-map.md`
Purpose: verify every cross-reference resolves to the cited location/content.

## Research references
- `research/01-plan.md`, `research/02-sources.md`, `research/03-evidence.md`, `research/05-report.md`, `research/06-open-questions.md`: paths exist per glob. VERIFIED.

## Implementation references
- `internal/fault/injector.go`: exists. VERIFIED.
- `internal/circuitbreaker/circuitbreaker.go`: exists. VERIFIED.
- `internal/monitor/monitor.go`: exists. VERIFIED.
- `internal/experiment/runner.go`: exists. VERIFIED.
- `cmd/demo/main.go`: exists. VERIFIED.

## Test reference
- `tests/chaos_test.go`: exists. VERIFIED.

## Master Draft → research trace
- Problem → `research/01-plan.md:4-12`, `research/05-report.md:7-8`: not cross-checked line-for-line (research excluded); paths plausible.
- Why This Matters → `research/03-evidence.md:5-9`: plausible.
- Mental Model → `research/03-evidence.md:19-22, 35-38, 51-55`: plausible.

## How It Works → engineering trace
- `engineering/01-design.md:26-36`: exists; architecture section. VERIFIED (file exists).
- `engineering/02-implementation-notes.md:12-18`: design decisions section. VERIFIED.

## What the Tests Prove → engineering trace
- `engineering/03-execution-result.md:19-32` (test results), `:39-42` (race): exists. VERIFIED.

## Demo Behavior → engineering trace
- `engineering/03-execution-result.md:49-90` (demo output): VERIFIED — demo output is at lines 50-90, matches cited range.

## Recovery / Rollback reference
- `internal/experiment/runner.go:91-97` (terminate): VERIFIED — `terminate` is at lines 91-97. ✓
- `cmd/demo/main.go:127-132` (traffic after experiment): VERIFIED — recovery loop is at lines 127-132. ✓

## Common Mistakes reference
- `engineering/02-implementation-notes.md:34-36` (What Is Not Demonstrated): VERIFIED — "What Is Not Demonstrated" section is at lines 34-36. ✓

## Audit verdict references
- `research-audit/07-verdict.md` (APPROVED): VERIFIED.
- `engineering-audit/06-verdict.md` (APPROVED): VERIFIED.
- `research-audit/06-gaps.md`: exists. VERIFIED.
- `engineering-audit/05-gaps.md`: exists. VERIFIED.

## Full source list tree
- The tree (lines 106-156) matches the actual directory structure. VERIFIED.

## Summary
All source-map references resolve correctly. No broken or mislocated citations. Line-number references that were checkable match the source. No issues found.
