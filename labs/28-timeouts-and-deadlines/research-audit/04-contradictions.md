# Contradictions & Trade-offs Audit: Research Stage — Timeouts and Deadlines

## Contradiction 1

Statement A:
"Timeouts should be generous (e.g., 30s) to avoid dropping requests during temporary network spikes or heavy downstream computation." (`research/04-contradictions.md`: Lines 6-7)

Statement B:
"Timeouts should be tight (e.g., 2x-3x P99 latency) to release resources immediately and prevent thread pool exhaustion." (`research/04-contradictions.md`: Line 7)

Type:
INTERNAL (Trade-off Analysis)

Impact:
Provides operational clarity on how to set timeout thresholds based on SLOs and workload types (interactive vs background).

Assessment:
PASS. Correctly analyzed and resolved via hierarchical deadline budgeting ($T_{timeout} = \min(\text{user\_deadline} - \text{elapsed} - \text{buffer}, 3 \times P99)$).

---

## Contradiction 2

Statement A:
"Propagating absolute timestamps across distributed systems (e.g. `deadline = 2026-09-28T12:00:05.000Z`)." (`research/04-contradictions.md`: Line 18)

Statement B:
"Propagating relative remaining timeout durations (e.g. `grpc-timeout = 1800m`)." (`research/04-contradictions.md`: Line 19)

Type:
SOURCE_CONFLICT / BEST_PRACTICE_RESOLUTION

Impact:
Identifies NTP clock skew vulnerability in distributed systems.

Assessment:
PASS. Correctly resolved by adopting relative remaining duration serialization (`timeout = deadline - now()`), as specified in gRPC wire protocol standards.

---

## Contradiction 3

Statement A:
"Automatically retry failed requests up to $N$ times immediately upon timeout." (`research/04-contradictions.md`: Line 29)

Statement B:
"Treat timeout as indeterminate state; never retry non-idempotent operations without an idempotency key; use exponential backoff with Full Jitter and a process-wide retry budget." (`research/04-contradictions.md`: Lines 30-31)

Type:
INTERNAL (Anti-pattern vs Resilience Pattern)

Impact:
Prevents retry storm outages and duplicate transactional state creation.

Assessment:
PASS. Highlighting naive immediate retries as an anti-pattern and establishing Practice B as mandatory for production architectures aligns with authoritative SRE standards.

---

## Overall Assessment

No material un-analyzed contradictions or invalid trade-offs exist across the research documentation. All identified tensions are explicitly resolved with evidence-based industry standards.
