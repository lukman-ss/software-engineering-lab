# Diagrams Audit

## File Reviewed
`content/04-diagrams.md` — 7 diagrams, 175 lines

## Verification

| Diagram | Source / Basis | Accuracy |
|---|---|---|
| 1 Pipeline Position | Pipeline metadata | PASS (informational) |
| 2 Race Condition: Check-Then-Act Failure | UnsafeStore store.go:24-47, test 210-239 | PASS |
| 3 Safe Insert: Atomic Constraint Enforcement | engine.go:45-99 InsertUser flow | PASS |
| 4 Constraint Evaluation Order | engine.go:45-99 + 122-148 combined | WARNING |
| 5 Partial Unique Index: Soft-Delete Flow | engine.go:71-96 + 101-120 | PASS |
| 6 Architecture | store.go + engine.go + dberr/errors.go + model.go | PASS (simplified) |
| 7 Error Code Mapping | dberr/errors.go:83-100 | PASS |

## Issue

### W-07 MEDIUM — Diagram 4 ordering claim exceeds research
Diagram 4 caption: "This ordering matches PostgreSQL behavior: NOT NULL first, then CHECK (alphabetical by name), then FK/UNIQUE (index-based)."

Research `04-contradictions.md` Open Question 3:
- Evidence supports: NOT NULL first, CHECK alphabetical by name
- **Missing**: "Explicit ordering for FK vs UNIQUE"
- Status: Assumption, not verified

The diagram presents the full 5-step order (NOT NULL → CHECK → UNIQUE → FK → COMMIT) as matching PostgreSQL, including the FK/UNIQUE ordering that research explicitly flagged as unknown. The lab's code never exercises both UNIQUE and FK in a single INSERT path (InsertUser has no FK; InsertOrder has no UNIQUE). The diagram idealizes the order.

Recommendation: Caption should clarify FK/UNIQUE ordering is lab's logical sequencing, not verified PostgreSQL behavior.

## Minor
- Diagram 1 includes pipeline-internal stage "Technical Writer (current)" — time-relative, will stale. Acceptable for internal docs.
- Diagram 6 shows single "Store" wrapping UnsafeStore/SafeStore; code has no common Store interface. Acceptable simplification.

## Verdict
VERIFIED_WITH_WARNINGS (W-07)