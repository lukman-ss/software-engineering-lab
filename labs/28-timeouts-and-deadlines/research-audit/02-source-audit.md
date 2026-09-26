# Source Audit: Lab 28 (Timeouts & Deadlines)

## Source 1
Claimed Title: Addressing Cascading Failures  
Claimed Publisher: Google SRE Book (Mike Ulrich)  
URL: https://sre.google/sre-book/addressing-cascading-failures/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Core chapter specifically explains cascading failure mechanics, thread exhaustion under bimodal latency, multi-layer retry amplification (4^3=64), deadline propagation, and backoff with jitter.  
Assessment: PASS  

---

## Source 2
Claimed Title: Deadlines  
Claimed Publisher: gRPC Documentation  
URL: https://grpc.io/docs/guides/deadlines/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Directly details client-side deadlines, server-side auto-cancellation, remaining timeout calculation across downstream calls, and clock skew mitigation.  
Assessment: PASS  

---

## Source 3
Claimed Title: PostgreSQL Documentation — Client Connection Defaults  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/runtime-config-client.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Contains authoritative definitions for `statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout`, and `transaction_timeout`.  
Assessment: PASS  

---

## Source 4
Claimed Title: Circuit Breaker Pattern  
Claimed Publisher: Microsoft Azure Architecture Center  
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/circuit-breaker  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Covers Closed, Open, Half-Open states, failure window tracking, and interaction with retry policies.  
Assessment: PASS  

---

## Source 5
Claimed Title: Idempotent Consumer Pattern  
Claimed Publisher: Microsoft Azure Architecture Center  
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/idempotent-consumer  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Covers at-least-once delivery duplication triggers, idempotency keys, and atomic persistence of deduplication markers.  
Assessment: PASS  

---

## Source 6
Claimed Title: Retry Pattern  
Claimed Publisher: Microsoft Azure Architecture Center  
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/retry  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Outlines retry policies (cancel, immediate, backoff), multi-layer retry anti-patterns, and fast-failure recommendations.  
Assessment: PASS  

---

## Source 7
Claimed Title: Exponential Backoff And Jitter  
Claimed Publisher: AWS Architecture Blog (Marc Brooker)  
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/  

Reachable: YES  
Source Type: SECONDARY (Engineering Blog)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Authoritative industry reference on jitter algorithms (Full Jitter, Equal Jitter, Decorrelated Jitter) preventing thundering herds.  
Assessment: PASS  

---

## Source 8
Claimed Title: Handling Overload  
Claimed Publisher: Google SRE Book  
URL: https://sre.google/sre-book/handling-overload/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Authoritative source for client-side throttling (K*accepts), retry budgeting (<10% ratio, 3 attempts), and criticality-based shedding.  
Assessment: PASS  

---

## Source 9
Claimed Title: Monitoring Distributed Systems  
Claimed Publisher: Google SRE Book  
URL: https://sre.google/sre-book/monitoring-distributed-systems/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Authoritative description of the Four Golden Signals (Latency, Traffic, Errors, Saturation) and tail latency distribution importance.  
Assessment: PASS  

---

## Source 10
Claimed Title: Health Endpoint Monitoring Pattern  
Claimed Publisher: Microsoft Azure Architecture Center  
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/health-endpoint-monitoring  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Explicitly details probe timeouts, liveness vs readiness separation, and preventing health checks from causing cascading failures.  
Assessment: PASS  

---

## Source 11
Claimed Title: Getting Started with Resilience4j  
Claimed Publisher: Resilience4j Documentation  
URL: https://resilience4j.readme.io/docs/getting-started-3  

Reachable: YES  
Source Type: SECONDARY (Library Documentation)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Verifies practical aspect composition (Retry -> CircuitBreaker -> RateLimiter -> TimeLimiter -> Bulkhead).  
Assessment: PASS  

---

## Source 12
Claimed Title: Reliability Guide  
Claimed Publisher: RabbitMQ Documentation  
URL: https://www.rabbitmq.com/docs/reliability  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Details consumer acknowledgements, publisher confirms, redelivered flag semantics, and idempotency needs in asynchronous queues.  
Assessment: PASS  

---

## Source 13
Claimed Title: PostgreSQL Documentation — Statement Timeout  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/queries-timeout.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Specifically details extended query protocol timing phases and log configuration.  
Assessment: PASS  
