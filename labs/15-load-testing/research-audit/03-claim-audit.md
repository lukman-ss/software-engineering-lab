# Claim Audit: Lab 15 (Load Testing)

## Claim 1
Claim: Load testing encompasses six primary test types: smoke, average-load, stress, soak/endurance, spike, and breakpoint.  
Location: `research/05-report.md: Finding 1`; `research/03-evidence.md: Evidence 1`  
Evidence Provided: Detailed citation of k6 documentation guides (Sources 1-7), Azure Well-Architected Framework (Source 11), Google SRE Book (Source 9).  
Source: Sources 1-7, 9, 11  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Solid consensus across primary cloud and testing vendor documentation.  

## Claim 2
Claim: Key metrics required for load test evaluation are P50, P95, P99 response times, error rate, throughput (RPS), and resource utilization (CPU, memory, disk I/O, network).  
Location: `research/05-report.md: Finding 2`; `research/03-evidence.md: Evidence 2`  
Evidence Provided: k6 metric reference & thresholds documentation, Azure Well-Architected Framework performance targets, ISO/IEC 25010 Quality Model.  
Source: Sources 10, 11, 15, 20  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: The emphasis on percentiles rather than averages is well-substantiated by both theoretical standards (ISO/IEC 25010) and practical documentation (k6, Azure).  

## Claim 3
Claim: k6 is JS-based for APIs, JMeter is GUI/XML for multi-protocol enterprise, Locust is Python/greenlet-based, Gatling is Scala/JVM DSL.  
Location: `research/05-report.md: Finding 3`; `research/03-evidence.md: Evidence 7`  
Evidence Provided: k6 and Locust official docs; general architectural characteristics for JMeter and Gatling.  
Source: Sources 1, 13, 14; general tool summaries  
Source Actually Supports Claim: PARTIAL  
Classification: FACT / INTERPRETATION  
Severity: MEDIUM  
Notes: k6 and Locust are supported by direct primary links. However, JMeter and Gatling were not backed by dedicated primary source entries in `02-sources.md` (only high-level website domain links given in `05-report.md`, and research report explicitly notes limitations in accessing their detailed documentation). Claim is technically accurate in the industry, but research citation depth is uneven.  

## Claim 4
Claim: Bottleneck isolation is achieved by decomposing HTTP request durations (e.g. `http_req_waiting` vs `http_req_connecting`) and correlating with server-side metrics.  
Location: `research/05-report.md: Finding 4`; `research/03-evidence.md: Evidence 3`  
Evidence Provided: k6 metric lifecycle breakdown (`http_req_duration`, `waiting`, `connecting`, etc.) and Azure hypothesis-driven testing guidelines.  
Source: Sources 11, 15, 20  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Accurate mapping of network/connection vs backend execution latency breakdown.  

## Claim 5
Claim: Common pitfalls include testing only `/health`, inadequate test data volume, testing on local developer laptops instead of production-mirror environments, and testing without explicit pass/fail SLA targets.  
Location: `research/05-report.md: Finding 5`; `research/03-evidence.md: Evidence 5`  
Evidence Provided: Azure Well-Architected Framework antipatterns and k6 threshold best practices.  
Source: Sources 9, 11, 15  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Well-documented industry failure modes accurately cited.  

## Claim 6
Claim: Concurrent users should be calculated using the formula: `Concurrent users = Hourly sessions * Average session duration (in seconds) / 3600`.  
Location: `research/03-evidence.md: Evidence 4`  
Evidence Provided: k6 guide "Calculate concurrent users for load tests".  
Source: Source 18  
Source Actually Supports Claim: YES  
Classification: FACT  
Severity: LOW  
Notes: Exact formula matches Grafana k6 sizing documentation (Little's Law application to session concurrency).  

## Claim 7
Claim: Target SLA baseline recommendations are P95 < 500ms, Error Rate < 1%, CPU < 75%, Memory < 80%.  
Location: `research/01-plan.md: Objective item 6`  
Evidence Provided: Stated as an objective in plan, but `06-open-questions.md` explicitly admits that no universal standard SLA thresholds exist and that these numbers are context/business specific.  
Source: None provided as universal standard  
Source Actually Supports Claim: NO (Recognized as context-dependent)  
Classification: HYPOTHESIS / EXAMPLE  
Severity: MEDIUM  
Notes: The research correctly refrains from claiming these arbitrary numbers as universal truths in `05-report.md` and isolates the limitation in `06-open-questions.md` (Question 2). However, in `01-plan.md`, it was phrased without clear caveat.  
