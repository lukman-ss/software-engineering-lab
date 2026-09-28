# Research Evidence: Timeouts and Deadlines

## Evidence 1: Slow Dependency Causes Resource Exhaustion & Cascading Failure

Claim:
A dependency that becomes slow (e.g. latency increases from 300ms to 60s) causes worker thread/process exhaustion and connection pool saturation, rapidly spreading failure across all endpoints in the service regardless of whether those endpoints touch the slow dependency.

Evidence:
According to Google SRE Book (Chapter 22 - Addressing Cascading Failures):
"Because requests take longer to handle, more requests are handled concurrently (up to a possible maximum capacity at which queuing may occur). This affects almost all resources, including memory, number of active threads (in a thread-per-request server model), number of file descriptors, and backend resources... Internal watchdogs in the server detect that the server isn't making progress, causing the servers to crash due to CPU starvation... Once a couple of servers crash on overload, the load on the remaining servers can increase, causing them to crash as well. The problem tends to snowball and soon all servers begin to crash-loop."

Source:
Google SRE Book - Chapter 22: Addressing Cascading Failures
URL: https://sre.google/sre-book/addressing-cascading-failures/
Publication Date: 2016
Confidence: HIGH
Corroborated By:
AWS Architecture Blog (Marc Brooker), Little's Law ($L = \lambda W$), where in-flight concurrency $L$ grows proportionally with latency $W$. If $W$ jumps 200x (300ms to 60s), concurrency demands exceed pool limits (e.g. 100 workers), halting all ingress traffic.

Notes:
Non-isolated thread pools / connection pools create failure domain leakage.

---

## Evidence 2: Granular Network Timeouts Across Connection Lifecycle

Claim:
Network timeouts are not a single monolithic number; HTTP/TCP connections require granular phases: Connect Timeout (DNS resolution + TCP handshake + TLS negotiation), Read/Response Timeout (waiting for server to send data/between packets), and Total/Request Timeout.

Evidence:
In network protocol mechanics and HTTP client implementations (e.g. Go `http.Client` with `Transport`, `net.Dialer.Timeout`, `TLSHandshakeTimeout`, `ResponseHeaderTimeout`), separate timers exist for connection establishment vs data reception. If a connect timeout is missing, a SYN packet to an unreachable or silently dropped IP can hang for the default OS TCP SYN retransmission timeout (75 to 127 seconds on Linux/BSD). If read timeout is missing, a server that accepts the connection but never sends bytes holds the client socket open indefinitely until TCP keepalive or OS timeout (often 2 hours).

Source:
Google SRE Book & Standard POSIX / BSD Socket Architecture
URL: https://sre.google/sre-book/addressing-cascading-failures/
Publication Date: 2016
Confidence: HIGH
Corroborated By:
gRPC Deadlines Guide & Go `net/http` documentation.

Notes:
Setting default library timeouts is dangerous because many default to infinity (`0` or none).

---

## Evidence 3: Timeout Budgeting vs Latency Distribution (P99)

Claim:
Timeout settings must not be arbitrary numbers (like "30 seconds just to be safe"); they must be derived from production latency distributions (P95/P99) and bounded by the client's overall deadline budget ($T_{total} \ge \sum T_{deps} + T_{safety}$). Setting timeouts too close to normal P95/P99 cuts off valid traffic, whereas setting them too high exhausts server worker capacity during bimodal latency degradation.

Evidence:
Google SRE Book:
"Having deadlines several orders of magnitude longer than the mean request latency is usually bad... With a 100-second deadline, 5% of requests would consume 5,000 threads (50 QPS * 100 seconds), but the frontend doesn’t have that many threads available. Assuming no other secondary effects, the frontend will only be able to handle 19.6% of the requests... resulting in an 80.4% error rate."
"As a server becomes overloaded, its responses to RPCs from its clients arrive later, which may exceed any deadlines those clients set. The work the server did to respond is then wasted, and clients may retry the RPCs, leading to even more overload."

Source:
Google SRE Book - Chapter 22: Addressing Cascading Failures
URL: https://sre.google/sre-book/addressing-cascading-failures/
Publication Date: 2016
Confidence: HIGH
Corroborated By:
AWS Architecture Blog (Marc Brooker), gRPC Deadlines Guide.

Notes:
A caller SLA of 3s with Postgres (P99=150ms), Inventory (P99=400ms), and Payment (P99=1200ms) requires strict budgeting: total latency P99 under normal conditions is $150 + 400 + 1200 = 1750\text{ ms}$, leaving $\approx 1250\text{ ms}$ buffer for overhead and network transit.

---

## Evidence 4: Naive Retries Trigger Retry Storms; Jittered Exponential Backoff Mitigates Load

Claim:
Retrying on timeout without exponential backoff and randomized jitter causes synchronized retry waves ("retry storms") that multiply downstream load ($N \times \text{retries}$), exacerbating outages. Full Jitter substantially reduces client work and server contention.

Evidence:
Marc Brooker (AWS Architecture Blog):
"With N clients contending, the total amount of work done by the system increases with $N^2$."
"The no-jitter exponential backoff approach is the clear loser. It not only takes more work, but also takes more time than the jittered approaches... The decision between 'Decorrelated Jitter' and 'Full Jitter' is less clear. The 'Full Jitter' approach uses less work, but slightly more time. Both approaches, though, present a substantial decrease in client work and server load."
Google SRE Book:
"Avoid amplifying retries by issuing retries at multiple levels: a single request at the highest layer may produce a number of attempts as large as the product of the number of attempts at each layer to the lowest layer (e.g. $4^3 = 64$ attempts)."

Source:
AWS Architecture Blog - Exponential Backoff And Jitter
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
Publication Date: 2015-03-04 (Updated 2023)
Confidence: HIGH
Corroborated By:
Google SRE Book - Chapter 22 (Ulrich, Sandler, Oshineye).

Notes:
Full Jitter sleep formulation: `sleep = rand_between(0, min(cap, base * 2 ** attempt))`.

---

## Evidence 5: Non-Idempotent Operations, Timeout Ambiguity, and Idempotency Keys

Claim:
A timeout does NOT mean that an operation failed; it merely indicates that the caller did not receive a response in time. In non-idempotent operations like payment processing (`POST /payment`), retrying blindly without an idempotency key can cause double charges. Correct architecture requires idempotency keys and out-of-band state reconciliation (query status or webhooks).

Evidence:
Stripe Official API Documentation:
"Treat the result of the API call as indeterminate. That is, don't assume that it succeeded or that it failed. Rely on webhooks for information about the outcome... You used an idempotency key for something unexpected, like replaying a request but passing different parameters... After you use an idempotency key, only reuse it for identical API calls."
If a network timeout occurs between client and payment processor, the charge may have been processed and funds debited, but the TCP ACK or HTTP response was dropped. Retrying with an Idempotency-Key guarantees the server returns the cached response rather than executing a second charge.

Source:
Stripe API Documentation: Error Handling & Idempotent Requests
URL: https://docs.stripe.com/error-handling.md?lang=go
Publication Date: 2026
Confidence: HIGH
Corroborated By:
Google SRE Book (Managing Critical State / Data Integrity), AWS Builders' Library.

Notes:
`timeout == unknown_state`, never `timeout == failed`.

---

## Evidence 6: Context and Deadline Propagation Prevents Wasted Downstream Work

Claim:
Propagating deadlines and cancellations downstream (e.g., using `grpc-timeout` or HTTP header translation) ensures downstream servers immediately abort or skip execution if the allocated deadline has elapsed or the client gave up.

Evidence:
gRPC Official Guide (Deadlines):
"A deadline is used to specify a point in time past which a client is unwilling to wait for a response from a server... A gRPC server deals with this situation by automatically cancelling a call (`CANCELLED` status) once a deadline set by the client has passed... gRPC converts the deadline to a timeout from which the already elapsed time is already deducted. This shields your system from any clock skew issues."
Google SRE Book:
"Had server B used deadline propagation, server C could immediately give up on the request because the 2-second deadline was exceeded. However, in this scenario, server C processes the request thinking it has 15 seconds to spare, but is not doing useful work... You don’t get credit for late assignments with RPCs."

Source:
gRPC Documentation - Deadlines
URL: https://grpc.io/docs/guides/deadlines/
Publication Date: 2025-07-07
Confidence: HIGH
Corroborated By:
Google SRE Book - Chapter 22: Addressing Cascading Failures.

Notes:
Propagating relative remaining duration avoids system clock drift / desynchronization errors between physical hosts.

---

## Evidence 7: Database-Level Timeout Guardrails

Claim:
Network API timeouts alone do not protect backend storage; databases require query statement timeouts (`statement_timeout`), lock acquisition timeouts (`lock_timeout`), and idle transaction timeouts (`idle_in_transaction_session_timeout`) to prevent lock starvation and connection pool exhaustion.

Evidence:
PostgreSQL 18 Documentation:
- `statement_timeout`: "Abort any statement that takes more than the specified amount of time... A value of zero (the default) disables the timeout... The timeout is measured from the time a command arrives at the server until it is completed by the server."
- `lock_timeout`: "Abort any statement that waits longer than the specified amount of time while attempting to acquire a lock on a table, index, row, or other database object."
- `idle_in_transaction_session_timeout`: "Terminate any session that has been idle (that is, waiting for a client query) within an open transaction for longer than the specified amount of time... prevents vacuuming away recently-dead tuples that may be visible only to this transaction; so remaining idle for a long time can contribute to table bloat."
- `transaction_timeout`: "Terminate any session that spans longer than the specified amount of time in a transaction."

Source:
PostgreSQL 18 Documentation - Section 19.11 Client Connection Defaults
URL: https://www.postgresql.org/docs/current/runtime-config-client.html
Publication Date: 2026-09-24
Confidence: HIGH
Corroborated By:
Database reliability engineering best practices and connection pool sizing specifications.

Notes:
Without `lock_timeout` or `statement_timeout`, a runaway query or lock wait holds a database connection, exhausting the pool for all other application threads.

---

## Evidence 8: Queue Worker Execution Limits and Poison Pills

Claim:
Background asynchronous jobs (e.g. report generation, webhooks) require explicit execution timeouts, maximum retry limits, and dead-letter queues (DLQ) to prevent infinite loops, resource leaks, and permanent worker thread lockups.

Evidence:
Google SRE Book:
"Internal watchdogs in the server detect that the server isn't making progress, causing the servers to crash due to CPU starvation... Limit retries per request. Don’t retry a given request indefinitely."
In queue engines (e.g. Celery `task_time_limit`, Sidekiq `timeout`, AWS SQS `VisibilityTimeout` + RedrivePolicy/DLQ), an un-timed job in an infinite loop or hung network call holds the worker thread permanently. Without a DLQ and retry limit, "poison pill" messages cause all worker processes to crash-loop continuously.

Source:
Google SRE Book - Chapter 22 & Queue Systems Architecture
URL: https://sre.google/sre-book/addressing-cascading-failures/
Publication Date: 2016
Confidence: HIGH
Corroborated By:
AWS Architecture Blog and standard worker concurrency models.

Notes:
Timeout on worker + max retries + DLQ is mandatory for queue safety.
