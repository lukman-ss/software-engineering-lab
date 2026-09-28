# Source Audit

## Source 1
Claimed Title: PostgreSQL 18 Documentation - Chapter 5. Data Definition: Table Partitioning
Claimed Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/ddl-partitioning.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Minor: PostgreSQL version in URL uses `current` alias rather than hardcoded `18`. PostgreSQL declarative table partitioning concepts apply from PostgreSQL 10 through current versions.

Assessment:
PASS

---

## Source 2
Claimed Title: MongoDB Manual - Sharding
Claimed Publisher: MongoDB, Inc.
URL: https://www.mongodb.com/docs/manual/sharding/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative reference for multi-instance physical sharding architecture, mongos routing, chunk allocation, and balancer mechanics.

Assessment:
PASS

---

## Source 3
Claimed Title: MongoDB Manual - Choose a Shard Key
Claimed Publisher: MongoDB, Inc.
URL: https://www.mongodb.com/docs/manual/core/sharding-choose-a-shard-key/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Comprehensive documentation regarding shard key cardinality, frequency, monotonicity write-hotspots, and query isolation.

Assessment:
PASS

---

## Source 4
Claimed Title: Vitess Documentation - Sharding
Claimed Publisher: Vitess.io / PlanetScale
URL: https://www.vitess.io/docs/24.0/reference/features/sharding/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Direct reference for MySQL horizontal sharding, keyspace IDs, key ranges, and resharding workflows.

Assessment:
PASS

---

## Source 5
Claimed Title: Vitess Documentation - Vindexes
Claimed Publisher: Vitess.io / PlanetScale
URL: https://www.vitess.io/docs/24.0/reference/features/vindexes/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative architectural guide for Primary Vindexes, Secondary Lookup Vindexes, functional vs lookup mapping, and scatter-gather avoidance.

Assessment:
PASS

---

## Source 6
Claimed Title: Vitess Documentation - Distributed Transactions
Claimed Publisher: Vitess.io / PlanetScale
URL: https://www.vitess.io/docs/24.0/reference/features/distributed-transaction/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Documents 2PC (Two-Phase Commit) implementation in Vitess, atomicity guarantees, and cross-shard transaction tradeoffs.

Assessment:
PASS

---

## Source 7
Claimed Title: Vitess Documentation - Sequences
Claimed Publisher: Vitess.io / PlanetScale
URL: https://www.vitess.io/docs/24.0/reference/features/vitess-sequences/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Explains centralized sequence tables with block allocation to replace distributed auto_increment limitations.

Assessment:
PASS

---

## Source 8
Claimed Title: RFC 9562 - Universally Unique IDentifiers (UUIDs)
Claimed Publisher: IETF (Internet Engineering Task Force)
URL: https://www.rfc-editor.org/rfc/rfc9562.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Standards-track RFC defining UUIDv7 timestamp-ordered specifications and B-tree index locality benefits over random UUIDv4.

Assessment:
PASS

---

## Source 9
Claimed Title: Consistent Hashing and Random Trees: Distributed Caching Protocols for Relieving Hot Spots on the World Wide Web
Claimed Publisher: ACM STOC '97 / MIT
URL: https://doi.org/10.1145/258533.258660

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Foundational paper establishing $O(K/N)$ key remapping bound for consistent hashing vs modulo hashing.

Assessment:
PASS
