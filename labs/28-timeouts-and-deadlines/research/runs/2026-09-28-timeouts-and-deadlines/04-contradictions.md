# Contradictions & Trade-offs: Timeouts and Deadlines

## 1. Timeout Value: Safety Margin vs Early Failure

### Point of Tension
- Position A: Timeouts should be generous (e.g., 30s) to avoid dropping requests during temporary network spikes or heavy downstream computation.
- Position B: Timeouts should be tight (e.g., 2x-3x P99 latency) to release resources immediately and prevent thread pool exhaustion.

### Analysis & Resolution
Setting timeouts too high turns slow dependencies into system-wide outages due to thread and connection starvation (Little's Law). Setting timeouts too tight cuts off legitimate slow queries, increasing error rates and triggering retry storms.
Resolution: Timeouts must be budgeted against the end-to-end SLO/deadline. For interactive user requests, timeout = min(user_deadline - elapsed - buffer, 3x P99). For non-critical background jobs, larger timeouts with separate worker pools and circuit breakers are appropriate.

---

## 2. Absolute Deadline vs Relative Timeout Propagation

### Point of Tension
- Method A: Propagating absolute timestamps (e.g., `deadline = 2026-09-28T12:00:05.000Z`).
- Method B: Propagating relative remaining timeout durations (e.g., `grpc-timeout = 1800m`).

### Analysis & Resolution
Absolute timestamps are vulnerable to NTP clock skew between distributed servers (even a 100ms clock difference can prematurely kill or needlessly prolong a request).
Resolution: gRPC and modern distributed frameworks convert absolute deadlines to relative remaining durations (`timeout = deadline - now()`) before transmitting across the wire, eliminating clock skew risks (as documented in the gRPC Deadlines specification).

---

## 3. Retrying on Timeout: Immediate vs Idempotent & Budgeted

### Point of Tension
- Practice A: Automatically retry failed requests up to $N$ times immediately upon timeout.
- Practice B: Treat timeout as indeterminate state; never retry non-idempotent operations without an idempotency key; use exponential backoff with Full Jitter and a process-wide retry budget.

### Analysis & Resolution
Practice A is the primary cause of retry storms and double charges in financial/order systems.
Resolution: Practice B is mandatory for production distributed systems. Retries must only be executed if:
1. The operation is safe/idempotent (or guarded by `Idempotency-Key`).
2. Jittered exponential backoff is applied.
3. The remaining request deadline budget allows another round trip.
4. Process-wide retry budget has not been exhausted.

---

## Summary of Material Contradictions

No material contradictions exist among authoritative engineering standards (Google SRE, AWS Architecture, gRPC, PostgreSQL). All agree that:
1. Timeouts and deadlines are resource protection mechanisms.
2. Latency without deadlines leads to cascading failure.
3. Jittered backoff and idempotency keys are non-negotiable for safe retries.
