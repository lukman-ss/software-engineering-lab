# Source Audit

**Target Lab:** `labs/22-n-plus-one-query-problem`  
**Audit Scope:** Research Sources (`research/02-sources.md`)  
**Audit Date:** 2026-09-26  

---

## Source 1

Claimed Title: Database access optimization | Django documentation | Django  
Claimed Publisher: Django Software Foundation  
URL: https://docs.djangoproject.com/en/5.1/topics/db/optimization/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Official authoritative Django documentation detailing query profiling, `select_related`, `prefetch_related`, `values`, `defer`, `only`, and bulk operations.  

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

Problems: None. Official Laravel documentation explicitly identifying the "N + 1" query problem, eager loading (`with`), `withCount`, and `chaperone`.  

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

Problems: None. Authoritative Microsoft Learn documentation covering eager, explicit, and lazy loading mechanisms in EF Core.  

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

Problems: None. Official documentation for `Include()`, `ThenInclude()`, filtered includes, and split queries warnings.  

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

Problems: None. Official Rails guide containing dedicated Section 16.1 ("N + 1 Queries Problem"), `includes`, `preload`, `eager_load`, and `strict_loading`.  

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

Problems: None. Authoritative SQLAlchemy documentation covering lazy loading default, `selectinload`, `joinedload`, `subqueryload`, `raiseload`, 500-key batching, and cartesian explosion warnings.  

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

Problems: None. Authoritative reference detailing `select_related()` (single-valued / JOINs) vs `prefetch_related()` (multi-valued / separate query with `IN`).  

Assessment: PASS  

---

## Source 7

Claimed Title: Topic specification — N+1 Query Problem lab (lukman-ss)  
Claimed Publisher: labs/22-n-plus-one-query-problem (source_repository: https://github.com/lukman-ss/software-engineering-lab)  
URL: labs/22-n-plus-one-query-problem  

Reachable: YES  
Source Type: COMMUNITY (Internal Specification)  
Relevant: YES  
Supports Claimed Topic: PARTIAL  

Problems:
- Source 7 is an internal lab requirements / topic specification, not an independent external authoritative standard.
- The research correctly identifies this source as Tier 3 and explicitly flags illustrative numbers (e.g. 712 queries, 2.4s) as scenario-specific rather than universal benchmarks.  

Assessment: PASS (Correctly classified and scoped)  
