# Contradictions Audit: SLO, SLI & Error Budget

## Contradiction 1: Month Availability Calculations in Topic Spec vs Standard SRE Tables
Statement A:
Topic specification stated 99% availability equates to ~7 hours 18 minutes downtime per month, and 99.99% equates to ~4 minutes 23 seconds per month.
Location:
Topic input specification.

Statement B:
Canonical availability table defines 99% as 7.2 hours/month (7 hours 12 minutes) and 99.99% as 4.32 minutes/month (4 minutes 19.2 seconds).
Location:
Google SRE Book Appendix A: https://sre.google/sre-book/availability-table/

Type:
SOURCE_CONFLICT / SPECIFICATION_ERROR

Impact:
Minor numeric discrepancy in input materials.

Assessment:
The research artifacts correctly noted and resolved this discrepancy using the primary canonical source while documenting the exact arithmetic difference in `research/03-evidence.md` (Evidence 21) and `research/04-contradictions.md`.

---

## Contradiction 2: SRE Conceptual Model vs Commercial Platform Mechanics
Statement A:
Google SRE Book describes error budget consumption and thresholding qualitatively and as policy-driven quarterly targets.
Location:
Google SRE Book Chapters 3 & 4.

Statement B:
Datadog documentation formalizes specific burn rate numeric tiers (1-6 elevated, >6 critical) over rolling 2-hour sliding windows.
Location:
Datadog SLO Documentation.

Type:
INTERNAL / SCOPE_DIFFERENTIATION

Impact:
Potential risk of treating vendor-specific alerting thresholds as immutable SRE laws.

Assessment:
The research correctly identifies this as complementary implementation detail rather than an irreconcilable conflict, appropriately isolating Datadog's specifics into lower confidence / implementation-specific tags.

---

## Summary
No material contradictions or unsupported conflicts remain unresolved across the research corpus.
