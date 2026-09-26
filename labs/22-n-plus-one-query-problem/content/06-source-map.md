# Source Map

## N+1 Query Problem Definition

Research:
- `research/runs/2026-09-25-n-plus-one-query-problem/05-report.md` (Finding 1: Definition consistent across Django, Rails, Laravel, SQLAlchemy, EF Core)
- `research-audit/07-verdict.md` — APPROVED

Implementation:
- `internal/blog/repository.go:13` — `GetAuthorsWithPostsNPlusOne()` naif method

Tests:
- `internal/blog/repository_test.go:8` — `TestGetAuthorsWithPostsNPlusOne` asserting query count == 4

---

## Eager Loading Solution (Batch Loading)

Research:
- `research/runs/2026-09-25-n-plus-one-query-problem/05-report.md` (Finding 2: Eager loading reduces N+1 to O(k))
- `research/runs/2026-09-25-n-plus-one-query-problem/05-report.md` (Finding 3: Trade-offs — cartesian explosion, memory overhead)

Implementation:
- `internal/blog/repository.go:32` — `GetAuthorsWithPostsEager()` batch loading method
- `internal/blog/store.go:63` — `GetPostsByAuthorIDs()` batch query implementation

Tests:
- `internal/blog/repository_test.go:27` — `TestGetAuthorsWithPostsEager` asserting query count == 2
- `internal/blog/repository_test.go:46-49` — Deep equivalence check (`reflect.DeepEqual`) ensuring data consistency

---

## Mock Data Store with Query Counter

Research:
- `research/runs/2026-09-25-n-plus-one-query-problem/05-report.md` (Finding 6: Detection requires tooling — query counting and profiling)

Implementation:
- `internal/blog/store.go:5-11` — Store struct with queryCount field
- `internal/blog/store.go:30-40` — `GetQueryCount()` and `ResetQueryCount()` methods
- `internal/blog/store.go:42-79` — Query methods that increment counter

---

## Lazy Loading as Default (Root Cause)

Research:
- `research/runs/2026-09-25-n-plus-one-query-problem/05-report.md` (Finding 5: Lazy loading is default in all major ORMs)

---

## Profiling-First Workflow

Research:
- `research/runs/2026-09-25-n-plus-one-query-problem/05-report.md` (Finding 7: Profile first → identify → eliminate unnecessary queries → scale infrastructure)

Engineering:
- `engineering/01-design.md:31-32` — Test strategy for query verification
- `engineering/02-implementation-notes.md:16-21` — Known limitations (lines 16-18) and trade-offs: deterministic query counting vs real I/O benchmarking (line 20-21)

---

## Column Selection and Aggregation

Research:
- `research/runs/2026-09-25-n-plus-one-query-problem/05-report.md` (Finding 4: `select()`, `pluck()`, `values()`, `withCount()` eliminate N+1 without full eager loading)

---

## Network/API N+1 Analogy

Research:
- `research/runs/2026-09-25-n-plus-one-query-problem/05-report.md` (Finding 8: N+1 extends to network/API calls — Confidence: MEDIUM)
- `research-audit/03-claim-audit.md` — Claims verified against ORM primary sources

---

## Mapping-Level vs Query-Level Eager Loading

Engineering:
- `engineering/02-implementation-notes.md:17-18` — Known limitations: memory bloat OoM not demonstrated, ORM mapping anomalies not shown
- `engineering-audit/06-verdict.md` — APPROVED, no warnings

---

## Verified Demo Output

Engineering:
- `engineering/03-execution-result.md:19-29` — Test results (all 3 tests PASS)
- `engineering/03-execution-result.md:31-40` — Race detector results (PASS)
- `engineering/03-execution-result.md:42-56` — Demo output showing 4 queries (N+1) vs 2 queries (eager)

---

## External References (Research Sources)

Research:
- `research/runs/2026-09-25-n-plus-one-query-problem/02-sources.md` — 7 sources + 1 internal spec
- `research-audit/02-source-audit.md` — Source integrity PASS

Documented at:
- Django Optimization Guide: https://docs.djangoproject.com/en/5.1/topics/db/optimization/
- Rails Active Record Query Interface: https://guides.rubyonrails.org/active_record_querying.html
- Laravel Eloquent Relationships: https://laravel.com/docs/13.x/eloquent-relationships
- SQLAlchemy Relationship Loading: https://docs.sqlalchemy.org/en/21/orm/queryguide/relationships.html
- EF Core Eager Loading: https://learn.microsoft.com/en-us/ef/core/querying/related-data/eager
