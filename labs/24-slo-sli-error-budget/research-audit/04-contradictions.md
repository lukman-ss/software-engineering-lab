# 04 - Contradictions Audit

Target Lab: `labs/24-slo-sli-error-budget`
Audit Date: 2026-09-28

---

## Contradiction 1: Target Setting — Current Performance vs User Expectations

Statement A:
"Don't pick a target based on current performance... may lock you into supporting a system that requires heroic efforts to meet its targets."
Location: Google SRE Book, Ch. 4

Statement B:
"While that advice is true, your current performance can be a good place to start if you don't have any other information, and if you have a good process for iterating in place."
Location: Google SRE Workbook, Ch. 2

Type:
SOURCE_CONFLICT / METHODOLOGICAL EVOLUTION

Impact:
Teams adopting SLOs may wonder whether to baseline against telemetry first or define SLO purely from business goals.

Assessment:
Not a logical contradiction, but a recognized pragmatic evolution. Both sources agree that baselining without subsequent iteration is dangerous.

---

## Contradiction 2: Measurement Window — Rolling vs Calendar-Aligned

Statement A:
"We have found a four-week rolling window to be a good general-purpose interval... closely aligned with user experience."
Location: Google SRE Workbook, Ch. 2

Statement B:
"We deliberately chose to bind our SLOs to a calendar month versus a rolling period to keep us focused and organized when running service reviews."
Location: Google SRE Workbook, Ch. 3 (Evernote Case Study)

Type:
SOURCE_CONFLICT / TRADE-OFF

Impact:
Affects alert dampening, budget reset dynamics, and alignment with business governance cycles.

Assessment:
Legitimate trade-off between user-perceived reliability (continuous rolling window) and corporate reporting cycles (calendar month). Both are documented in official SRE literature.

---

## Contradiction 3: Availability Paradigm — Request-Based vs Time-Based

Statement A:
"At Google, however, a time-based metric for availability is usually not meaningful... we define availability in terms of the request success rate."
Location: Google SRE Book, Ch. 3

Statement B:
Evernote used synthetic prober-based node uptime checks to define availability.
Location: Google SRE Workbook, Ch. 3

Type:
INTERNAL / ARCHITECTURAL VARIANCE

Impact:
Request-based availability weights impact by traffic volume; time-based availability treats an outage at 3 AM the same as noon peak.

Assessment:
Google's multi-region distributed setup renders time-based binary uptime irrelevant, whereas single-region services often start with prober checks before maturing to request-ratio SLIs.

---

## Contradiction 4: Alerting Thresholds — Google SRE Workbook vs Vendor Implementations

Statement A:
Google SRE Workbook Ch. 5 specifies multi-window multi-burn-rate alerting across 1h, 6h, and 3d windows.

Statement B:
Datadog implements a 2-hour rolling window burn rate indicator (thresholds 1x and 6x).

Type:
INTERNAL / IMPLEMENTATION_VARIANCE

Impact:
Engineers implementing SLOs with commercial tools might confuse UI dashboard status indicators with production paging alert logic.

Assessment:
Documented variation between theoretical multi-window alert engineering and vendor UI operationalization.

---

## Contradiction 5: Downtime per Month Calculation Conventions

Statement A:
Google SRE Book Appendix A: 99% availability = 7.2 hours/month, 99.99% = 4.32 minutes/month (based on a 30-day month).

Statement B:
Standard astronomical/Gregorian average month calculation (365.25 / 12 = 30.4375 days): 99% = 7 hours 18 minutes/month, 99.99% = 4 minutes 23 seconds/month.

Type:
INTERNAL / NUMERICAL ASSUMPTION

Impact:
Slight numerical drift (~1.5%) depending on which monthly baseline is assumed.

Assessment:
Minor convention mismatch. Research explicitly identified this variance in `04-contradictions.md`.
