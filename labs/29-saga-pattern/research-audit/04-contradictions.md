# Contradictions Audit: Saga Pattern Research

## Contradiction 1

Statement A:
"Choreography is 'easier to implement, at least initially' but 'orchestration is often easier to build when one uses it from the start.'" — Temporal Blog (2023-07-13)

Location:
`research/04-contradictions.md:6` and `research/03-evidence.md:61-70`

Statement B:
"Good for simple workflows that have few services and don't need a coordination logic." (Choreography benefit) — Microsoft Azure Architecture Center

Location:
`research/04-contradictions.md:8-9`

Type:
INTERNAL

Impact:
LOW — the research agent correctly assessed this as a difference in framing, not a factual conflict. Both sources agree choreography suits simpler workflows; only emphasis differs.

Assessment:
RESOLVED by research agent. Acceptable as documented.

---

## Contradiction 2

Statement A:
Microsoft Azure Architecture Center lists 6 specific isolation countermeasures: semantic lock, commutative updates, pessimistic view, reread values, version files, value-based concurrency.

Location:
`research/04-contradictions.md:23-24`, `research/03-evidence.md:203-215`

Statement B:
Microservices.io references "countermeasures" in Chapter 4 of the Manning book without enumerating them on the public page.

Location:
`research/04-contradictions.md:27`

Type:
SOURCE_CONFLICT (minor, by incompleteness)

Impact:
MEDIUM — the 6-countermeasure taxonomy rests on a single public source (Microsoft). Richardson's complete list is paywalled. Research agent correctly reported MEDIUM confidence for this claim.

Assessment:
NOT FULLY RESOLVED. Single-source claim for the specific 6-item enumeration. Requires acknowledgment in gaps.

---

## Contradiction 3

Statement A:
"Standard workflows have exactly-once execution." — AWS Step Functions documentation

Location:
`research/04-contradictions.md:70-71`, `research/03-evidence.md:221-227`

Statement B:
Temporal claims "exactly-once" at workflow level via deterministic replay, but activities (participants) are at-least-once.

Location:
`research/04-contradictions.md:74`

Statement C:
Debezium pipeline is at-least-once; consumers must detect duplicates.

Location:
`research/04-contradictions.md:77`

Type:
SOURCE_CONFLICT (layered terminology)

Impact:
MEDIUM — conflating workflow-engine exactly-once with end-to-end participant delivery could mislead lab readers into thinking no idempotency is needed when orchestrated by Step Functions Standard. Research agent correctly resolved this as a layered distinction, not a fundamental contradiction.

Assessment:
RESOLVED adequately. Required note in final publication.

---

## Note

No material contradictions in core claims. All potential conflicts identified in `04-contradictions.md` are correctly classified as terminological or emphasis differences, not factual disagreements. Research agent exercised appropriate critical reasoning.
