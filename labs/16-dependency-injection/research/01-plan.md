# 01 — Research Plan

## Research Topic
Dependency Injection (DI) & Inversion of Control (IoC) — principles, patterns, trade-offs, and practices for testable and evolvable software. Target lab: `labs/16-dependency-injection`.

## Objective
Investigate DI/IoC fundamentals, container mechanics, interface-based design, testing benefits, lifecycle management, and anti-patterns. Collect authoritative evidence, cross-check claims, and produce audit-ready report without implementation.

## Research Questions
1. What are DI and IoC definitions and their relationship?
2. What are the recognized forms of DI (constructor, setter, interface) and when to use each?
3. How do IoC containers work (binding, resolution, auto-wiring, lifetimes)?
4. Why does DI improve testability and evolvability (vendor swap, mock)?
5. What is Interface vs implementation principle and PSR-11 standard?
6. How does Service Locator differ from DI and why is it discouraged?
7. When NOT to use DI and what are over-injection / Service Locator smells?
8. What do Laravel, Spring, .NET containers demonstrate in practice?

## Search Strategy
- Tier 1 priority: Fowler canonical article, Spring Framework docs, Laravel Container docs, PSR-11 spec, Microsoft .NET DI docs, PHP interfaces manual
- Tier 2: Wikipedia for terminology baseline (cross-checked with Fowler/Spring)
- Fetch full source pages (not snippets), extract verbatim evidence with URLs
- Cross-check each significant claim with 2+ independent sources

## Expected Primary Sources
- martinfowler.com/articles/injection.html (Fowler 2004)
- martinfowler.com/bliki/InversionOfControl
- docs.spring.io/spring-framework/reference/core/beans/introduction.html
- laravel.com/docs/container
- php-fig.org/psr/psr-11
- learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection
- php.net/manual/en/language.oop5.interfaces.php
- en.wikipedia.org/wiki/Dependency_injection, en.wikipedia.org/wiki/Inversion_of_control

## Risks / Unknowns
- Empirical impact of DI on defect rate / velocity not in docs (likely NOT VERIFIED)
- Performance overhead of containers not quantified
- TypeScript/NestJS evidence weaker (fetch failure)
- 12-parameter heuristic from topic spec not in primary sources
- Value-object vs service boundary is heuristic, not formal standard
