# Claim Audit

## Claim 1
Claim: Industrial standard recognizes six primary performance test types: smoke, average-load, stress, soak/endurance, spike, and breakpoint.  
Location: `research/03-evidence.md:3-19`, `research/05-report.md:13-30`  
Evidence Provided: Detailed definitions from k6, Azure Well-Architected Framework, and Google SRE Book.  
Source: Sources 1-7, 9, 11  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Solid consensus across primary sources.

---

## Claim 2
Claim: Senior engineers monitor P50, P95, P99 response times, error rate, RPS, and resource utilization (CPU, memory, disk I/O, network).  
Location: `research/03-evidence.md:21-37`, `research/05-report.md:31-52`  
Evidence Provided: ISO/IEC 25010 subcharacteristics, k6 metrics definitions, Azure performance targets.  
Source: Sources 10, 11, 20  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Well supported by ISO standards and cloud architecture guidelines.

---

## Claim 3
Claim: Bottleneck layer (application vs database vs external API) can be identified via TTFB (`http_req_waiting`), connection time (`http_req_connecting`), and external dependency metrics breakdown.  
Location: `research/03-evidence.md:38-55`, `research/05-report.md:73-91`  
Evidence Provided: k6 HTTP metric breakdown and Azure hypothesis-driven testing guidelines.  
Source: Sources 11, 20  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Standard request lifecycle metric correlation methodology.

---

## Claim 4
Claim: Concurrent virtual users are calculated as: Peak sessions per second $\times$ Average session duration.  
Location: `research/03-evidence.md:56-72`  
Evidence Provided: k6 calculate-concurrent-users guide and Azure throughput guidance.  
Source: Sources 2, 18, 11  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Mathematically equivalent to Little's Law applied to user concurrency.

---

## Claim 5
Claim: Common pitfalls include health-endpoint-only testing, insufficient data volume, testing without server monitoring, lack of defined targets, and testing on developer laptops.  
Location: `research/03-evidence.md:73-89`, `research/05-report.md:92-110`  
Evidence Provided: Azure Well-Architected Framework performance testing anti-patterns and Google SRE reliability principles.  
Source: Sources 9, 11  
Source Actually Supports Claim: YES  
Classification: FACT / BEST_PRACTICE  
Severity: LOW  
Notes: Directly aligns with cloud well-architected anti-patterns.

---

## Claim 6
Claim: Tool selection (k6 vs JMeter vs Locust vs Gatling) should be based on language/protocol fit rather than popularity.  
Location: `research/03-evidence.md:107-123`, `research/05-report.md:53-72`  
Evidence Provided: Locust official architecture docs (greenlet/gevent), k6 documentation, Gatling/JMeter general properties.  
Source: Sources 13, 14  
Source Actually Supports Claim: PARTIAL  
Classification: INTERPRETATION  
Severity: MEDIUM  
Notes: JMeter and Gatling citations are weak/secondary (unreachable URLs noted in Limitations section of 05-report.md). However, Locust and k6 claims are firmly backed.

---

## Claim 7
Claim: Ramp-up duration should be 5-15% of total test duration.  
Location: `research/02-sources.md:21`, `research/06-open-questions.md:61`  
Evidence Provided: k6 load testing guide.  
Source: Source 2  
Source Actually Supports Claim: YES  
Classification: IMPLEMENTATION-SPECIFIC  
Severity: LOW  
Notes: Specific to k6 staging guidance, correctly flagged in open questions as system-dependent.
