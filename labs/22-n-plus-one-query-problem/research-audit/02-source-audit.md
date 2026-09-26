# Source Audit

## Source 1
Claimed Title: Database access optimization | Django documentation | Django
Claimed Publisher: Django Software Foundation
URL: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Documentation explicitly covers "Profile first", `select_related()`, `prefetch_related()`, `values()`, `defer()`, `only()`, `count()`, `exists()`, and bulk operations.
Assessment: PASS

---

## Source 2
Claimed Title: Eloquent: Relationships | Laravel 13.x
Claimed Publisher: Laravel
URL: https://laravel.com/docs/13.x/eloquent-relationships
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Covers `with()`, `withCount()`, `chaperone()`, lazy loading behavior, and explicit mentions of the N+1 problem.
Assessment: PASS

---

## Source 3
Claimed Title: Loading Related Data - EF Core | Microsoft Learn
Claimed Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/ef/core/querying/related-data/
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Defines eager, explicit, and lazy loading mechanisms in EF Core.
Assessment: PASS

---

## Source 3a
Claimed Title: Eager Loading of Related Data - EF Core | Microsoft Learn
Claimed Publisher: Microsoft
URL: https://learn.microsoft.com/en-us/ef/core/querying/related-data/eager
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Covers `Include()`, `ThenInclude()`, and single vs split queries tradeoff.
Assessment: PASS

---

## Source 4
Claimed Title: Active Record Query Interface | Ruby on Rails Guides
Claimed Publisher: Ruby on Rails Core Team
URL: https://guides.rubyonrails.org/active_record_querying.html
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Section 16.1 explicitly details the "N + 1 Queries Problem", `includes()`, `preload()`, `eager_load()`, `strict_loading`, and `pluck()`.
Assessment: PASS

---

## Source 5
Claimed Title: Relationship Loading Techniques — SQLAlchemy 2.1 Documentation
Claimed Publisher: SQLAlchemy authors
URL: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Details lazy loading default, selectinload, joinedload cartesian explosion warnings, and raiseload.
Assessment: PASS

---

## Source 6
Claimed Title: QuerySet API reference | Django documentation
Claimed Publisher: Django Software Foundation
URL: https://docs.djangoproject.com/en/5.1/ref/models/querysets/
Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None. Authoritative API reference for `select_related`, `prefetch_related`, and queryset evaluation caching.
Assessment: PASS

---

## Source 7
Claimed Title: Topic specification — N+1 Query Problem lab (lukman-ss)
Claimed Publisher: labs/22-n-plus-one-query-problem (source_repository: https://github.com/lukman-ss/software-engineering-lab)
URL: labs/22-n-plus-one-query-problem
Reachable: YES (Local repository file)
Source Type: COMMUNITY / INTERNAL (Tier 3)
Relevant: YES
Supports Claimed Topic: PARTIAL
Problems:
- Internal task specification, not an authoritative external engineering reference.
- Used to seed specific numbers (712 queries = 2.4s, 180ms target) and local narrative claims.
- Research correctly classified this as Tier 3 and noted numbers are illustrative rather than universal facts.
Assessment: WARNING
