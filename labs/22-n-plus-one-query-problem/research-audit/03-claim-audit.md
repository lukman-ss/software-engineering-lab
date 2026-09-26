# Claim Audit

**Target Lab:** `labs/22-n-plus-one-query-problem`  
**Audit Scope:** Major claims in `research/05-report.md` and `research/03-evidence.md`  
**Audit Date:** 2026-09-26  

---

## Claim 1

Claim: Fetching $N$ related objects via lazy loading triggers $N+1$ queries: 1 for the parent collection + $N$ for each child relationship access.

Location: `research/05-report.md:17-34`, `research/03-evidence.md:3-20`

Evidence Provided: Rails Section 16.1 example (10 books -> 11 queries), Django docs on lazy relationship access, Laravel Eloquent relationship docs, SQLAlchemy "N plus one problem" docs.

Source: Rails Guide, Django Optimization Docs, Laravel Eloquent Docs, SQLAlchemy Docs.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Universal ORM behavior across all surveyed ecosystems.

---

## Claim 2

Claim: Eager loading eliminates per-object lazy loads by fetching all related data in a batch query, reducing query count to $O(k)$ where $k$ is the number of relationship levels.

Location: `research/05-report.md:35-56`, `research/03-evidence.md:23-42`

Evidence Provided: Rails `includes()` generates 2 queries (`WHERE id IN (...)`), Django `prefetch_related()`, Laravel `with()`, SQLAlchemy `selectinload()`, EF Core `Include()`.

Source: EF Core docs, Django QuerySet API, SQLAlchemy docs, Rails guides.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Thoroughly substantiated with direct syntax and query output semantics.

---

## Claim 3

Claim: Eager loading has trade-offs including memory overhead and cartesian product explosion when joining multiple collections.

Location: `research/05-report.md:57-73`, `research/03-evidence.md:44-62`

Evidence Provided: SQLAlchemy warning on `joinedload()` collection multiplication and Result.unique(), Django reverse relation multiplication warning, EF Core split queries recommendation.

Source: SQLAlchemy docs, Django docs, EF Core docs.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Accurately identifies that eager loading is not a cost-free silver bullet.

---

## Claim 4

Claim: Column projection (`select`, `pluck`, `values`, `only`, `defer`) and aggregation (`withCount`, `annotate(Count)`) eliminate N+1 without loading full relationship objects when only scalar values/subsets are needed.

Location: `research/05-report.md:74-95`, `research/03-evidence.md:64-83`

Evidence Provided: Rails `pluck()`, Django `values()` and `annotate(Count())`, Laravel `withCount()`, SQLAlchemy `load_only()`.

Source: Django optimization docs, Rails guide, Laravel docs, SQLAlchemy docs.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Supported across all frameworks.

---

## Claim 5

Claim: Lazy loading is the default in all major ORMs, causing N+1 to be invisible in development with small datasets but problematic in production.

Location: `research/05-report.md:96-116`, `research/03-evidence.md:103-125`

Evidence Provided: SQLAlchemy default loading style declaration, Django lazy query evaluation, Rails lazy association default, EF Core default behavior.

Source: SQLAlchemy, Django, Rails, EF Core docs.

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Validated against official architectural descriptions.

---

## Claim 6

Claim: N+1 detection requires profiling tools (query logging, debug toolbars, strict loading / raiseload modes) and a profile-first workflow.

Location: `research/05-report.md:117-152`, `research/03-evidence.md:128-175`

Evidence Provided: Django "Profile first" guide & `connection.queries`, Rails `strict_loading`, SQLAlchemy `raiseload()`, Laravel Debugbar.

Source: Django optimization guide, Rails guide, SQLAlchemy docs.

Source Actually Supports Claim: YES

Classification: FACT / INTERPRETATION

Severity: LOW

Notes: "Profile first" is an explicit recommendation in Django official documentation.

---

## Claim 7

Claim: The N+1 pattern extends to network and microservice API calls and GraphQL resolvers.

Location: `research/05-report.md:153-170`, `research/03-evidence.md:178-198`

Evidence Provided: Conceptual equivalence of $N$ HTTP `GET` calls vs 1 bulk endpoint; GraphQL DataLoader pattern.

Source: Lab spec / Architectural literature.

Source Actually Supports Claim: YES

Classification: INTERPRETATION

Severity: LOW

Notes: Research explicitly notes this as an architectural generalization rather than a framework-specific feature.

---

## Claim 8

Claim: Specific performance numbers (e.g. 712 queries = 2.4s, 180ms targets) and connection pool cascade thresholds are universal benchmarks.

Location: `research/05-report.md:195-202`, `research/05-report.md:216-217`

Evidence Provided: Explicitly disclaimed in Limitations section as scenario-specific illustrative data, not universal facts.

Source: Source 7 (Topic Specification).

Source Actually Supports Claim: YES (Properly contextualized)

Classification: IMPLEMENTATION-SPECIFIC

Severity: LOW

Notes: The research report correctly refused to overgeneralize these numbers and treated them strictly as illustrative scenario examples.
