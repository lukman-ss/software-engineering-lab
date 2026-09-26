# Content Brief

Topic: Backward Compatibility Strategy — Implementing Expand → Migrate → Contract (Parallel Change) Pattern for Database Schema and API Evolution

Target Reader: Backend engineer, infrastructure architect, DevOps engineer implementing zero-downtime migrations in production systems

Problem: Legacy systems face breaking changes when evolving schema (e.g., 1:1 → 1:N relationships) or API contracts (e.g., single field → array). Traditional migration strategies cause downtime or break existing clients. Need safe, incremental transition without service interruption.

Core Mental Model: Treat breaking changes as reversible by introducing dual representations (legacy + new) simultaneously, gradually migrating traffic and data, then safely retiring legacy representation only after zero-traffic verification.

Approved Research Status: APPROVED

Approved Engineering Status: APPROVED

Main Concepts:
- Parallel Change / Expand → Migrate → Contract pattern (3-phase non-breaking change)
- Dual-write: application writes to both legacy and new data stores simultaneously
- Backfill: batched, idempotent, resumable data migration worker
- Fallback reading (dual-read): read new schema first, fall back to legacy if empty
- Feature flags: runtime control of WriteMode and ReadMode for phased rollout
- Observability metrics: track legacy vs new traffic for safe cutover
- Rollback safety: dual-write ensures data preserved if version N+1 rolls back to N

Verified Behaviors:
- Expand phase adds new schema without affecting legacy clients
- Dual-write preserves data in both representations during migration
- Backfill migrates historical data in chunks with checkpoint resume
- Fallback read prevents data starvation before backfill completes
- Contract phase safely drops legacy only when legacy traffic = 0
- Rollback to older version during dual-write causes zero data loss
- Feature flags decouple deployment from activation for canary rollout

Available Case Studies:
- Customer phone (1:1 → 1:N) — Case Study A
- Invoice mechanics (1:1 → N:M) — Case Study B  
- Multi-currency pricing (single field → multiple currencies) — Case Study C

Warnings:
- Contract phase is IRREVERSIBLE if applied before legacy traffic drops to zero
- Dual-write increases write latency and complexity; use only when necessary
- Batch backfill must be throttled to avoid database lockup on large tables
- Fallback reading adds read-path complexity; prefer full backfill before switching to ReadNewOnly
- Research gaps: no published benchmarks for transformation pipeline overhead; PostgreSQL-specific DDL examples dominate (limited coverage of other databases)
- Engineering gaps: in-memory store used instead of PostgreSQL for demo; DDL lock timeout emulation not executed (documented but not runtime)
