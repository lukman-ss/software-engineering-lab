# Claim Audit

## Claim 1
Claim: Table partitioning divides a table into smaller physical pieces within a single database instance, enabling partition pruning and fast bulk dropping, but does not distribute CPU/RAM scaling across multiple independent servers.
Location: `research/05-report.md` (Finding 1), `research/03-evidence.md` (Evidence 1)
Evidence Provided: PostgreSQL 18 Documentation Section 5.12 quotes on partition pruning and drop/detach partition performance vs bulk DELETE.
Source: PostgreSQL Documentation (DDL Partitioning)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately delineates single-node logical partitioning from multi-node distributed sharding.

---

## Claim 2
Claim: Sharding distributes data horizontally across multiple independent servers/instances to scale CPU, memory, disk I/O, and storage capacity beyond the limits of a single machine.
Location: `research/05-report.md` (Finding 1), `research/03-evidence.md` (Evidence 2)
Evidence Provided: MongoDB and Vitess documentation definitions of horizontal scaling and multi-server dataset distribution.
Source: MongoDB Manual (Sharding Overview) & Vitess Documentation (Sharding)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Supported by primary sources across document and relational paradigms.

---

## Claim 3
Claim: Monotonically increasing or decreasing keys (such as timestamps or autoincrement IDs) create severe write hotspots when used as range-based sharding keys because inserts concentrate on the single active shard/chunk.
Location: `research/05-report.md` (Finding 2), `research/03-evidence.md` (Evidence 3)
Evidence Provided: MongoDB shard key guidelines highlighting chunk bottlenecks for monotonic keys, Vitess Primary Vindex routing.
Source: MongoDB Manual (Choose a Shard Key)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Well-documented distributed database bottleneck mechanism.

---

## Claim 4
Claim: Sharding keys with high cardinality and low frequency (e.g., `user_id` or `tenant_id`) ensure balanced write distribution and enable direct point-lookup routing.
Location: `research/05-report.md` (Finding 2), `research/03-evidence.md` (Evidence 4)
Evidence Provided: MongoDB Manual cardinality/frequency definitions and Vitess Vindex single-shard targeted query routing.
Source: MongoDB Manual (Choose a Shard Key) & Vitess Documentation (Vindexes)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Validated against industry best practices and vendor documentation.

---

## Claim 5
Claim: Consistent Hashing minimizes key redistribution during cluster resizing to $O(K/N)$ keys on average, whereas Hash Modulo ($N \pmod M$) forces nearly all keys to relocate.
Location: `research/05-report.md` (Finding 3), `research/03-evidence.md` (Evidence 5)
Evidence Provided: Karger et al. (1997) STOC paper formal proof and bounds for consistent hashing vs modulo hashing.
Source: ACM STOC '97 (Karger et al.)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Mathematical guarantee verified against primary academic source.

---

## Claim 6
Claim: Queries executed without the sharding key require Scatter-Gather (broadcast) operations across all shards, increasing latency and limiting cluster scalability.
Location: `research/05-report.md` (Finding 5), `research/03-evidence.md` (Evidence 6)
Evidence Provided: MongoDB and Vitess query router behavior when sharding key is absent from the WHERE predicate.
Source: MongoDB Manual (Sharding) & Vitess Documentation (Vindexes)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Direct architectural consequence of horizontal data partitioning.

---

## Claim 7
Claim: Secondary Lookup Vindexes (Global Secondary Indexes) provide point-lookup routing for non-sharding-key queries at the cost of additional write overhead and cross-shard consistency maintenance.
Location: `research/05-report.md` (Finding 5), `research/03-evidence.md` (Evidence 7)
Evidence Provided: Vitess documentation on Lookup Vindex table maintenance, read lookup vs extra insert/delete writes, and transactional synchronization.
Source: Vitess Documentation (Vindexes)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately captures the read-optimization vs write-amplification trade-off.

---

## Claim 8
Claim: Standard database auto-increment cannot generate globally unique identifiers across independent database instances without central coordination; distributed ID generation (UUIDv7, Snowflake, sequences) is required.
Location: `research/05-report.md` (Finding 4), `research/03-evidence.md` (Evidence 8)
Evidence Provided: RFC 9562 Section 2.1/5.7 on UUIDv7 time-ordered indexing and Vitess sequence tables with block allocation.
Source: RFC 9562 & Vitess Documentation (Sequences)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately contrasts random UUIDv4, time-ordered UUIDv7, and sequence block allocation mechanisms.
