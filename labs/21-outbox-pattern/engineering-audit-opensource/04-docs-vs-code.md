# Docs vs Code Comparison

Target Lab: labs/21-outbox-pattern

---

## DOC_CODE_MISMATCH

Location: engineering/01-design.md:Execution Plan (line 69) and Architecture diagram (line 29-47)
Claims: "Create SQLite DB schema (`orders`, `outbox_events`, `consumer_log`)"
Diagram: Shows SQLite Database component
Observed: Implementation uses in-memory Go maps (`orders map[string]Order`, `outbox map[string]OutboxMessage`)
No SQLite or `database/sql` dependency in go.mod
Consumer uses `map[string]bool`, not `consumer_log` table
Assessment: MISMATCH
Severity: MEDIUM
Notes: Design doc describes production SQLite implementation, but actual code is in-memory simulation as documented in engineering/02-implementation-notes.md and content/02-master-draft.md. The design doc is aspirational; implementation is a deliberate simplification for portability.

## DOC_CODE_MISMATCH

Location: engineering/01-design.md:Outbox Table Structure
Claims: Table structure includes `id`, `aggregateid`, `aggregatetype`, `type`, `payload`
Observed: internal/outbox/model.go defines OutboxMessage with fields: ID, EventType, Payload, Status, CreatedAt
No aggregateid/aggregatetype columns
Assessment: MISMATCH
Severity: LOW
Notes: Research/05-report.md Finding 5 specifies these columns (per Debezium/microservices.io), but implementation uses a simplified model for demo. The payload contains full order state. Assessment: Implementation omits aggregate identifiers but retains payload. Acceptable simplification for lab context.

## DOC_CODE_MISMATCH

Location: engineering/01-design.md:Relay Mechanism
Claims: "OutboxRelay: Background worker polling pending events, delivering to broker, marking processed, with retry / cleanup support."
Observed: Relay implements polling but has NO retry support (no backoff, no DLQ) and NO cleanup support (no automatic purge)
Assessment: PARTIAL_MISMATCH
Severity: MEDIUM
Notes: The relay polls and retries on broker failure by re-attempting publish, but has no retry policy configuration, exponential backoff, or dead-letter queue. No cleanup worker mentioned in design — only referenced in Success Criteria #5 ("Cleanup worker successfully purges processed events") and Production Considerations section. Implementation lacks both automated retry policy and cleanup worker.

## DOC_CODE_MISMATCH

Location: engineering/01-design.md:Success Criteria #5
Claims: "Cleanup worker successfully purges processed events"
Observed: No cleanup worker exists. PurgeProcessedOutbox() is a manual function tested but never called in demo or relay.
Assessment: MISMATCH
Severity: MEDIUM
Notes: Design claims automated cleanup worker, but implementation only provides a manual function. Production considerations section correctly notes cleanup is required, but implementation does not demonstrate it.

## DOC_CODE_MISMATCH

Location: engineering/02-implementation-notes.md:Known Limitations
Claims: "Log-tailing (CDC/Debezium) is not implemented; polling publisher pattern is used instead."
Observed: This is accurate — implementation uses polling publisher.
Assessment: MATCH
Severity: LOW
Notes: Honest documentation of implementation choice.

## DOC_CODE_MISMATCH

Location: engineering/03-execution-result.md:Build/Test/Race Detector/Demo
Claims: All show PASS results
Observed: Verified via actual execution — all commands pass.
Assessment: MATCH
Severity: LOW
Notes: Execution results accurately reflect actual code behavior.

## DOC_CODE_MISMATCH

Location: content/02-master-draft.md:Architecture table (lines 210-218)
Claims: Lists components with correct file mappings
Observed: Matches actual implementation files
Assessment: MATCH
Severity: LOW
Notes: Master draft accurately documents the as-implemented architecture.

## DOC_CODE_MISMATCH

Location: content/02-master-draft.md:Production Considerations (lines 554-576)
Claims: Lists production considerations (persistence, backoff, cleanup, CDC, monitoring, etc.)
Observed: None of these are implemented in the code — correctly labeled as considerations.
Assessment: MATCH
Severity: LOW
Notes: Documentation properly flags what is *not* implemented as per research/engineering findings.

## TEST_CLAIM_MISMATCH

Location: research/05-report.md:Finding 6 (Operational Requirements)
Claims: Outbox systems require cleanup and monitoring of key metrics (unprocessed event count, oldest unprocessed event age, publish failure rate, etc.)
Observed: No monitoring or metrics collection in code. PurgeProcessedOutbox exists but not automated.
Assessment: MISMATCH
Severity: MEDIUM
Notes: Research identifies requirement, but implementation does not demonstrate it — only provides manual purge function.

## RESEARCH_IMPLEMENTATION_MISMATCH

Location: research/05-report.md:Finding 2 (Outbox Pattern Guarantees Atomicity)
Claims: Guarantees messages sent iff DB transaction commits
Observed: Implementation achieves this via Tx.Commit() writing both order and outbox atomically.
Assessment: MATCH
Severity: LOW
Notes: Core guarantee correctly implemented.

## RESEARCH_IMPLEMENTATION_MISMATCH

Location: research/05-report.md:Finding 4 (Idempotent Consumer Requirement)
Claims: At-least-once delivery requires idempotent consumer
Observed: Consumer.Handle() provides idempotency via ID tracking.
Assessment: MATCH
Severity: LOW
Notes: Correctly implemented.

## SUMMARY

Matches: Core atomicity, dual-write problem demonstration, relay polling, consumer idempotency, test/demo execution results, production considerations acknowledgment.
Mismatches: Design doc's SQLite/CDC claims vs in-memory polling implementation, missing automated cleanup/retry policy, missing monitoring, simplified outbox schema.