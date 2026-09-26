# Contradictions

No material contradictions discovered.

All authoritative sources (Django, Laravel, Rails, EF Core, SQLAlchemy) agree on:
1. Definition of N+1 problem (1 query + N additional queries for N related objects)
2. Eager loading as primary solution (fetching related data in fewer queries)
3. Lazy loading as default behavior in most ORMs
4. Trade-offs: joined loading can cause cartesian explosion, select-in is safer for collections
5. Detection methods: query logging, profiling tools, strict loading modes
6. Mitigation strategies: eager loading, column selection, pagination, aggregation

Minor differences in implementation details:
- Rails uses `includes()` which intelligently chooses between `preload()` and `eager_load()`
- Django separates `select_related()` (JOIN, single-valued) from `prefetch_related()` (separate query, collections)
- SQLAlchemy provides multiple strategies: joinedload, selectinload, subqueryload with different trade-offs
- EF Core uses `Include()` with `ThenInclude()` for nested loading

These differences reflect ORM design philosophy, not factual disagreements about the N+1 problem itself.