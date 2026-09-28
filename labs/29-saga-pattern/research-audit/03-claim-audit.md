# Claim Audit: Saga Pattern Research

## Claim 1: 2PC is infeasible / ill-suited for microservices with database-per-service

Location:
`03-evidence.md`: Evidence 1
`05-report.md`: Finding 1

Evidence Provided:
Direct quotes from Microsoft Azure Architecture Center and Chris Richardson noting that 2PC coordinator creates availability bottlenecks, long lock durations, and lack of support across disparate modern storage engines.

Source:
Microsoft Azure Architecture Center & microservices.io

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Supported by consensus in distributed systems literature.

---

## Claim 2: Saga executes as a sequence of local transactions with saga-level atomicity

Location:
`03-evidence.md`: Evidence 2, Evidence 10
`05-report.md`: Finding 2

Evidence Provided:
Step-by-step breakdown: each step executes atomically within local service DB and emits an event/command to trigger the subsequent step; atomicity of whole workflow is achieved via forward progress or compensation.

Source:
Microsoft Azure Architecture Center & microservices.io

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately distinguishes local ACID boundaries from distributed saga completion.

---

## Claim 3: Choreography vs Orchestration trade-offs

Location:
`03-evidence.md`: Evidence 4
`05-report.md`: Finding 3

Evidence Provided:
Comparison table detailing loose coupling vs coordinator overhead, risk of cyclic dependencies and debugging difficulty in choreography vs centralized logic in orchestration.

Source:
Microsoft Azure Architecture Center & microservices.io

Source Actually Supports Claim:
YES

Classification:
FACT / INTERPRETATION

Severity:
LOW

Notes:
Standard architectural trade-offs corroborated across both sources.

---

## Claim 4: Compensating transactions are semantic corrective actions, not automatic ACID rollbacks

Location:
`03-evidence.md`: Evidence 3
`05-report.md`: Finding 4

Evidence Provided:
Explains that local changes are committed and visible to concurrent actors; rolling back requires explicit application logic with opposing effect.

Source:
Chris Richardson & Microsoft Azure Architecture Center

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Critical concept properly documented.

---

## Claim 5: Transaction decomposition into Compensable, Pivot, and Retryable

Location:
`03-evidence.md`: Evidence 5
`05-report.md`: Finding 5

Evidence Provided:
Defines Compensable (can be reversed), Pivot (point of no return), and Retryable (guaranteed to succeed eventually through idempotency).

Source:
Microsoft Azure Architecture Center

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Directly adopted from Garcia-Molina/Richardson taxonomy.

---

## Claim 6: Lack of Isolation causes data anomalies (lost updates, dirty reads, fuzzy reads)

Location:
`03-evidence.md`: Evidence 6, Evidence 11
`05-report.md`: Finding 6

Evidence Provided:
Saga lacks the "I" of ACID across service boundaries because intermediate states are visible immediately after local commit.

Source:
Microsoft Azure Architecture Center & microservices.io

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Well formulated and highlighted as a primary drawback.

---

## Claim 7: Recommended countermeasures for isolation anomalies

Location:
`03-evidence.md`: Evidence 7
`05-report.md`: Finding 7

Evidence Provided:
Lists semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency.

Source:
Microsoft Azure Architecture Center

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION / IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
While Microsoft lists these countermeasures conceptually, the research notes that deep implementation recipes for each countermeasure require further design in the code stage.

---

## Claim 8: Idempotency is mandatory for reliable saga execution

Location:
`03-evidence.md`: Evidence 8
`05-report.md`: Finding 8

Evidence Provided:
Network retries and compensations can execute multiple times; participants must guarantee idempotent handling.

Source:
Microsoft Azure Architecture Center & microservices.io

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Essential operational requirement.

---

## Claim 9: Atomically updating state AND publishing messages requires patterns like Transactional Outbox

Location:
`03-evidence.md`: Evidence 13
`05-report.md`: Areas of Disagreement / Complementary differences

Evidence Provided:
Local DB update and message broker publish cannot span a 2PC boundary without re-introducing 2PC overhead; Transactional Outbox pattern is needed.

Source:
Chris Richardson / Microservices.io

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately identifies the dual-write problem inherent in saga event propagation.

---

## Claim 10: Compensating transactions may fail, leaving system in inconsistent state

Location:
`03-evidence.md`: Evidence 9
`05-report.md`: Finding 9

Evidence Provided:
Explicit warning that compensation failure requires manual intervention, alert monitoring, or escalation.

Source:
Microsoft Azure Architecture Center

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Honest appraisal of failure modes; avoids "zero downtime" / "perfect recovery" falsehoods.
