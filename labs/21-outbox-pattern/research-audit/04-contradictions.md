# Contradiction Audit: Transactional Outbox Pattern

Target Lab: `labs/21-outbox-pattern`
Research Set Under Audit: `research/2026-09-26-outbox-pattern/`
Date: 2026-09-26

---

## Contradiction 1: Exactly-Once Processing vs At-Least-Once Outbox Delivery

Statement A:
"Transactional Outbox pattern guarantees at-least-once delivery; consumers MUST be idempotent to prevent duplicate processing."
Location: `05-report.md` (Finding 6), `03-evidence.md` (Evidence 7)

Statement B:
"Kafka supports Exactly-Once Semantics (EOS) with transactions across topics."
Location: `05-report.md` (Finding 7), `02-sources.md` (Source 6)

Type: SCOPE_BOUNDARIES (Resolved Nuance)

Impact: LOW (Clarified in report)

Assessment:
No contradiction exists in the research files. The report explicitly clarifies that Kafka EOS applies exclusively to Kafka-to-Kafka read-process-write streams. Outbox spans relational DB to Kafka, where network partitions or relay restarts can produce duplicate publishes, maintaining the requirement for at-least-once + consumer idempotency.

---

## Contradiction 2: Table Housekeeping & Cleanup Mechanics

Statement A:
"Outbox tables must be periodically pruned or archived to prevent unbounded table growth and database bloat."
Location: `05-report.md` (Finding 8), `03-evidence.md` (Evidence 12)

Statement B:
"Debezium outbox implementation uses INSERT and DELETE within the same business transaction; CDC captures the WAL INSERT, leaving the outbox table permanently empty."
Location: `05-report.md` (Finding 8), `04-contradictions.md` (Nuance 3)

Type: IMPLEMENTATION_VARIANT (Resolved Nuance)

Impact: LOW

Assessment:
No conflict. These represent two legitimate implementation strategies dependent on the message relay architecture:
1. Persistent Outbox with Polling Relay -> requires retention purge daemon.
2. Ephemeral Outbox with CDC Log Tailing -> table stays empty because delete is committed in the same transaction after WAL logging.

---

## Summary of Findings

No un-reconciled material contradictions found within the research document set or between cited sources.
