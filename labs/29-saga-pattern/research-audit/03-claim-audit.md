# Claim Audit: Saga Pattern Research

Target Lab: `labs/29-saga-pattern`
Date: 2026-09-29

---

## Claim 1

Claim:
Saga Pattern was introduced by Garcia-Molina & Salem (ACM SIGMOD 1987) for long-lived database transactions using compensating transactions; later adapted to microservices.

Location:
`research/05-report.md` → Finding 1; `research/03-evidence.md` → Evidence 1

Evidence Provided:
Cornell mirror PDF; DOI 10.1145/62224.62226; Temporal Blog footnote; citation chain across Microsoft and Richardson.

Source:
- https://www.cs.cornell.edu/andru/cs711/2002fa/reading/sagas.pdf

Source Actually Supports Claim: PARTIAL
The Cornell PDF is reachable and is the 1987 paper. Direct text extraction failed due to LZW/scanned image encoding. Attribution and DOI are cross-confirmed by multiple Tier-1 and Tier-2 sources. The claim is accepted on citation chain verification, not direct PDF text parsing.

Classification: FACT

Severity: LOW
The historical attribution is academically established and consistent across all sources.

Notes:
The verbatim definition from the 1987 paper was acknowledged as unextracted in `06-open-questions.md`. This is transparent. No fabrication detected.

---

## Claim 2

Claim:
A saga is a sequence of local transactions; each updates its own database and publishes an event/message to trigger the next; on failure, compensating transactions undo prior steps semantically (not ACID rollback).

Location:
`research/05-report.md` → Finding 2; `research/03-evidence.md` → Evidence 2, 6

Evidence Provided:
Verbatim quotation from microservices.io: "A saga is a sequence of local transactions. Each local transaction updates the database and publishes a message or event to trigger the next local transaction..."
Microsoft: "breaks them into a sequence of local transactions... performs a series of compensating transactions to reverse the changes."

Source:
- https://microservices.io/patterns/data/saga.html
- https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes:
Quotations are accurate. Both sources confirmed reachable and contain the cited text.

---

## Claim 3

Claim:
Two coordination styles: Choreography (each service publishes domain events triggering the next) and Orchestration (central orchestrator sends commands to participants).

Location:
`research/05-report.md` → Finding 2; `research/03-evidence.md` → Evidence 3

Evidence Provided:
Microservices.io verbatim: "There are two ways of coordination sagas: Choreography... Orchestration..."
Microsoft: "The two typical saga implementation approaches are choreography and orchestration."

Source:
- https://microservices.io/patterns/data/saga.html
- https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Fully supported by both primary-tier sources verified as reachable.

---

## Claim 4

Claim:
Choreography benefits: loose coupling, no coordinator, no SPOF, good for simple workflows. Drawbacks: spaghetti events, cyclic dependencies, hard to track/debug.
Orchestration benefits: clear flow, no cyclic dependencies. Drawbacks: SPOF, coordination complexity.

Location:
`research/05-report.md` → Finding 2; `research/03-evidence.md` → Evidence 4, 5

Evidence Provided:
Microsoft benefit/drawback tables quoted verbatim. Temporal analogy and framing corroborated.

Source:
- https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- https://temporal.io/blog/to-choreograph-or-orchestrate-your-saga-that-is-the-question

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Microsoft benefit/drawback tables verified from live page fetch.

---

## Claim 5

Claim:
Transactions in a saga are classified as: Compensable (can be undone), Pivot (point of no return), and Retryable (idempotent post-pivot, must complete).

Location:
`research/05-report.md` → Finding 4; `research/03-evidence.md` → Evidence 7

Evidence Provided:
Microsoft: "Pivot transactions serve as the point of no return... After a pivot transaction succeeds, compensable transactions are no longer relevant." "Retryable transactions follow the pivot transaction. Retryable transactions are idempotent..."

Source:
- https://learn.microsoft.com/en-us/azure/architecture/patterns/saga (ONLY)

Source Actually Supports Claim: YES (single source)

Classification: FACT (implementation-specific taxonomy detail)

Severity: MEDIUM
This taxonomy is detailed exclusively in the Microsoft Azure Architecture source. Microservices.io does not enumerate it on its public page; Richardson's book (behind paywall) is not directly cited. The research correctly reported confidence as MEDIUM.

Notes:
Not a fabricated claim. Source confirmed containing this taxonomy. However, the research report could note this is Microsoft's framing and may not be universally standardized terminology.

---

## Claim 6

Claim:
2PC blocks resources across services, reduces availability, and creates tight coupling; Sagas trade ACID isolation for availability and eventual consistency.

Location:
`research/05-report.md` → Finding 5; `research/03-evidence.md` → Evidence 8

Evidence Provided:
Microservices.io Forces: "2PC is not an option."
Microsoft: "architectures that rely on... traditional transaction models like two-phase commit protocol, are often better suited for the Saga pattern."
AWS: "There are long-lived transactions and you don't want other microservices to be blocked..."

Source:
- https://microservices.io/patterns/data/saga.html
- https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes:
All three primary sources confirmed reachable and verified.

---

## Claim 7

Claim:
Dual-write problem: atomically updating DB and publishing to broker is impossible without distributed XA transactions (which Kafka does not support). Transactional Outbox solves this by writing business data + outbox record in one local ACID transaction; CDC relays to broker.

Location:
`research/05-report.md` → Finding 6; `research/03-evidence.md` → Evidence 9

Evidence Provided:
Debezium verbatim: "So how can this situation be avoided? The answer is to only modify one of the two resources... The idea of this approach is to have an 'outbox' table..."
Outbox schema: `id uuid`, `aggregatetype`, `aggregateid`, `type`, `payload jsonb`.
Microservices.io: "it cannot use the traditional mechanism of a distributed transaction... Instead, it must use one of the patterns listed below."

Source:
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
- https://microservices.io/patterns/data/saga.html

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes:
Debezium outbox schema and narrative confirmed from live page fetch. Table structure matches what research reported.

---

## Claim 8

Claim:
Idempotency is mandatory for saga participants; mechanism: PROCESSED_MESSAGES table keyed by (subscriberId, messageID); or event UUID in Kafka header (Debezium); or idempotency key (clientId) for external APIs (Stripe-style).

Location:
`research/05-report.md` → Finding 7; `research/03-evidence.md` → Evidence 10

Evidence Provided:
Microservices.io Idempotent Consumer: "INSERT will fail if the message has been already processed."
Debezium: event UUID propagated as Kafka header for consumer duplicate detection.
Temporal: programmer must ensure Activity idempotency.

Source:
- https://microservices.io/patterns/data/idempotent-consumer.html
- https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
- https://temporal.io/blog/saga-pattern-made-easy

Source Actually Supports Claim: YES (partially for microservices.io idempotent consumer — that URL was cited but not directly verified in this audit session)

Classification: FACT

Severity: LOW

Notes:
Debezium idempotency mechanism confirmed directly. Microservices.io Idempotent Consumer URL not directly opened in this session, but the pattern text is consistent with and corroborated by the Debezium and Microsoft sources. Minor: marking as partially verified on microservices.io idempotent-consumer URL specifically.

---

## Claim 9

Claim:
Sagas lack isolation, causing data anomalies: lost updates, dirty reads, fuzzy/non-repeatable reads. Countermeasures: semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency.

Location:
`research/05-report.md` → Finding 8; `research/03-evidence.md` → Evidence 11, 12

Evidence Provided:
Microsoft enumerates 3 anomaly types + 6 countermeasures verbatim.
Microservices.io: "Lack of isolation (the 'I' in ACID) — the lack of isolation means that there's risk that the concurrent execution of multiple sagas and transactions can use data anomalies."

Source:
- https://learn.microsoft.com/en-us/azure/architecture/patterns/saga (detailed enumeration)
- https://microservices.io/patterns/data/saga.html (conceptual reference)

Source Actually Supports Claim: YES

Classification: FACT (single source for the specific 6-countermeasure taxonomy)

Severity: MEDIUM
The 6 specific countermeasures are detailed only in Microsoft's architecture page. Microservices.io references the concept but refers to Chapter 4 of Richardson's book. The research correctly reports this as MEDIUM confidence.

Notes:
One note on terminology: Microsoft uses "risk-based concurrency" for the sixth countermeasure; the research report and evidence use "value-based concurrency" or "risk-based concurrency" alternately. This is a minor inconsistency in research labeling:
- `03-evidence.md` Evidence 12: "version file, value-based concurrency"
- `05-report.md` Finding 8: "version files, value-based concurrency"
- Microsoft's actual page: "Risk-based concurrency based on value"

This is a naming variation, not a factual error.

---

## Claim 10

Claim:
Compensating transactions can themselves fail, leaving the system inconsistent; monitoring, retry, and operator intervention are required.

Location:
`research/05-report.md` → Finding 9; `research/03-evidence.md` → Evidence 14

Evidence Provided:
Microsoft: "Compensating transactions might not always succeed, which can leave the system in an inconsistent state."
AWS: "The saga pattern is difficult to debug and its complexity increases with the number of microservices."

Source:
- https://learn.microsoft.com/en-us/azure/architecture/patterns/saga
- https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes:
Both sources confirmed reachable and contain the cited statements.

---

## Claim 11

Claim:
AWS Step Functions Standard workflows provide exactly-once execution; Express workflows provide at-least-once; the saga implementation uses Catch/Retry states.

Location:
`research/05-report.md` → Finding 10; `research/03-evidence.md` → Evidence 13

Evidence Provided:
AWS Step Functions Docs: "Standard workflows have exactly-once workflow execution... Express workflows have at-least-once workflow execution."
AWS Prescriptive Guidance illustrates Step Functions saga with Catch/Retry.

Source:
- https://docs.aws.amazon.com/step-functions/latest/dg/welcome.html
- https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes:
Confirmed from live fetch of both AWS pages. "Exactly-once" is at the workflow-execution level, not at individual Lambda invocation level — this nuance is correctly reflected in `04-contradictions.md` (Contradiction 5).

---

## Claim 12

Claim:
Lab scenario (Order Pending → Debit Payment → Reserve Inventory FAIL → Refund Payment → Cancel Order) is a classical orchestration saga with LIFO compensation.

Location:
`research/05-report.md` → Finding 11

Evidence Provided:
Richardson orchestration example mapping; Microsoft compensable/pivot model; Temporal Java `saga.addCompensation` + `saga.compensate()`.

Source:
- Microservices.io orchestration example
- Temporal Java Saga class example
- Microsoft Azure Architecture Center

Source Actually Supports Claim: YES

Classification: EXAMPLE

Severity: LOW

Notes:
The pattern described is a straightforward application of known saga concepts corroborated by multiple sources. The lab scenario is not attributing a unique technical claim; it's an exercise design example.

---

## Summary Table

| Claim | Classification | Severity | Assessment |
|-------|---------------|----------|------------|
| 1 — Origin 1987 paper | FACT | LOW | PASS (citation chain) |
| 2 — Local transaction + compensation definition | FACT | LOW | PASS |
| 3 — Two coordination styles | FACT | LOW | PASS |
| 4 — Choreography vs Orchestration tradeoffs | FACT | LOW | PASS |
| 5 — Compensable/Pivot/Retryable taxonomy | FACT | MEDIUM | PASS (single source, correctly flagged) |
| 6 — 2PC vs Saga tradeoffs | FACT | LOW | PASS |
| 7 — Dual-write problem + Outbox solution | FACT | LOW | PASS |
| 8 — Idempotency mechanisms | FACT | LOW | PASS |
| 9 — Isolation anomalies + 6 countermeasures | FACT | MEDIUM | PASS (single source for list, correctly flagged) |
| 10 — Compensation failure risk | FACT | LOW | PASS |
| 11 — AWS Step Functions exactly-once / at-least-once | FACT | LOW | PASS |
| 12 — Lab scenario LIFO compensation | EXAMPLE | LOW | PASS |
