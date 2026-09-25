Target Lab: labs/14-circuit-breaker
Files Reviewed:
- labs/14-circuit-breaker/README.md
- labs/14-circuit-breaker/research/02-sources.md
- labs/14-circuit-breaker/research/03-core-concepts.md
- labs/14-circuit-breaker/research/04-cascade-failure.md
- labs/14-circuit-breaker/research/05-circuit-states.md
- labs/14-circuit-breaker/research/06-timeout-retry-backoff.md
- labs/14-circuit-breaker/research/07-fallback-bulkhead.md
- labs/14-circuit-breaker/research/08-observability.md
- labs/14-circuit-breaker/research/09-failure-modes.md
- labs/14-circuit-breaker/research/10-final-research.md

Claims To Verify:
1. Circuit Breakers isolate network blast radius and prevent cascading failures.
2. The caller fails in microseconds instead of blocking.
3. Asynchronous non-critical flows must be decoupled.
4. Silent fallbacks must never be used for critical mutations.
5. Circuit states (CLOSED, OPEN, HALF_OPEN) behave as described.

Code To Execute:
None (PIPELINE OVERRIDE: Do not audit implementation/code).

Primary Risks:
- Overgeneralized claims using "never" or "must" (e.g., fallbacks, decoupling).
- Uncited architectural patterns (PPOB/CMMS examples).

Audit Strategy:
- Validate URLs.
- Cross-reference major claims against sources.
- Identify missing sources for architectural claims.
