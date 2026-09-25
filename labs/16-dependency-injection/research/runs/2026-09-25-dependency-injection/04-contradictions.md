# Contradictions and Divergences

## Disagreement 1: Service Locator Viability vs Anti-Pattern

### SOURCE A
Martin Fowler (2004) — "Inversion of Control Containers and the Dependency Injection pattern"
Claim: Service Locator and Dependency Injection are roughly equivalent for application-internal classes. Fowler states: "When building application classes the two are roughly equivalent, but I think Service Locator has a slight edge due to its more straightforward behavior... I don't see the injector's inversion as providing anything compelling [for simple app code]."

### SOURCE B
PHP-FIG PSR-11 Meta Document & Modern Architectural Consensus
Claim: Service Locator is explicitly discouraged for business objects. Passing the container into an object obscures dependencies, couples the class to the container interface, and harms interoperability and unit testability.

### ASSESSMENT
The consensus evolved significantly from 2004 to present. Early in container adoption, runtime reflection and complex XML configurations made DI feel opaque ("magic"), leading authors like Fowler to see Service Locator as more intuitive. However, modern IDE tooling, static typing, static analysis, autowiring, and clean testing paradigms proved that hidden dependencies in Service Locators cause severe long-term maintenance degradation. Modern software architecture uniformly classifies Service Locator inside domain/business services as an anti-pattern.

---

## Disagreement 2: Constructor Injection vs Setter Injection

### SOURCE A
Early Spring Framework Community (circa 2003-2005)
Claim: Favored Setter Injection to avoid massive constructors and allow flexible optional configuration.

### SOURCE B
Modern Frameworks (.NET, modern Spring 5/6/7, modern PHP)
Claim: Favors Constructor Injection as the primary mechanism because it enforces immutability, guarantees complete object initialization, and naturally reveals violating the Single Responsibility Principle when parameter counts grow too high.

### ASSESSMENT
Constructor injection has won the debate for mandatory dependencies. Setter injection is now reserved almost exclusively for truly optional dependencies or resolving rare circular references.
