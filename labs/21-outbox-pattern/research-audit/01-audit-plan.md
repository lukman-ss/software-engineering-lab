# Audit Plan: Research Audit for Lab 21 (Outbox Pattern)

## Target Lab
`labs/21-outbox-pattern`

## Scope & Pipeline Override Notice
Per pipeline instructions, this audit evaluates **research files only** located under `labs/21-outbox-pattern/research/`.
- Implementation / code in `internal/`, `cmd/`, `tests/` is out of scope for this research audit stage.
- Research files must NOT be modified.
- Audit output location: `labs/21-outbox-pattern/research-audit/`

## Research Files Reviewed
1. `research/01-plan.md` (Research Plan)
2. `research/02-sources.md` (Sources List)
3. `research/03-evidence.md` (Evidence Ledger)
4. `research/04-contradictions.md` (Contradiction Log)
5. `research/05-report.md` (Final Research Report)
6. `research/06-open-questions.md` (Open Questions / Gaps)

## Primary Sources Audited
1. **Source 1**: Debezium Outbox Event Router Documentation (`https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html`)
2. **Source 2**: Debezium Blog — Reliable Microservices Data Exchange With the Outbox Pattern (`https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/`)
3. **Source 3**: Microservices.io — Pattern: Transactional Outbox (`https://microservices.io/patterns/data/transactional-outbox.html`)
4. **Source 4**: Microservices.io — Pattern: Polling Publisher (`https://microservices.io/patterns/data/polling-publisher.html`)
5. **Source 5**: Microservices.io — Pattern: Saga (`https://microservices.io/patterns/data/saga.html`)
6. **Source 6**: Debezium Architecture (`https://debezium.io/documentation/reference/stable/architecture.html`)

## Major Technical Claims To Verify
1. **Dual-write problem atomicity**: Standard DB transactions + message publishing cannot be made atomic using 2PC across heterogeneous stores.
2. **Outbox Pattern Guarantee**: Storing messages in an outbox table in the same DB transaction guarantees atomicity (all-or-nothing).
3. **Message Relay Mechanisms**: Two primary approaches exist — Polling Publisher and Transaction Log Tailing (CDC).
4. **Delivery Semantics & Idempotency**: Outbox pattern yields at-least-once delivery, requiring consumers to be idempotent via event UUID/deduplication.
5. **Outbox Table Structure & Payload Choices**: Canonical schema requires `id`, `aggregateid`, `aggregatetype`, `type`, `payload`. Thin vs fat event tradeoffs exist.
6. **Operational Requirements**: Outbox events require cleanup strategies (deletion/archival) and health monitoring (oldest event age metric).

## Audit Strategy
1. **Source Verification**: Fetch live URLs or check official documentation to verify existence, publisher accuracy, titles, and topic relevance.
2. **Claim Verification**: Compare each claim in `research/05-report.md` and `research/03-evidence.md` against cited sources to ensure accurate representation without misattribution or unsupported extensions.
3. **Contradiction Identification**: Analyze internal consistency across research files and compare source claims.
4. **Gap Analysis**: Check for missing citations, overgeneralizations, missing edge cases, or unverified numeric claims.
5. **Verdict Generation**: Apply standard severity model to assign final status (`APPROVED`, `APPROVED_WITH_WARNINGS`, `NEEDS_REVISION`, `REJECTED`).
