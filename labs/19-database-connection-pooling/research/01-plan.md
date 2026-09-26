# Research Plan: Database Connection Pooling

## Research Topic
Database Connection Pooling: Architecture, Sizing Formulas, Connection Leaks, Resource Contention, and Production Diagnostics.

## Objective
Investigate how connection pooling operates at the protocol, process, and application layers; determine authoritative sizing formulas and limits; understand failure modes (connection exhaustion, leaks, lock contention); and establish production diagnostic playbooks for PostgreSQL deployments.

## Research Questions
1. **Connection Lifecycle & Overhead**: What exact resource overhead (memory, process creation, TLS/handshake) does a direct database connection incur in process-per-connection architectures like PostgreSQL?
2. **Pool Sizing & Saturation Dynamics**: Why does throughput degrade when pool sizes exceed CPU core limits ("the knee" in the TPS curve), and what mathematical formulas govern optimal sizing?
3. **Failure Modes (Leaks & Exhaustion)**: What mechanisms trigger connection leaks (`idle in transaction`, unclosed connections, hanging external I/O), and how does connection starvation cascade into application 500 errors?
4. **Architecture Sizing**: In multi-instance distributed deployments (e.g., 4 instances × 16 workers × 10 connections = 640 potential connections vs `max_connections = 200`), how should client-side pooling vs intermediary proxy pooling (e.g., PgBouncer) be structured?
5. **Monitoring & Diagnostic Sequence**: What metrics and database catalog views (`pg_stat_activity`, wait events) must be inspected before altering server-side `max_connections`?

## Search Strategy
- Query primary database engine documentation (PostgreSQL 15-18 documentation).
- Consult authoritative connection pool engineering literature (HikariCP / Brett Wooldridge, Oracle Real-World Performance Group).
- Review database proxy architecture specifications (PgBouncer documentation).
- Review community wiki and benchmark analyses (PostgreSQL Wiki on "Number Of Database Connections").

## Expected Primary Sources
- PostgreSQL Official Documentation: Connections and Authentication (`max_connections`, `reserved_connections`), Cumulative Statistics System (`pg_stat_activity`).
- HikariCP Architecture Documentation: Pool Sizing Principles & Formulas.
- PostgreSQL Wiki: Number Of Database Connections (Contention analysis, scaling bottlenecks).
- PgBouncer Documentation: Architecture & Pooling Modes (Session, Transaction, Statement).

## Risks / Unknowns
- SSD vs spinning disk variance in the classical PostgreSQL connection formula `((core_count * 2) + effective_spindle_count)`.
- Transaction pooling incompatibilities with application-level session states (`SET/RESET`, prepared statements, advisory locks).
- Discrepancy between single-tenant internal connection pooling and distributed microservice pooling over multi-tier orchestrators.
