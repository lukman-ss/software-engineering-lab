# Source List

## Source 1

Title: N+1 query problem with JPA and Hibernate
Publisher: Vlad Mihalcea
URL: https://vladmihalcea.com/n-plus-1-query-problem/
Published: 2020-03-17
Accessed: 2026-09-26
Source Tier: Tier 2
Relevance: Authoritative technical article defining N+1 for JPA/Hibernate, covering FetchType.EAGER default, FetchType.LAZY with lazy load triggers, JOIN FETCH fixes, second-level cache N+1, and the stealthiness of N+1 in slow query logs.

## Source 2

Title: Eloquent: Relationships | Laravel 13.x
Publisher: Laravel
URL: https://laravel.com/docs/13.x/eloquent-relationships
Published: 2026
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Official Laravel documentation for `with()` (eager loading), `withCount()` (aggregation), `chaperone()` (parent hydration), and N+1 prevention via `preventLazyLoading()`.

## Source 3

Title: Database access optimization | Django documentation | Django
Publisher: Django Software Foundation
URL: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
Published: 2026
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Official Django ORM documentation covering `select_related()`, `prefetch_related()`, `values()`, `only()`, `defer()`, `count()`, `exists()`, bulk methods, and the "profile first" approach.

## Source 4

Title: Active Record Query Interface | Ruby on Rails Guides
Publisher: Ruby on Rails Core Team
URL: https://guides.rubyonrails.org/active_record_querying.html
Published: 2026
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Official Rails documentation with dedicated "N + 1 Queries Problem" section (16.1), `includes()`, `preload()`, `eager_load()`, `strict_loading`, `pluck()`, `find_each` batch processing.

## Source 5

Title: Relationship Loading Techniques — SQLAlchemy 2.1 Documentation
Publisher: SQLAlchemy authors
URL: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
Published: 2026-09-25 (release 2.1.1)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Official SQLAlchemy documentation for lazy loading default, `selectinload()`, `joinedload()`, `subqueryload()`, `raiseload()`, N+1 problem definition, select-in batch size of 500 keys.

## Source 6

Title: Eager Loading of Related Data - EF Core | Microsoft Learn
Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/ef/core/querying/related-data/eager
Published: 2026-07-31 (updated)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Official EF Core documentation for `Include()`, `ThenInclude()`, filtered includes, auto-include configuration, and the caution about collection navigation performance.

## Source 7

Title: Solving the N+1 Problem for GraphQL through Batching
Publisher: Shopify Engineering
URL: https://shopify.engineering/solving-the-n-1-problem-for-graphql-through-batching
Published: 2018-04-24
Accessed: 2026-09-26
Source Tier: Tier 2
Relevance: Shopify case study on GraphQL N+1 (network N+1), DataLoader pattern, graphql-batch Ruby library, resolver batching.
