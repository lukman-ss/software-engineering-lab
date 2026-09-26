# Research Audit Plan: Transactional Outbox Pattern

## Target Lab
`labs/21-outbox-pattern`

## Scope
Audit of research artifacts only (per pipeline override):
- `research/01-plan.md`
- `research/02-sources.md`
- `research/03-evidence.md`
- `research/04-contradictions.md`
- `research/05-report.md`
- `research/06-open-questions.md`

## Files Reviewed
1. `labs/21-outbox-pattern/research/01-plan.md`
2. `labs/21-outbox-pattern/research/02-sources.md`
3. `labs/21-outbox-pattern/research/03-evidence.md`
4. `labs/21-outbox-pattern/research/04-contradictions.md`
5. `labs/21-outbox-pattern/research/05-report.md`
6. `labs/21-outbox-pattern/research/06-open-questions.md`

## Claims To Verify
1. **Dual-write problem**: Inability of standard database transactions to roll back external message broker writes (Redis, RabbitMQ, Kafka) or ensure atomicity without 2PC.
2. **Outbox pattern atomicity guarantee**: Inserting message into outbox table in the same DB transaction as aggregate update guarantees atomicity.
3. **Relay implementations**: Polling publisher and transaction log tailing are the two canonical relay patterns, with their respective trade-offs.
4. **Idempotency requirement**: Pattern provides at-least-once delivery; consumer must be idempotent via message ID / UUID tracking.
5. **Outbox table schema**: Canonical columns are `id`, `aggregatetype`, `aggregateid`, `type`, `payload`.
6. **Debezium SMT Event Router**: Specific field routing, aggregateid partitioning key, header placement of event UUID.
7. **Operational guidance & metrics**: Need for archival/cleanup of processed outbox entries, monitoring oldest unprocessed event age and queue lag.

## Primary Risks
- Overgeneralization of numeric SLA thresholds (e.g. 2s normal vs 47m alert).
- Scope conflation between generic pattern and Debezium-specific CDC implementation.
- Verifiability of external canonical documentation URLs.

## Audit Strategy
1. Live fetch and verify every cited URL in `02-sources.md`.
2. Map claims from `03-evidence.md` and `05-report.md` to retrieved source text.
3. Analyze contradictions across sources and within research documents.
4. Identify gaps, unsupported assertions, or overgeneralized statements.
5. Formulate evidence-based verdict and quality gate assessments.
