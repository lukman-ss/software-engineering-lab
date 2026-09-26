# Evidence

## Evidence 1

Claim: N+1 query problem occurs when fetching N related objects triggers N+1 database queries (1 for the main query + N for each related object accessed via lazy loading)

Evidence: 
- Laravel Eloquent example: fetching 100 invoices with `$invoice->customer->name` in a loop causes 101 queries (1 SELECT * FROM invoices + 100 SELECT * FROM customers WHERE id = X)
- Rails documentation states: "Retrieving a list of records N (where N is a number greater that 1) in a single query can sometimes trigger N extra queries; one for each record"
- Django documentation describes N+1 as accessing related objects via lazy loading after fetching a queryset
- SQLAlchemy identifies this as the "N plus one problem" where accessing lazy-loaded attributes on N objects triggers N additional SELECT statements

Source: Laravel Eloquent Relationships documentation
URL: https://laravel.com/docs/13.x/eloquent-relationships
Confidence: HIGH

Corroborated By: Rails Active Record Query Interface, Django ORM optimization guide, SQLAlchemy relationship loading documentation

Notes: All major ORMs describe the same N+1 pattern: fetch collection → loop → access relationship → trigger additional query per object

---

## Evidence 2

Claim: Eager loading solves N+1 by fetching related data in fewer queries (typically 2 queries instead of N+1)

Evidence:
- Rails `includes()` method generates 2 queries: one for main table, one with `IN` clause for related table
- Laravel `with()` method similarly generates separate queries instead of joining
- Django `prefetch_related()` executes a second query with `WHERE id IN (...)` to load all related objects at once
- SQLAlchemy `selectinload()` emits a second SELECT with `WHERE foreign_key IN (primary_key_list)`
- Entity Framework `Include()` uses either JOIN or separate queries depending on configuration

Source: Rails Active Record Query Interface guide (Section 16.2)
URL: https://guides.rubyonrails.org/active_record_querying.html
Confidence: HIGH

Corroborated By: Django optimization guide (select_related/prefetch_related), Laravel docs (with method), SQLAlchemy docs (selectinload), EF Core docs (Include)

Notes: All ORMs provide eager loading mechanisms. Strategy differs: Rails/Django/SQLAlchemy typically use separate queries with IN clauses for collections; joined loading available but has tradeoffs

---

## Evidence 3

Claim: Eager loading trade-offs include cartesian product explosion, memory usage, and over-fetching

Evidence:
- SQLAlchemy docs warn: "joinedload() creates anonymous aliases to avoid affecting query results" and when joining collections, rows are multiplied
- Django docs note: "Because ManyToManyField attributes and reverse relations can have multiple related rows, including these can have a multiplier effect on the size of your result set"
- Rails docs: "Using `includes()` with conditions can generate LEFT OUTER JOIN whereas `joins()` uses INNER JOIN"
- EF Core docs: "Eager loading a collection navigation in a single query may cause performance issues"

Source: SQLAlchemy Relationship Loading Techniques - Joined Eager Loading
URL: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
Confidence: HIGH

Corroborated By: Django QuerySet API (values() warning), Rails guides (includes vs joins comparison), EF Core docs (single vs split queries)

Notes: The "Zen of Joined Eager Loading" concept from SQLAlchemy emphasizes that eager loading should not change query results, only loading strategy — this requires workarounds that can impact performance

---

## Evidence 4

Claim: Column selection (`select()`, `pluck()`, `values()`, `only()`, `defer()`) reduces data transfer and can replace relationship loading

Evidence:
- Rails `pluck()` directly returns array of values without constructing model instances
- Django `values()` returns dictionaries instead of model instances for subsets of fields
- Django `only()` and `defer()` control which fields are loaded
- Laravel `select()` limits columns retrieved
- SQLAlchemy `load_only()` loads only specific columns

Source: Django ORM optimization guide
URL: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
Confidence: HIGH

Corroborated By: Rails pluck documentation, Laravel query builder docs, SQLAlchemy column loading options

Notes: These methods are often more efficient than eager loading when you only need specific fields or counts rather than full objects

---

## Evidence 5

Claim: Pagination is critical to prevent N+1 from scaling catastrophically with dataset size

Evidence:
- All ORMs support limit/offset or cursor-based pagination
- Rails: "The `find_each` method retrieves records in batches and then yields each one to the block"
- Django: "Use `iterator()` when you have a lot of objects and caching causes memory issues"
- Laravel: "The `chunk` method is useful for processing many records when processing the data in batches"
- SQLAlchemy: "The strategy emits a SELECT for up to 500 parent primary key values at a time"

Source: Rails Active Record Query Interface - Retrieving Multiple Records in Batches
URL: https://guides.rubyonrails.org/active_record_querying.html
Confidence: HIGH

Corroborated By: Django iterator() documentation, Laravel chunk() method, SQLAlchemy batch size configuration

Notes: Without pagination or batching, N+1 on 10,000 records = 10,001 queries. With pagination (50 per page), impact is limited per request

---

## Evidence 6

Claim: ORMs provide tools to detect N+1: strict loading, logging, query counting

Evidence:
- Rails `strict_loading` option raises error when lazy loading occurs
- Django `connection.queries` shows executed SQL statements
- Laravel with Laravel Debugbar or Query Counter packages
- EF Core logging can track query count
- SQLAlchemy `echo=True` or event listeners

Source: Rails Guide - Eager Loading Associations (strict_loading section)
URL: https://guides.rubyonrails.org/active_record_querying.html
Confidence: HIGH

Corroborated By: Django database query logging, EF Core logging configuration

Notes: Detection is essential — N+1 often doesn't manifest with small test datasets, only in production with real data volume

---

## Evidence 7

Claim: `withCount()` or aggregation solves N+1 when only count is needed

Evidence:
- Laravel: `Post::withCount('comments')->get()` generates 2 queries instead of N+1
- Django: `Blog.objects.annotate(num_entries=Count('entry'))` adds count as annotation
- Rails: `Author.includes(:books).count` works efficiently
- SQLAlchemy: `func.count()` in query projection

Source: Laravel Eloquent Relationships - Aggregating Related Models
URL: https://laravel.com/docs/13.x/eloquent-relationships
Confidence: HIGH

Corroborated By: Django aggregation guide, SQLAlchemy func documentation

Notes: This is often the most efficient solution when you only need a scalar value like count, not the full collection

---

## Evidence 8

Claim: N+1 concept extends beyond database to API/microservice calls

Evidence:
- The user's lab specification mentions: "GET /orders returns 100 order. Backend does GET Payment API /payment/1 through /payment/100"
- GraphQL uses DataLoader pattern to batch resolvers
- Microservices often batch operations via bulk endpoints

Source: Original lab specification (topic description)
URL: labs/22-n-plus-one-query-problem
Confidence: MEDIUM

Corroborated By: GraphQL DataLoader documentation (external source)

Notes: The pattern is identical: batch operations instead of per-object requests. Database: IN clause. API: bulk endpoint with JSON array

---

## Evidence 9

Claim: Recommended troubleshooting workflow: measure first, identify hotspots, optimize queries before scaling

Evidence:
- Django docs: "Profile first... Find out what queries you are doing and what they are costing you"
- User's lab specification outlines workflow: "Berapa query yang dijalankan? Query apa yang paling mahal? Ada N+1? Index benar? Data yang diambil terlalu banyak? Pagination benar? Baru pertimbangkan caching"
- Rails guides emphasize profiling before optimization

Source: Django optimization guide - "Profile first" section
URL: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
Confidence: HIGH

Corroborated By: User's lab specification (identical workflow pattern), general performance engineering best practices

Notes: The "rule of thumb" from the lab specification reflects industry consensus: measure → identify pattern → fix root cause → only then consider infrastructure scaling

---

## Evidence 10

Claim: Common mistakes: over-eager loading, testing with small datasets, using SELECT *, ignoring pagination

Evidence:
- Laravel docs warn about loading "10,000 invoice × 20 item × several relationship" into memory
- Django: "Don't overuse `contains()`, `count()`, and `exists()`... if you are going to need other data from the QuerySet, evaluate it immediately"
- Rails: "The `take` method returns `nil` if no record is found and no exception will be raised. Since `take` doesn't specify an ORDER BY clause, the retrieved record may vary depending on the database engine"

Source: User's lab specification - "Kesalahan Umum" section
URL: labs/22-n-plus-one-query-problem
Confidence: MEDIUM

Corroborated By: General performance anti-patterns documented across ORMs

Notes: These anti-patterns are consistent across all ORMs despite different syntax

---

## Evidence 11

Claim: Different ORMs use different terminology but same concepts

Evidence:
- Rails: `includes()`, `preload()`, `eager_load()`, `strict_loading`
- Django: `select_related()` (for single-valued), `prefetch_related()` (for collections)
- Laravel: `with()`, `withCount()`, `chaperone()`
- EF Core: `Include()`, `ThenInclude()`, `IgnoreAutoIncludes()`
- SQLAlchemy: `joinedload()`, `selectinload()`, `subqueryload()`, `lazyload()`, `raiseload()`

Source: Cross-referencing official documentation for each ORM
URLs:
- Rails: https://guides.rubyonrails.org/active_record_querying.html
- Django: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
- Laravel: https://laravel.com/docs/13.x/eloquent-relationships
- EF Core: https://learn.microsoft.com/en-us/ef/core/querying/related-data/eager
- SQLAlchemy: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
Confidence: HIGH

Corroborated By: Direct comparison of documentation sections on eager loading

Notes: Despite different API names, the concepts map consistently:
Rails includes ≈ Django prefetch_related ≈ SQLAlchemy selectinload
Rails eager_load ≈ Django select_related (for single) ≈ SQLAlchemy joinedload

---

## Evidence 12

Claim: Lazy loading is the default in most ORMs, causing N+1 if not understood

Evidence:
- SQLAlchemy: "Lazy loading is the default loading style for all relationship constructs"
- Django: QuerySets are lazy — not executed until iterated
- Rails: associations load on first access by default
- EF Core: navigation properties load on access
- Laravel: relationships are accessed as dynamic properties (lazy by default)

Source: SQLAlchemy Relationship Loading Techniques - Lazy Loading section
URL: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
Confidence: HIGH

Corroborated By: All other ORM docs confirm lazy loading is default

Notes: This is why N+1 is so common — code looks correct with small data, fails in production

---