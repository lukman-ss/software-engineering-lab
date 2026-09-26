# Contradictions

## Contradiction 1: Pool Sizing Formula Applicability to SSDs

**Source A (PostgreSQL Wiki / HikariCP Wiki):**
Claim: "There hasn't been any analysis so far regarding how well the formula works with SSDs." The formula `connections = ((core_count * 2) + effective_spindle_count)` was developed for spinning disks where `effective_spindle_count` matters.

**Source B (HikariCP Wiki):**
Claim: "SSDs perform better with *fewer* threads... less blocking and therefore fewer threads [closer to core count] will perform better than more threads."

**Assessment:**
These are not contradictory but represent a knowledge gap. The formula was designed for HDDs where disk seeks create blocking opportunities. With SSDs, the `effective_spindle_count` concept doesn't directly apply (no physical spindles). HikariCP argues SSDs need fewer connections, which aligns with setting effective_spindle_count = 0. However, no empirical SSD-specific formula has been established by the PostgreSQL project or HikariCP. The community acknowledges this explicitly as an unanalyzed area.

**Status:** NO MATERIAL CONTRADICTION — Acknowledged gap in research, not a disagreement between sources.

---

## Contradiction 2: Optimal Pool Size — Fixed vs Dynamic

**Source A (HikariCP Wiki):**
Claim: "We recommend *not* setting [minimumIdle] and instead allowing HikariCP to act as a *fixed size* connection pool."

**Source B (HikariCP README):**
Claim: `minimumIdle` default is "same as maximumPoolSize" which effectively creates a fixed-size pool. But the configuration allows dynamic sizing if `minimumIdle < maximumPoolSize`.

**Source C (HikariCP Wiki — Spike Demand Analysis):**
Claim: "The customer's environment imposed a high cost of new connection acquisition, and a requirement for a dynamically-sized pool, but yet a need for responsiveness to request spikes."

**Assessment:**
These are not contradictions but different recommendations for different workloads. HikariCP's primary recommendation is fixed-size pools for maximum performance. However, they acknowledge spike-demand scenarios where dynamic sizing is necessary. The default configuration (minimumIdle = maximumPoolSize) implements the primary recommendation. The configuration knobs exist for edge cases.

**Status:** NO MATERIAL CONTRADICTION — Different recommendations for different use cases, clearly documented.

---

## Contradiction 3: PgBouncer Transaction Mode vs Application Compatibility

**Source A (PgBouncer Config):**
Claim: Transaction pooling breaks several PostgreSQL features: SET/RESET, LISTEN, WITH HOLD CURSOR, PREPARE/DEALLOCATE, session-level advisory locks, LOAD, and temp tables with PRESERVE/DELETE ROWS.

**Source B (PgBouncer Features):**
Claim: "This mode breaks a few session-based features of PostgreSQL. You can use it only when the application cooperates by not using features that break."

**Source C (Azure Docs):**
Claim: "We recommend that you use PgBouncer... in transaction mode."

**Assessment:**
Azure recommends transaction mode as a general best practice, while PgBouncer documentation clearly states it requires application cooperation. This is not a contradiction — Azure is making a recommendation assuming the application is compatible or can be adapted. The PgBouncer documentation provides the compatibility matrix so teams can assess feasibility. Cloud providers may have more prescriptive defaults than the generic tool documentation.

**Status:** NO MATERIAL CONTRADICTION — Different audiences (cloud provider vs tool documentation) with different levels of prescription.

---

## Contradiction 4: max_connections Calculation — Cloud Provider vs On-Premise

**Source A (PostgreSQL Wiki):**
Claim: "max_connections should be a bit bigger than the number of connections you enable in your connection pool. That way there are always a few slots available for direct connections for system maintenance and monitoring."

**Source B (Azure):**
Claim: Reserves 15 connections for replication/monitoring automatically. max_connections default is calculated based on vCores. "Any subsequent changes of product selection... won't have any effect on the default value for max_connections... We recommend that whenever you change the product assigned to an instance, you also adjust the value for max_connections."

**Source C (AWS):**
Claim: max_connections = `LEAST(DBInstanceClassMemory/9531392, 5000)` — calculated automatically from instance memory.

**Assessment:**
Different deployment models have different approaches. On-premise/self-managed: set max_connections slightly above pool size + buffer. Cloud-managed: default is calculated from instance size (vCores or memory), with reserved connections built in. The principle is consistent (ensure headroom), but the mechanism differs. Cloud providers automate the calculation; self-managed requires manual tuning.

**Status:** NO MATERIAL CONTRADICTION — Different operational models with same underlying principle.

---

## Contradiction 5: Connection Leak Detection — HikariCP vs PostgreSQL JDBC Driver

**Source A (HikariCP):**
Claim: `leakDetectionThreshold` (default 0, min 2000ms) logs a message when a connection is out of pool beyond threshold.

**Source B (PostgreSQL JDBC Driver):**
Claim: `logUnclosedConnections` parameter captures stack trace at connection open time for leak debugging.

**Assessment:**
These are complementary, not contradictory. HikariCP detects leaks at the pool level (connection not returned). The JDBC driver logs at connection creation to help trace where leaked connections originated. They operate at different layers and solve different parts of the problem. Using both together provides better coverage.

**Status:** NO MATERIAL CONTRADICTION — Complementary mechanisms at different layers.

---

## Summary

**No material contradictions discovered** across the authoritative sources examined. The sources are remarkably consistent on:

1. The fundamental principle: more connections past saturation = worse performance
2. The pool sizing formula (with the explicit caveat about SSD uncertainty)
3. The necessity of external pooling (PgBouncer or application-level) for production
4. The danger of over-provisioning pool sizes
5. The need for connection leak detection and monitoring
6. Cloud providers' preference for conservative limits with external pooling over raw max_connections increases

The only "gaps" are:
- SSD-specific pool sizing formula (explicitly acknowledged as unresearched)
- Exact quantitative benchmark data for specific workloads (context-dependent)

All sources agree on principles; quantitative specifics depend on hardware, workload, and deployment model.