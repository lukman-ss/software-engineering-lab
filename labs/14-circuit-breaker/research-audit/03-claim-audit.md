## Claim 1

Claim:
Circuit Breakers isolate network blast radius and prevent cascading failures.

Location:
labs/14-circuit-breaker/research/10-final-research.md (item 1), README.md (Cascade Failure)

Evidence Provided:
Azure Architecture Center documentation, Martin Fowler blog post.

Source:
Martin Fowler; Microsoft Learn

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Core justification for Circuit Breaker pattern.


## Claim 2

Claim:
Caller fails in microseconds instead of multi-second HTTP timeout hangs.

Location:
labs/14-circuit-breaker/research/10-final-research.md (item 2), README.md (With Circuit Breaker)

Evidence Provided:
Circuit breaker blocks the call locally without network dispatch.

Source:
Martin Fowler; Microsoft Learn

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Directly derived from the fail-fast mechanism of the OPEN state.


## Claim 3

Claim:
Asynchronous non-critical flows (e.g. notifications) must be decoupled using queues + idempotency so third-party downtime never aborts core transaction persistence.

Location:
labs/14-circuit-breaker/research/10-final-research.md (item 3)

Evidence Provided:
None cited in 02-sources.md.

Source:
Uncited.

Source Actually Supports Claim:
NO

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
While an industry best practice, it is an overgeneralized assertion with no direct reference in the cited sources.


## Claim 4

Claim:
Never use silent fallback for critical state-altering mutations (e.g., balance debit).

Location:
labs/14-circuit-breaker/README.md (Fallback)

Evidence Provided:
None cited.

Source:
Uncited.

Source Actually Supports Claim:
NO

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
Strong universal rule ("Never") not explicitly mentioned or justified by the primary sources.


## Claim 5

Claim:
Circuit states (CLOSED, OPEN, HALF_OPEN) transition automatically with probes and cooldowns.

Location:
labs/14-circuit-breaker/research/05-circuit-states.md, README.md

Evidence Provided:
Finite State Machine definitions in Martin Fowler and Azure Architecture Center.

Source:
Martin Fowler; Microsoft Learn

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Standard three-state model matches authoritative texts.
