# Claim Audit

## Claim 1

Claim:
Increasing the connection pool size beyond a certain physical capacity limit degrades database throughput and increases response times.

Location:
research/03-evidence.md (Evidence 1) and research/05-report.md (Finding 2)

Evidence Provided:
Quotes regarding contention overhead and reducing pool size to improve response times (from ~100ms to ~2ms).

Source:
HikariCP Wiki & PostgreSQL Wiki

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Both sources explicitly describe performance degradation when concurrent connection counts exceed available processing resources, highlighting the "knee" in throughput graphs.

---

## Claim 2

Claim:
The optimal formula for baseline database connection sizing is `((core_count * 2) + effective_spindle_count)`.

Location:
research/03-evidence.md (Evidence 2) and research/05-report.md (Finding 3)

Evidence Provided:
Verbatim formula citation.

Source:
HikariCP Wiki & PostgreSQL Wiki

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Formula is present in both sources as a recommended starting point for tuning, not a rigid absolute, which the research report accurately contextualizes regarding SSDs.

---

## Claim 3

Claim:
High connection counts degrade performance due to disk contention, memory exhaustion (work_mem), lock contention, context switching, and CPU cache line eviction.

Location:
research/03-evidence.md (Evidence 3)

Evidence Provided:
Direct lists of bottlenecks.

Source:
PostgreSQL Wiki

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Source precisely lists these five mechanical reasons for throughput collapse.

---

## Claim 4

Claim:
Distributed deployments without an intermediate pooler easily exceed backend connection limits, resulting in starvation. Proxy pooling (PgBouncer) collapses virtual connections into a small real pool with low memory footprint (2kB).

Location:
research/03-evidence.md (Evidence 4) and research/05-report.md (Finding 4)

Evidence Provided:
Description of PgBouncer transaction pooling mode and 2kB overhead.

Source:
PgBouncer Features & PostgreSQL Documentation

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
PgBouncer features page explicitly states transaction pooling mode and 2kB per connection memory requirement.

---

## Claim 5

Claim:
Connection deadlocks (pool locking) can be avoided using the resource allocation formula: `pool size = Tn x (Cm - 1) + 1`.

Location:
research/03-evidence.md (Evidence 5)

Evidence Provided:
Verbatim formula and variables mapping.

Source:
HikariCP Wiki

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Formula and examples are taken verbatim from the HikariCP wiki section on "Pool-locking".

---

## Claim 6

Claim:
Idle connections held open due to application leaks or long-running non-database tasks waste slots, observable via `pg_stat_activity` states like `idle in transaction` and wait events like `ClientRead`.

Location:
research/03-evidence.md (Evidence 6) and research/05-report.md (Finding 5)

Evidence Provided:
References to state and wait event tracking.

Source:
PostgreSQL Documentation (pg_stat_activity)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
`idle in transaction` state and `ClientRead` wait event are explicitly documented in PostgreSQL statistics documentation as indicators of backend waiting on client activity.
