# Audit Plan: Research on Transactional Outbox Pattern

Target Lab: `labs/21-outbox-pattern`
Research Set Under Audit: `research/2026-09-26-outbox-pattern/`
Date: 2026-09-26

## Target Scope & Pipeline Override
- Stage: Research Audit ONLY
- Implementation / code execution: EXCLUDED in this phase per pipeline override
- Research artifacts reviewed:
  - `01-plan.md`
  - `02-sources.md`
  - `03-evidence.md`
  - `04-contradictions.md`
  - `05-report.md`
  - `06-open-questions.md`

## Files Reviewed
1. `research/2026-09-26-outbox-pattern/01-plan.md`
2. `research/2026-09-26-outbox-pattern/02-sources.md`
3. `research/2026-09-26-outbox-pattern/03-evidence.md`
4. `research/2026-09-26-outbox-pattern/04-contradictions.md`
5. `research/2026-09-26-outbox-pattern/05-report.md`
6. `research/2026-09-26-outbox-pattern/06-open-questions.md`
7. Prior revision notes in `research-revision/` for context

## Major Claims To Verify
1. **Dual-write atomicity impossibility & 2PC infeasibility**: DB + message broker cannot be coordinated via 2PC realistically without severe coupling or lack of XA support.
2. **Transactional Outbox core mechanism**: Inserting domain entity + outbox record in a single database transaction guarantees atomic intent.
3. **Message Relay approaches**: Polling publisher vs Transaction log tailing (CDC), trade-offs in complexity, ordering, latency, and duplicates.
4. **Delivery guarantee**: Outbox produces at-least-once semantics, requiring consumer idempotency.
5. **Schema design**: Standard Debezium table structure (`id`, `aggregatetype`, `aggregateid`, `type`, `payload`) vs generic schema.
6. **Kafka Transactions scope**: Kafka EOS solves Kafka-to-Kafka read-process-write, not DB-to-Kafka dual-write.
7. **Operational practices**: Table cleanup (ephemeral vs retention), payload sizing, DLQ, and oldest unprocessed event age monitoring.

## Primary Risks
- **Overgeneralization of canonical schema**: Presenting Debezium-specific table schemas as universal requirements.
- **Unverified operational thresholds**: Numeric latency or monitoring thresholds presented without authoritative benchmarks.
- **Conflation of Kafka EOS with DB-to-Broker atomicity**: Claiming exactly-once across database and messaging layers.
- **Source Reachability & Faithfulness**: Verifying whether external URLs exist, HTTP status is 200, and quotes/claims faithfully reflect source contents.

## Audit Strategy
1. Live network verification of all cited source URLs.
2. Cross-reference quotes and paraphrases against canonical texts from Microservices.io, Debezium, Confluent, Martin Fowler, and Enterprise Integration Patterns.
3. Check for unsupported claims, overgeneralizations, or conflated architectural semantics.
4. Verify internal consistency across `01-plan.md` through `06-open-questions.md`.
5. Identify any remaining research gaps or unverified operational assertions.
