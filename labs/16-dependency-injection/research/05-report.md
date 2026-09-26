# Research Report

## Research Question
What are the principles, patterns, benefits, trade-offs, and practices of Dependency Injection (DI) and Inversion of Control (IoC) for writing testable, evolvable software? How do these apply to the topic specification (PaymentService/PPOB/NotificationService examples, interface-based design, container bindings, testing with mocks)?

## Executive Summary
DI is a specialized form of IoC for service assembly: objects declare dependencies on abstractions (constructor/setter), and an external injector/container supplies implementations at runtime. Consensus across Fowler 2004, Spring 7.x, Microsoft .NET 2026, Laravel 12.x/13.x, PSR-11, and PHP interfaces manual: DI reduces coupling, makes configuration swappable per deployment, and enables mock/stub testing. Service Locator achieves similar decoupling but creates a locator dependency in every consumer and is discouraged by PSR-11. Modern practice: constructor injection default, setter fallback, interface injection obsolete. Containers manage lifetimes (singleton/scoped/transient). Anti-patterns: container-as-locator, over-injection (design smell), injecting value objects unnecessarily.

## Findings

### Finding 1
Claim: DI separates construction from use; object receives dependencies from external injector instead of creating them; yields loose coupling.
Evidence: Wikipedia DI definition verbatim; Spring "objects define dependencies only through constructor/factory args or properties; container injects at bean creation"; .NET three-step (abstract, register, inject).
Sources: https://en.wikipedia.org/wiki/Dependency_injection, https://docs.spring.io/spring-framework/reference/core/beans/introduction.html, https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection, https://martinfowler.com/articles/injection.html
Confidence: HIGH

### Finding 2
Claim: IoC is broader (framework calls user code; Hollywood Principle); DI is one IoC form for dependency implementations. "IoC container" conflates the two senses.
Evidence: Wikipedia IoC; Fowler settled on "Dependency Injection" name because IoC too generic (UI main loop vs plugin lookup inversion).
Sources: https://en.wikipedia.org/wiki/Inversion_of_control, https://martinfowler.com/articles/injection.html
Confidence: HIGH

### Finding 3
Claim: Three historical DI forms: constructor, setter, interface (type 3/2/1). Modern frameworks support constructor + setter; interface injection obsolete.
Evidence: Fowler full taxonomy with PicoContainer/Spring/Avalon examples; Spring lists only constructor-based and setter-based as two major variants; .NET/Laravel focus constructor.
Sources: https://martinfowler.com/articles/injection.html, https://docs.spring.io/spring-framework/reference/core/beans/introduction.html
Confidence: HIGH

### Finding 4
Claim: Prefer constructor injection (valid object at birth, immutable fields); switch to setter when many params, multiple valid combos, string params needing names, inheritance explosion.
Evidence: Fowler detailed constructor-vs-setter analysis with Kent Beck reference; .NET constructor primary with selection rules.
Sources: https://martinfowler.com/articles/injection.html, https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection
Confidence: HIGH

### Finding 5
Claim: Both DI and Service Locator decouple from concrete impl; DI avoids locator dependency; DI preferred for externally-consumed components; locator OK for app-internal code with known locator API.
Evidence: Fowler comparison + segregated-interface/dynamic-locator variants; PSR-11 SHOULD NOT pass container into objects.
Sources: https://martinfowler.com/articles/injection.html, https://www.php-fig.org/psr/psr-11/
Confidence: HIGH

### Finding 6
Claim: DI enables mock/stub unit testing in isolation without real external calls; often first benefit noticed. (Locator equally stub-able if well-designed, per Fowler.)
Evidence: Wikipedia testing section; .NET "use mock or stub... isn't possible with hard-coded"; Laravel "easily mock when testing"; Fowler nuance.
Sources: https://en.wikipedia.org/wiki/Dependency_injection, https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection, https://laravel.com/docs/container
Confidence: HIGH

### Finding 7
Claim: Containers support singleton (app-wide), scoped/request (per lifecycle), transient/prototype (per resolve). Misuse (scoped-from-root, scoped-into-singleton) caught by validation in dev.
Evidence: .NET AddSingleton/AddScoped/AddTransient + scope validation; Laravel singleton()/scoped()/instance(); Spring bean scopes.
Sources: https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection, https://laravel.com/docs/container, https://docs.spring.io/spring-framework/reference/core/beans/introduction.html
Confidence: HIGH

### Finding 8
Claim: Program to interfaces: PHP interfaces define contracts; multiple impls interchangeable (payment gateways, DB, cache); consumer unchanged on swap.
Evidence: PHP manual explicitly lists "multiple payment gateways... swapped without changes to code that uses them"; Laravel bind(Interface, Impl); Fowler MovieFinder interface + plugin.
Sources: https://www.php.net/manual/en/language.oop5.interfaces.php, https://laravel.com/docs/container, https://martinfowler.com/articles/injection.html
Confidence: HIGH

### Finding 9
Claim: PSR-11 standardizes get/has; unknown id MUST throw NotFoundException; passing container into objects (Service Locator) discouraged.
Evidence: PSR-11 sections 1.1.2/1.2/1.3 verbatim; Laravel PSR-11 compliance section.
Sources: https://www.php-fig.org/psr/psr-11/, https://laravel.com/docs/container
Confidence: HIGH

### Finding 10
Claim: DI costs: config burden, harder tracing, reflection hurts IDE, upfront effort, framework lock-in risk. IoC hard to understand/debug; justify over simpler alternative.
Evidence: Wikipedia disadvantages list; Fowler "inversion comes at a price... hard to understand... leads to problems debugging... prefer to avoid unless needed."
Sources: https://en.wikipedia.org/wiki/Dependency_injection, https://martinfowler.com/articles/injection.html
Confidence: MEDIUM (criticisms listed, single-source each)

### Finding 11
Claim: Over-injection signals SRP violation; industry guidance is qualitative (Fowler: "lot of params... sign of over-busy object"); the 12-parameter threshold is a lab-specific heuristic, not an industry standard.
Evidence: Fowler "lot of params... sign of over-busy object that should be split"; topic spec 12-param rule.
Sources: https://martinfowler.com/articles/injection.html + topic spec (Lab Requirement/Heuristic)
Confidence: MEDIUM (number 12 not in primary source)

### Finding 12
Claim: Don't inject value objects (stateful domain data); inject external-boundary services (DB/HTTP/gateway/cache/queue/email/storage). The specific examples (DateTime/Money/Address) are lab heuristics, not a universal standard.
Evidence: Topic spec list; Fowler distinction between entities and value objects (no concrete list); Wikipedia new-keyword diminished except value objects.
Sources: Topic spec (Lab Requirement/Heuristic); https://martinfowler.com/articles/injection.html (principle)
Confidence: MEDIUM (list from spec, principle corroborated)

## Areas of Agreement
1. DI definition and construction/use separation
2. Constructor + setter as practical pair; interface injection historical
3. DI vs Locator distinction (locator dependency)
4. Testability via substitution
5. Lifetimes singleton/scoped/transient
6. PSR-11 discourages locator
7. IoC broader; DI specific

## Areas of Disagreement
- Testing exclusivity (Fowler: locator equally stub-able; vendors present DI as fix)
- Interface injection (Fowler includes; modern docs omit — evolution)
- Config mechanism (Fowler programmatic preference vs modern annotation/attribute hybrids)
- 12-param threshold (spec heuristic vs Fowler qualitative)

## Limitations
- No empirical studies on defect/velocity impact found — NOT VERIFIED
- Container performance overhead not quantified — NOT VERIFIED
- NestJS/TS evidence weaker (fetch failure)
- Spec heuristics (value-object list, 12 params) not independently verified
- PHP interface property hooks (8.4) tangential, not DI-specific

## Conclusion
DI is mature cross-platform consensus: depend on abstractions via constructor, container resolves lifetimes, separate bindings from consumers, avoid locator, watch over-injection. Topic spec examples (PaymentService→gateway interface→container bind→mock test; PPOB provider swap; notification provider swap) align directly with findings. Open quantification and language-nuance questions remain.
