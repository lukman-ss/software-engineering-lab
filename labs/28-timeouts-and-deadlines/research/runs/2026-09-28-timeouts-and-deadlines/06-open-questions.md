# Open Questions: Timeouts and Deadlines

## 1. Adaptive & Dynamic Timeout Mechanisms
- Unanswered Question: How do dynamic timeout algorithms (e.g. dynamically calculating timeout as $P99_{moving\_window} + k \cdot MAD$) perform under sudden traffic spikes compared to static budgeted timeouts?
- Next Step: Evaluate Netflix Hystrix / Resilience4j / Envoy adaptive timeout implementations.

## 2. Distributed Tracing & W3C Header Standardization for HTTP Deadlines
- Unanswered Question: While gRPC uses `grpc-timeout`, HTTP/1.1 and HTTP/2 lack a single universally adopted W3C standard header for deadline propagation (e.g., `Timeout-Access`, `Request-Timeout`, `X-Request-Deadline`).
- Next Step: Investigate OpenTelemetry Context Propagation specification for custom baggage propagation of request deadlines across REST endpoints.

## 3. Database Connection Pool Wait Timeout Tuning
- Unanswered Question: What is the optimal ratio between pool size, connection acquisition timeout (`pool_timeout`), and SQL statement timeout (`statement_timeout`) to prevent thread starvation under database node degradation?
- Next Step: Conduct benchmark testing in PostgreSQL with HikariCP/pgxpool under synthetic lock contention and latency injection.
