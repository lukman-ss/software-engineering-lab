# Research Plan: Timeouts and Deadlines — Resiliency and Resource Protection in Distributed Systems

## Research Topic

Timeouts and Deadlines — Jangan Biarkan Satu Dependency Lambat Menahan Seluruh Sistem. Why unconstrained downstream latencies cause thread/worker exhaustion and cascading outages, how timeout budgets, context/deadline propagation, backoff jitter, circuit breaking, idempotency keys, database statement timeouts, and queue execution limits prevent system collapse.

## Objective

Investigate and verify the technical claims made in the lab topic specification:

1. Slow dependencies (high latency) cause worker/connection pool exhaustion leading to cascading failure.
2. HTTP/Network timeouts consist of multiple phases (DNS, TCP Connect, TLS Handshake, Write, Wait/Read Response, Total Timeout).
3. Timeout budget calculation based on overall SLA/SLO and P99 latency distribution prevents wasting resources.
4. Retry storms occur when client timeouts trigger aggressive retries without jitter, backoff, or budgets.
5. Non-idempotent operations (e.g. POST /payment) subject to timeout require idempotency keys and state reconciliation to prevent double-charging/double-processing.
6. Deadline propagation across microservices (e.g., gRPC metadata headers / HTTP `Timeout-Access` or `grpc-timeout`) prevents waste on expired requests.
7. Database query guards (PostgreSQL `statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout`) protect database CPU, locks, and connection pools.
8. Queue workers require execution timeouts, retry limits, and dead-letter queues (DLQ) to prevent infinite loops and worker lockup.

## Research Questions

### RQ1: Mechanisms of Resource Exhaustion & Cascading Failures
- How does latency (e.g. 60s response instead of 300ms) exhaust worker pools (Little's Law: $L = \lambda W$)?
- What is the cascading failure mechanism where a slow dependency causes unrelated endpoints to fail?

### RQ2: Network Timeout Anatomy & Granular Controls
- What are the distinct phases of an HTTP/TCP connection lifecycle?
- What are the specific roles and configuration parameters of Connect Timeout, Read Timeout, Write Timeout, and Total Request Timeout?

### RQ3: Timeout Budgeting & Latency Distributions (P50/P95/P99)
- How should a timeout budget be constructed across downstream services?
- Why using fixed large timeouts (e.g. 30s for all) or unaligned timeouts breaks SLA/SLO expectations?

### RQ4: Retry Storms, Exponential Backoff, and Jitter
- What causes retry storms when downstream services experience latency spikes?
- How do Exponential Backoff, Full Jitter, Equal Jitter, Decorrelated Jitter, and Retry Budgets mitigate retry amplification?

### RQ5: Non-Idempotent Operations & Timeout Ambiguity
- Why does a timeout NOT imply failure (the "State Ambiguity Problem")?
- How do Idempotency Keys and status reconciliation patterns prevent double-charging on payment gateway timeouts?

### RQ6: Context & Deadline Propagation
- How does deadline propagation work in distributed tracing and gRPC (`grpc-timeout`) or HTTP headers?
- What happens when a client cancels or times out while downstream microservices are processing?

### RQ7: Database & Worker Timeout Controls
- What PostgreSQL timeout settings (`statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout`, `transaction_timeout`) guard the storage layer?
- How do worker queue systems implement execution timeouts, retry limits, and Dead Letter Queues (DLQ)?

## Search Strategy

- Primary: Official Documentation (PostgreSQL Docs, gRPC Docs, Go `context` package, AWS Builders' Library / Architecture Blog, Google SRE Book).
- Secondary: Industry engineering publications (Martin Fowler, Netflix TechBlog, Uber Engineering).
- Keywords: `timeouts`, `deadlines`, `grpc-timeout`, `statement_timeout`, `cascading failure`, `retry storm`, `exponential backoff jitter`, `idempotency key`, `circuit breaker`.

## Expected Primary Sources

1. Google SRE Book - Chapter 22: Addressing Cascading Failures (https://sre.google/sre-book/addressing-cascading-failures/)
2. AWS Architecture Blog - Exponential Backoff And Jitter (Marc Brooker, 2015/2023)
3. gRPC Documentation - Deadlines & Deadline Propagation (https://grpc.io/docs/guides/deadlines/)
4. PostgreSQL 18 Documentation - Client Connection Defaults & Statement Behavior (`statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout`)

## Risks / Unknowns

- Varied implementation specifics across programming languages (e.g., Go `context.WithDeadline` vs Java `gRPC` context propagation vs HTTP client timeouts).
- Differentiating between connection pool timeouts and socket read timeouts.
