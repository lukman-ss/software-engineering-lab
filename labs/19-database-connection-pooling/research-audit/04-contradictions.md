# Contradiction Audit

## Contradiction 1

Statement A:
"for optimal throughput the number of active connections should be somewhere near ((core_count * 2) + effective_spindle_count)... There hasn't been any analysis so far regarding how well the formula works with SSDs."

Location:
PostgreSQL Wiki & HikariCP Wiki (Classical Sizing Formula)

Statement B:
"Don't be tricked into thinking, 'SSDs are faster and therefore I can have more threads'. That is exactly 180 degrees backwards. Faster, no seeks, no rotational delays means less blocking and therefore fewer threads [closer to core count] will perform better than more threads."

Location:
HikariCP Wiki (Modern Flash Storage Guidance)

Type:
SOURCE_CONFLICT

Impact:
Practitioners reading only the historical formula might attempt to guess an arbitrary non-zero number for `effective_spindle_count` on SSD arrays, inadvertently over-provisioning pools.

Assessment:
PASS. The research explicitly highlights this divergence, resolves it correctly based on hardware physics (zero rotational latency eliminates I/O wait opportunities for thread context switching), and guides modern configurations toward `core_count * 2` or `core_count`.

---

## Contradiction 2

Statement A:
Application frameworks (HikariCP, Java EE) recommend direct client-side pooling for lowest borrow latency and zero proxy hops.

Location:
HikariCP Architecture Documentation

Statement B:
PostgreSQL architecture documentation recommends external connection poolers (PgBouncer) outside the core server to handle horizontally scaled clients.

Location:
PostgreSQL Wiki ("The Need for an External Pool") & PgBouncer Documentation

Type:
SOURCE_CONFLICT

Impact:
Architecture confusion regarding whether client-side pooling alone is adequate in cloud/microservice deployments.

Assessment:
PASS. The research reconciles the two models: application-side pooling suffices for monolithic/single-instance architectures, while intermediary proxy pooling (e.g., PgBouncer in transaction mode) is necessary when aggregate application workers exceed backend `max_connections`.
