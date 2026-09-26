# Source List

## Source 1

Title: Database access optimization | Django documentation | Django
Publisher: Django Software Foundation
URL: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
Published: 2026
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Core ORM optimization patterns including N+1 detection, `select_related()`, `prefetch_related()`, `values()`, `defer()`, `only()`, `count()`, `exists()`, bulk operations, profiling-first approach

## Source 2

Title: Eloquent: Relationships | Laravel 13.x
Publisher: Laravel
URL: https://laravel.com/docs/13.x/eloquent-relationships
Published: 2026
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: `with()`, `withCount()`, `chaperone()`, eager loading patterns, lazy loading, the explicit term "N + 1", aggregating related models (Counting Related Models section)

## Source 3

Title: Loading Related Data - EF Core | Microsoft Learn
Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/ef/core/querying/related-data/
Published: 2026-06-24
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Three loading patterns defined: eager loading, explicit loading, lazy loading (lazy = "transparently loaded... when the navigation property is accessed")

## Source 3a

Title: Eager Loading of Related Data - EF Core | Microsoft Learn
Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/ef/core/querying/related-data/eager
Published: 2026-07-31
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: `Include()`, `ThenInclude()`, filtered includes, auto-include configuration, single vs split queries caution

## Source 4

Title: Active Record Query Interface | Ruby on Rails Guides
Publisher: Ruby on Rails Core Team
URL: https://guides.rubyonrails.org/active_record_querying.html
Published: 2026
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Explicit "N + 1 Queries Problem" section (16.1), `includes()`, `preload()`, `eager_load()`, `strict_loading`, `pluck()`, batch processing (`find_each`)

## Source 5

Title: Relationship Loading Techniques — SQLAlchemy 2.1 Documentation
Publisher: SQLAlchemy authors
URL: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
Published: 2026
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Lazy loading default, "N plus one problem" reference, `selectinload()`, `joinedload()`, `subqueryload()`, `raiseload()`, Select IN loading batch size (500 keys), cartesian explosion warning for joined loading

## Source 6

Title: QuerySet API reference | Django documentation
Publisher: Django Software Foundation
URL: https://docs.djangoproject.com/en/5.1/ref/models/querysets/
Published: 2026
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: `select_related()` vs `prefetch_related()` distinction, `values()`/`values_list()`, `only()`/`defer()`, `iterator()`, query evaluation/caching behavior, N+1 via lazy evaluation of related collections

## Source 7

Title: Topic specification — N+1 Query Problem lab (lukman-ss)
Publisher: labs/22-n-plus-one-query-problem (source_repository: https://github.com/lukman-ss/software-engineering-lab)
URL: labs/22-n-plus-one-query-problem
Published: N/A
Accessed: 2026-09-26
Source Tier: Tier 3
Relevance: Source of the lab scenario's specific claims (712 queries = 2.4s, 180ms target, connection pool cascade, network N+1 analogy, troubleshooting workflow, common mistakes)