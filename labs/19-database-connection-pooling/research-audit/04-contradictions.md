# Contradictions Audit

## Contradiction 1

Statement A:
Active connections should be sized near `((core_count * 2) + effective_spindle_count)`.
Location:
`04-contradictions.md`: Contradiction 1 (PostgreSQL Wiki / HikariCP Wiki)

Statement B:
SSDs eliminate seeks and rotational delays, meaning fewer threads (closer to `core_count`) perform better than higher connection numbers.
Location:
`04-contradictions.md`: Contradiction 1 (Modern Cloud / HikariCP Commentary)

Type:
SOURCE_CONFLICT

Impact:
Practitioners might oversize pools if assuming `effective_spindle_count` applies to flash storage.

Assessment:
The research resolves this conflict accurately: modern NVMe/SSD storage effectively sets spindle wait times to zero, making `core_count` or `core_count * 2` the realistic upper ceiling rather than an inflated count.

---

## Contradiction 2

Statement A:
Connection pooling should live inside the client application (in-memory, lowest latency, zero proxy hops).
Location:
`04-contradictions.md`: Contradiction 2 (HikariCP / Java EE view)

Statement B:
PostgreSQL explicitly chose not to implement internal pooling and recommends intermediate proxy poolers (PgBouncer) for large, distributed deployments.
Location:
`04-contradictions.md`: Contradiction 2 (PostgreSQL Architecture / PgBouncer view)

Type:
INTERNAL / ARCHITECTURAL

Impact:
Engineers may believe client-side pooling and intermediate proxy pooling are incompatible or competing alternatives.

Assessment:
The research correctly delineates topology: single monolithic or small instances benefit from in-process client pooling; multi-instance auto-scaling worker tiers require proxy-based pooling (transaction pooling) to protect database backend resources from multiplication exhaustion.
