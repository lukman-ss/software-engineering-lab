# Gap Analysis: Lab 28 (Timeouts & Deadlines)

## Gap 1
Type: SCOPE_ERROR / IMPLEMENTATION_GAP
Severity: LOW
Location: `06-open-questions.md: Section 1`
Problem: Specific timeout budget numbers for Indonesian e-commerce benchmark scenarios (PG=150ms, Inventory=400ms, Payment=1.2s) are heuristic estimates (e.g. PG=200ms, Inventory=500ms, Payment=1.5s, Total=2.5s within 3s SLA) and require empirical load test validation.
Required Revision: Validate during implementation lab benchmarking phase.
Can Be Approved Without Fix: YES (Acknowledged as open question in research).

---

## Gap 2
Type: MISSING_CASE
Severity: LOW
Location: `06-open-questions.md: Section 3`
Problem: Hedged requests (sending parallel requests to 2nd replica after P95 delay) vs simple timeouts trade-offs are identified but not deeply quantified.
Required Revision: Include optional discussion item in lab presentation/article.
Can Be Approved Without Fix: YES.

---

## Summary
All research gaps have been documented in `06-open-questions.md`. No critical or blocking research gaps remain.
