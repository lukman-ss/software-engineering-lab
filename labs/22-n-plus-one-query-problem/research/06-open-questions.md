# Open Questions

## Unanswered Questions

1. **What are the quantitative performance benefits of each eager loading strategy?**
   - The documentation describes joined vs. select-in trade-offs but lacks specific benchmarks showing latency improvements across different dataset sizes and database backends.
   - Source: SQLAlchemy docs state "selectin loading is usually the most simple and efficient way to eagerly load collections" but provides no comparative metrics.

2. **How does the 500-row batch size in selectin loading interact with database query planning?**
   - SQLAlchemy docs mention batching to 500 parent keys but note that some databases (Oracle) have hard limits on IN clause size. What is the actual performance impact of different batch sizes on different databases?

3. **What is the interaction between N+1 detection and database indexes?**
   - The lab specification lists "Index benar?" as a step in the debugging workflow, but the documentation doesn't explain how missing indexes exacerbate N+1 or whether N+1 creates specific index requirements.

4. **How do ORMs handle N+1 in complex polymorphic or single-table inheritance scenarios?**
   - The lab scenario involves work orders, mechanics, and service items - potentially involving polymorphic relationships. The documentation mentions polymorphic eager loading (e.g., in EF Core) but focuses on simple foreign key cases.

5. **What is the memory usage profile of different eager loading strategies?**
   - The docs warn about cartesian explosion and memory overhead but lack specific guidance on estimating memory consumption for large result sets with multiple relationships.

6. **Are there automated tools that can detect N+1 in CI/CD pipelines?**
   - The docs recommend profiling tools for development, but is there a standardized approach to automatically flag N+1 in automated tests?

## Weak Evidence

1. **Performance numbers in lab specification** (e.g., "1 request = 712 queries = 2.4 seconds") are presented as realistic monitoring signals but are not cross-referenced with independent benchmarks. These appear to be scenario-specific examples rather than universally applicable figures.

2. **The "Rule of Thumb" workflow** in the lab specification matches Django's "Profile first" guidance, but this is the only source presenting this specific 6-step debugging sequence. While the steps are individually sound, their ordering as presented is not explicitly documented in other ORM guides.

## Claims Needing Deeper Research

1. **Connection pool exhaustion mechanics**: The lab describes a cascade (too many queries → connection pool full → request queuing → latency → timeout). This is plausible but requires deeper investigation into how connection pooling implementations handle long-running multi-query requests.

2. **Batch size optimization**: The 500-row batch limit in selectin loading appears arbitrary. Research into database IN clause performance and memory characteristics could provide better guidance.

3. **Cache as concealment**: The lab warns that "cache applied on top of bad queries only hides the problem, not solves it." This is conceptually sound but would benefit from quantitative evidence showing how caching masks N+1 vs. actually fixing it.

## Possible Next Research Directions

1. **Compare eager loading strategies with TPC-H-style workloads**: How do joinedload vs. selectinload vs. subqueryload perform on normalized schemas with multiple relationships per table?

2. **Benchmark cache + N+1 vs. query optimization**: Measure actual throughput improvements from fixing N+1 vs. adding Redis cache layer.

3. **Investigate bulk insert/update patterns**: While the Django docs mention `bulk_create()` and `bulk_update()`, how do these interact with eager loading and do they have their own N+1 anti-patterns?

4. **GraphQL data loader patterns**: Research how Facebook's DataLoader pattern handles N+1 in GraphQL resolvers and whether similar batching approaches exist outside GraphQL.

5. **Database-specific query optimization**: How do PostgreSQL's planner decisions differ from MySQL/MariaDB when handling large IN clauses versus JOINs for eager loading?