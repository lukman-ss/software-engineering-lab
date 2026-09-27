# Research Audit Plan: Outbox Pattern

Target Lab: `labs/21-outbox-pattern`
Audit Stage: Research Audit (Pipeline Override: Research Only)

## Target Lab Summary
The target lab explores the **Transactional Outbox Pattern** to solve the dual-write problem in distributed systems and event-driven architectures.

## Files Reviewed
- `labs/21-outbox-pattern/research/01-plan.md`
- `labs/21-outbox-pattern/research/02-sources.md`
- `labs/21-outbox-pattern/research/03-evidence.md`
- `labs/21-outbox-pattern/research/04-contradictions.md`
- `labs/21-outbox-pattern/research/05-report.md`
- `labs/21-outbox-pattern/research/06-open-questions.md`

## Claims To Verify
1. **Dual-write problem atomicity**: Standard database transactions (without 2PC) cannot atomically commit local DB updates and external message publishing (Redis, RabbitMQ, Kafka, HTTP).
2. **Outbox pattern atomicity guarantee**: Writing messages to an outbox table within the same DB transaction guarantees that messages are persisted if and only if business data commits.
3. **Message Relay alternatives**: Message relay can be implemented via Polling Publisher or Transaction Log Tailing (CDC).
4. **Delivery semantics & Consumer Idempotency**: Outbox pattern provides at-least-once delivery; consumers must track processed event IDs to handle duplicate messages.
5. **Outbox table schema**: Common columns are `id` (UUID), `aggregateid`, `aggregatetype`, `type`, `payload` (JSON/JSONB).
6. **Fat vs Thin payload trade-offs**: Thin events reduce DB bloat but require consumer RPC callbacks; fat events avoid RPC but increase storage/network payload.
7. **Operational requirements**: Outbox tables grow indefinitely without cleanup/archival (e.g. 1M events/day); monitoring must track oldest unprocessed event age and unprocessed count.

## Code To Execute
*Pipeline Override:* Code execution skipped for research-only audit.

## Primary Risks
- **Overgeneralization of SLAs/monitoring numbers**: Claiming specific metric thresholds (e.g. 2s normal vs 47min critical) as universal facts rather than contextual/illustrative examples.
- **Source scope mismatch**: Using Debezium/CDC documentation as sole proof for general outbox mechanics.
- **Citation integrity**: Claiming URLs exist or sources support claims when web fetching/inspection reveals discrepancies or unverified content.

## Audit Strategy
1. Perform web fetch checks on all 6 cited URLs in `02-sources.md` and `03-evidence.md`.
2. Inspect every finding in `05-report.md` against evidence items in `03-evidence.md` and primary sources.
3. Evaluate classification (FACT / INTERPRETATION / EXAMPLE / HYPOTHESIS) and severity for each claim.
4. Record research gaps, contradictions, and non-blocking vs blocking issues.
5. Issue final verdict in `labs/21-outbox-pattern/research-audit/07-verdict.md`.
