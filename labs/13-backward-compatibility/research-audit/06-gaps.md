# Research Gap Analysis

## Gap 1

Type:
UNVERIFIED_CLAIM

Severity:
LOW

Location:
`research/08-failure-modes.md` (Failure Mode 4), `research/11-final-research.md` (Question 11)

Problem:
The 30-day heuristic for declaring legacy endpoint traffic at 0 before deleting code lacks authoritative Tier 1 evidence. The research correctly flags this as `NOT VERIFIED`, but it remains an uncalibrated recommendation.

Required Revision:
None blocking. Research authors already flagged the heuristic explicitly. Downstream labs should advise teams to rely on actual SLA/monitoring thresholds.

Can Be Approved Without Fix:
YES

---

## Gap 2

Type:
MISSING_CASE

Severity:
MEDIUM

Location:
`research/10-open-questions.md` (Question 1)

Problem:
Dual-write consistency patterns across distributed microservices (without 2PC) are outlined conceptually via Outbox/Saga, but detailed implementation patterns specifically tailored for database schema migration across independent services are marked as open questions.

Required Revision:
None for this phase. Accurately captured as an open research question in `10-open-questions.md`.

Can Be Approved Without Fix:
YES

---

## Gap 3

Type:
WEAK_SOURCE

Severity:
LOW

Location:
`research/10-open-questions.md` (Question 4)

Problem:
The quantitative CPU and memory overhead of runtime response transformation layers (such as Stripe's date-based version change modules) lacks public benchmarks.

Required Revision:
None for this phase. Recorded in open questions.

Can Be Approved Without Fix:
YES
