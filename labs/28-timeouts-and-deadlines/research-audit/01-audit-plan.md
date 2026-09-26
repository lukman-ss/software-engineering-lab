# Audit Plan: Research Audit for Lab 28 (Timeouts & Deadlines)

## Target Lab
`labs/28-timeouts-and-deadlines`

## Scope & Pipeline Override
- Research audit only.
- Implementation and code are deferred to Engineering phase.
- Audit output written directly to `labs/28-timeouts-and-deadlines/research-audit/`.

## Files Reviewed
- `labs/28-timeouts-and-deadlines/research/01-plan.md`
- `labs/28-timeouts-and-deadlines/research/02-sources.md`
- `labs/28-timeouts-and-deadlines/research/03-evidence.md`
- `labs/28-timeouts-and-deadlines/research/04-contradictions.md`
- `labs/28-timeouts-and-deadlines/research/05-report.md`
- `labs/28-timeouts-and-deadlines/research/06-open-questions.md`

## Claims To Verify
1. Cascading failures stemming from slow dependencies and resource exhaustion (threads, sockets, queues).
2. Retry storm amplification (mathematical amplification, multi-layer retry multiplication 4^3=64).
3. Deadline propagation (absolute deadlines vs fixed timeouts, clock skew handling via remaining budget).
4. Database timeout mechanisms (PostgreSQL `statement_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout`).
5. Circuit breaker state machine integration with retry policies.
6. Idempotency guarantees required when retrying timed-out requests (deduplication keys, atomic commits).
7. Backoff with jitter (preventing retry ripples / thundering herds).
8. Overload handling and criticality-based shedding (Google SRE K*accepts, 4 criticality levels).
9. Four golden signals and tail latency (P99 vs mean) for timeout detection.
10. Health check probe timeout separation (liveness vs readiness).

## Code To Execute
- None in this stage (pipeline override: research only; no runnable code exists in target lab).

## Primary Risks
- Inaccurate URL citations or dead links.
- Unsupported numeric rules (e.g. arbitrary 3-attempt caps without qualification, arbitrary latency multipliers).
- Misattribution of platform-specific features (e.g., gRPC deadline propagation details) as universal distributed systems behavior.
- Failure to document crash-window edge cases in idempotent retries.

## Audit Strategy
1. Perform reachability and content verification of all 13 sources cited in `02-sources.md`.
2. Cross-examine claims in `03-evidence.md` and `05-report.md` against primary documentation.
3. Review contradiction analysis in `04-contradictions.md` for completeness and validity.
4. Evaluate open questions and research gaps in `06-gaps.md`.
5. Issue final verdict in `07-verdict.md`.
