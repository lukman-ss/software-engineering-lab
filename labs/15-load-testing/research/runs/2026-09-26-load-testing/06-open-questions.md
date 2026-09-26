# Open Questions

## Unanswered Questions

### 1. What are realistic P95/P99 thresholds for Booking Bengkel workshop booking workflows?

**Current Evidence**: Threshold examples exist (`p(95)<200`, `p(95)<500`) but no industry-standard baselines for booking applications. Azure mentions 400ms API + 150ms DB budgets, but Booking Bengkel has multi-step workflows.

**Confidence**: LOW — Methodology established but specific values unknown.

**Next Steps**: 
- Analyze production analytics for current 95th percentile booking completion times
- Benchmark user tolerance: industry average 2-5 seconds for web transactions
- Set P95 < 500ms for individual API calls, P95 < 3000ms for end-to-end booking flow

---

### 2. How to calculate concurrent users for Booking Bengkel's specific workflows?

**Current Evidence**: Formula `concurrent_users = hourly_sessions × avg_session_duration / 3600` verified (Source 2) but no Booking Bengkel production data.

**Confidence**: HIGH — Method is correct; specific values require business data.

**Next Steps**:
- Instrument production to measure: peak sessions/second × average booking flow duration
- Factors: login (5s), select branch (3s), select service (4s), choose time (2s), payment (10s), generate invoice (3s), WhatsApp confirmation (5s)
- Total estimated duration: 32 seconds per booking session
- If peak traffic is 60 bookings/minute = 1 booking/second → 1 × 32 = 32 concurrent users baseline

---

### 3. How to handle WhatsApp API latency testing when external API is slow/unavailable?

**Current Evidence**: Resource-based bottleneck identification method (Source 4, 8) but no guidance on decoupling external dependency latency.

**Confidence**: LOW — External API dependency testing not fully addressed.

**Next Steps**:
- Research mock server patterns vs stubbed responses
- Implement graceful degradation: queue WhatsApp for batch processing
- Monitor WhatsApp delivery success separately from booking flow latency
- Consider asynchronous WhatsApp sending (return HTTP 202 immediately)

---

### 4. What constitutes a memory leak in production vs acceptable memory growth during soak tests?

**Current Evidence**: Soak tests detect memory leaks (Source 6) but no threshold defined. Google SRE discusses GC pressure (Source 8) but not quantitative thresholds.

**Confidence**: MEDIUM — Problem exists, specific thresholds undefined.

**Next Steps**:
- Monitor heap growth rate: < 5% growth/hour during steady state = acceptable
- > 10% growth/hour = investigation required
- Set alerts using memory leak detection tools (Go pprof, JVM profilers)
- Consider using HdrHistogram for long-running percentile tracking

---

### 5. Correlation methodology: At what P95 degradation level should investigation be triggered?

**Current Evidence**: P95 as metric is recommended (Source 3), but degradation thresholds undefined. Azure thresholds are binary pass/fail; no guidance on acceptable degradation ranges.

**Confidence**: LOW — No Tier 1 source defines degradation investigation thresholds.

**Next Steps**:
- Industry practice: 20-50% increase from baseline triggers investigation
- For Booking Bengkel: if P95 degrades 30% from established baseline, initiate root cause analysis
- Create SLO around acceptable degradation: e.g., "P95 should not exceed 1.5x baseline"
- Document in performance SLA

---

### 6. Optimal ramp-up duration relative to application warm-up time?

**Current Evidence**: k6 suggests "ramp-up 5-15% of total duration" (Source 1) but no guidance on JIT compilation, cache warming, connection pool initialization.

**Confidence**: MEDIUM — General guidance exists; system-specific optimal unknown.

**Next Steps**:
- Warm-up phase: 2-5 minutes at low load (10-20% target) before main test
- Connect to Google SRE's "executor load average" concept: ramp-up to allow resource allocation
- Measure warm-up by monitoring: CPU ramp-up, cache hit rate improvement, connection pool initialization

---

### 7. Database-specific metrics for connection pool saturation during load testing?

**Current Evidence**: k6 http_req_blocked indicates client-side pool saturation; Azure mentions connection pool problems (Source 10) but no specific metrics.

**Confidence**: LOW — Client-side metrics available; DB-server metrics not specified.

**Next Steps**:
- Database metrics to track:
  - Active connections vs pool size
  - Connection acquisition wait time
  - Pool saturation rate (% time pool at capacity)
  - Query response time percentiles per query type
- For Booking Bengkel: Pool size 5, so monitor when active connections = 5 consistently

---

### 8. How to coordinate distributed load injectors and aggregate results?

**Current Evidence**: Locust mentions distributed testing (Source 12); Azure Load Testing mentions CI/CD integration (Source 14) but no coordination patterns.

**Confidence**: LOW — Topic acknowledged but not detailed.

**Next Steps**:
- Research k6 vus/arrival-rate coordination across multiple executors
- Load testing patterns: hub-and-spoke (1 controller + N injectors) vs peer-to-peer
- Result aggregation: use InfluxDB/Prometheus for centralized metrics
- Network overhead modeling: subtract injectors' network impact from service metrics

---

### 9. How does authentication (JWT tokens, sessions) affect load test realism?

**Current Evidence**: Booking Bengkel requires login; JMeter and Gatling have auth docs but not integrated into load testing methodology.

**Confidence**: LOW — Session management impact not addressed by primary sources.

**Next Steps**:
- Research token acquisition rate limiting
- Session storage: in-memory vs Redis vs DB — affects scalability testing
- Token lifetime: need for refresh during long soak tests
- Synthetic sessions vs real user behavior: how to model realistic session patterns

---

### 10. CI/CD pipeline integration patterns for automated load testing gates?

**Current Evidence**: k6 Threshold CI/CD mentioned (Source 1), Azure Load Testing CI/CD integration (Source 14), but no standardized gate criteria.

**Confidence**: MEDIUM — Tools support CI/CD; standards not established.

**Next Steps**:
- Gate criteria examples:
  - Block deployment if P95 > 500ms for 5 consecutive runs
  - Block if error rate > 1% on main workflow
  - Block if no new degradation detected vs previous build
- Frequency: Pre-merge for critical systems, nightly baseline, before tagged releases
- Flaky test handling: require 2/3 consecutive failures to block, use statistical significance

---

### 11. What is the optimal threshold for Booking Bengkel's specific endpoints?

**Current Evidence**: General thresholds exist but no Booking Bengkel-specific metrics documented.

**Confidence**: MEDIUM — Need to establish Booking Bengkel-specific SLA.

**Next Steps**:
- Per-endpoint thresholds based on user impact:
  - Login: P95 < 2000ms (critical for subsequent actions)
  - Booking creation: P95 < 1000ms (user wait time)
  - Payment processing: P95 < 3000ms (external API dependent)
  - Invoice generation: P95 < 500ms
  - WhatsApp confirmation: P95 < 2000ms (external API)
- Overall flow: P95 < 5000ms for complete booking process

---

### 12. How to model connection pool saturation as a "closed system" when web layer is "open"?

**Current Evidence**: Gatling workload models (Source 13) distinguish open vs closed; Booking Bengkel has queue at DB pool level.

**Confidence**: MEDIUM — Theoretically sound; practical implementation unclear.

**Next Steps**:
- Hybrid approach: web layer uses open model (arrival rate), DB layer uses closed (pool capacity = 5)
- Monitor HTTP_req_blocked to detect when web layer exceeds DB capacity
- k6 scenario: use VU model for realistic booking session pacing
- Document saturation point: when 5 connections busy, remaining requests queue

---

### 13. What are the specific thresholds for aborting load tests on threshold failure?

**Current Evidence**: k6 cloud evaluates thresholds every 60 seconds (Source 3); local evaluates continuously.

**Confidence**: MEDIUM — Platform behavior documented; optimal interval unclear.

**Next Steps**:
- For CI/CD: abort immediately (local execution preferred for gates)
- For exploratory testing: use 60-second delay in cloud to avoid early abort from warm-up noise
- Define minimum sample requirement before threshold evaluation (e.g., >100 samples)

---

## Summary Table

| # | Question | Confidence | Priority |
|---|----------|------------|----------|
| 1 | P95/P99 thresholds for Booking workflows | LOW | High |
| 2 | Concurrent user calculation (practice) | HIGH | High |
| 3 | WhatsApp API latency testing | LOW | Medium |
| 4 | Memory leak thresholds | MEDIUM | Medium |
| 5 | P95 degradation investigation trigger | LOW | High |
| 6 | Ramp-up vs warm-up duration | MEDIUM | Medium |
| 7 | DB connection pool metrics | LOW | High |
| 8 | Distributed injector coordination | LOW | Low |
| 9 | Auth/session impact on realism | LOW | Medium |
| 10 | CI/CD load test gates | MEDIUM | High |
| 11 | Endpoint-specific thresholds | MEDIUM | High |
| 12 | Hybrid open/closed workload modeling | MEDIUM | Medium |
| 13 | Threshold evaluation intervals | MEDIUM | Low |

---

## Recommended Next Research

1. **Industry Benchmarks**: Survey Booking Bengkel domain peers (e-commerce, booking systems) for P95/P99 benchmarks
2. **Production Data**: Instrument Booking Bengkel for actual traffic patterns using Analytics/GA4
3. **DB Monitoring**: Investigate PostgreSQL/MySQL-specific connection pool metrics and alerting
4. **External API Patterns**: Document WhatsApp API SLA and implement circuit breaker pattern
5. **CI/CD Integration**: Prototype load test gate in existing pipeline with rollback strategy