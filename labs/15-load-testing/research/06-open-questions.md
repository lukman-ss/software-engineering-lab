# Open Questions

## Unanswered Questions

### 1. How to calculate realistic concurrent users for the specific Booking Bengkel workflow?

The research confirms the general methodology (peak sessions per second × average session duration) but does not provide specific figures for a typical Booking Bengkel scenario. The exact parameters depend heavily on:
- Number of workshops/branches per session
- Average booking process steps (login → select branch → select service → choose time → payment → confirmation)
- Average duration of each booking workflow
- Peak traffic patterns for a workshop booking application

**Confidence:** LOW — Methodology is verified, but no concrete scenario parameters available.

**Next Steps:** Instrument production-like data or consult business stakeholders for realistic concurrent user projections.

---

### 2. What is the optimal threshold for each test type's pass/fail criteria?

The research confirms that thresholds (P95 < 200ms, error rate < 1%, etc.) are essential but does not provide industry-standard baselines.

**Confidence:** MEDIUM — Best practice is to define thresholds, but no universal "correct" thresholds exist.

**Next Steps:** Establish thresholds based on business SLAs and user experience research.

---

### 3. How to handle WhatsApp API latency in load testing?

The original lab scenario mentions "WhatsApp Konfirmasi" as a bottleneck source. The research identifies external API latency as a bottleneck category, but:

**Confidence:** LOW — Need specific guidance on testing with external API dependencies.

**Next Steps:** Research dependency testing strategies and fallback patterns for third-party APIs during load testing.

---

### 4. Memory leak detection during soak testing — specific metric thresholds?

The research confirms soak tests detect memory leaks over 3-72 hours but does not specify how much memory growth constitutes a leak that needs investigation.

**Confidence:** MEDIUM — Methodology known, threshold criteria not defined.

**Next Steps:** Look into memory profiling tooling and industry standards for acceptable memory growth rates.

---

### 5. Correlation methodology: When does P95 degradation indicate an actual problem?

The research shows P95 should be monitored but does not specify at what degradation point (e.g., 20% increase, 2x increase) requires investigation.

**Confidence:** LOW — No specific thresholds provided for degradation detection.

**Next Steps:** Correlate with business impact metrics and user experience data.

---

### 6. Optimal ramp-up duration relative to system warm-up time?

The k6 documentation mentions ramp-up should be "5-15% of total test duration" but does not define the relationship to application warm-up time (JIT compilation, cache warming, connection pool initialization).

**Confidence:** MEDIUM — General guidance exists (5-15%) but system-specific optimal values unknown.

**Next Steps:** Measure warm-up characteristics of the specific application and infrastructure.

---

### 7. Database-specific load testing metrics for connection pooling?

While the research mentions "Database Connection" as a metric to monitor, it does not detail specific metrics like:
- Active connections vs. pooled connections
- Connection acquisition time
- Pool saturation rate
- Query response time percentiles per query type

**Confidence:** LOW — Covered conceptually but no specific DB metrics defined.

**Next Steps:** Research database-specific performance counters for the DB backend in use (MySQL, PostgreSQL, etc.).

---

### 8. Distributed load testing coordination and result aggregation?

The research mentions distributed load testing as a concept but does not detail:
- How to coordinate multiple load injectors
- How to aggregate and interpret results across distributed runners
- Network overhead considerations in distributed load testing

**Confidence:** LOW — Concept acknowledged, detailed methodology not available in current sources.

**Next Steps:** Research distributed/load testing orchestration patterns for k6, Locust, or JMeter.

---

### 9. Impact of authentication and session management on load test accuracy?

For the Booking Bengkel app with a Login feature, how session management (JWT tokens, session stores, auth database lookups) affects load test results and whether synthetic sessions accurately represent real user behavior.

**Confidence:** LOW — Not specifically addressed in sources.

**Next Steps:** Research session-aware load testing patterns and authentication token management in load testing tools.

---

### 10. CI/CD pipeline integration patterns for automated load testing?

The research mentions k6 automated performance testing in CI/CD but does not detail:
- Gate criteria for blocking deployments
- Frequency of automated tests in the pipeline
- Handling of flaky test results in automated contexts

**Confidence:** MEDIUM — Mentioned but not detailed adequately.

**Next Steps:** Research industry best practices for CI/CD load testing gates and false-positive reduction.
