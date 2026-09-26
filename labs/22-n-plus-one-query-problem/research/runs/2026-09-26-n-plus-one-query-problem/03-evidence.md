# Evidence

## Evidence 1

Claim: The N+1 query problem occurs when an application executes 1 primary query to fetch N records, then N additional queries (one per record) to fetch related data.

Evidence:
- Mihalcea: "The N+1 query problem happens when the data access framework executes N additional SQL statements to fetch the same data that could have been retrieved when executing the primary SQL query."
- Rails docs: "Retrieving a list of records N (where N is a number greater than 1) in a single query can sometimes trigger N extra queries; one for each record."
- SQLAlchemy docs define lazy loading as default: accessing a relationship on N objects triggers N additional SELECTs.
- Laravel, Django, and EF Core all describe the same pattern: lazy loading triggers additional queries per object.

Source: Vlad Mihalcea, Rails Active Record Query Interface, SQLAlchemy Relationship Loading Techniques
URLs:
- https://vladmihalcea.com/n-plus-1-query-problem/
- https://guides.rubyonrails.org/active_record_querying.html
- https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
Confidence: HIGH

Corroborated By: Laravel Eloquent docs, Django ORM optimization guide, EF Core related data docs

---

## Evidence 2

Claim: N+1 queries bypass slow query logs because each individual query is fast; the aggregate latency is what causes production degradation.

Evidence:
- Mihalcea: "Unlike the slow query log that can help you find slow-running queries, the N+1 issue won't be spotted because each individual additional query runs sufficiently fast to not trigger the slow query log."
- The lab specification calculates: 201 queries × 5ms each ≈ 1 second total, plus network latency and contention.

Source: Vlad Mihalcea
URL: https://vladmihalcea.com/n-plus-1-query-problem/
Confidence: HIGH

Corroborated By: Lab specification performance analysis (2026-09-26)

---

## Evidence 3

Claim: Lazy loading is the default behavior in most ORMs, making N+1 a pervasive risk.

Evidence:
- SQLAlchemy: "Lazy loading is the default loading style for all relationship constructs that don't otherwise indicate the lazy option."
- Django: "QuerySets are lazy — the creation of a QuerySet does not hit the database."
- Rails: "Associations are loaded the first time the method is used."
- EF Core: Navigation properties load on access by default.
- Laravel: Relationships are accessed as dynamic properties (lazy by default).

Source: SQLAlchemy Relationship Loading Techniques
URL: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
Confidence: HIGH

Corroborated By: All 5 ORM documentation sets confirm lazy loading as default

---

## Evidence 4

Claim: Eager loading solves N+1 by fetching related data in fewer queries (typically 2 instead of N+1).

Evidence:
- Rails `includes()` generates 2 queries: one for the main table, one with `IN` clause for related data.
- Laravel `with()` similarly uses separate queries instead of joining.
- Django `prefetch_related()` executes a second query with `WHERE id IN (...)`.
- SQLAlchemy `selectinload()` emits `WHERE foreign_key IN (primary_key_list)` with a batch size of up to 500.
- EF Core `Include()` uses JOIN or separate queries depending on configuration.

Source: Rails Active Record Query Interface guide (Section 16)
URL: https://guides.rubyonrails.org/active_record_querying.html
Confidence: HIGH

Corroborated By: Django optimization guide, Laravel docs, SQLAlchemy docs, EF Core docs

---

## Evidence 5

Claim: Eager loading has trade-offs: cartesian product explosion (joined loading), memory bloat (over-fetching), and SELECT * overloading.

Evidence:
- SQLAlchemy warns: joined loading creates anonymous aliases; joining collections multiplies rows (cartesian product).
- Django docs: "ManyToManyField attributes and reverse relations can have multiple related rows, including these can have a multiplier effect."
- EF Core: "Eager loading a collection navigation in a single query may cause performance issues."
- Mihalcea: "Using FetchType.EAGER either implicitly or explicitly for your JPA associations is a bad idea because you are going to fetch way more data that you need."

Source: SQLAlchemy Relationship Loading Techniques
URL: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
Confidence: HIGH

Corroborated By: Django QuerySet API, EF Core eager loading docs, Vlad Mihalcea

---

## Evidence 6

Claim: `withCount()` or aggregation eliminates N+1 when only scalar values (count, sum, etc.) are needed.

Evidence:
- Laravel: `Invoice::withCount('items')` generates 2 queries instead of N+1.
- Django: `Blog.objects.annotate(num_entries=Count('entry'))` adds count as annotation.
- SQLAlchemy: `func.count()` in query projection avoids separate queries.
- Rails: `Author.includes(:books).count` with efficient aggregation.

Source: Laravel Eloquent Relationships — Aggregating Related Models
URL: https://laravel.com/docs/13.x/eloquent-relationships
Confidence: HIGH

Corroborated By: Django aggregation guide, SQLAlchemy func documentation

---

## Evidence 7

Claim: Pagination prevents N+1 from scaling catastrophically with dataset size.

Evidence:
- Rails `find_each` retrieves records in batches (default 1,000).
- Laravel `chunk()` processes records in batches.
- Django `iterator()` avoids caching large result sets.
- SQLAlchemy select-in loading: "emits a SELECT for up to 500 parent primary key values at a time."
- Without pagination: 10,000 records + N+1 = 10,001 queries. With pagination (50/page), each request limited.

Source: Rails Active Record Query Interface — Retrieving Multiple Records in Batches
URL: https://guides.rubyonrails.org/active_record_querying.html
Confidence: HIGH

Corroborated By: Django iterator(), Laravel chunk(), SQLAlchemy batch size

---

## Evidence 8

Claim: ORMs provide tools to detect N+1: strict loading, logging, query counting.

Evidence:
- Rails `strict_loading` raises error when lazy loading occurs.
- Django `connection.queries` shows executed SQL.
- Laravel with Debugbar or `preventLazyLoading()`.
- EF Core logging tracks query count.
- SQLAlchemy `echo=True` or event listeners.

Source: Rails Guide — Eager Loading Associations
URL: https://guides.rubyonrails.org/active_record_querying.html
Confidence: HIGH

Corroborated By: Django, Laravel, EF Core, SQLAlchemy documentation

---

## Evidence 9

Claim: The N+1 concept extends beyond database to API/microservice calls (network N+1).

Evidence:
- Shopify Engineering: "The n+1 problem means that the server executes multiple unnecessary round trips to datastores for nested data."
- DataLoader pattern batches resolver calls in GraphQL.
- Lab specification: "GET /orders returns 100 order. Backend does GET Payment API /payment/1 through /payment/100."
- Solution: bulk endpoints like `POST /payments/batch` with `{ "order_ids": [...] }`.

Source: Shopify Engineering — Solving the N+1 Problem for GraphQL through Batching
URL: https://shopify.engineering/solving-the-n-1-problem-for-graphql-through-batching
Confidence: HIGH

Corroborated By: Lab specification network N+1 example

---

## Evidence 10

Claim: Column selection (`select()`, `pluck()`, `values()`, `only()`, `defer()`) reduces data transfer and is often more efficient than eager loading when only specific fields are needed.

Evidence:
- Rails `pluck()` returns array of values without constructing model instances.
- Django `values()` returns dictionaries instead of model instances.
- Django `only()` and `defer()` control which fields are loaded.
- Laravel `select()` limits columns.
- SQLAlchemy `load_only()` loads only specific columns.

Source: Django ORM optimization guide
URL: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
Confidence: HIGH

Corroborated By: Rails pluck(), Laravel query builder, SQLAlchemy column loading options

---

## Evidence 11

Claim: Recommended troubleshooting workflow: measure first, identify hotspots, optimize queries before scaling infrastructure.

Evidence:
- Django docs: "Profile first... Find out what queries you are doing and what they are costing you."
- Lab specification workflow: "Berapa query → Query paling mahal → Ada N+1 → Index benar → Data terlalu banyak → Pagination → Baru caching."
- All ORM docs emphasize profiling before optimization.

Source: Django optimization guide — "Profile first" section
URL: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
Confidence: HIGH

Corroborated By: Lab specification troubleshooting workflow, general performance engineering consensus

---

## Evidence 12

Claim: Common mistakes include: over-eager loading, testing with small datasets only, SELECT *, and ignoring pagination.

Evidence:
- Lab specification: "Testing hanya menggunakan 10 record. N+1 sering tidak terlihat sampai dataset cukup besar."
- Mihalcea: "Using FetchType.EAGER... you are going to fetch way more data that you need."
- Django: "Don't overuse contains(), count(), and exists()... if you are going to need other data from the QuerySet, evaluate it immediately."
- Lab: "Pagination dianggap opsional. Endpoint list yang bisa mengembalikan seluruh tabel adalah bom waktu."

Source: Lab specification — "Kesalahan Umum" section
URL: labs/22-n-plus-one-query-problem
Confidence: MEDIUM

Corroborated By: Django ORM docs, Vlad Mihalcea, general ORM best practices
