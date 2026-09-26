# Contradictions Audit

## Contradiction 1

Statement A: Status tracking uses `Proposed`, `Accepted`, `Deprecated`, `Superseded`.
Location: Nygard (2011)

Statement B: Status tracking adds `Rejected` to preserve institutional memory of discarded options.
Location: AWS Prescriptive Guidance / MADR

Type: SOURCE_CONFLICT
Impact: Minor implementation variance.
Assessment: Correctly identified by the research agent. Resolved cleanly by the agent adopting a superset of statuses (`Proposed`, `Accepted`, `Rejected`, `Deprecated`, `Superseded`) which preserves the intent of both schools of thought.

## Contradiction 2

Statement A: Accepted ADRs are strictly immutable. Updates require new ADRs that supersede older ones.
Location: Nygard / AWS / Microsoft Azure Well-Architected

Statement B: Practical teams often prefer "living documents" with date-stamped additions instead of immutability.
Location: Joel Parker Henderson community repository (README Teamwork Advice)

Type: SOURCE_CONFLICT
Impact: Operational tension between auditable history and maintenance convenience.
Assessment: Correctly flagged as a practice-level tension. The research agent rightly favored the Tier 1 strict immutability rule, as it enforces the referential integrity logic the lab aims to validate (supersession lineage).

## Contradiction 3

Statement A: Start new projects with a monolith, even if confident it will scale.
Location: Martin Fowler ("Monolith First" core thesis)

Statement B: Starting with microservices allows teams to get used to the rhythm... viable for system replacements with stable boundaries.
Location: Martin Fowler ("Monolith First" counter-argument section)

Type: INTERNAL
Impact: Nuanced architectural trade-off.
Assessment: Correctly identified by the research agent. This is not a contradiction of fact, but an acknowledgment of differing engineering contexts. The research maps the monolith-first consensus appropriately to the specific constraints of the lab scenario (small team, volatile boundaries, greenfield).

## Overall Assessment
No material contradictions were hidden or ignored. All identified contradictions are well-analyzed and appropriately resolved within the context of the engineering lab.
