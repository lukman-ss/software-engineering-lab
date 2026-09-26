# Claim Audit

## Claim 1
Claim: Fetching N related objects via lazy loading triggers N+1 database queries (1 main query + N child queries).
Location: `research/05-report.md:Finding 1`, `research/03-evidence.md:Evidence 1`
Evidence Provided: Documentation quotes from Rails Guides, Django Docs, Laravel Eloquent, SQLAlchemy.
Source: Sources 1, 2, 4, 5
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Authoritative universal consensus across all ORM implementations.

---

## Claim 2
Claim: Eager loading eliminates per-object lazy queries, reducing query count to O(1) or O(k) queries.
Location: `research/05-report.md:Finding 2`, `research/03-evidence.md:Evidence 2`
Evidence Provided: `includes()` in Rails, `prefetch_related()` in Django, `with()` in Laravel, `selectinload()` in SQLAlchemy, `Include()` in EF Core.
Source: Sources 1, 2, 3a, 4, 5
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Confirmed across all primary ORM documentation sources.

---

## Claim 3
Claim: Eager loading has trade-offs including memory overhead, over-fetching, and cartesian product explosion.
Location: `research/05-report.md:Finding 3`, `research/03-evidence.md:Evidence 3`
Evidence Provided: SQLAlchemy joinedload warnings, Django prefetch reverse/m2m multiplier notes, EF Core split query recommendations.
Source: Sources 3a, 5, 6
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Well documented in official ORM performance caveats.

---

## Claim 4
Claim: Column selection (`select`, `pluck`, `values`, `defer`, `only`) and aggregation (`withCount`) replace relationship loading when scalar values or subset fields are required.
Location: `research/05-report.md:Finding 4`, `research/03-evidence.md:Evidence 4, Evidence 7`
Evidence Provided: Django `values`/`annotate`, Rails `pluck`, Laravel `withCount`.
Source: Sources 1, 2, 4, 6
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Supported by primary sources as optimal alternative to full entity hydration.

---

## Claim 5
Claim: Lazy loading is the default loading strategy across major ORMs and is the root cause of N+1 problems in production.
Location: `research/05-report.md:Finding 5`, `research/03-evidence.md:Evidence 12`
Evidence Provided: Direct quotes from SQLAlchemy, Django, Rails, EF Core, Laravel.
Source: Sources 1, 2, 3, 4, 5
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Core architectural design choice across mainstream ORMs.

---

## Claim 6
Claim: N+1 cannot be reliably detected by code review or small test datasets; it requires profiling tools, query counting, or strict loading modes.
Location: `research/05-report.md:Finding 6`, `research/03-evidence.md:Evidence 6`
Evidence Provided: Rails `strict_loading`, SQLAlchemy `raiseload`, Django `connection.queries`/`django-debug-toolbar`.
Source: Sources 1, 4, 5
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Explicitly advised in ORM guides.

---

## Claim 7
Claim: The recommended troubleshooting workflow follows "measure first → profile → identify N+1 → fix root cause → scale infrastructure only if needed".
Location: `research/05-report.md:Finding 7`, `research/03-evidence.md:Evidence 9`
Evidence Provided: Django "Profile first" doc section, aligned with standard performance engineering steps.
Source: Source 1, Source 7
Source Actually Supports Claim: YES
Classification: INTERPRETATION
Severity: LOW
Notes: Methodological guidance matching authoritative advice.

---

## Claim 8
Claim: N+1 pattern extends beyond databases to network/microservice API calls and GraphQL resolvers.
Location: `research/05-report.md:Finding 8`, `research/03-evidence.md:Evidence 8`
Evidence Provided: Specification analogy, GraphQL DataLoader concept reference.
Source: Source 7, External DataLoader documentation reference
Source Actually Supports Claim: PARTIAL
Classification: EXAMPLE
Severity: MEDIUM
Notes: Analogy is conceptually valid, but evidence relies partly on internal lab spec rather than primary microservices benchmarking paper. Research report correctly flags confidence as MEDIUM.

---

## Claim 9
Claim: "1 request = 712 queries = 2.4 seconds, target 180ms" represents realistic production monitoring signals.
Location: `research/05-report.md:Executive Summary & Limitations`, `research/06-open-questions.md:Weak Evidence 1`
Evidence Provided: Lab topic specification.
Source: Source 7
Source Actually Supports Claim: PARTIAL
Classification: IMPLEMENTATION-SPECIFIC
Severity: MEDIUM
Notes: Research correctly notes that numeric values are scenario-specific illustrative figures, not universal benchmarks.
