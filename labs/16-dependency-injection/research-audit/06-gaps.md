# Research Gaps

## Gap 1

Type: UNVERIFIED_CLAIM
Severity: LOW
Location: `research/05-report.md` (Finding 10, Finding 11)
Problem: The heuristics regarding the specific threshold for over-injection (12 parameters) and the exact list of objects that require vs don't require DI (Value objects vs Services) are derived from the lab specification, not authoritative primary sources.
Required Revision: None required. The research agent already correctly identified these as heuristics/spec-driven rules and downgraded confidence to MEDIUM. This shows excellent critical analysis.
Can Be Approved Without Fix: YES

## Gap 2

Type: MISSING_SOURCE
Severity: LOW
Location: `research/02-sources.md` (Source 7)
Problem: NestJS documentation was not fully machine-readable (returned a JS shell redirect) so the specific decorators/modules claims could not be deeply verified via text extraction.
Required Revision: None required. The research agent flagged this in `04-contradictions.md` and `06-open-questions.md`, noting the weak evidence. The overall conclusions do not hinge on NestJS.
Can Be Approved Without Fix: YES
