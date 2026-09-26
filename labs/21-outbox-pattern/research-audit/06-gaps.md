# Research Gap Analysis

## Gaps Review

The research report in `labs/21-outbox-pattern/research/05-report.md` and evidence ledger in `research/03-evidence.md` provide solid coverage of the core transactional outbox pattern mechanics. The research agent appropriately captured open questions in `research/06-open-questions.md`.

Below is the auditor assessment of documented gaps and potential nuances:

---

## Gap 1: Schema Evolution for Long-Lived Outbox Payloads

Type:
MISSING_CASE

Severity:
LOW

Location:
`research/06-open-questions.md` (Question 4)

Problem:
When payload models change over time, downstream consumers processing backlogged or replayed outbox events may encounter serialization mismatches unless schema registry (Avro/Protobuf) or forward/backward compatible JSON schemas are enforced.

Required Revision:
None required for base lab approval. Adequately identified in open questions and Debezium Avro docs.

Can Be Approved Without Fix:
YES

---

## Gap 2: High-Throughput Partitioning & Outbox Table Contention

Type:
MISSING_CASE

Severity:
LOW

Location:
`research/06-open-questions.md` (Question 2)

Problem:
In high-throughput relational databases (10k+ tx/sec), a single monotonically incrementing or polled outbox table can suffer from lock contention, write hotspots, and vacuum/bloat overhead.

Required Revision:
None required for conceptual understanding and lab scope. Identified as an open question for scale-up scenarios.

Can Be Approved Without Fix:
YES

---

## Gap 3: Dead-Letter Queue (DLQ) and Relay Retry Policies

Type:
MISSING_CASE

Severity:
LOW

Location:
`research/06-open-questions.md` (Question 1)

Problem:
If the broker is unreachable or a message payload fails validation permanently, the relay must specify retry backoff and poison pill routing to avoid stalling the publisher queue.

Required Revision:
Covered as an open research question; not blocking for core pattern foundation.

Can Be Approved Without Fix:
YES
