# Research Gaps Analysis: SLO, SLI & Error Budget

## Gap 1: Lack of Empirical Data for Cost of Nines
Type:
WEAK_SOURCE

Severity:
LOW

Location:
`research/05-report.md` (Finding 11), `research/06-open-questions.md`

Problem:
The claim that each additional "nine" costs ~100x more is accepted as a rule of thumb from the Google SRE Book, but lacks empirical financial case studies or formulaic proofs.

Required Revision:
None required for approval; research transparently flagged this in `06-open-questions.md`.

Can Be Approved Without Fix:
YES

---

## Gap 2: Specificity of Burn Rate Alerting Thresholds
Type:
OVERGENERALIZATION / IMPLEMENTATION_GAP

Severity:
MEDIUM

Location:
`research/05-report.md` (Finding 8), `research/03-evidence.md` (Evidence 12)

Problem:
Numeric burn rate thresholds (1-6 elevated, >6 critical over 2 hours) rely exclusively on Datadog documentation rather than an open standard (e.g., OpenTelemetry or Google SRE Workbook multi-window multi-burn rate alerts).

Required Revision:
None required for approval; research clearly separated vendor implementation from canonical SRE theory.

Can Be Approved Without Fix:
YES

---

## Gap 3: User Perception of Variance vs Speed Citation
Type:
MISSING_SOURCE

Severity:
LOW

Location:
`research/03-evidence.md` (Evidence 10), `research/06-open-questions.md`

Problem:
Google SRE Book asserts user studies show preference for consistent latency over lower mean latency with high variance, but does not explicitly cite the underlying academic paper/study in that section.

Required Revision:
None required for basic research approval; noted in open questions.

Can Be Approved Without Fix:
YES
