# Source Audit

## Source 1

Claimed Title: PostgreSQL Documentation: Connections and Authentication
Claimed Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/runtime-config-connection.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. URL is active and contains the claimed documentation on `max_connections`, `superuser_reserved_connections`, and shared memory allocations.

Assessment:
PASS

---

## Source 2

Claimed Title: About Pool Sizing
Claimed Publisher: HikariCP Wiki / Brett Wooldridge
URL: https://github.com/brettwooldridge/HikariCP/wiki/About-Pool-Sizing

Reachable:
YES

Source Type:
PRIMARY / COMMUNITY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. URL is active. Quotes ("50x improvement", formula `Tn x (Cm - 1) + 1`, and `((core_count * 2) + effective_spindle_count)`) are exact matches to the page content.

Assessment:
PASS

---

## Source 3

Claimed Title: Number Of Database Connections
Claimed Publisher: PostgreSQL Wiki
URL: https://wiki.postgresql.org/wiki/Number_Of_Database_Connections

Reachable:
YES

Source Type:
SECONDARY / COMMUNITY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Content perfectly corroborates claims around lock contention, context switches, cache line contention, and disk thrashing.

Assessment:
PASS

---

## Source 4

Claimed Title: PgBouncer Features
Claimed Publisher: PgBouncer Authors
URL: https://www.pgbouncer.org/features.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Document explicitly cites "Transaction pooling" and "Low memory requirements (2 kB per connection by default)".

Assessment:
PASS

---

## Source 5

Claimed Title: PostgreSQL Documentation: The Cumulative Statistics System (pg_stat_activity)
Claimed Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/monitoring-stats.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Document lists state parameters (`idle in transaction`) and discusses client wait events.

Assessment:
PASS
