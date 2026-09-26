# Evidence Collection: Timeouts — Jangan Biarkan Satu Dependency Lambat Menahan Seluruh Sistem

## Evidence 1: Cascading Failure Mechanism

**Claim:** A single slow dependency can cause cascading failures by exhausting resources across multiple services in a call chain.

**Evidence:** When cluster B fails in a load-balanced setup, cluster A receives increased traffic (1,200 QPS instead of 1,000 QPS), causing resource exhaustion. As servers in A crash or miss deadlines, the rate of successfully handled requests drops below normal, potentially triggering overload in other clusters through load balancing redistribution.

**Source:** Google SRE Book - Addressing Cascading Failures (https://sre.google/sre-book/addressing-cascading-failures/)

**URL:** https://sre.google/sre-book/addressing-cascading-failures/

**Confidence:** HIGH

**Corroborated By:** Source 8 (Handling Overload)

**Notes:** This establishes the fundamental problem that timeouts aim to prevent — when one service slows down, it can trigger resource exhaustion that spreads through the system.

## Evidence 2: Retry Storm Amplification

**Claim:** Retries without proper backoff can amplify load during partial outages, creating retry storms that overwhelm already-struggling dependencies.

**Evidence:** Backend service has a limit of 10,000 QPS per task. Frontend calls at 10,100 QPS, overloading by 100 QPS which backend rejects. Failed requests are retried every second, adding to the request stream. Backend now receives 10,200 QPS (200 QPS failing), leading to more retries. The retry volume grows: 100 → 200 → 300 QPS of failed requests. With retries at multiple layers (database, backend, frontend, JavaScript), a single user action can create 4^3=64 attempts on the database.

**Source:** Google SRE Book - Addressing Cascading Failures (https://sre.google/sre-book/addressing-cascading-failures/)

**URL:** https://sre.google/sre-book/addressing-cascading-failures/

**Confidence:** HIGH

**Corroborated By:** Source 7 (Exponential Backoff And Jitter)

**Notes:** Demonstrates mathematically how naive retry logic can exponentially increase load on struggling services, turning a minor issue into a major outage.

## Evidence 3: Deadline Propagation Benefits

**Claim:** Propagating deadlines downstream prevents services from performing useless work after the original request deadline has passed.

**Evidence:** Client sends request with 2-second deadline to User Server. User Server spends 0.5 seconds before calling Billing Server. With deadline propagation, Billing Server receives request with 1.5-second timeout (2.0 - 0.5 = 1.5). If Billing Server takes 2 seconds to process, it times out at 1.5 seconds rather than doing useless work that arrives after the original 2-second deadline has passed.

**Source:** gRPC Documentation - Deadlines (https://grpc.io/docs/guides/deadlines/)

**URL:** https://grpc.io/docs/guides/deadlines/

**Confidence:** HIGH

**Corroborated By:** Source 1 (Addressing Cascading Failures) discusses deadline propagation concept

**Notes:** This prevents resource waste on work that users have already abandoned due to timeout, improving overall system efficiency.

## Evidence 4: Database Statement Timeouts

**Claim:** Database statement timeouts prevent runaway queries from consuming resources indefinitely, protecting against resource exhaustion from long-running queries.

**Evidence:** PostgreSQL statement_timeout parameter aborts any statement taking more than the specified amount of time. The timeout is measured from command arrival at server to completion by server. If multiple SQL statements appear in a single query, timeout applies to each separately. Setting statement_timeout prevents queries from consuming CPU, memory, locks, and connections when they would otherwise run for extended periods due to bugs or unexpected data volumes.

**Source:** PostgreSQL Documentation - Client Connection Defaults (https://www.postgresql.org/docs/current/runtime-config-client.html)

**URL:** https://www.postgresql.org/docs/current/runtime-config-client.html

**Confidence:** HIGH

**Corroborated By:** Source 13 (PostgreSQL Statement Timeout)

**Notes:** Essential guardrail for database reliability — without it, a single slow query can tie up connection pool resources and cause application-wide timeouts.

## Evidence 5: Circuit Breaker Integration with Retries

**Claim:** Circuit breakers should be combined with retry logic, where retry policies are sensitive to circuit breaker state to prevent retrying when failures are not transient.

**Evidence:** Circuit breaker has three states: Closed (normal operation), Open (immediate failure), Half-Open (limited trial requests). The Retry pattern should be implemented such that when circuit breaker indicates a fault isn't transient (Open state), retry logic stops attempting retries. An application can combine patterns by using Retry to invoke operations through a circuit breaker, but retry logic must respect circuit breaker exceptions.

**Source:** Microsoft Azure Architecture Center - Circuit Breaker Pattern (https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker)

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker

**Confidence:** HIGH

**Corroborated By:** Source 6 (Retry Pattern) mentions circuit breaker integration

**Notes:** Prevents retry storms by stopping retry attempts when circuit breaker determines service is likely to fail, rather than blindly retrying based on time-based policies alone.

## Evidence 6: Idempotency Requirement for Safe Retries

**Claim:** Retrying non-idempotent operations without deduplication can cause incorrect results like duplicate charges, requiring idempotency keys for safe retry mechanisms.

**Evidence:** In at-least-once delivery systems, consumers may receive the same message multiple times due to producer retries (acknowledgment lost), consumer crashes before acknowledgment, or crashes after processing but before acknowledgment. Without idempotency protection, retrying a payment request could cause double charges. Idempotent consumers use deduplication stores to track processed messages and skip duplicates, with the key insight being that the deduplication marker and business side effects must be committed atomically to avoid crash windows where side effects apply but marker isn't recorded.

**Source:** Microsoft Azure Architecture Center - Idempotent Consumer Pattern (https://learn.microsoft.com/en-us/azure/architecture/patterns/idempotent-consumer)

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/idempotent-consumer

**Confidence:** HIGH

**Corroborated By:** Source 11 (Resilience4j) shows retry + circuit breaker integration

**Notes:** Critical correctness consideration — timeouts alone don't solve the problem of not knowing whether an operation succeeded when timing out; idempotency is required for safe retries.

## Evidence 7: Exponential Backoff with Jitter

**Claim:** Retry delays should use exponential backoff with jitter to prevent synchronized retry storms that amplify load on struggling services.

**Evidence:** Without jitter, retry attempts from multiple clients can synchronize, causing "retry ripples" where all clients retry at the same moment, amplifying overload. With jitter, delays are randomized (e.g., base * 2^attempt * random(0,1)), spreading retry attempts over time. Amazon's AWS SDKs have used this approach for 8+ years as a standard pattern for resilient remote service calls.

**Source:** AWS Architecture Blog - Exponential Backoff And Jitter (https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/)

**URL:** https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/

**Confidence:** HIGH

**Corroborated By:** Source 1 (Addressing Cascading Failures) recommends randomized exponential backoff

**Notes:** Key insight is that randomization prevents thundering herd problems where coordinated retries overwhelm recovering services.

## Evidence 8: Overload Handling and Criticality

**Claim:** During overload, services should implement client-side throttling and criticality-based request rejection to maintain availability for important traffic while shedding less critical load.

**Evidence:** Rather than modeling capacity as "queries per second," measure actual resource consumption (CPU, memory). Implement per-customer quotas so misbehaving customers are throttled while others unaffected. Use client-side throttling where clients self-regulate when requests exceed K×accepts (K=2 default). Implement four criticality levels: CRITICAL_PLUS (most critical), CRITICAL (production jobs), SHEDDABLE_PLUS (batch jobs), SHEDDABLE (frequent partial/full unavailability). Services reject lower criticality requests first during overload.

**Source:** Google SRE Book - Handling Overload (https://sre.google/sre-book/handling-overload/)

**URL:** https://sre.google/sre-book/handling-overload/

**Confidence:** HIGH

**Corroborated By:** Source 1 (Addressing Cascading Failures) discusses overload causes

**Notes:** Provides proactive overload management strategy that works alongside timeouts to prevent cascading failures before they start.

## Evidence 9: Four Golden Signals for Timeout Monitoring

**Claim:** Monitoring latency, traffic, errors, and saturation (the four golden signals) provides essential visibility into timeout-related issues and system health.

**Evidence:** Latency should distinguish between successful and failed request latency. Tracking error latency is important because slow errors are worse than fast errors. Traffic measures demand (HTTP requests/second). Errors count failed requests (explicit like HTTP 500s, implicit like wrong content, policy violations). Saturation measures how constrained resources are (CPU in CPU-constrained systems, memory in memory-constrained systems). For detecting saturation tails, measure latency distribution rather than mean — if 1% of requests take 5 seconds while average is 100ms, the rest are actually faster than average suggests.

**Source:** Google SRE Book - Monitoring Distributed Systems (https://sre.google/sre-book/monitoring-distributed-systems/)

**URL:** https://sre.google/sre-book/monitoring-distributed-systems/

**Confidence:** HIGH

**Corroborated By:** Source 10 (Health Endpoint Monitoring) discusses monitoring integration

**Notes:** Essential for detecting emerging timeout issues before they cause outages — rising P99 latency often precedes increased timeout rates.

## Evidence 10: Health Endpoints with Dependency Timeouts

**Claim:** Health check endpoints should separate liveness (process responsiveness) from readiness (dependency availability) and apply timeouts to dependency checks to prevent monitoring from causing cascading failures.

**Evidence:** Implement two health endpoints:
- Liveness (`/healthz/live`): Returns 200 if process can respond. No dependency queries. Detects unresponsive process.
- Readiness (`/healthz/ready`): Checks initialization and required dependencies. Returns 200 when application can serve requests, 503 if any required dependency check fails.
Each dependency check within the readiness probe should be bounded with a timeout to prevent a slow dependency from causing the health check itself to hang and trigger false unhealthy reports. For unreliable dependencies, apply circuit-breaker logic within health checks to prevent the monitoring mechanism from exacerbating issues.

**Source:** Microsoft Azure Architecture Center - Health Endpoint Monitoring Pattern (https://learn.microsoft.com/en-us/azure/architecture/patterns/health-endpoint-monitoring)

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/health-endpoint-monitoring

**Confidence:** HIGH

**Corroborated By:** Source 2 (gRPC Deadlines) shows similar timeout concept

**Notes:** Prevents monitoring tools from creating false alarms and unintentionally causing cascading failures when checking on already-stressed systems.

## Evidence 11: Practical Library Implementation (Resilience4j)

**Claim:** Production-ready libraries provide integrated implementations of timeouts, retries, circuit breakers, and related patterns with configurable policies.

**Evidence:** Resilience4j provides:
- CircuitBreaker: Configurable failure rate threshold (default 50%), sliding window size, wait duration in open state
- Retry: Max attempts, wait duration, exponential backoff multiplier, retryable/ignorable exception lists
- Bulkhead: Max concurrent calls (semaphore-based) or max wait duration (thread pool-based)
- TimeLimiter: Timeout duration, cancel-running-future flag
- Automatic metrics publishing to `/actuator/metrics` including circuit breaker calls, state, failure rate
- Health indicator integration: CircuitBreaker OPEN → DOWN, CLOSED → UP, HALF-OPEN → UNKNOWN
- Aspect ordering: Retry → CircuitBreaker → RateLimiter → TimeLimiter → Bulkhead → Function

**Source:** Resilience4j Documentation - Getting Started (https://resilience4j.readme.io/docs/getting-started-3)

**URL:** https://resilience4j.readme.io/docs/getting-started-3

**Confidence:** HIGH

**Corroborated By:** Sources 4, 5, 6 show individual patterns

**Notes:** Shows how patterns integrate in practice — timeouts (TimeLimiter) work with circuit breakers and retries to provide layered protection.

## Evidence 12: Message Queue Reliability and Timeouts

**Claim:** Message queue systems require consumers to handle duplicate deliveries and use acknowledgements with timeouts to prevent resource exhaustion from unacknowledged messages.

**Evidence:** Most message brokers (RabbitMQ, Kafka, etc.) provide at-least-once delivery, meaning messages can be redelivered. Consumers should use acknowledgements to guarantee message processing. Without acknowledgements, only at-most-once delivery is guaranteed. Publisher confirms provide broker-to-producer acknowledgment. If consumer crashes after processing but before acknowledgement, message will be redelivered. Idempotent consumer design is recommended over explicit deduplication. The redelivered flag on messages is a hint (not guarantee) that message may have been processed before.

**Source:** RabbitMQ Documentation - Reliability Guide (https://www.rabbitmq.com/docs/reliability)

**URL:** https://www.rabbitmq.com/docs/reliability

**Confidence:** HIGH

**Corroborated By:** Source 5 (Idempotent Consumer) shows pattern application

**Notes:** Extends timeout/retry/idempotency concepts to asynchronous messaging systems where timing out while waiting for acknowledgements creates duplicate processing risks.

## Evidence 13: Statement Timeout Implementation Details

**Claim:** PostgreSQL statement_timeout provides precise control over maximum query execution time, measured from command arrival to completion, with logging integration for debugging.

**Evidence:** The statement_timeout parameter aborts any statement taking more than specified time. If value specified without units, it's taken as milliseconds. Timeout measured from time command arrives at server until completed by server. In extended query protocol, timeout starts when Parse/Bind/Execute/Describe message arrives and is cancelled by completion of Execute or Sync message. If log_min_error_statement set to ERROR or lower, timed-out statements also logged. Setting in postgresql.conf not recommended as it affects all sessions — use ALTER ROLE or ALTER DATABASE for session-specific configuration.

**Source:** PostgreSQL Documentation - Statement Timeout (https://www.postgresql.org/docs/current/runtime-config-client.html#runtime-config-client-statement-timeout)

**URL:** https://www.postgresql.org/docs/current/runtime-config-client.html

**Confidence:** HIGH

**Corroborated By:** Source 4 (Circuit Breaker) and Source 6 (Retry) show complementary patterns

**Notes:** Essential database-level timeout that complements application-level timeouts to prevent resource exhaustion from long-running queries.