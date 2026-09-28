# Source Audit

**Target Lab**: `labs/33-read-replicas-and-replication-lag`  
**Audit Date**: 2026-09-28

## Source Verification Matrix

| Content Claim | Source File | Claimed Location | Verified? | Notes |
|---------------|-------------|-------------------|-----------|-------|
| Async replication default | `research/05-report.md` Finding 1 | "Async Replication Causes Variable Lag" | ✅ | Matches content section "Asynchronous vs Synchronous Replication" |
| Sticky routing 5s heuristic | `research/06-open-questions.md` Question 1 | "Optimal Sticky Routing Duration" | ✅ | Content correctly labels as convention, not guarantee |
| Causal token / LSN routing | `research/05-report.md` Finding 3 | "Session Guarantees via LSN/ClusterTime" | ✅ | Matches `ReadWithToken` implementation |
| Lag-aware + fallback | `research/05-report.md` Finding 4 | "Middleware Supports Automatic Splitting" | ⚠️ | Indirect — findings describe middleware splitting, not ΔLSN filter specifically |
| Monitoring metrics | `research/05-report.md` Finding 5 | lag metrics across platforms | ✅ | Referenced in Production Considerations |
| Sync replication trade-off | `research/05-report.md` Finding 5 | "Synchronous Replication Eliminates Lag" | ✅ | Matches `SyncReplication` path in cluster.go |
| Terry et al. session guarantees | `research/02-sources.md` Source 10 | "Replicated Data Consistency Explained Through Baseball" | ✅ | Referenced in Causal Token section |
| MongoDB causal sessions | `research/02-sources.md` Source 4 | "MongoDB — Read Isolation, Consistency, Recency" | ✅ | Mentioned in research report section |
| PostgreSQL `remote_apply` | `research/02-sources.md` Source 1 | "PostgreSQL 18 Docs — Log-Shipping Standby" | ✅ | Accurately described in async vs sync section |
| Vitess semi-sync ≠ freshness | `research/02-sources.md` Source 7 | "Vitess — MySQL Replication" | ✅ | Correctly distinguished in master draft |
| GORM DBResolver | `research/02-sources.md` Source 8 | "GORM DBResolver" | ✅ | Mentioned in production recommendations |

## Research-Gap Cross-Reference

| Research Gap | Severity | Content Coverage | Disclosed? |
|--------------|----------|------------------|------------|
| Gap 1: Sticky window is heuristic, not guarantee | MEDIUM | `01-content-brief.md` Warning, `05-key-takeaways.md` #3 | ✅ Yes |
| Gap 2: LSN polling overhead per-query | LOW | `01-content-brief.md` Warning, `02-master-draft.md` Production Considerations | ✅ Yes |
| Gap 3: Multi-region p99 lag unavailable | LOW | `01-content-brief.md` Warning | ✅ Yes |
