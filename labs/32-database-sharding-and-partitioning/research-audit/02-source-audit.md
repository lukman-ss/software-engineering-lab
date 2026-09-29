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
- URL uses `/docs/current/` which resolves dynamically to the current active release. In late 2026, PG 18 is current.

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
- None. Standard official documentation for MongoDB multi-instance sharding.

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
- None. Direct primary source for cardinality, frequency, and monotonicity tradeoffs.

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
- Specific release version (24.0) pinned. Valid primary documentation for MySQL horizontal sharding.

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
- None. Primary source for primary vs secondary lookup vindexes.

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
- None. Covers TwoPC latency and isolation trade-offs.

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
- None. Documents centralized table-backed sequence allocation for sharded MySQL.

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
- None. Official IETF standard for UUIDv7 time-ordered identifiers.

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
- Full text behind paywall/landing page; paper restatement was cross-checked via secondary reference (Wikipedia Consistent Hashing), which is explicitly recorded in research metadata.

Assessment:
PASS
