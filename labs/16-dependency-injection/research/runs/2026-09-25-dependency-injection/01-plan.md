# Research Topic
Dependency Injection (DI) & Inversion of Control (IoC) — Writing Testable and Extensible Code

# Objective
To investigate Dependency Injection (DI) and Inversion of Control (IoC) concepts, their impact on software design (coupling, testability, extensibility), and to distinguish between appropriate DI usage and anti-patterns like pervasive Service Locator.

# Research Questions
1. What is the fundamental difference between Inversion of Control (IoC) and Dependency Injection (DI)?
2. How does DI reduce coupling and improve testability compared to direct object instantiation?
3. What is the role of an IoC Container, and what are the standard practices for configuring it?
4. Why is using a Service Locator generally discouraged, and how does it differ from DI?
5. When is it appropriate *not* to use DI (e.g., simple data objects/value objects)?

# Search Strategy
- Query Martin Fowler's seminal article on IoC and DI.
- Review official documentation for major frameworks (.NET, Spring).
- Analyze PHP-FIG standards (PSR-11) for container interfaces and its Meta Document regarding Service Locator.

# Expected Primary Sources
- "Inversion of Control Containers and the Dependency Injection pattern" by Martin Fowler
- .NET Documentation: "Dependency injection - .NET" (Microsoft)
- Spring Framework Documentation: "Introduction to the Spring IoC Container and Beans"
- PSR-11: Container Interface and PSR-11 Meta Document (PHP-FIG)

# Risks / Unknowns
- Terminology overlap between IoC, DI, and Service Locator may cause confusion.
- Nuances between constructor injection and setter injection preference across different communities.
