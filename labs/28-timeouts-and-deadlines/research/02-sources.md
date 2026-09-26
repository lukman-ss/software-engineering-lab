# Research Sources: Timeouts — Jangan Biarkan Satu Dependency Lambat Menahan Seluruh Sistem

## Source 1

Title: Addressing Cascading Failures  
Publisher: Google SRE Book  
URL: https://sre.google/sre-book/addressing-cascading-failures/  
Published: 2017  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Core research question — cascading failures caused by timeouts and retries  

**Key evidence:** Explains how server overload is the most common cause of cascading failures. When cluster B fails, cluster A receives 1,200 QPS instead of 1,000 QPS, causing resources to be exhausted. The reduction in successfully handled requests spreads to other failure domains.

**Cascading failure scenario:** Frontend with 1,000 QPS → cluster B fails → cluster A receives 1,200 QPS → resources exhausted → servers crash → latency increases → missed deadlines → clients retry → further overload.

**Retry storm example:** Backend limit 10,000 QPS, frontend sends 10,100 QPS → 100 QPS fail → retries add 100 QPS → total 10,200 QPS → more failures → retry volume grows.

**Recommendation:** Use randomized exponential backoff when scheduling retries. Limit retries per request. Consider server-wide retry budget (e.g., only 60 retries per minute per process). Avoid retrying at multiple layers — a single request can create 4^3=64 attempts if database, backend, frontend, and JavaScript all retry.

## Source 2

Title: Deadlines  
Publisher: gRPC Documentation  
URL: https://grpc.io/docs/guides/deadlines/  
Published: 2025-07-07  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Deadline propagation and client/server timeout implementation  

**Key evidence:** By default, gRPC does not set a deadline which means clients can wait for responses effectively forever. Clients should always explicitly set realistic deadlines. Servers automatically cancel calls once the client deadline has passed.

**Deadline propagation:** When server receives RPC with 2-second deadline from client, spends 0.5s processing, then calls downstream service, the downstream RPC gets 1.5-second timeout (deadline minus elapsed time). This prevents downstream servers from doing work that won't benefit the original client.

**Clock skew protection:** gRPC converts deadlines to timeouts and deducts already-elapsed time, protecting against clock synchronization issues between servers.

**Server cancellation responsibility:** Server application must stop any activity spawned to service the RPC when deadline is exceeded. For long-running processes, periodically check if RPC has been cancelled.

## Source 3

Title: PostgreSQL Documentation — Client Connection Defaults  
Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/runtime-config-client.html  
Published: 2026-09-24  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Database timeouts — statement_timeout, lock_timeout, transaction_timeout  

**Key evidence:**
- `statement_timeout`: Abort statements taking more than specified time. Measured from command arrival to completion. Default 0 (disabled). Can be set per-session.
- `lock_timeout`: Abort statements waiting longer than specified time to acquire locks. Only triggers while waiting for locks, not during actual execution.
- `transaction_timeout`: Terminate sessions spanning longer than specified time in a transaction.
- `idle_in_transaction_session_timeout`: Terminate sessions idle within open transactions. Prevents long-running locks and vacuum issues.

**Recommendation:** Setting these in postgresql.conf affects all sessions — use ALTER ROLE or ALTER DATABASE for more granular control.

## Source 4

Title: Circuit Breaker Pattern  
Publisher: Microsoft Azure Architecture Center  
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker  
Published: 2025-02-05  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Circuit breaker integration with timeouts and retry logic  

**Key evidence:** Circuit breaker states:
- **Closed:** Request routed to operation. Counts failures. If threshold exceeded within time period, switches to Open.
- **Open:** Request fails immediately. No attempt to invoke service.
- **Half-Open:** Limited requests allowed. If successful, switches to Closed. If any fail, reverts to Open.

**Failure counter:** Time-based and automatically resets at periodic intervals to prevent overreaction to occasional failures.

**Success counter in Half-Open:** Records consecutive successful invocations. Reverts to Closed after specified number of successes.

**Recommended integration:** Retry pattern should be sensitive to circuit breaker exceptions and stop retrying when circuit breaker indicates fault isn't transient.

**Adaptive technique:** Modern approaches use AI/ML to dynamically adjust thresholds based on real-time traffic patterns and historical failure rates.

## Source 5

Title: Idempotent Consumer Pattern  
Publisher: Microsoft Azure Architecture Center  
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/idempotent-consumer  
Published: 2026-08-13  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Correctness implications of timeouts with non-idempotent operations  

**Key evidence:** At-least-once delivery guarantees mean messages can be delivered multiple times. Duplicates arise from:
- Producer retries (acknowledgment lost)
- Consumer crashes before acknowledgment
- Consumer crashes after processing but before acknowledgment

**Idempotent processing steps:**
1. Read message and extract deduplication key
2. Check deduplication store for key
3. If key exists, treat as duplicate, acknowledge, skip processing
4. If key doesn't exist, process message and record key in atomic operation

**Atomic commit:** Write deduplication marker and business side effects in same transaction to avoid crash window where side effects applied but key unrecorded.

**Natural idempotency:** Prefer naturally idempotent operations (upserts, absolute value writes, HTTP PUT) instead of deduplication bookkeeping.

## Source 6

Title: Retry Pattern  
Publisher: Microsoft Azure Architecture Center  
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/retry  
Published: 2024-07-18  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Retry strategies, performance impact, and idempotency considerations  

**Key evidence:**
- Retry strategies: Cancel, Retry immediately, Retry after delay (with exponential backoff)
- Aggressive retry policy with minimal delay and many retries can degrade busy services and affect application responsiveness
- If request still fails after significant retries, prevent further requests and fail immediately (use circuit breaker)
- Implement retry logic only where full context of failing operation is understood
- Avoid retrying at multiple layers — lower-level tasks should fail fast and report back
- Consider exception type — adjust delays based on failure nature
- Log early failures as informational, last failure as actual error

## Source 7

Title: Exponential Backoff And Jitter  
Publisher: AWS Architecture Blog  
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/  
Published: 2015-03-04  
Accessed: 2026-09-26  
Source Tier: 2  
Relevance: Retry backoff with jitter to prevent synchronized retries  

**Key evidence:** Without jitter, retries from multiple clients can synchronize, causing retry ripples that amplify overload. Randomization spreads retries over time.

**Amazon's approach:** Used for 8+ years as pillar for AWS SDK retry behavior. Most AWS SDKs support exponential backoff with jitter as part of standard retry behavior.

**Formula example:** Wait time = min(cap, random(0, base * 2^attempt)) where base is typically 100ms and cap prevents excessively long waits.

## Source 8

Title: Handling Overload  
Publisher: Google SRE Book  
URL: https://sre.google/sre-book/handling-overload/  
Published: 2017  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Load balancing, throttling, and request rejection strategies  

**Key evidence:**
- Model capacity in available resources (CPU, memory) rather than "queries per second" which doesn't account for request cost variation
- Per-customer limits allow misbehaving customers to be throttled while others remain unaffected
- Client-side throttling: Clients stop self-regulating when requests exceed K×accepts (default K=2). Prevents backend overload from rejecting many requests.

**Criticality levels:**
- CRITICAL_PLUS: Most critical, user-visible impact if fail
- CRITICAL: Default for production jobs, user-visible but less severe impact
- SHEDDABLE_PLUS: Partial unavailability expected (batch jobs)
- SHEDDABLE: Frequent partial or full unavailability expected

**Retry behavior:** When backend overloaded, clients can retry immediately only if small subset overloaded. If large subset overloaded, don't retry — return error immediately to caller. Per-request retry budget of 3 attempts. Per-client retry budget: only retry as long as retry ratio < 10%.

**Request retry counter:** Include counter in request metadata (starts at 0, increments on each retry). Backends can determine if other tasks are also overloaded by checking histogram of retry counts.

## Source 9

Title: Monitoring Distributed Systems  
Publisher: Google SRE Book  
URL: https://sre.google/sre-book/monitoring-distributed-systems/  
Published: 2017  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Four golden signals and monitoring timeout-related metrics  

**Key evidence:** Four golden signals:
- **Latency:** Time to service request. Track error latency separately from successful request latency.
- **Traffic:** Measure demand (e.g., HTTP requests/second for web services)
- **Errors:** Rate of failed requests (HTTP 500s, implicit failures, policy violations)
- **Saturation:** How "full" service is (CPU, memory, bandwidth)

**Tails matter:** If 1% of requests take 5 seconds, the rest are twice as fast as average. Measure distribution, not just mean.

**Monitoring philosophy:** Alert on symptoms, not causes. Page only when:
- Condition is urgent and actionable
- User is being negatively affected
- Can take meaningful action
- Not duplicating pages from other systems

**White-box vs black-box:** White-box detects imminent problems (e.g., slow database causing slow website). Black-box detects active problems (system isn't working right now).

## Source 10

Title: Health Endpoint Monitoring Pattern  
Publisher: Microsoft Azure Architecture Center  
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/health-endpoint-monitoring  
Published: 2026-09-24  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Health checks as part of timeout and reliability strategy  

**Key evidence:** Two types of health endpoints:
- **Liveness (`/healthz/live`):** Returns 200 if process can respond. No dependency queries. Detects unresponsive process.
- **Readiness (`/healthz/ready`):** Checks initialization and required dependencies. Returns 200 when application can serve requests, 503 if required checks fail.

**Timeouts in probes:** Each dependency check within probe should be bounded with a timeout to prevent cascading failure from slow dependencies.

**Circuit-breaker in probes:** Use circuit-breaker logic for unreliable dependencies during health checks to prevent probe itself from causing failure when dependency is slow.

**Caching trade-off:** Cache health status to avoid expensive checks on every probe, but cache TTL creates staleness window. Size TTL against acceptable detection delay.

## Source 11

Title: Getting Started with Resilience4j  
Publisher: Resilience4j Documentation  
URL: https://resilience4j.readme.io/docs/getting-started-3  
Published: 2026  
Accessed: 2026-09-26  
Source Tier: 2  
Relevance: Practical library implementation of timeouts, retries, circuit breakers  

**Key evidence:** Resilience4j provides:
- CircuitBreaker: Configurable failure rate threshold, sliding window size, wait duration in open state
- Retry: Max attempts, exponential backoff, configurable exception types
- Bulkhead: Max concurrent calls or max wait duration
- TimeLimiter: Timeout duration, configurable whether to cancel running future

**Aspect order:** Retry → CircuitBreaker → RateLimiter → TimeLimiter → Bulkhead → Function

**Metrics:** Automatic publishing to `/actuator/metrics` including circuit breaker calls, state, failure rate.

**Health indicators:** CircuitBreaker OPEN → DOWN, CLOSED → UP, HALF-OPEN → UNKNOWN

## Source 12

Title: Reliability Guide  
Publisher: RabbitMQ Documentation  
URL: https://www.rabbitmq.com/docs/reliability  
Published: 2026  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Message queue reliability and timeout mechanisms  

**Key evidence:** At-least-once delivery guarantees mean consumers must handle duplicate messages. Redelivered messages have redelivered flag set as hint (not guaranteed).

**Acknowledgements:** Use consumer acknowledgements to guarantee at-least-once delivery. Without acknowledgements, only at-most-once delivery is guaranteed.

**Publisher confirms:** Broker confirms messages once it has taken responsibility. Consumers may receive duplicate messages if acknowledgment lost during network failure.

**Consumer recommendations:** Design consumer to be idempotent rather than performing explicit deduplication. If redelivered flag not set, message has definitely not been seen before.

## Source 13

Title: PostgreSQL Documentation — Statement Timeout  
Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/queries-timeout.html  
Published: 2026-09-24  
Accessed: 2026-09-26  
Source Tier: 1  
Relevance: Statement-level timeout for long-running queries  

**Key evidence:** statement_timeout aborts any statement taking more than specified time. Measured from command arrival to completion. In extended query protocol, timeout starts when Parse/Bind/Execute/Describe arrives and is cancelled on Execute/Sync completion.

**Use cases:** Prevent runaway queries from consuming CPU, memory, locks, and connections. Essential for query that due to bug might take 4 minutes — timeout ensures resources freed.

**Log integration:** If log_min_error_statement set to ERROR or lower, timed-out statements also logged for debugging.

**Per-session setting:** Not recommended to set in postgresql.conf (affects all sessions). Use ALTER ROLE or ALTER DATABASE for granular control.