# Claim Audit

## Claim 1
Claim: Table partitioning divides a table into smaller physical pieces within a single database instance, enabling partition pruning and fast dropping of data, but does not distribute CPU/memory across multiple machines.
Location: `research/03-evidence.md:Evidence 1`, `research/05-report.md:Finding 1`
Evidence Provided: PostgreSQL 18 Documentation (Section 5.12).
Source: PostgreSQL 18 Documentation (DDL Partitioning).
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurate distinction between logical partitioning and physical multi-machine scaling.

---

## Claim 2
Claim: Sharding distributes data horizontally across multiple independent servers/instances to scale CPU, memory, disk I/O, and storage capacity.
Location: `research/03-evidence.md:Evidence 2`, `research/05-report.md:Finding 1`
Evidence Provided: MongoDB Sharding Overview, Vitess Sharding Overview.
Source: MongoDB Manual, Vitess Documentation.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Correctly cites both document and relational horizontal sharding engines.

---

## Claim 3
Claim: Monotonically increasing or decreasing keys (such as `created_at` or autoincrement IDs) create severe write hotspots when used as range-based sharding keys; Hashed Sharding is the primary mitigation.
Location: `research/03-evidence.md:Evidence 3`, `research/05-report.md:Finding 2 & Finding 6`
Evidence Provided: MongoDB Choose a Shard Key, MongoDB Hashed Sharding.
Source: MongoDB Documentation.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Nuance correctly captured that monotonic keys are not forbidden, but require hashed routing rather than ranged routing.

---

## Claim 4
Claim: Sharding keys with high cardinality and low frequency (e.g., `user_id` or `tenant_id`) ensure balanced write distribution and enable direct point-lookup routing.
Location: `research/03-evidence.md:Evidence 4`, `research/05-report.md:Finding 2`
Evidence Provided: MongoDB Choose a Shard Key, Vitess Vindexes (Primary Vindex).
Source: MongoDB Documentation, Vitess Documentation.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Accurately explains how cardinality caps the maximum number of useful shards.

---

## Claim 5
Claim: Consistent Hashing minimizes key redistribution during cluster resizing to $O(K/N)$ keys on average ($1/n$ fraction when adding a node), whereas Hash Modulo (`key % M`) forces almost all keys to relocate.
Location: `research/03-evidence.md:Evidence 5`, `research/05-report.md:Finding 3`
Evidence Provided: ACM STOC '97 (Karger et al.), Wikipedia Consistent Hashing.
Source: ACM STOC '97 DOI, Wikipedia.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Math verified; $1/n$ applies to single-node addition/removal while $n/m$ describes general resize where $n$ is keys and $m$ is slots.

---

## Claim 6
Claim: Queries executed without the sharding key require Scatter-Gather (broadcast) operations across all shards, increasing latency and limiting cluster scalability (with exceptions for large aggregation queries).
Location: `research/03-evidence.md:Evidence 6`, `research/05-report.md:Finding 4 & 5`
Evidence Provided: MongoDB Sharding Overview, Vitess Vindexes.
Source: MongoDB Documentation, Vitess Documentation.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Correctly identifies p99 tail latency risk and notes the aggregation exception acknowledged by MongoDB docs.

---

## Claim 7
Claim: Secondary Lookup Vindexes (Global Secondary Indexes) provide point-lookup routing for non-sharding-key queries at the cost of additional write overhead and cross-shard consistency maintenance.
Location: `research/03-evidence.md:Evidence 7`, `research/05-report.md:Finding 5`
Evidence Provided: Vitess Vindexes (Lookup Vindex types).
Source: Vitess Documentation.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Correctly details Lookup Unique Owned and consistent lookup trade-offs.

---

## Claim 8
Claim: Standard auto-increment cannot generate globally unique identifiers across independent database instances without central coordination; distributed ID generation (UUIDv7, Snowflake, Vitess Sequences) is required.
Location: `research/03-evidence.md:Evidence 8`, `research/05-report.md:Finding 4`
Evidence Provided: RFC 9562 §2.1 & §5.7, Vitess Sequences Documentation.
Source: RFC 9562, Vitess Documentation.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Distinguishes coordination-free temporal IDs (UUIDv7) from centralized block-allocation sequences (Vitess Sequences).

---

## Claim 9
Claim: Cross-shard transactions via TwoPC trade commit latency for atomicity without providing full ACID isolation across shards.
Location: `research/03-evidence.md:Evidence 9`, `research/05-report.md:Finding 4`
Evidence Provided: Vitess Distributed Transactions, MongoDB Transactions.
Source: Vitess Documentation, MongoDB Documentation.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fractured reads and cross-shard visibility nuances are accurately captured.

---

## Claim 10
Claim: Resharding reorganizes data to match a new sharding scheme while serving live traffic, with brief read-only cutover downtime.
Location: `research/03-evidence.md:Evidence 10`, `research/05-report.md:Finding 7`
Evidence Provided: Vitess Sharding (Resharding), MongoDB Sharding Balancer.
Source: Vitess Documentation, MongoDB Documentation.
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Distinguishes continuous chunk migration (MongoDB) from operator-triggered workflows (Vitess).
