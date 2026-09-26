# Content Brief

Topic: Deadlock in Concurrent Systems
Target Reader: Software Engineer
Problem: Deadlocks form permanent circular waits blocking transactions, requiring system-level aborts and application-level handling to resolve or prevent.
Core Mental Model: Deadlock is a circular dependency. The database aborts one as a "victim". Lock ordering prevents cycles from forming. Application retry recovers aborted victims.
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts: Circular Wait, Deadlock Monitor/Victim, Lock Ordering, Transaction Duration, Application-Level Retry.
Verified Behaviors: Naive bidirectional resource access causes deadlocks. Consistent lock ordering completely prevents deadlocks. Retry mechanisms can successfully recover from deadlock aborts. Longer transactions increase deadlock frequency.
Available Case Studies: None. The "Sistem PPOB" framing in research is illustrative only (unbacked by primary domain literature — see research-audit WARNING, non-blocking issue 3).
Warnings:
- Research source integrity WARNING: Coffman Conditions / Two-Phase Locking sourced via Wikipedia (Tier 3); Oracle omitted (source unreachable); MySQL documentation accessed via Wayback Machine snapshots.
- Engineering is an in-memory application-level simulation using Go channels + context timeouts. It does NOT implement a real Wait-For Graph (WFG); the `ErrDeadlock` timeout models the database's deadlock victim selection rather than true cycle detection.
- Wait values (`10ms`, `20ms`, `50ms`) and retry backoff (`2ms`, fixed) are arbitrary and tuned only for fast, reproducible tests — not production guidance.
- `TransferWithRetry` uses fixed 2ms backoff (not exponential backoff with jitter) — a documented deliberate simplification for readability.
