# 03 — Evidence

## Evidence 1
Claim: DI is an object receiving dependencies from external code (injector) rather than creating them internally; separates construction from use, yields loose coupling.
Evidence: "Dependency injection is a programming technique in which an object receives other objects that it requires, as opposed to creating them internally... aims to separate the concerns of constructing objects and using them, leading to loosely coupled programs."
Source: Wikipedia — Dependency injection
URL: https://en.wikipedia.org/wiki/Dependency_injection
Confidence: HIGH
Corroborated By: Fowler 2004, Spring 7.0.9 docs
Notes: None

## Evidence 2
Claim: DI implements "inverting control over implementations of dependencies"; Java frameworks generically name this IoC (distinct from inversion of control flow).
Evidence: "Dependency injection implements the idea of 'inverting control over the implementations of dependencies', which is why certain Java frameworks generically name the concept 'inversion of control' (not to be confused with inversion of control flow)."
Source: Wikipedia — Dependency injection
URL: https://en.wikipedia.org/wiki/Dependency_injection
Confidence: HIGH
Corroborated By: Fowler (alternative meaning section), Spring docs
Notes: Distinguishes two IoC senses.

## Evidence 3
Claim: DI involves four roles: services, clients, interfaces, injectors (assembler/container/provider/factory). Injector must not be the client (avoids circular dependency).
Evidence: "Dependency injection involves four roles: services, clients, interfaces, and injectors... The injector... introduces services to the client... must not be the client, as this would create a circular dependency."
Source: Wikipedia — Dependency injection
URL: https://en.wikipedia.org/wiki/Dependency_injection
Confidence: HIGH
Corroborated By: Fowler injector/assembler concept
Notes: None

## Evidence 4
Claim: DI reduces coupling; improves reuse, testability, maintainability, flexibility; reduces boilerplate since one component handles creation; enables concurrent development.
Evidence: "A basic benefit of dependency injection is decreased coupling... By removing a client's knowledge of how its dependencies are implemented, programs become more reusable, testable and maintainable... dependency injection reduces boilerplate code, since all dependency creation is handled by a singular component."
Source: Wikipedia — Dependency injection
URL: https://en.wikipedia.org/wiki/Dependency_injection
Confidence: HIGH
Corroborated By: Spring, Microsoft .NET, Laravel docs
Notes: None

## Evidence 5
Claim: DI enables unit testing with stubs/mocks in isolation; often first benefit noticed.
Evidence: "This makes clients more independent and are easier to unit test in isolation, using stubs or mock objects, that simulate other objects not under test. This ease of testing is often the first benefit noticed when using dependency injection."
Source: Wikipedia — Dependency injection
URL: https://en.wikipedia.org/wiki/Dependency_injection
Confidence: HIGH
Corroborated By: Microsoft .NET ("use a mock or stub... isn't possible with [hard-coded] approach"), Spring, Laravel (mock service in tests)
Notes: None

## Evidence 6
Claim: DI has costs: demanding configuration, harder traceability (behavior vs construction separated), reflection-based implementations hinder IDE automation, more upfront effort, framework dependence risk.
Evidence: Critics argue DI "creates clients that demand configuration details... makes code difficult to trace... typically implemented with reflection hindering IDE automation... requires more upfront development effort... encourages dependence on a framework."
Source: Wikipedia — Dependency injection
URL: https://en.wikipedia.org/wiki/Dependency_injection
Confidence: MEDIUM
Corroborated By: Fowler (IoC hard to understand, debug problems)
Notes: Trade-off section.

## Evidence 7
Claim: Three main DI forms: Constructor injection, Setter injection, Interface injection (historically type 3/2/1 IoC).
Evidence: Fowler: "There are three main styles... Constructor Injection, Setter Injection, and Interface Injection... referred to as type 1 IoC (interface), type 2 (setter), type 3 (constructor)."
Source: Martin Fowler — IoC Containers and DI pattern
URL: https://martinfowler.com/articles/injection.html
Confidence: HIGH
Corroborated By: Wikipedia (constructor/method/setter/interface), Spring (constructor + setter as two major variants)
Notes: Modern docs omit interface injection.

## Evidence 8
Claim: Start with constructor injection; switch to setter when: many params, multiple valid construction combos, simple string params needing names, inheritance constructor explosion.
Evidence: Fowler: "My preference is to start with constructor injection, but be ready to switch to setter injection as soon as the problems I've outlined above start to become a problem." Constructor gives valid object at birth; immutable fields hidden via no setter.
Source: Martin Fowler — IoC Containers and DI pattern
URL: https://martinfowler.com/articles/injection.html
Confidence: HIGH
Corroborated By: Microsoft .NET (constructor primary), Spring (both supported)
Notes: PHP 8 named args / C# reduce some setter advantages.

## Evidence 9
Claim: Service Locator vs DI: both decouple from concrete impl; difference is how impl is provided (explicit request to locator vs appears via injection). Locator makes every consumer depend on locator.
Evidence: Fowler: "The key difference is that with a Service Locator every user of a service has a dependency to the locator... With injection there is no explicit request, the service appears in the application class - hence the inversion of control."
Source: Martin Fowler — IoC Containers and DI pattern
URL: https://martinfowler.com/articles/injection.html
Confidence: HIGH
Corroborated By: Wikipedia, PSR-11
Notes: None

## Evidence 10
Claim: DI better for components used in multiple external apps; Service Locator acceptable for app-internal classes where locator API is known.
Evidence: Fowler: "If you are building classes to be used in multiple applications then Dependency Injection is a better choice... building an application with various classes... using a service locator works quite well."
Source: Martin Fowler — IoC Containers and DI pattern
URL: https://martinfowler.com/articles/injection.html
Confidence: HIGH
Corroborated By: None directly; interpretation consistent with PSR-11 discouragement
Notes: Judgment call, not absolute rule.

## Evidence 11
Claim: Core principle outweighing pattern choice: separate service configuration from use within application.
Evidence: Fowler concluding: "The choice between Service Locator and Dependency Injection is less important than the principle of separating service configuration from the use of services within an application." Also: separate config enables plugin substitution per deployment.
Source: Martin Fowler — IoC Containers and DI pattern
URL: https://martinfowler.com/articles/injection.html
Confidence: HIGH
Corroborated By: Spring (configuration metadata), Laravel (service providers)
Notes: None

## Evidence 12
Claim: IoC general principle: custom code receives flow of control from external framework; "Don't call us, we'll call you" (Hollywood Principle). Applies to UI loops, callbacks, schedulers, template method, event loops.
Evidence: "Inversion of control is a design principle in which custom-written portions of a program receive the flow of control from an external source (e.g. a framework)... 'Don't call us, we'll call you'."
Source: Wikipedia — Inversion of control
URL: https://en.wikipedia.org/wiki/Inversion_of_control
Confidence: HIGH
Corroborated By: Fowler bliki InversionOfControl, Spring docs
Notes: None

## Evidence 13
Claim: Term "inversion of control" separately used by Java community to mean DI patterns with IoC containers (Spring); refers to framework controlling dependency implementations, not control flow.
Evidence: "The phrase 'inversion of control' has separately also come to be used in the community of Java programmers to refer specifically to the patterns of dependency injection... that occur with 'IoC containers' in Java frameworks such as the Spring Framework."
Source: Wikipedia — Inversion of control
URL: https://en.wikipedia.org/wiki/Inversion_of_control
Confidence: HIGH
Corroborated By: Fowler 2004 (settled on name DI because IoC too generic), Spring docs
Notes: Explains terminology confusion.

## Evidence 14
Claim: Spring defines DI as objects declaring dependencies only via constructor args, factory method args, or post-construction properties; container injects them at bean creation. Inverse of bean locating dependencies itself (direct construction or Service Locator).
Evidence: "Dependency injection (DI) is a specialized form of IoC, whereby objects define their dependencies... only through constructor arguments, arguments to a factory method, or properties... The IoC container then injects those dependencies when it creates the bean. This process is fundamentally the inverse... of the bean itself controlling the instantiation or location of its dependencies by using direct construction of classes or a mechanism such as the Service Locator pattern."
Source: Spring Framework docs 7.0.9
URL: https://docs.spring.io/spring-framework/reference/core/beans/introduction.html
Confidence: HIGH
Corroborated By: Fowler, Microsoft, Laravel
Notes: None

## Evidence 15
Claim: Laravel container manages class dependencies and performs DI via constructor or setter; enables mocking in tests; zero-config resolution for concrete classes; interface bindings registered in service providers.
Evidence: "The Laravel service container is a powerful tool for managing class dependencies and performing dependency injection... class dependencies are 'injected' into the class via the constructor or, in some cases, 'setter' methods... Since the service is injected, we are able to easily 'mock'... when testing." Plus: "$this->app->bind(EventPusher::class, RedisEventPusher::class); ... container should inject RedisEventPusher when a class needs EventPusher."
Source: Laravel docs — Service Container
URL: https://laravel.com/docs/container
Confidence: HIGH
Corroborated By: PSR-11 (Laravel implements it), Microsoft, Spring
Notes: None

## Evidence 16
Claim: Laravel supports singleton, scoped (per request/job lifecycle, e.g. Octane), instance bindings; contextual binding (different impls per consumer); tagging; method invocation injection via App::call.
Evidence: Docs sections: singleton(), scoped(), instance(), when()->needs()->give(), tag/tagged, App::call([...]) for method injection.
Source: Laravel docs — Service Container
URL: https://laravel.com/docs/container
Confidence: HIGH
Corroborated By: Microsoft (AddSingleton/scoped/transient), Spring (bean scopes)
Notes: None

## Evidence 17
Claim: PSR-11 standardizes ContainerInterface with get($id) and has($id); unknown id MUST throw NotFoundExceptionInterface; users SHOULD NOT pass container into objects (Service Locator discouraged).
Evidence: Spec 1.1.2, 1.2, 1.3: "get takes one mandatory parameter... throw NotFoundExceptionInterface if identifier not known"; "Users SHOULD NOT pass a container into an object so that the object can retrieve its own dependencies. This means the container is used as a Service Locator which is generally discouraged."
Source: PHP-FIG PSR-11
URL: https://www.php-fig.org/psr/psr-11/
Confidence: HIGH
Corroborated By: Fowler, Laravel PSR-11 section
Notes: RFC 2119 keywords.

## Evidence 18
Claim: .NET has built-in DI: abstract via interface/base class, register in IServiceCollection at startup, inject via constructor; framework creates and disposes instances. Hard-coded `new` causes swap pain, scattered config, untestability.
Evidence: ".NET supports the dependency injection (DI) software design pattern, which is a technique for achieving Inversion of Control (IoC)... use of an interface or base class to abstract implementation; registration in a service container; injection into constructor... Hard-coded dependencies... To replace MessageWriter you must modify Worker... difficult to unit test... should use a mock or stub."
Source: Microsoft Learn — Dependency injection .NET
URL: https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection
Confidence: HIGH
Corroborated By: Spring, Laravel, Fowler naive MovieLister example
Notes: None

## Evidence 19
Claim: DI service lifetimes: Singleton, Scoped, Transient; scope validation catches scoped-from-root and scoped-into-singleton errors in dev.
Evidence: ".NET docs: AddSingleton/AddScoped/AddTransient... Validating service scopes catches these situations... Scoped services aren't resolved from root... aren't injected into singletons."
Source: Microsoft Learn — Dependency injection .NET
URL: https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection
Confidence: HIGH
Corroborated By: Laravel singleton/scoped, Spring bean scopes
Notes: None

## Evidence 20
Claim: PHP interfaces define method contracts without implementation; classes implement via `implements`; enable interchangeable implementations (multiple payment gateways, caching strategies) swappable without changing consuming code.
Evidence: PHP manual: "Object interfaces allow you to create code which specifies which methods a class must implement, without having to define how... To allow developers to create objects of different classes that may be used interchangeably because they implement the same interface... A common example is multiple database access services, multiple payment gateways, or different caching strategies. Different implementations may be swapped out without requiring any changes to the code that uses them."
Source: PHP manual — Object Interfaces
URL: https://www.php.net/manual/en/language.oop5.interfaces.php
Confidence: HIGH
Corroborated By: Topic spec (PaymentGatewayInterface), Laravel interface binding, Fowler MovieFinder interface
Notes: None
