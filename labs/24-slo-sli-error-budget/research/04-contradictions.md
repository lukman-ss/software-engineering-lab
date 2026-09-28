# Contradictions

## Contradiction 1: Pick SLO from current performance or not?

SOURCE A (Google SRE Book, Ch.4 "Choosing Targets"):
"Don't pick a target based on current performance. While understanding the merits and limits of a system is essential, adopting values without reflection may lock you into supporting a system that requires heroic efforts to meet its targets."

SOURCE B (Google SRE Workbook, Ch.2):
"In our first book, we advise against picking an SLO based upon current performance... While that advice is true, your current performance can be a good place to start if you don't have any other information, and if you have a good process for iterating in place."

ASSESSMENT:
Not a factual contradiction but a documented evolution of guidance (2016 → 2018). Source A warns against anchoring; Source B softens it for teams with no data, with iteration as the condition. Both agree current performance alone must not be the final answer. Uncertainty: which approach is preferable in practice is not resolved by either source; context-dependent.

## Contradiction 2: Time window — rolling vs calendar

SOURCE A (Google SRE Workbook, Ch.2 "Choosing an Appropriate Time Window"):
"We have found a four-week rolling window to be a good general-purpose interval... Rolling windows are more closely aligned with user experience. We recommend defining this period as an integral number of weeks so it always contains the same number of weekends."

SOURCE B (Evernote case study, same workbook):
"We deliberately chose to bind our SLOs to a calendar month versus a rolling period to keep us focused and organized when running service reviews."

SOURCE C (Home Depot case study, same workbook):
"Like many companies adopting error budgets, we're weighing the pros and cons of rolling windows versus fixed windows."

ASSESSMENT:
Google recommends rolling windows (user experience alignment); Evernote explicitly chose calendar month (review cadence alignment); Home Depot undecided. Sources may differ because rolling windows optimize for user experience while calendar windows optimize for business planning. What remains uncertain: neither source demonstrates a data-driven conclusion that one is universally better; both trade-offs are acknowledged.

## Contradiction 3: Availability measurement — time-based vs request-based

SOURCE A (Google SRE Book, Ch.3 "Embracing Risk"):
"At Google, however, a time-based metric for availability is usually not meaningful... we define availability in terms of the request success rate." (For globally distributed services, they are "at least partially 'up' at all times.")

SOURCE B (Evernote case study, Google SRE Workbook):
Evernote used prober-based uptime checks: "If a prober check fails, the node is marked as Unconfirmed Down and then a second geographically separate prober performs a check."

ASSESSMENT:
Google prefers aggregate request-based availability; Evernote used time-based prober uptime (a black-box measure). They may differ because Evernote had a single-region service with a status page endpoint, while Google services are globally distributed so time-based availability is always "up". Uncertainty: Evernote acknowledged limitations (moved toward client-side/API-level SLIs in later versions). Both are valid depending on system topology; neither source claims universality.

## Contradiction 4: Burn rate indicator implementation differs across vendors

SOURCE A (Google SRE Workbook, Ch.5):
Recommends multi-window, multi-burn-rate alerting: page at 14.4x burn rate over 1h+5m windows, 6x over 6h+30m; ticket at 1x over 3d+6h.

SOURCE B (Datadog documentation):
"Burn rate indicators use a rolling 2-hour window... A red icon indicating a critical burn rate above 6 in the past 2 hours. A yellow icon indicating an elevated burn rate between 1 and 6 in the past 2 hours."

ASSESSMENT:
Google's approach is multi-window with different thresholds per window; Datadog uses a single 2-hour rolling window with two thresholds (1 and 6). They differ because Datadog provides a simplified built-in indicator for dashboards, while Google's approach is tuned for alerting precision/recall trade-offs. Uncertainty: Datadog does not document how its indicator maps to page/ticket decisions; Google's parameters are explicitly stated as "starting point" requiring tuning. Core burn rate math agrees across both.

## Contradiction 5: Lab material downtime figures vs Google Availability Table

SOURCE A (Google SRE Book, Appendix A "Availability Table"):
99% → 7.2 hours per month; 99.9% → 43.2 minutes per month; 99.99% → 4.32 minutes per month (assuming 30-day month).

SOURCE B (Lab topic specification):
"99% → ~7 jam 18 menit/bulan; 99,9% → ~43 menit/bulan; 99,99% → ~4 menit 23 detik/bulan"

ASSESSMENT:
99.9% matches (43.2 min ≈ 43 min). 99% and 99.99% differ slightly: lab figures (7h18m, 4m23s) correspond to an average month of 365.25/12 ≈ 30.44 days, while Google's table uses a 30-day month (7.2h, 4m19.2s). Assessment: both are arithmetically correct under different month-length assumptions; not a factual error but an undisclosed assumption. Lab material does not state which month definition it uses. Google's table is the authoritative primary source for "per month" figures under its stated 30-day convention.

## Contradiction 6: Google Cloud blog overview unavailable

SOURCE A (Intended): Google Cloud blog "SRE basics: SLI, SLO, SLA" — fetch returned 404 (2026-09-28).

ASSESSMENT:
Cannot be used as corroboration. Definitions of SLI/SLO/SLA are instead verified via Google SRE Book Ch.4 (Tier 1) plus Datadog docs (Tier 2), which agree. No substantive disagreement; the missing source only reduces redundancy of cross-check for the definitional claims (still 2 independent sources: Google SRE Book + Datadog).