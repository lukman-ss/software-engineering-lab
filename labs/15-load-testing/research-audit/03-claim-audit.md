# Claim Audit

## Claim 1

Claim: Performance testing encompasses six standard types: Smoke, Average-Load, Stress, Soak/Endurance, Spike, and Breakpoint testing.

Location: `research/05-report.md:Finding 1` & `research/03-evidence.md:Evidence 1`

Evidence Provided: Grafana k6 docs (Sources 1-7), Azure Well-Architected Framework (Source 11), Google SRE Book (Source 9).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Standard industry categorization documented across multiple cloud frameworks and testing tool manuals.

---

## Claim 2

Claim: Key metrics for evaluation include latency percentiles (P50, P95, P99), error rate, throughput (RPS), and server resource utilization (CPU, memory, disk I/O, network).

Location: `research/05-report.md:Finding 2` & `research/03-evidence.md:Evidence 2`

Evidence Provided: k6 Metrics Reference (Source 20), k6 Thresholds (Source 15), Azure Performance Testing (Source 11), ISO/IEC 25010 (Source 10).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Percentiles prevent average latency masking tail performance degradation.

---

## Claim 3

Claim: Tool capabilities differ by architectural design: k6 (JS/API), Locust (Python gevent greenlets), JMeter (GUI/XML/Multi-protocol), Gatling (Scala DSL/JVM high performance).

Location: `research/05-report.md:Finding 3` & `research/03-evidence.md:Evidence 7`

Evidence Provided: Official docs for k6 (Source 1), Locust (Source 13), JMeter (Source 16), Gatling (Source 21).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Tool positioning and language/runtime model accurately represented.

---

## Claim 4

Claim: Server-side vs network vs dependency bottlenecks can be isolated by analyzing metric breakdowns like `http_req_waiting` (TTFB) vs `http_req_connecting`/`http_req_tls_handshaking`.

Location: `research/05-report.md:Finding 4` & `research/03-evidence.md:Evidence 3`

Evidence Provided: k6 Built-in Metrics Reference (Source 20) and Azure Performance Testing guidance (Source 11).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: High `http_req_waiting` specifically isolates server-side request processing time.

---

## Claim 5

Claim: Concurrent users formula for sizing virtual users: `Concurrent Users = Hourly Sessions * Average Session Duration (seconds) / 3600`.

Location: `research/03-evidence.md:Evidence 4` & `research/02-sources.md:Source 18`

Evidence Provided: Grafana k6 guide "Calculate concurrent users for load tests" (Source 18).

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes: Valid standard baseline formula for throughput estimation, though multi-step transactional workflows require session think-time adjustments.

---

## Claim 6

Claim: Breakpoint testing in elastic cloud environments risks finding account bill limits rather than application limits unless auto-scaling is disabled.

Location: `research/04-contradictions.md:Contradiction 3` & `research/02-sources.md:Source 6`

Evidence Provided: Grafana k6 Breakpoint Testing Documentation (Source 6).

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Practical warning for elastic cloud environments.
