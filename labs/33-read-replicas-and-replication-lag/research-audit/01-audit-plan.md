# Audit Plan: Research on Read Replicas & Replication Lag

**Target Lab**: `labs/33-read-replicas-and-replication-lag`  
**Audit Date**: 2026-09-28  
**Auditor**: Independent Technical Research Auditor  
**Scope**: Research-only audit (`research/` files, sources, claims, contradictions, open questions).

---

## Files Reviewed

- `labs/33-read-replicas-and-replication-lag/research/01-plan.md`
- `labs/33-read-replicas-and-replication-lag/research/02-sources.md`
- `labs/33-read-replicas-and-replication-lag/research/03-evidence.md`
- `labs/33-read-replicas-and-replication-lag/research/04-contradictions.md`
- `labs/33-read-replicas-and-replication-lag/research/05-report.md`
- `labs/33-read-replicas-and-replication-lag/research/06-open-questions.md`

---

## Claims To Verify

1. **Default Replication Mode**: PostgreSQL streaming replication, MySQL replication, and MongoDB replica sets operate asynchronously by default.
2. **Replication Lag Behavior**: Async replication creates an eventual consistency window where replica reads can be stale (lag ranging from sub-second to hours under load).
3. **Synchronous Replication Semantics & Trade-offs**: PostgreSQL `synchronous_commit` modes (`remote_write` vs `remote_apply`) and MySQL semi-sync guarantee durability/relay log storage, but only `remote_apply` or primary reads guarantee immediate read visibility.
4. **Consistency Models & Linearizability**: Followers returning stale data violate linearizability (Kleppmann 2015); session guarantees (read-your-writes, monotonic reads) require coordination or sticky primary routing (Terry et al. 2011).
5. **Storage-level vs Log-shipping Architecture**: Amazon Aurora separates compute and storage, allowing up to 15 Aurora Replicas to share the same underlying storage volume.
6. **Application Routing Mechanisms**: ORM/proxy middleware (e.g. GORM DBResolver, Vitess) can split reads/writes and support override clauses (`dbresolver.Write`).

---

## Code To Execute

- **None** (Pipeline override: research audit only; implementation/code audit excluded for this phase).

---

## Primary Risks

- Unverified or dead citation URLs.
- Overgeneralization of database-specific features (e.g., assuming PostgreSQL has built-in causal consistency sessions like MongoDB).
- Arbitrary numeric recommendations presented as facts without empirical or mathematical grounding (e.g., sticky routing timeout durations).
- Source tier misattribution or reliance on search snippets.

---

## Audit Strategy

1. **Source Integrity**: Check all 10 cited sources for accessibility, domain authority, accuracy of titles/publishers, and fidelity of claims against actual contents.
2. **Claim Verification**: Map all major factual claims in `03-evidence.md` and `05-report.md` back to primary sources; classify claims (FACT, INTERPRETATION, IMPLEMENTATION-SPECIFIC, HEURISTIC).
3. **Contradiction Analysis**: Assess internal consistency between research documents and source citations.
4. **Gaps Identification**: Record missing sources, unverified assumptions, and heuristic boundaries in `06-gaps.md`.
5. **Final Verdict**: Issue status based on evidence standards in `07-verdict.md`.
