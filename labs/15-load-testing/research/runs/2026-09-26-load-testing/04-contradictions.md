# Contradictions

## Contradiction 1: Number of Test Types Recognized

**SOURCE A** (k6): Documents six test types (smoke, average-load, stress, soak, spike, breakpoint) in the "test-type cheat sheet" (Source 1).

**SOURCE B** (Azure WAF): Lists four test types (load, stress, spike, endurance/soak) in the "Apply multiple performance test types" table (Source 10).

**ASSESSMENT**: This is not a factual disagreement but a taxonomy granularity difference. k6 includes smoke and breakpoint as first-class types; Azure treats them as implementation details within other categories or as part of capacity planning. The four Azure types are a subset of the six k6 types. All sources agree on the core concepts (stress = above-average load, soak = extended duration, spike = sudden burst). 

**Resolution**: Both are correct within their context. k6's 6-type framework is more comprehensive for tool-level classification; Azure's 4-type is operationally focused. For Booking Bengkel, use k6's taxonomy for completeness.

**Confidence**: LOW (apparent disagreement, actually a categorization difference)

---

## Contradiction 2: Stress Test Load Increase Percentage

**SOURCE A** (k6 Stress Testing): Explicitly states: "Some testers might have default targets for stress tests—say an increase upon average load by 50 or 100 percent—there's no fixed percentage. The load simulated in a Stress test depends on the stressful situations that the system may be subject to." (Source 5)

**SOURCE B** (Common industry practice / older tutorials): Frequently cite "50-100% above average load" as a stress test guideline.

**ASSESSMENT**: k6 (Tier 1) is authoritative and explicitly rejects fixed percentages. The "50-100%" figure from secondary sources is an oversimplification. Google SRE (Source 7) supports this by using risk-based stress scenarios ("paydays, rush hours, end of workweek") rather than fixed multipliers. The risk-profile-driven approach is confirmed by both Tier 1 sources.

**Resolution**: No fixed percentage is correct. Stress test load should be derived from the system's operational risk profile (peak events, business cycles), not arbitrary multipliers.

**Confidence**: HIGH (k6 is authoritative, industry blogs are secondary)

---

## Contradiction 3: Breakpoint Testing in Elastic Cloud Environments

**SOURCE A** (k6): Explicitly warns: "Avoid breakpoint tests in elastic cloud environments. The elastic environment may grow as the test moves further, finding only the limit of your cloud account bill. If this test runs on a cloud environment, turning off elasticity on all the affected components is strongly recommended." (Source 1)

**SOURCE B** (Azure): Lists breakpoint-style capacity testing as valid approach for "finding maximum capacity" and "failure modes" (Source 10) without mentioning the elasticity risk.

**ASSESSMENT**: Both sources agree on the purpose of breakpoint testing (find capacity limits). k6 provides operational warning about auto-scaling masking true limits and incurring costs; Azure's guidance is general-purpose and doesn't surface this implementation-specific risk. k6's warning reflects real operational concern: auto-scaling can mask the true system limit by adding resources, potentially leading to unbounded costs.

**Resolution**: k6's warning should be followed — disable auto-scaling before breakpoint tests in cloud environments. This is an implementation detail, not a disagreement on methodology.

**Confidence**: LOW (operational nuance, not a factual disagreement)

---

## Contradiction 4: Threshold Evaluation Frequency in k6 Cloud vs Local

**SOURCE A** (k6 Thresholds): "When k6 runs in the cloud, thresholds are evaluated every 60 seconds. Therefore, the abortOnFail feature may be delayed by up to 60 seconds." (Source 3)

**SOURCE B**: Implicitly, local execution evaluates thresholds continuously.

**ASSESSMENT**: This is by-design platform-dependent behavior, not a contradiction between sources. Local CLI evaluates thresholds frequently during execution; cloud evaluates every 60 seconds for aggregation efficiency. Engineers should be aware that cloud-based k6 tests may continue running up to 60 seconds past threshold failure when using abortOnFail.

**Resolution**: No contradiction — this is documented platform behavior. Use local execution for faster abort feedback; cloud for large-scale distributed testing.

**Confidence**: N/A (documented platform behavior, not a source disagreement)

---

## Contradiction 5: Test Environment — Staging vs Production Testing

**SOURCE A** (Azure): Recommends "Run controlled production testing. Schedule tests during off-peak hours" for "the most realistic results" (Source 10).

**SOURCE B** (Lab specification context): "Load test sebaiknya dilakukan pada lingkungan yang mendekati production" (test in environment close to production), warns against laptop testing.

**ASSESSMENT**: Azure's guidance is nuanced — staging for most tests, controlled production for realistic validation with safeguards. The lab's advice aligns with Azure's "mirror production as closely as practical." The apparent tension (staging vs production) resolves under Azure's framework: use staging for routine tests, production only with controls (off-peak, extra capacity, rollback plans).

**Resolution**: Both agree: environment should mirror production. Azure adds production testing as an option for noncritical workloads with controls. For Booking Bengkel: staging for routine load testing, production testing during off-peak for validation.

**Confidence**: LOW (apparent tension resolved by source context)

---

## Contradiction 6: Workload Model Choice — Concurrent Users vs Arrival Rate

**SOURCE A** (Gatling): "Don't reason in terms of concurrent users if your system can't push excess traffic into a queue." Recommends open model (arrival rate) for most websites.

**SOURCE B** (k6 Load Testing Guide): Uses Virtual Users (VU) model with concurrent user semantics in default examples.

**ASSESSMENT**: k6 supports both open (arrival-rate) and closed (VU) models via scenarios (Source 1: "k6 can model load by either number of VUs or by number of iterations per second (open vs. closed)"). k6's default quickstart uses VU model for simplicity, but ramping-arrival-rate executor supports open model. Azure's cheat sheet and k6's test types use "VUs/Throughput" which conflates the two concepts.

**Resolution**: Both tools support both models. Gatling is more opinionated about matching model to system behavior. For Booking Bengkel (web app with queue at DB connection), the system is effectively closed at the bottleneck (DB connection pool of 5) but open for web arrivals — both should be tested to understand full behavior.

**Confidence**: MEDIUM (tool positioning difference, not factual disagreement)

---

## Contradiction 7: Memory Leak Detection Specificity

**SOURCE A** (k6 Soak Testing): States soak tests detect "memory or other resource leaks, data saturation, and storage depletion" but only lists test duration ranges (3-72 hours) (Source 6).

**SOURCE B** (Azure): States soak tests detect "memory leaks, resource exhaustion, connection pool problems" but does not specify memory growth thresholds (Source 10).

**ASSESSMENT**: Both sources agree soak tests detect memory leaks and resource issues. Neither provides specific thresholds for "when does memory growth indicate a leak requiring investigation." The existing open-questions document (Source from prior run) also notes this gap. This is an evidence gap, not a contradiction.

**Resolution**: No contradiction, but a documented gap in threshold specificity. Memory leak detection requires additional tooling (profilers, growth rate analysis) not covered by these sources.

**Confidence**: LOW (evidence gap noted)

---

## Contradiction 8: P95/P99 Degradation Investigation — What Constitutes Degradation

**SOURCE A** (k6): Documents p(95) and p(99) threshold syntax (e.g., `p(95)<200`, `p(99)<400`) but does not define what percentage degradation from baseline triggers investigation (Source 3).

**SOURCE B** (Azure): Emphasizes setting thresholds based on SLOs, but thresholds are binary pass/fail — no guidance on what gradual increase warrants investigation (Source 10).

**ASSESSMENT**: Both sources agree P95/P99 should be monitored with thresholds. Neither defines specific degradation indicators (e.g., "20% increase from baseline" or "doubling of P95"). The investigation methodology exists implicitly (compare against baseline) but thresholds for "investigation-worthy" degradation aren't standardized. This appears to be a gap in both sources.

**Resolution**: Degradation criteria are organization-specific. Industry practice suggests 20-50% P95 increase or 2x P95 as investigation triggers, but not documented by Tier 1 sources.

**Confidence**: LOW (evidence gap, not contradiction)

---

## Summary

**Total contradictions identified**: 8  
**Factual contradictions**: 0  
**Taxonomy/nuance differences**: 5 (Sources 1, 3, 4, 5, 6)  
**Evidence gaps**: 3 (Sources 7, 8, and P95 investigation methodology)

**Assessment**: No material factual contradictions exist in the load testing domain. All Tier 1 sources agree on fundamental principles: test type definitions, metric importance (percentiles), bottleneck identification methodology, common pitfalls, and timing.

**Apparent disagreements resolve as**:
1. Granularity of categorization (k6 6 types vs Azure 4 types)
2. Platform-specific operational warnings (k6 cloud threshold evaluation, elasticity)
3. Implementation defaults vs advanced features (k6 VU model default vs arrival-rate option)
4. Scope of coverage (general frameworks vs tool-specific details)

**Evidence gaps identified requiring further investigation**:
- Specific thresholds for P95/P99 degradation triggering investigation
- Memory leak detection thresholds (how much growth = leak)
- P95/P99 degradation investigation methodology (step-by-step)
