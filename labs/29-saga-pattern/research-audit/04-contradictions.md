# Contradiction Audit: Saga Pattern Research

Target Lab: `labs/29-saga-pattern`
Date: 2026-09-29

---

## Contradiction 1: Choreography vs Orchestration Complexity Threshold

Statement A:
Temporal (2023-07-13) states: "Orchestration is often easier to build when one uses it from the start." Choreography appears simple initially, but becomes hard quickly.

Location:
`research/04-contradictions.md` → Contradiction 1; `research/02-sources.md` (Source 6)

Statement B:
Microsoft Azure Architecture Center & lab specification frame Choreography as "Good for simple workflows that have few services and don't need a coordination logic."

Location:
`research/04-contradictions.md` → Contradiction 1; `research/02-sources.md` (Source 2)

Type:
SOURCE_CONFLICT

Impact: LOW
Both sources agree that Choreography degrades as complexity grows. The disagreement is about whether orchestration should be the default from day one (Temporal's perspective, influenced by their product which is a durable execution engine) or whether choreography is acceptable for 2-3 service linear workflows (Microsoft/standard industry view).

Assessment:
The research correctly identifies this as a perspective difference rather than a factual error, noting that "No quantitative threshold (e.g. number of services or branching factor) exists; community consensus is qualitative." This is an accurate assessment.

---

## Contradiction 2: "Exactly-Once" Semantics vs Mandatory Participant Idempotency

Statement A:
AWS Step Functions Standard Workflows claim "exactly-once workflow execution."

Location:
`research/04-contradictions.md` → Contradiction 5; `research/02-sources.md` (Source 9)

Statement B:
Debezium outbox pattern and Temporal documentation explicitly state that participant services receive "at-least-once" message delivery and MUST be idempotent.

Location:
`research/04-contradictions.md` → Contradiction 5; `research/02-sources.md` (Sources 5, 7)

Type:
INTERNAL / SOURCE_CONFLICT

Impact: MEDIUM
If a developer relies on "exactly-once workflow execution" without implementing idempotency at the participant level, duplicate payments or inventory reservations can occur under network retry conditions.

Assessment:
The research correctly resolves this layered distinction: the orchestrator workflow state machine guarantees single execution of the workflow *definition*, but individual network calls/Lambda invocations to participant services can still experience retries, making participant idempotency non-negotiable. The research's analysis is technically accurate.

---

## Contradiction 3: Dual-Write Problem Attribution to Garcia-Molina (1987)

Statement A:
The research plan (`01-plan.md` Objective 5) lists the Dual-Write Problem as a core saga research topic alongside the 1987 paper.

Location:
`research/01-plan.md` → Objective 5

Statement B:
Garcia-Molina & Salem (1987) addressed long-lived transactions within a *single* database system. The Dual-Write Problem (updating a DB + sending a Kafka event without distributed transactions) is a modern event-driven microservices concern.

Location:
`research/04-contradictions.md` → Contradiction 4

Type:
HISTORICAL_CONTEXT_DISTINCTION

Impact: LOW
Dual-write is a modern prerequisite for implementing reliable saga messaging, not a feature of the 1987 paper.

Assessment:
The research explicitly clarifies this in `04-contradictions.md` (Contradiction 4): "Dual-write is a modern problem that sagas inherited when adapted to event-driven microservices... Terminology difference, not factual contradiction." This demonstrates clear historical discernment.

---

## Summary

- Total Contradictions Identified by Research: 5
- Total Contradictions Verified by Audit: 5
- Material Contradictions Unhandled by Research: 0
- Silent Fixes or Concealed Contradictions: None found.
- Assessment: Internal consistency is PASS with HIGH technical accuracy.
