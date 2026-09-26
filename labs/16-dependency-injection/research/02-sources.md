# Research Sources

## Research Topic
Dependency Injection (DI) & Inversion of Control (IoC)

## Source List

### Source 1
Title: Inversion of Control Containers and the Dependency Injection pattern
Publisher: Martin Fowler (martinfowler.com)
URL: https://martinfowler.com/articles/injection.html
Published: 23 January 2004
Accessed: 26 September 2026
Source Tier: Tier 1 (Foundational Authority)
Relevance: Definitive article coining the "Dependency Injection" term; defines three forms (constructor, setter, interface injection), contrasts DI with Service Locator, discusses configuration vs. use separation.

### Source 2
Title: Inversion Of Control (Bliki)
Publisher: Martin Fowler (martinfowler.com)
URL: https://martinfowler.com/bliki/InversionOfControl.html
Published: 26 June 2005
Accessed: 26 September 2026
Source Tier: Tier 1 (Foundational Authority)
Relevance: Distinguishes Inversion of Control as general principle (frameworks) from the specific DI pattern used by IoC containers; provides etymology (Johnson & Foote 1988, Gang of Four, Hollywood Principle).

### Source 3
Title: Dependency injection - .NET
Publisher: Microsoft Learn
URL: https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection/overview
Published: 26 January 2026 (last updated 23 April 2026)
Accessed: 26 September 2026
Source Tier: Tier 1 (Official Documentation)
Relevance: Practical .NET implementation of DI via IServiceCollection/IServiceProvider; defines hard-coded dependency problems, constructor injection, service lifetimes (singleton/scoped/transient), scope validation, keyed services.

### Source 4
Title: Service Container
Publisher: Laravel (laravel.com)
URL: https://laravel.com/docs/12.x/container
Published: Laravel 12.x documentation
Accessed: 26 September 2026
Source Tier: Tier 1 (Official Documentation)
Relevance: Laravel's PSR-11 compliant IoC container; describes binding interfaces to implementations, contextual binding, tagging, zero-configuration resolution, singleton/scoped lifetimes, automatic injection.

### Source 5
Title: Dependency Injection
Publisher: Spring Framework Documentation (spring.io)
URL: https://docs.spring.io/spring-framework/reference/core/beans/dependencies/factory-collaborators.html
Published: Spring Framework 7.0.9
Accessed: 26 September 2026
Source Tier: Tier 1 (Official Documentation)
Relevance: Java ecosystem DI via Spring Container; constructor-based and setter-based DI, autowiring, bean scopes, factory-method injection; emphasizes object does not know location/class of dependencies.

### Source 6
Title: PSR-11: Container interface
Publisher: PHP Framework Interop Group (PHP-FIG)
URL: https://www.php-fig.org/psr/psr-11/
Published: PSR-11 Standard
Accessed: 26 September 2026
Source Tier: Tier 1 (Standards)
Relevance: Standardizes container interface (`get`, `has`); explicitly RECOMMENDS NOT using container as Service Locator ("Users SHOULD NOT pass a container into an object so that the object can retrieve its own dependencies").

### Source 7
Title: NestJS Documentation - Providers
Publisher: NestJS (nestjs.com)
URL: https://docs.nestjs.com/providers
Published: NestJS v11 documentation
Accessed: 26 September 2026
Source Tier: Tier 2 (Reputable Technical Publication)
Relevance: TypeScript/Node.js DI framework based on Angular DI; demonstrates providers, modules, scoped providers — shows cross-language DI patterns.