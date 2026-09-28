# Claim Audit: Research Stage — Timeouts and Deadlines

## Claim 1

Claim:
Slow downstream dependencies cause worker thread/process exhaustion and connection pool saturation (Little's Law: $L = \lambda W$), propagating failures to unrelated endpoints sharing the same process/pool.

Location:
`research/03-evidence.md`: Lines 5-21; `research/05-report.md`: Lines 24-37.

Evidence Provided:
Google SRE Book Chapter 22 direct quotations regarding thread saturation, queue accumulation, watchdog crash-loops, and cascading server failure under increased request duration.

Source:
Google SRE Book - Chapter 22: Addressing Cascading Failures (`https://sre.google/sre-book/addressing-cascading-failures/`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately grounded in queuing theory and verified SRE failure mechanics.

---

## Claim 2

Claim:
Network timeouts must be configured granularly across the connection lifecycle (Connect Timeout, Response/Read Timeout, Total Request Timeout); missing connect or read timeouts can cause sockets to hang for OS defaults (75–127s SYN retransmit or hours on idle TCP sockets).

Location:
`research/03-evidence.md`: Lines 25-43; `research/05-report.md`: Lines 40-54.

Evidence Provided:
POSIX socket behavior, TCP SYN retransmission specifications, and Go `net/http` / `Transport` timeout properties.

Source:
Google SRE Book & Standard POSIX / BSD Socket Architecture; Go `net/http` Docs.

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Highlights critical default behavior where standard libraries default socket timeouts to 0 (infinite).

---

## Claim 3

Claim:
Timeouts must be budgeted hierarchically based on production latency distributions ($T_{total} \ge \sum T_{deps\_P99} + T_{safety}$); arbitrary large timeouts allow bimodal latency spikes to exhaust thread capacity.

Location:
`research/03-evidence.md`: Lines 48-66; `research/05-report.md`: Lines 56-75.

Evidence Provided:
Google SRE Chapter 22 calculations showing that a 100s deadline on a 50 QPS service consumes 5,000 threads during tail latency spikes, causing an 80.4% failure rate on upstream frontends.

Source:
Google SRE Book - Chapter 22: Addressing Cascading Failures.

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Mathematical examples clearly explain why SLA budgeting is superior to static arbitrary numbers.

---

## Claim 4

Claim:
Retrying on timeout without exponential backoff and randomized jitter causes synchronized retry waves ($O(N^2)$ contention); Full Jitter substantially minimizes client contention and server load.

Location:
`research/03-evidence.md`: Lines 69-91; `research/05-report.md`: Lines 78-92.

Evidence Provided:
Marc Brooker's AWS Architecture Blog mathematical simulation and formulas (`sleep = rand_between(0, min(cap, base * 2 ** attempt))`), corroborated by Google SRE retry budget principles.

Source:
AWS Architecture Blog - Exponential Backoff And Jitter (`https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Mathematical formulation and comparison between unjittered, Full Jitter, and Decorrelated Jitter are verified.

---

## Claim 5

Claim:
A network timeout is an indeterminate state (`timeout != failure`); mutating operations (`POST /payment`) require idempotency keys and asynchronous status reconciliation to avoid duplicate execution.

Location:
`research/03-evidence.md`: Lines 94-114; `research/05-report.md`: Lines 94-107.

Evidence Provided:
Stripe developer documentation explaining that network drops after charge processing leave client state unknown, requiring the reuse of `Idempotency-Key` or webhook reconciliation.

Source:
Stripe API Documentation: Error Handling & Idempotent Requests.

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Essential distributed data integrity pattern correctly classified.

---

## Claim 6

Claim:
Propagating remaining deadline durations downstream across RPC/HTTP hops (e.g. `grpc-timeout`) prevents downstream servers from executing doomed work for cancelled/expired requests and eliminates clock skew issues.

Location:
`research/03-evidence.md`: Lines 117-138; `research/05-report.md`: Lines 110-123.

Evidence Provided:
gRPC Deadlines Guide detailing automatic server cancellation and the transformation of absolute deadlines into relative remaining timeouts before serialization over the wire.

Source:
gRPC Documentation - Deadlines (`https://grpc.io/docs/guides/deadlines/`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately clarifies the necessity of relative duration transmission over absolute timestamps.

---

## Claim 7

Claim:
Backend database layers require `statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout`, and `transaction_timeout`, and queue workers require execution timeouts and DLQs to prevent connection pool exhaustion and infinite worker lockups.

Location:
`research/03-evidence.md`: Lines 141-185; `research/05-report.md`: Lines 126-140.

Evidence Provided:
PostgreSQL 18 Chapter 19.11 runtime configuration parameters and Google SRE Chapter 22 worker watchdog specifications.

Source:
PostgreSQL 18 Documentation & Google SRE Book.

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Comprehensive multi-layer defense coverage beyond standard HTTP client timeouts.
