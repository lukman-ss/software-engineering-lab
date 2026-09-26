# Content Brief

**Topic:** Backward Compatibility - Expand → Migrate → Contract Pattern for Database Schema Evolution

**Target Reader:** Backend engineers, SRE, technical leads responsible for production database migrations and API versioning.

**Problem:** How to evolve database schemas (1:1 to 1:N) and API contracts without causing downtime or breaking legacy clients during rolling deployments.

**Core Mental Model:** Parallel Change - decouple structural changes into three independent, reversible phases: Expand (additive schema), Migrate (dual write + backfill), Contract (legacy removal).

**Approved Research Status:** APPROVED

**Approved Engineering Status:** APPROVED

**Main Concepts:**
- Backward Compatibility vs Forward Compatibility
- Breaking Change Identification
- Expand → Migrate → Contract Pattern
- Dual-Write Data Consistency
- Batch Idempotent Resumable Backfill
- Fallback Read (Dual-Read)
- Observability-Driven Contract Enforcement
- Feature Flags for Safe Rollout/Rollback
- HTTP Deprecation and Sunset Headers (RFC 8594)

**Verified Behaviors:**
1. Legacy V1 clients receive `phone` string field; V2 clients receive `phones` array
2. Dual-write atomically writes to both legacy and modern tables
3. Backfill migrates historical data in configurable batches
4. Fallback read enables reading from legacy when modern field is empty
5. Rollback during Dual-Write preserves data integrity
6. Contract drops legacy field only after zero legacy traffic observed
7. Concurrency-safe operations pass `-race` detector

**Available Case Studies:**
- Case Study A: Customer Phone (1:1 to 1:N)
- Case Study B: CMMS Invoice Mechanics (1:1 to N:M)
- Case Study C: Multi Currency (Schema Splitting)

**Warnings:**
- 30-day legacy observation window is a heuristic (not universally verified)
- In-memory storage simplified; production requires database-level transaction handling
- Batch backfill mechanics not explicitly backed by Tier 1 sources
- Asynchronous CDC patterns (Debezium) acknowledged but not demonstrated