# Revision Plan: Optimistic vs Pessimistic Locking

Target Lab: `labs/23-optimistic-vs-pessimistic-locking`
Pipeline Override: Research revision only. No code changes.
Revision Directory: `labs/23-optimistic-vs-pessimistic-locking/research-revision/`

Previous Audit Status: **APPROVED_WITH_WARNINGS**

## Blocking Issues

None.

## Non-Blocking Issues

1. **Source 14 (Hibernate ORM 6.x)**: Reachability recorded but direct HTTP fetch not executed in audit session. Status: LOW severity.
2. **Atomic Decrement Pattern**: Conditional `UPDATE ... SET stock = stock - N WHERE stock >= N` supported by synthesized evidence (statement atomicity + conditional guard) but no single Tier 1 source quotes recipe verbatim. Status: MEDIUM severity.
3. **MySQL Mirror Dependency**: MySQL 8.0 docs fetched via Oracle CDN mirror due to 403 on direct mysql.com endpoint. Status: WARNING (not a revision blocker, documented in report limitations).

## Files To Modify

- `research/02-sources.md`: Clarify Source 14 verification note; clarify Source 13 as contrastive NoSQL example.
- `research/03-evidence.md`: No change needed; Evidence 10 correctly marked MEDIUM confidence with explanation.
- `research/05-report.md`: No change needed; Limitations section correctly documents MySQL mirror, Hibernate, and atomic-pattern gaps.
- `research/06-open-questions.md`: No change needed; OQ-1 and OQ-2 correctly identify the gaps.

## Verification Plan

- [x] Inspect audit verdict
- [x] Inspect claim audit (7 claims, 0 unsupported)
- [x] Inspect source audit (14 sources, 1 MEDIUM warning)
- [x] Inspect contradictions (3 vendor divergences, properly documented)
- [x] Inspect gaps (4 gaps, 1 MEDIUM, 3 LOW, all non-blocking)
- [x] Confirm research files correctly classify confidence and document limitations
- [x] Confirm open questions correctly identify upgrade paths

Verification status: **PASS** — Research is internally consistent and correctly classifies confidence. Warnings do not block approval but should be tracked for future revision cycles.
