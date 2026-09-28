# 02 - Source Audit

Target Lab: `labs/24-slo-sli-error-budget`
Audit Date: 2026-09-28

---

## Source 1

Claimed Title: Service Level Objectives (Chapter 4, Google SRE Book)
Claimed Publisher: Google (O'Reilly Media, CC BY-NC-ND 4.0)
URL: `https://sre.google/sre-book/service-level-objectives/`

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None.

Assessment:
PASS

---

## Source 2

Claimed Title: Implementing SLOs (Chapter 2, Google SRE Workbook)
Claimed Publisher: Google (O'Reilly Media, CC BY-NC-ND 4.0)
URL: `https://sre.google/workbook/implementing-slos/`

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None.

Assessment:
PASS

---

## Source 3

Claimed Title: Embracing Risk (Chapter 3, Google SRE Book)
Claimed Publisher: Google (O'Reilly Media, CC BY-NC-ND 4.0)
URL: `https://sre.google/sre-book/embracing-risk/`

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None.

Assessment:
PASS

---

## Source 4

Claimed Title: Alerting on SLOs (Chapter 5, Google SRE Workbook)
Claimed Publisher: Google (O'Reilly Media, CC BY-NC-ND 4.0)
URL: `https://sre.google/workbook/alerting-on-slos/`

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None.

Assessment:
PASS

---

## Source 5

Claimed Title: Availability Table (Appendix A, Google SRE Book)
Claimed Publisher: Google (O'Reilly Media, CC BY-NC-ND 4.0)
URL: `https://sre.google/sre-book/availability-table/`

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Note: Calculates "per month" using standard 30-day month convention.

Assessment:
PASS

---

## Source 6

Claimed Title: SLO Engineering Case Studies (Chapter 3, Google SRE Workbook)
Claimed Publisher: Google (O'Reilly Media, CC BY-NC-ND 4.0)
URL: `https://sre.google/workbook/slo-engineering-case-studies/`

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Covers Evernote & Home Depot case studies accurately.

Assessment:
PASS

---

## Source 7

Claimed Title: Practical Alerting from Time-Series Data (Chapter 10, Google SRE Book)
Claimed Publisher: Google (O'Reilly Media, CC BY-NC-ND 4.0)
URL: `https://sre.google/sre-book/practical-alerting/`

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Focuses partly on Borgmon legacy tooling, but principles (symptom vs cause) remain authoritative.

Assessment:
PASS

---

## Source 8

Claimed Title: OpenSLO Specification
Claimed Publisher: OpenSLO Community (Apache 2.0)
URL: `https://openslo.github.io/OpenSLO/`

Reachable:
YES

Source Type:
PRIMARY (for OpenSLO specification)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Industry-wide adoption metrics outside of OpenSLO repository are not documented.

Assessment:
PASS

---

## Source 9

Claimed Title: Alerting Best Practices
Claimed Publisher: Prometheus (CNCF)
URL: `https://prometheus.io/docs/practices/alerting/`

Reachable:
YES

Source Type:
PRIMARY (for Prometheus monitoring best practices)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None.

Assessment:
PASS

---

## Source 10

Claimed Title: Service Level Objectives (Datadog)
Claimed Publisher: Datadog
URL: `https://docs.datadoghq.com/service_level_objectives/`

Reachable:
YES

Source Type:
SECONDARY / VENDOR

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Reflects Datadog-specific implementation details (e.g., 2-hour burn rate window indicator).

Assessment:
PASS

---

## Source 11

Claimed Title: Google Cloud Blog: SRE Basics (404)
Claimed Publisher: Google Cloud
URL: `https://cloud.google.com/blog/products/devops-sre/sre-basics-service-level-indicators-service-level-objectives-and-service-level-agreements`

Reachable:
NO (404 Not Found)

Source Type:
UNKNOWN

Relevant:
NO (Unavailable)

Supports Claimed Topic:
NO

Problems:
- URL dead/removed. Research agent correctly recorded it as 404 and did not base evidence on it.

Assessment:
WARNING

---

## Source 12

Claimed Title: Grafana SLO docs (404)
Claimed Publisher: Grafana
URL: `https://grafana.com/docs/grafana-cloud/monitor-applications/slos/`

Reachable:
NO (404 Not Found)

Source Type:
UNKNOWN

Relevant:
NO (Unavailable)

Supports Claimed Topic:
NO

Problems:
- URL dead. Correctly flagged by research agent as unusable.

Assessment:
WARNING

---

## Source 13

Claimed Title: OpenTelemetry SLO semconv (404)
Claimed Publisher: OpenTelemetry
URL: `https://opentelemetry.io/docs/specs/semconv/service-level-objectives/`

Reachable:
NO (404 Not Found)

Source Type:
UNKNOWN

Relevant:
NO (Unavailable)

Supports Claimed Topic:
NO

Problems:
- URL dead / semantic convention does not exist at path. Correctly flagged by research agent.

Assessment:
WARNING
