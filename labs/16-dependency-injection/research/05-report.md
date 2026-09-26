# Research Report

## Research Question
What are the principles, patterns, benefits, trade-offs, and best practices of Dependency Injection (DI) and Inversion of Control (IoC) for writing testable, evolvable software?

## Executive Summary
Dependency Injection is a specific application of the broader Inversion of Control principle for service assembly. DI enables loose coupling by inverting dependency creation: objects declare dependencies (via constructor, setter, or interface), and a container supplies implementations at runtime. All authoritative sources (Fowler 2004, Spring 7.x, Microsoft .NET 2026, Laravel 12.x, PSR-11) agree on the core mechanics and benefits: decoupling, testability via mock/stub substitution, and explicit separation of configuration from use. The Service Locator pattern is an alternative that achieves similar decoupling but is explicitly discouraged by PSR-11 and considered inferior for components consumed by external applications. Modern practice converges on constructor injection as primary, setter injection as fallback, and has largely abandoned interface injection. Service lifetimes (singleton, scoped/request, transient/prototype) are universally supported. The main anti-patterns are: using the container as a Service Locator (passing container into classes), over-injecting dependencies (signaling SRP violations), and injecting value objects that have no substitution need.

## Findings

### Finding 1 — DI Definition and Core Mechanics
**Claim**: DI is a pattern where objects define dependencies only through constructor arguments, factory method arguments, or properties set after construction; the container injects dependencies when creating the bean. The object does not look up its dependencies or know their location/class.
**Evidence**: Spring Framework 7.0.9 explicitly defines DI this way: "Dependency injection (DI) is a process whereby objects define their dependencies... only through constructor arguments, arguments to a factory method, or properties that are set on the object instance after it is constructed... The container then injects those dependencies when it creates the bean."
**Sources**: Spring (Source 5), Fowler (Source 1), Microsoft .NET (Source 3)
**Confidence**: HIGH

### Finding 2 — Three Forms of DI
**Claim**: The three recognized forms are Constructor Injection, Setter Injection, and Interface Injection. Modern frameworks primarily use Constructor and Setter; Interface Injection has fallen out of practice.
**Evidence**: Fowler (2004) names all three with code examples. Spring docs 7.0.9 document "Constructor-based dependency injection" and "Setter-based dependency injection" as "two major variants" without mentioning interface injection. Microsoft .NET focuses on constructor injection.
**Sources**: Fowler (Source 1), Spring (Source 5), Microsoft (Source 3)
**Confidence**: HIGH

### Finding 3 — DI vs Service Locator
**Claim**: Both achieve decoupling from concrete implementations. DI avoids a dependency on the locator; Service Locator requires every consumer to depend on the locator API. For components used by external applications, DI is preferred. PSR-11 explicitly discourages Service Locator ("Users SHOULD NOT pass a container into an object...").
**Evidence**: Fowler (2004) provides detailed comparison; PSR-11 (Source 6) mandates "SHOULD NOT" per RFC 2119.
**Sources**: Fowler (Source 1), PSR-11 (Source 6)
**Confidence**: HIGH

### Finding 4 — DI Improves Testability
**Claim**: DI enables unit testing by allowing replacement of real implementations with stubs/mocks without modifying the class under test.
**Evidence**: Spring: "classes become easier to test... allow for stub or mock implementations." Microsoft .NET: "the app should use a mock or stub... which isn't possible with this approach" (hard-coded deps). Fowler: "both [DI and SL] are very amenable to stubbing" — testing benefit is shared but DI is simpler default.
**Sources**: Spring (Source 5), Microsoft (Source 3), Fowler (Source 1)
**Confidence**: HIGH

### Finding 5 — Constructor vs Setter Injection Trade-offs
**Claim**: Prefer constructor injection for valid-object-at-construction-time and immutable fields; switch to setter injection when: multiple valid construction configurations exist, many constructor parameters (no keyword args in Java historically), primitive/string parameters needing named disambiguation, or inheritance creates constructor forwarding complexity.
**Evidence**: Fowler (2004) provides detailed analysis. Modern languages (PHP 8, C#, Kotlin) have named constructor args, reducing some setter advantages.
**Sources**: Fowler (Source 1), Microsoft (Source 3 — constructor primary)
**Confidence**: HIGH

### Finding 6 — Service Lifetimes (Scopes)
**Claim**: All major containers support at least three lifetimes: Singleton (app-wide), Scoped/Request (per lifecycle boundary), Transient/Prototype (new instance per resolve).
**Evidence**: .NET: `AddSingleton`, scoped per request. Laravel: `singleton()`, `scoped()` per request/job. Spring: bean scopes singleton, prototype, request, session.
**Sources**: Microsoft (Source 3), Laravel (Source 4), Spring (Source 5)
**Confidence**: HIGH

### Finding 7 — IoC is Broader than DI
**Claim**: Inversion of Control is a general framework principle (framework calls user code: UI events, template methods, EJB lifecycle). DI is one specific form of IoC used by containers for service assembly. "IoC Container" is a misnomer conflating the two.
**Evidence**: Fowler (2005): "Inversion of Control is a common phenomenon... the specific styles of inversion of control (such as dependency injection) that these containers use."
**Sources**: Fowler (Source 2)
**Confidence**: HIGH

### Finding 8 — PSR-11 Standard
**Claim**: PSR-11 standardizes ContainerInterface with `get($id)` and `has($id)`. A non-existent id MUST throw `NotFoundExceptionInterface`. The standard explicitly discourages Service Locator usage.
**Evidence**: PSR-11 document (Source 6) ratified by PHP-FIG.
**Sources**: PHP-FIG (Source 6)
**Confidence**: HIGH

### Finding 9 — Anti-Pattern: Service Locator Overuse
**Claim**: Passing container into classes so they call `container.get()` or `app()->make()` hides dependencies, makes them unclear from signatures, and violates separation of configuration from use.
**Evidence**: PSR-11 "SHOULD NOT"; Fowler "separation of configuration from use"; Topic spec "❌ Menggunakan Service Locator di Mana-mana".
**Sources**: PSR-11 (Source 6), Fowler (Source 1), Topic Spec
**Confidence**: HIGH

### Finding 10 — Anti-Pattern: Constructor Over-injection
**Claim**: Excessive constructor parameters (e.g., 12+) signal design problems (too many responsibilities, SRP violation).
**Evidence**: Topic spec explicitly states "Kalau constructor berisi 12 parameter... Biasanya ada masalah desain." Fowler: "If you have a lot of constructor parameters things can look messy... often a sign of an over-busy object that should be split."
**Sources**: Topic Spec, Fowler (Source 1)
**Confidence**: MEDIUM (specific number "12" from topic spec, not independently verified)

### Finding 11 — When NOT to Use DI
**Claim**: Value objects without external dependencies or substitution needs (DateTime, Money, Address) can be created directly. DI is appropriate for services with external dependencies: Database, HTTP Client, Payment Gateway, Cache, Queue, Email, Storage, cross-module services.
**Evidence**: Topic spec provides this guidance explicitly. Fowler discusses principle (separation of configuration from use) but does not give a concrete checklist.
**Sources**: Topic Spec
**Confidence**: MEDIUM (guidance from spec, not independently verified primary source)

## Areas of Agreement
All authoritative sources agree on:
1. DI definition: container injects dependencies declared by object
2. Three DI forms historically; constructor + setter are the practical pair
3. DI vs Service Locator distinction (locator dependency vs no locator dependency)
4. DI enables testability via mock/stub substitution
5. Service lifetimes (singleton, scoped, transient) exist in all containers
6. PSR-11 explicitly discourages Service Locator
7. IoC is broader; DI is specific IoC for service assembly

## Areas of Disagreement / Divergence
- **Testing benefit exclusivity**: Fowler states Service Locator is equally amenable to stubbing if well-designed; Spring/Microsoft present DI as the solution to testing difficulties. (Divergence in emphasis, not fact.)
- **Interface Injection**: Fowler includes as third form; modern framework docs omit entirely. (Evolution, not contradiction.)
- **Configuration mechanism**: Fowler (2004) advocates programmatic builders over XML; modern frameworks use annotations/attributes that blur the line.
- **"12 parameter" threshold**: Topic spec cites specific number; Fowler says "a lot" qualitatively. (Heuristic vs principle.)

## Limitations
- Primary sources are documentation and Fowler's articles; no empirical studies on DI impact on defect rates, velocity, or maintainability were found in this research.
- NestJS documentation fetch failed to retrieve code examples; TypeScript ecosystem evidence is weaker.
- The topic specification's practical heuristics (value objects list, 12-parameter threshold) are not independently cross-checked against primary sources.
- DI performance overhead (container resolution cost) not quantified in sources reviewed.

## Conclusion
Dependency Injection is a mature, cross-platform pattern for achieving loose coupling and testability. The consensus across .NET, Java Spring, PHP Laravel, and PHP-FIG standards is:
1. Declare dependencies on abstractions (interfaces) via constructor injection as default.
2. Use an IoC container for automatic resolution and lifetime management.
3. Avoid Service Locator (passing container into classes).
4. Separate configuration (bindings) from use (consumer classes).
5. Watch for constructor over-injection as a design smell.

The topic specification's practical exercises (PPOB system, identifying injectable dependencies, creating interfaces, mocking for tests) directly align with these findings. The remaining open questions concern empirical validation and language-specific nuances.