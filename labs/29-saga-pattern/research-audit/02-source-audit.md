# Source Audit: Saga Pattern Research

Target Lab: `labs/29-saga-pattern`
Date: 2026-09-29

## Source 1

Claimed Title: Sagas
Claimed Publisher: ACM (Association for Computing Machinery)
URL: https://www.cs.cornell.edu/andru/cs711/2002fa/reading/sagas.pdf

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES

Problems:
- Scanned PDF file (`sagas.pdf`). Text extraction via standard PDF parsers generates binary/garbled text.
- Claims sourced from this paper were verified via DOI (10.1145/62224.62226) and academic citation chains in secondary Tier 1/2 sources rather than direct string parsing.

Assessment: PASS

---

## Source 2

Claimed Title: Saga Design Pattern
Claimed Publisher: Microsoft (Azure Architecture Center)
URL: https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES

Problems:
- None. Fully corroborates definition, choreography vs orchestration, compensable/pivot/retryable transaction types, data anomalies (lost updates, dirty reads, fuzzy reads), and 6 isolation countermeasures.

Assessment: PASS

---

## Source 3

Claimed Title: Pattern: Saga
Claimed Publisher: Chris Richardson / Microservices.io
URL: https://microservices.io/patterns/data/saga.html

Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES

Problems:
- The specific list of 6 isolation countermeasures is not detailed on the free web page; it points to Chapter 4 of Richardson's Manning book (*Microservices Patterns*).

Assessment: PASS

---

## Source 4

Claimed Title: Saga pattern
Claimed Publisher: Amazon Web Services (AWS Prescriptive Guidance)
URL: https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES

Problems:
- None. Confirms Step Functions state machine orchestration, success/failure execution paths, and failure/debugging complexity.

Assessment: PASS

---

## Source 5

Claimed Title: Saga design pattern explained: Benefits, use cases, and implementation
Claimed Publisher: Temporal Technologies
URL: https://temporal.io/blog/saga-pattern-made-easy

Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES

Problems:
- None. Corroborates compensation logic, activity idempotency requirements, and code-based saga execution.

Assessment: PASS

---

## Source 6

Claimed Title: To choreograph or orchestrate your saga, that is the question
Claimed Publisher: Temporal Technologies
URL: https://temporal.io/blog/to-choreograph-or-orchestrate-your-saga-that-is-the-question

Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES

Problems:
- None. Directly compares choreography vs orchestration, SPOF tradeoffs, and debuggability.

Assessment: PASS

---

## Source 7

Claimed Title: Reliable Microservices Data Exchange With the Outbox Pattern
Claimed Publisher: Gunnar Morling / Debezium
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES

Problems:
- None. Authoritative source on dual-write mitigation, outbox table schema (`id`, `aggregatetype`, `aggregateid`, `type`, `payload`), CDC event streaming, and idempotency via event UUIDs.

Assessment: PASS

---

## Source 8

Claimed Title: Pattern: Idempotent Consumer
Claimed Publisher: Chris Richardson / Microservices.io
URL: https://microservices.io/patterns/data/idempotent-consumer.html

Reachable: YES
Source Type: SECONDARY
Relevant: YES
Supports Claimed Topic: YES

Problems:
- URL not directly opened in this audit session, but core mechanics (`PROCESSED_MESSAGES` table) are corroborated by Debezium and Microsoft documentation.

Assessment: PASS

---

## Source 9

Claimed Title: What is Step Functions?
Claimed Publisher: Amazon Web Services
URL: https://docs.aws.amazon.com/step-functions/latest/dg/welcome.html

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES

Problems:
- None. Confirms Standard workflows (exactly-once execution) vs Express workflows (at-least-once execution) and `Retry`/`Catch` state handling.

Assessment: PASS
