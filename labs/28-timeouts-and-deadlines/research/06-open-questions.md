# Open Questions: Timeouts — Jangan Biarkan Satu Dependency Lambat Menahan Seluruh Sistem

## Unanswered Questions

### 1. Optimal Timeout Values for Indonesian/P99-Specific Scenario
**Question:** For the lab exercise (PostgreSQL P99=150ms, Inventory P99=400ms, Payment P99=1.2s, user deadline ~3s), what exact timeout budget allocation is optimal?
**Status:** NOT VERIFIED — No source provides exact budget for this specific combination. Sources provide principles (P99 + margin, total budget ≤ deadline minus processing overhead) but not concrete numbers for this case.
**Next research:** Load testing with actual traffic patterns to validate budget. Example heuristic: PostgreSQL 200ms, Inventory 500ms, Payment 1.5s, safety margin 300ms = 2.5s total (within 3s deadline) — but needs empirical validation.

### 2. Payment Retry Safety with Double-Spend Prevention
**Question:** When payment API times out after charge succeeded but response lost, how to guarantee exactly-once charging across retry?
**Status:** Partially verified — Idempotent Consumer pattern covers deduplication via idempotency keys, but payment provider specifics (bank API idempotency support, reconciliation window) not covered in sources.
**Next research:** Investigate specific payment provider idempotency key implementation (Stripe, Midtrans, Xendit) and reconciliation mechanisms.

### 3. Hedged Requests vs Timeouts Trade-off
**Question:** When should hedged requests (send to second replica after short delay) be used instead of simple timeouts?
**Status:** NOT VERIFIED — Google SRE mentions hedged requests briefly, but no detailed comparison found.
**Next research:** Review Dean & Barroso (2013) Tail at Scale paper for hedged request analysis.

### 4. Adaptive Timeout Configuration
**Question:** Should timeouts be static or dynamically adjusted based on current latency?
**Status:** Partially verified — Azure Circuit Breaker mentions AI/ML adaptive techniques, but no concrete adaptive timeout implementation details found.
**Next research:** Investigate Envoy proxy adaptive timeout, Netflix Hystrix adaptive configuration.

## Weak Evidence

### 1. Queue Worker Timeout Best Practices
**Evidence strength:** MEDIUM — RabbitMQ docs cover acknowledgements but not specific execution timeout configuration for job workers.
**Gap:** No authoritative source found for queue job execution timeout values, retry limit tuning, dead-letter handling specifics.
**Related sources:** RabbitMQ Reliability Guide mentions monitoring but not worker timeout policy details.

### 2. Lock Timeout vs Statement Timeout Interaction
**Evidence strength:** MEDIUM — PostgreSQL docs describe both but don't explain when to use which or how they interact in practice.
**Gap:** Need practical guidance on setting lock_timeout < statement_timeout to distinguish lock contention from slow queries.

### 3. Connection Pool Timeout Configuration
**Evidence strength:** LOW — Not directly covered in sources. Inferred from thread exhaustion discussion but no specific connection pool timeout guidance.
**Gap:** Connection acquisition timeout, pool exhaustion handling not addressed.

## Claims Needing Deeper Research

### 1. "Timeout is also resource protection mechanism"
**Current evidence:** Supported by SRE discussion of thread/memory exhaustion, but quantitative data on resource savings from timeouts is limited.
**Needed:** Benchmarks showing worker utilization with vs without timeouts under slow dependency.

### 2. "Timeout too small causes retry storm"
**Current evidence:** Supported by retry storm example, but threshold for "too small" not quantified.
**Needed:** Empirical data on timeout vs retry rate relationship at different percentile cutoffs.

### 3. "Deadline propagation prevents wasted work"
**Current evidence:** Strong conceptual evidence from gRPC, but quantitative impact on resource savings not provided.
**Needed:** Measurements of CPU/memory saved by early cancellation via deadline propagation.

## Possible Next Research Directions

### 1. Language-Specific Timeout Implementation
Investigate timeout configuration in:
- Go: context.WithTimeout, http.Client timeout fields
- Java: Resilience4j TimeLimiter, Hystrix, OkHttp timeouts
- Node.js: AbortController, axios timeout
- Python: requests timeout, aiohttp timeout

### 2. Observability Implementation
Research OpenTelemetry trace propagation for deadline tracking:
- Trace context propagation
- Timeout span attributes
- Dashboard examples for timeout monitoring

### 3. Chaos Engineering for Timeout Validation
Research methods to test timeout configurations:
- Latency injection tools (Toxiproxy, Chaos Monkey)
- Load testing timeout behavior
- Circuit breaker testing strategies

### 4. Database-Specific Deep Dive
For PostgreSQL: statement_timeout vs idle_in_transaction_session_timeout tuning
For MySQL: wait_timeout, interactive_timeout, innodb_lock_wait_timeout
For Redis: command timeout, connection timeout

### 5. Cost-Benefit Analysis
Research trade-offs between:
- Aggressive timeouts (better protection, more false positives)
- Conservative timeouts (fewer false positives, less protection)
- Circuit breaker thresholds (failure rate %, window size)

## Research Limitations Acknowledged

1. **Temporal limitation:** Research conducted 2026-09-26; practices evolve, especially cloud-native patterns
2. **Access limitation:** Some authoritative sources (AWS Builders Library specific articles) returned 404, suggesting URL changes or access restrictions
3. **Depth limitation:** Research focused on principles; implementation specifics require code-level investigation deferred to Engineering phase
4. **Context limitation:** Timeout optimal values are highly context-dependent (traffic patterns, resource constraints, SLO requirements) — no universal numbers exist