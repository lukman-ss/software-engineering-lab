# Claim Audit

## Claim 1

Claim: Direct database connection establishment imposes significant latency and backend process memory overhead; PostgreSQL allocates shared memory structures based on `max_connections`.

Location:
`05-report.md`: Finding 1 & `03-evidence.md`: Evidence 3

Evidence Provided:
PostgreSQL fork overhead per connection, `max_connections` shared memory sizing, compared to PgBouncer 2 kB connection memory footprint.

Source:
PostgreSQL Documentation (`runtime-config-connection`) & PgBouncer Features

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects engine-level mechanics for process-per-connection architectures.

---

## Claim 2

Claim: Increasing pool size beyond hardware resource saturation degrades transaction throughput and spikes response latency ("the knee").

Location:
`05-report.md`: Finding 2 & `03-evidence.md`: Evidence 1, 3

Evidence Provided:
HikariCP / Oracle benchmark references (decreasing pool size dropped response times from ~100ms to ~2ms) and PostgreSQL Wiki analysis of contention points (RAM/work_mem, lock contention, context switching, cache line contention).

Source:
HikariCP Wiki & PostgreSQL Wiki

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Supported by empirical tests cited in both primary/secondary sources.

---

## Claim 3

Claim: Baseline connection sizing formula is `((core_count * 2) + effective_spindle_count)`, simplifying toward `core_count * 2` (or closer to `core_count`) on modern SSD/NVMe storage.

Location:
`05-report.md`: Finding 3 & `03-evidence.md`: Evidence 2 & `04-contradictions.md`: Contradiction 1

Evidence Provided:
Formulas documented in HikariCP wiki and PostgreSQL Wiki with nuance regarding zero rotational latency on SSDs.

Source:
HikariCP Wiki & PostgreSQL Wiki

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
LOW

Notes:
The formula serves as a baseline starting point for load testing, not an absolute invariant. The research correctly identifies this caveat.

---

## Claim 4

Claim: In horizontally scaled distributed deployments, aggregate application pools multiply linearly and can exhaust database backend slots unless decoupled by a proxy pooler like PgBouncer.

Location:
`05-report.md`: Finding 4 & `03-evidence.md`: Evidence 4

Evidence Provided:
Multiplication example (4 instances * 16 workers * 10 connections = 640 vs max_connections 200) and PgBouncer transaction pooling mechanics.

Source:
PgBouncer Features & PostgreSQL Documentation

Source Actually Supports Claim:
YES

Classification:
EXAMPLE

Severity:
LOW

Notes:
Clear demonstration of pool multiplication risk in distributed architectures.

---

## Claim 5

Claim: Pool deadlocks occur when threads hold multiple connections simultaneously; the theoretical minimum pool size to avoid deadlock is `pool size = Tn x (Cm - 1) + 1`.

Location:
`03-evidence.md`: Evidence 5

Evidence Provided:
Resource allocation formula from HikariCP documentation where `Tn` is maximum thread count and `Cm` is maximum simultaneous connections per thread.

Source:
HikariCP Wiki

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Verified directly in HikariCP documentation.

---

## Claim 6

Claim: Connection leaks and starvation are observable via `pg_stat_activity` when sessions remain in `idle in transaction` while waiting on client-side I/O (`ClientRead`).

Location:
`05-report.md`: Finding 5 & `03-evidence.md`: Evidence 6

Evidence Provided:
`pg_stat_activity` state definitions and wait event descriptions from PostgreSQL documentation.

Source:
PostgreSQL Documentation (`monitoring-stats.html`)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Matches PostgreSQL catalog documentation and production diagnostic standards.
