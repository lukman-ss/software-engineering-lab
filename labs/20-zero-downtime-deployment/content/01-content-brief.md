# Content Brief

Topic: Zero-Downtime Deployment (ZDD)
Target Reader: Software Engineers, Backend Developers, DevOps Engineers
Problem: Application updates frequently drop in-flight HTTP requests, interrupt active background jobs, and cause database schema conflicts, leading to user-facing downtime.
Core Mental Model: "Start new before stopping old." Version N and Version N-1 must coexist, traffic must be intelligently routed using probes, and state/schema changes must be backward compatible.
Approved Research Status: APPROVED_WITH_WARNINGS
Approved Engineering Status: APPROVED_WITH_WARNINGS
Main Concepts:
- Health Probes (Liveness vs. Readiness)
- Graceful HTTP Shutdown and Connection Draining
- Kubernetes Asynchronous Routing Detachment (preStop hook)
- Database Expand and Contract (Parallel Change) Pattern
- Worker Graceful Termination
Verified Behaviors:
- Server rejects traffic when unready (HTTP 503) and accepts when ready (HTTP 200).
- Server preStop hook delays shutdown execution to simulate routing table detachment.
- Server connection draining allows active in-flight requests to complete before exit.
- Background worker completes the currently active job upon receiving a stop signal.
- Database fallback logic successfully reads legacy records and writes modern dual-state records.
Warnings:
- Research Warning: Kubernetes pod termination has asynchronous endpoint propagation; a `preStop` sleep hook is strictly required to avoid 502/504 errors.
- Research Warning: DDL table locks (e.g., PostgreSQL `ACCESS EXCLUSIVE`) can cause downtime during schema expansion; connection lock timeouts (e.g., `lock_timeout`) are necessary in production.
- Research Warning: Laravel-specific cache-based worker signals are used in some ecosystems, but POSIX signals (`SIGTERM`/`SIGQUIT`) are the universal standard.
- Engineering Warning: The worker implementation successfully finishes the *active* job but explicitly abandons remaining *buffered* jobs in the channel due to immediate context cancellation.
