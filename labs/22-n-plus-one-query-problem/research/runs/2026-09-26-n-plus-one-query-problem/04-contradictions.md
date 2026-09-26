# Contradictions

No material contradictions discovered between the authoritative sources investigated.

All primary sources (Laravel, Django, Rails, SQLAlchemy, EF Core) and secondary sources (Vlad Mihalcea, Shopify Engineering) agree on:

1. Definition of the N+1 problem (1 query + N additional queries for N related objects)
2. Lazy loading as default behavior causing N+1
3. Eager loading as the primary solution (batching related data)
4. Trade-offs of eager loading (cartesian product, memory bloat, over-fetching)
5. Detection methods (query logging, profiling, strict loading)
6. Mitigation strategies (eager loading, column selection, pagination, aggregation)

Minor implementation differences exist between ORMs, but these reflect design philosophy, not factual disagreements:

- Rails `includes()` intelligently chooses between `preload()` and `eager_load()` based on conditions.
- Django separates `select_related()` (JOIN, single-valued) from `prefetch_related()` (separate query, collections).
- SQLAlchemy provides more fine-grained strategies: `joinedload`, `selectinload`, `subqueryload`, `raiseload`.
- EF Core uses `Include()` with `ThenInclude()` for nested loading, with a caution about single-query collection performance.
- Laravel uses `with()` for eager loading, `withCount()` for aggregation, and `chaperone()` for parent hydration.

All sources agree that the core problem is architectural and framework-agnostic: any data access pattern that triggers N separate queries where 1 batched query would suffice is an N+1.
