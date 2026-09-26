# Claim Audit

## Claim 1
Claim: Load testing encompasses six primary test types: smoke, average-load, stress, soak/endurance, spike, and breakpoint.
Location: `research/05-report.md` (Finding 1), `research/03-evidence.md` (Evidence 1)
Evidence Provided: Detailed breakdown of each test type with load patterns and objectives.
Source: Grafana k6 docs (Sources 1-7), Microsoft Azure Well-Architected Framework (Source 11), Google SRE Book (Source 9).
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully supported across multiple authoritative Tier 1 sources.

---

## Claim 2
Claim: Standard performance evaluation relies on percentiles (P50, P95, P99), error rates, throughput (RPS), and system resource metrics (CPU, memory, disk, network) rather than arithmetic averages alone.
Location: `research/05-report.md` (Finding 2), `research/03-evidence.md` (Evidence 2)
Evidence Provided: k6 metric percentiles definition, Azure performance targets, ISO/IEC 25010 Quality Model subcharacteristics.
Source: Grafana k6 Metrics Reference (Source 20), Azure Performance Testing (Source 11), ISO/IEC 25010 (Source 10).
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Soundly evidenced; correctly emphasizes percentile evaluation over unweighted averages to catch tail latency.

---

## Claim 3
Claim: Tools serve distinct operational contexts: k6 is specialized for API and script-driven testing in JS; Locust provides lightweight greenlet coroutines in Python; JMeter supports multi-protocol enterprise workflows via GUI/XML; Gatling offers high performance on the JVM.
Location: `research/05-report.md` (Finding 3), `research/03-evidence.md` (Evidence 7)
Evidence Provided: Architecture descriptions for Locust (gevent greenlets), k6 execution model, JMeter and Gatling capabilities.
Source: Locust docs (Source 13), k6 docs (Source 1), Azure docs (Source 11).
Source Actually Supports Claim: YES
Classification: INTERPRETATION
Severity: LOW
Notes: Appropriately qualified with MEDIUM confidence where specific benchmark comparisons are absent.

---

## Claim 4
Claim: Bottlenecks across application, database, and external APIs can be isolated by decomposing the HTTP lifecycle (e.g. `http_req_waiting` vs `http_req_connecting`) and correlating latency shifts with component-level metrics.
Location: `research/05-report.md` (Finding 4), `research/03-evidence.md` (Evidence 3)
Evidence Provided: k6 metric component breakdown equation and Azure hypothesis-driven testing guidelines.
Source: Azure Performance Testing (Source 11), k6 Built-in Metrics Reference (Source 20).
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Core diagnosis strategy is backed by standard request timeline decomposition principles.

---

## Claim 5
Claim: Calculating virtual users for transactional scenarios uses peak sessions per second multiplied by average session duration (Little's Law application).
Location: `research/03-evidence.md` (Evidence 4)
Evidence Provided: k6 VU calculation guidelines and Azure throughput estimation formulas.
Source: k6 VU calculation guide (Source 18), Azure Performance Testing (Source 11).
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Standard queuing theory and performance engineering calculation.

---

## Claim 6
Claim: Common performance testing anti-patterns include testing `/health` exclusively, using miniature data fixtures, running tests against unmonitored infrastructure, testing without strict SLO thresholds, and using developer laptops as production proxies.
Location: `research/05-report.md` (Finding 5), `research/03-evidence.md` (Evidence 5)
Evidence Provided: Azure anti-pattern documentation and Google SRE reliability principles.
Source: Azure Performance Testing (Source 11), Google SRE Book (Source 9).
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Highly practical, directly cited from Azure Well-Architected and Google SRE guidance.

---

## Claim 7
Claim: Load testing is a continuous SDLC activity required before go-live, ahead of major traffic events, and following major database, cloud, or architectural changes.
Location: `research/05-report.md` (Finding 6), `research/03-evidence.md` (Evidence 6)
Evidence Provided: Google SRE equivalence testing principles, Azure continuous testing guidance, CI/CD automated testing references.
Source: Google SRE Book (Source 9), Azure Performance Testing (Source 11), k6 docs (Source 1).
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully aligned with modern CI/CD and SRE practices.
