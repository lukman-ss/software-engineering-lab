# Open Questions

## Unanswered Questions
- What are the real-world performance benchmarks for N+1 vs eager loading vs aggregation across different databases (MySQL, PostgreSQL, SQL Server)? The sources provide theoretical analysis but few concrete cross-database benchmarks.
- What specific APM tools or middleware are most effective at detecting N+1 in production at scale? Sources mention query logging and Datadog/New Relic but no authoritative comparison.
- What is the industry-standard query-count threshold for flagging a potential N+1 issue? (e.g., "more than X queries per HTTP request is suspicious")

## Weak Evidence
- The lab specification cites specific figures: 712 queries = 2.4 seconds, target P95 < 400ms, 180ms after optimization. These appear illustrative rather than from a published benchmark. No independent verification source found.
- The network N+1 section relies heavily on the lab specification's example (`GET /payment/1` through `/payment/100`) rather than a widely-cited production case study. Shopify's article focuses on GraphQL N+1 specifically.
- The "201 query × 5ms ≈ 1 second" latency extrapolation is from the lab spec (Tier 3 source), not independently cross-checked.

## Claims Needing Deeper Research
- The precise conditions under which `select_related()` vs `prefetch_related()` (Django) or `preload()` vs `eager_load()` (Rails) each produce better performance for different relationship cardinalities.
- Real production case studies (beyond Shopify's GraphQL example) of N+1 remediation ROI metrics — e.g., cost savings from query reduction vs. server scaling.
- The interaction of N+1 with connection pooling under concurrent load — the lab spec describes the cascade (pool exhaustion → queueing → timeout) but this deserves quantitative analysis.

## Possible Next Research Directions
- Benchmark N+1 vs eager loading with realistic data volumes across MySQL, PostgreSQL, and SQLite using the lab's Go demo implementation.
- Investigate ORM-specific N+1 detection extensions (e.g., Laravel's `strict` mode, Rails' `Bullet` gem, SQLAlchemy's `nplus1` detection).
- Survey production observability platforms (Datadog APM, New Relic) for built-in N+1 detection capabilities.
- Explore whether "automatic eager loading" features (Laravel 13's automatic eager loading) effectively prevent N+1 without developer intervention.
