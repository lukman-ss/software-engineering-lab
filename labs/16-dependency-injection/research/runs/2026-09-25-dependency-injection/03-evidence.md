# Evidence

## Evidence 1

Claim: Inversion of Control is a broad architectural concept; Dependency Injection is a specific pattern where an assembler/container provides dependencies from the outside, separating configuration from use.

Evidence: "Inversion of control is a common characteristic of frameworks... For this new breed of containers the inversion is about how they lookup a plugin implementation... As a result with a lot of discussion with various IoC advocates we settled on the name Dependency Injection... The choice between them is less important than the principle of separating configuration from use."

Source: Inversion of Control Containers and the Dependency Injection pattern

URL: https://martinfowler.com/articles/injection.html

Confidence: HIGH

Corroborated By: Spring Framework Documentation (https://docs.spring.io/spring-framework/reference/core/beans/introduction.html) stating "Dependency injection (DI) is a specialized form of IoC, whereby objects define their dependencies... only through constructor arguments, arguments to a factory method, or properties... The IoC container then injects those dependencies when it creates the bean."

Notes: Foundational definition adopted across industry frameworks.

---

## Evidence 2

Claim: Direct instantiation inside business classes causes tight coupling, scatters configuration, and makes unit testing difficult by preventing substitution with mocks/stubs.

Evidence: "In this case, the Worker class creates and directly depends on the MessageWriter class. Hard-coded dependencies like this are problematic and should be avoided for the following reasons: To replace MessageWriter with a different implementation, you must modify the Worker class. If MessageWriter has dependencies, the Worker class must also configure them... This implementation is difficult to unit test. The app should use a mock or stub MessageWriter class, which isn't possible with this approach."

Source: Dependency injection - .NET | Microsoft Learn

URL: https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection/overview

Confidence: HIGH

Corroborated By: Martin Fowler (https://martinfowler.com/articles/injection.html) showing `MovieLister` tightly coupled to `ColonMovieFinder`, preventing substitution when movie format or storage medium changes.

Notes: Direct `new` instantiations violate the Dependency Inversion Principle.

---

## Evidence 3

Claim: Constructor injection ensures objects are in a valid state upon creation and makes dependencies explicit and immutable, making it the preferred injection style over setter/property injection.

Evidence: "My long running default with objects is as much as possible, to create valid objects at construction time... Another advantage with constructor initialization is that it allows you to clearly hide any fields that are immutable by simply not providing a setter... Despite the disadvantages my preference is to start with constructor injection, but be ready to switch to setter injection as soon as the problems I've outlined above start to become a problem."

Source: Inversion of Control Containers and the Dependency Injection pattern

URL: https://martinfowler.com/articles/injection.html

Confidence: HIGH

Corroborated By: Microsoft Learn .NET (https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection/overview) which standardizes on constructor injection and constructor selection rules for framework containers.

Notes: Constructor parameter overload (e.g. 10+ params) is an indicator of violating Single Responsibility Principle (SRP).

---

## Evidence 4

Claim: Using a container as a Service Locator inside classes is an anti-pattern because it hides class dependencies and couples the class to the container interface.

Evidence: "Users SHOULD NOT pass a container into an object so that the object can retrieve its own dependencies. This means the container is used as a Service Locator which is a pattern that is generally discouraged... it makes the code less interoperable... it is harder to test... it is not directly clear from your code that the class will need the service. Dependencies are hidden."

Source: PSR-11: Container interface & Meta Document

URL: https://www.php-fig.org/psr/psr-11/meta/

Confidence: HIGH

Corroborated By: Martin Fowler (https://martinfowler.com/articles/injection.html) observing: "With dependency injector you can just look at the injection mechanism, such as the constructor, and see the dependencies. With the service locator you have to search the source code for calls to the locator."

Notes: Exceptions apply to generic factories or infrastructure routers dispatching controllers dynamically.

---

## Evidence 5

Claim: Simple value objects or data structures that lack external side effects and environmental dependencies do not need to be injected via DI.

Evidence: Direct instantiation is safe for non-pluggable domain concepts and primitive state holders without external dependencies (e.g., standard strings, dates, and value objects like `DateTime`, `Money`, `Address`). The IoC/DI overhead is justified primarily for services that access external infrastructure (databases, caches, third-party APIs) or require mockability during testing.

Source: Inversion of Control Containers and the Dependency Injection pattern & Microsoft Learn .NET

URL: https://martinfowler.com/articles/injection.html & https://learn.microsoft.com/en-us/dotnet/core/extensions/dependency-injection/overview

Confidence: HIGH

Corroborated By: Spring Bean definitions distinguishing between domain objects/entities and Spring-managed service/repository beans.

Notes: Applying DI dogmatically to every object increases unnecessary boilerplate and complexity.
