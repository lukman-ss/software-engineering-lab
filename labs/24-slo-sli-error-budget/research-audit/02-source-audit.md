# 02 — Source Audit

Audit of all external sources cited in `labs/24-slo-sli-error-budget/research/02-sources.md`.

---

## Source 1

Claimed Title: Chapter 4 - Service Level Objectives
Claimed Publisher: Google SRE Book (Google Inc., O'Reilly Media)
URL: https://sre.google/sre-book/service-level-objectives/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative primary source for SLI/SLO/SLA definitions, percentiles vs averages, and user-facing metrics.

Assessment:
PASS

---

## Source 2

Claimed Title: Chapter 3 - Embracing Risk
Claimed Publisher: Google SRE Book (Google Inc., O'Reilly Media)
URL: https://sre.google/sre-book/embracing-risk/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- "Cost increases 100x per nine" cited inside text is a qualitative rule of thumb / internal heuristic without published empirical dataset.

Assessment:
PASS

---

## Source 3

Claimed Title: Appendix A - Availability Table
Claimed Publisher: Google SRE Book (Google Inc., O'Reilly Media)
URL: https://sre.google/sre-book/availability-table/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Basis is strict 30-day month (43,200 min), leading to a 6-minute discrepancy when comparing with average Gregorian month calculations (30.44 days = 43,833.6 min).

Assessment:
PASS

---

## Source 4

Claimed Title: Chapter 2 - Implementing SLOs
Claimed Publisher: Google SRE Workbook (Google Inc.)
URL: https://sre.google/workbook/implementing-slos/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Primary source establishing good/total event ratios, rolling 4-week windows, and reasons 100% reliability is wrong.

Assessment:
PASS

---

## Source 5

Claimed Title: Chapter 5 - Alerting on SLOs
Claimed Publisher: Google SRE Workbook (Google Inc.)
URL: https://sre.google/workbook/alerting-on-slos/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Multi-burn-rate parameters (14.4x/1h, 6x/6h, 1x/3d, short window = 1/12) are specific Google operational recommendations, not formal mathematical standards.

Assessment:
PASS

---

## Source 6

Claimed Title: Appendix B - Example Error Budget Policy
Claimed Publisher: Google SRE Workbook (Google Inc.)
URL: https://sre.google/workbook/error-budget-policy/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Contains claim "Changes are a major source of instability, representing roughly 70% of our outages" which is an internal Google retrospective observation without external corroboration.

Assessment:
WARNING

---

## Source 7

Claimed Title: Concepts in Service Monitoring
Claimed Publisher: Google Cloud Documentation (Google Inc.)
URL: https://cloud.google.com/stackdriver/docs/solutions/slo-monitoring

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Official Cloud documentation explaining metric-based and request-based SLO monitoring.

Assessment:
PASS

---

## Source 8

Claimed Title: Alerting | Prometheus
Claimed Publisher: Prometheus Authors / The Linux Foundation
URL: https://prometheus.io/docs/practices/alerting/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Canonical documentation on symptom-based alerting and golden signals.

Assessment:
PASS

---

## Source 9

Claimed Title: Service Level Objectives
Claimed Publisher: Datadog
URL: https://docs.datadoghq.com/service_level_objectives/

Reachable:
YES

Source Type:
SECONDARY (Platform Documentation)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Error budget remaining formula `100 * (current_status - target) / (100 - target)` and burn rate thresholds (1-6 elevated, >6 critical over 2h) are Datadog-specific SaaS implementations.

Assessment:
PASS
