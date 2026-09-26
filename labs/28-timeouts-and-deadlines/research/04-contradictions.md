# Contradictions and Disagreements: Timeouts — Jangan Biarkan Satu Dependency Lambat Menahan Seluruh Sistem

## Contradiction 1: Retry Immediately vs. Retry with Delay

**SOURCE A (Microsoft Azure Retry Pattern):**
"Retry immediately. If the specific fault reported is unusual or rare, like a network packet becoming corrupted while it was being transmitted, the best course of action might be to immediately retry the request."

**SOURCE B (Google SRE Handling Overload):**
When backend tasks are overloaded, requests should not be retried immediately but errors should bubble up. Only retry immediately if small subset of tasks are overloaded, not when large subset is involved.

**ASSESSMENT:**
The sources address different failure scenarios. Azure refers to rare/unusual faults (corrupted packets, transient network blips) where immediate retry is safe because the fault is unlikely to be systemic. Google SRE refers to overload scenarios where retrying immediately can amplify an existing problem.

**RESOLUTION:**
Both are correct in their contexts. The key is distinguishing between:
- Rare transient faults → immediate retry may be safe
- Overload/rejection faults → retrying immediately worsens the problem
- Large-scale overload → don't retry at all; bubble up error

Application should implement exception-type-aware retry policies that differentiate failure modes.

## Contradiction 2: Health Check Depth and Frequency

**SOURCE A (Microsoft Azure Health Endpoint Monitoring):**
Recommends separate liveness and readiness endpoints. Liveness should be cheap (no dependency queries). Readiness should check dependencies but each check must be bounded with timeouts.

**SOURCE B (General SRE Practice from Sources 8, 9):**
Emphasizes monitoring four golden signals, some requiring dependency interaction. Alerting philosophy suggests detailed monitoring for diagnosis.

**ASSESSMENT:**
No actual contradiction — the sources complement each other. Azure provides implementation detail for health endpoints, while Google SRE provides monitoring philosophy. Both agree that:
- Liveness checks must be cheap and dependency-free
- Dependency checks need timeouts to prevent false alarms
- Cached health status trades freshness for reduced load

**RESOLUTION:**
Implement both patterns:
- Shallow liveness endpoints for orchestrator decisions (restart unhealthy containers)
- Dependency-aware readiness for load balancer decisions (route around unhealthy instances)
- Full metrics collection via monitoring agents (separate from health probes)

## Contradiction 3: Retry Budget Size

**SOURCE A (Google SRE Handling Overload):**
Per-request retry budget: 3 attempts max. Per-client retry budget: only retry when retry ratio < 10%. These limits prevent retry storms.

**SOURCE B (Resilience4j Documentation):**
Default maxAttempts: 3 for retry instances. Configurable to higher values.

**SOURCE C (Google SRE Cascading Failures):**
Example code shows up to 10 retry attempts in frontend RPC calls.

**ASSESSMENT:**
The 10-attempt example is explicitly described as naive code that contributes to cascading failures. The handling overload chapter provides the recommended budgets (3 attempts, 10% ratio). Resilience4j's default of 3 aligns with Google SRE recommendations.

**RESOLUTION:**
No real disagreement — the 10-attempt code is an anti-pattern illustration. All authoritative sources recommend 3 attempts as sensible default, with additional client-level budgets to prevent amplification.

## Contradiction 4: Where to Implement Retries in Stack

**SOURCE A (Microsoft Azure Retry Pattern):**
"Implement retry logic only where full context of failing operation is understood. If lower-level task also has retry policy, extra retries can add long delays. Better to configure lower-level task to fail fast."

**SOURCE B (Google SRE Handling Overload):**
"Requests should only be retried at layer immediately above the layer that is rejecting them. When request can't be served, return 'overloaded; don't retry' error to avoid combinatorial retry explosion."

**ASSESSMENT:**
Both sources agree on the core principle: don't retry at every layer. They differ slightly on implementation:
- Azure: Lower levels fail fast, higher levels handle retry
- Google: Retry only at immediate parent layer, propagate "don't retry" up the stack

**RESOLUTION:**
Both approaches achieve the same goal — preventing multiplicative retry amplification. Google's approach is more specific: each layer retries only failures from its direct dependency, then signals upstream whether further retry is appropriate.

## Contradiction 5: Timeout Value Philosophy

**SOURCE A (Topic Specification Example):**
Suggests timeouts calculated from latency budgets: PostgreSQL P99=150ms, Inventory P99=400ms, Payment P99=1.2s. User deadline ~3 seconds. Timeouts must fit within deadline.

**SOURCE B (gRPC Documentation):**
"Determine appropriate deadline from educated guess based on network latency, server processing time, validated by load testing."

**SOURCE C (Google SRE Addressing Cascading Failures):**
"Having deadlines several orders of magnitude longer than mean latency is usually bad. Example: 100ms mean with 100-second deadline causes thread exhaustion when 5% of requests hit deadline."

**ASSESSMENT:**
All sources agree — timeouts must be based on actual latency distributions plus safety margins, not arbitrary large values. The topic specification's example (P99 + margin) aligns with Google SRE guidance. The only nuance is how to validate: load testing vs production data.

**RESOLUTION:**
Consistent guidance across all sources:
1. Measure production latency distribution (P50, P95, P99)
2. Set timeout at roughly P99 + margin (not orders of magnitude above)
3. Validate via load testing at expected traffic levels
4. Account for total budget across all dependencies
5. Monitor and adjust as latency characteristics change

## Areas with NO Material Contradictions

The following claims show strong agreement across multiple independent sources:

1. **Timeouts prevent resource exhaustion** — All sources (SRE books, Microsoft patterns, PostgreSQL docs, gRPC) agree timeouts protect finite resources (threads, connections, memory) from being held indefinitely by slow operations.

2. **Retries require backoff with jitter** — Google SRE (randomized exponential backoff), AWS (exponential backoff and jitter), Azure (increasing delays, spread requests) all agree on this pattern.

3. **Idempotency required for safe retries** — Microsoft (Idempotent Consumer), RabbitMQ (idempotent consumer design), Azure Retry Pattern (consider idempotency) all emphasize that non-idempotent operations cannot be safely retried without deduplication.

4. **Circuit breakers complement retries** — Microsoft Circuit Breaker Pattern, Azure Retry Pattern, Resilience4j documentation, Google SRE (deadline propagation + cancellation) all show these patterns work together.

5. **Deadline propagation is essential** — gRPC documentation and Google SRE (Addressing Cascading Failures) both emphasize propagating absolute deadlines downstream, not using fixed timeouts at each layer.

6. **Database/queue timeouts are necessary** — PostgreSQL docs (statement_timeout, lock_timeout), RabbitMQ docs (consumer acknowledgements, heartbeats) both establish that network API timeouts alone are insufficient.