# Open Questions and Research Frontiers

## Unanswered Questions
1. **Handling Distributed Transactions**: How do teams manage dual-write consistency across microservices safely without employing full two-phase commits (2PC) or sacrificing throughput?
2. **Schema Migration Tooling**: Do existing migration tools (e.g., Flyway, Liquibase, Prisma Migrate) natively support automated checks that block "destructive" operations during the expand phase?
3. **Consumer Lag**: For external APIs, how can API providers safely contract an interface if a small percentage of critical B2B consumers refuse to migrate off the deprecated version for years?

## Weak Evidence
- Strict quantitative evidence on the overhead of internal API transformation layers (e.g., CPU/Memory impact of running dozens of backward-walking transformation modules per request) is sparse outside of Stripe's initial descriptions.

## Claims Needing Deeper Research
- **Lazy Backfilling vs. Batch Backfilling**: Under what precise dataset sizes and latency thresholds does lazy backfilling (fallback read + write) become strictly superior to batch background worker backfilling?

## Next Research Directions
- Investigating the CDC (Change Data Capture) pattern (e.g., Debezium) as a mechanism for handling the dual-write/migration phase asynchronously instead of modifying application-level code.
- How GraphQL interfaces naturally handle Expand/Contract patterns compared to REST API versioning.
