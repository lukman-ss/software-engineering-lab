# Source Audit: SLO, SLI & Error Budget

## Source 1
Claimed Title: Service Level Objectives (Google SRE Book)
Claimed Publisher: Google, Inc. / O'Reilly Media
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
- None. Authoritative origin of modern SRE SLI/SLO concepts.

Assessment:
PASS

---

## Source 2
Claimed Title: Embracing Risk (Google SRE Book)
Claimed Publisher: Google, Inc. / O'Reilly Media
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
- None. Primary source for error budget mechanics, velocity tradeoff, and risk tolerance.

Assessment:
PASS

---

## Source 3
Claimed Title: Monitoring Distributed Systems (Google SRE Book)
Claimed Publisher: Google, Inc. / O'Reilly Media
URL: https://sre.google/sre-book/monitoring-distributed-systems/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Establishes the Four Golden Signals (Latency, Traffic, Errors, Saturation) and symptom-based alerting.

Assessment:
PASS

---

## Source 4
Claimed Title: Availability Table (Google SRE Book Appendix)
Claimed Publisher: Google, Inc. / O'Reilly Media
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
- None. Canonical reference table for downtime calculations per time period.

Assessment:
PASS

---

## Source 5
Claimed Title: Alerting (Prometheus Documentation)
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
- None. Industry standard guidance on symptom-based alerting.

Assessment:
PASS

---

## Source 6
Claimed Title: Service Level Objectives (Datadog Documentation)
Claimed Publisher: Datadog
URL: https://docs.datadoghq.com/service_level_objectives/

Reachable:
YES

Source Type:
SECONDARY (Commercial platform documentation)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Burn rate thresholds (1-6 elevated, >6 critical) and remaining error budget formula are proprietary Datadog implementation details, correctly classified in research as such.

Assessment:
PASS
