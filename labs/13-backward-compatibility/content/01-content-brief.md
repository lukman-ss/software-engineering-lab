# Content Brief

Topic: Backward Compatibility in Database Schema & API Evolution (Expand -> Migrate -> Contract)
Target Reader: Backend Engineers, System Architects, Platform Engineers
Problem: Evolving production schemas and API contracts without downtime, data loss, or breaking existing consumers.
Core Mental Model: Expand -> Migrate -> Contract (Parallel Change) decouples breaking changes into independent, reversible phases.
Approved Research Status: APPROVED_WITH_WARNINGS
Approved Engineering Status: APPROVED
Main Concepts:
- Backward Compatibility vs Forward Compatibility
- Breaking Changes (structural vs semantic)
- Expand -> Migrate -> Contract Pattern
- Dual-Write vs Dual-Read (Fallback Read)
- Resumable & Idempotent Backfill
- Zero-Downtime Rolling Deployments
- Feature Flag Gating & Observability
- Contract Enforcement & Safe Deprecation

Verified Behaviors:
- Un-upgraded V1 consumers parse legacy fields (`phone`) without failure across Expand and Migrate phases.
- V2 consumers parse modern array structures (`phones`).
- Dual-write writes atomically to both legacy and modern tables.
- Fallback read prevents data starvation before backfill completes.
- Backfill executes in batches with checkpointing and idempotency.
- Feature flag rollback safely reverts traffic to V1 without data loss.
- Contract phase safely disables legacy paths after traffic reaches zero.

Available Case Studies:
- Transition from 1:1 user phone relation (`users.phone`) to 1:N structure (`user_phones` table) with zero downtime.

Warnings:
- The 30-day zero-traffic metric is an operational guideline, not a universal standard.
- In-memory store is an illustrative simulation; production environments require ACID database transactions or outbox patterns.
