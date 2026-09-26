# Content Brief

Topic: Dependency Injection (DI) and Inversion of Control (IoC)
Target Reader: Software Engineers
Problem: Hard-coded dependencies create tightly coupled code, making isolated testing and swapping implementations difficult.
Core Mental Model: Externalize object creation. Instead of a class building its own dependencies, an external container or assembler provides them (usually via constructors).
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts: Separation of Configuration from Use, Constructor Injection, Service Locator Anti-Pattern, Value Objects, Interface-based design, Test doubles for isolation.
Verified Behaviors: Dependency mocking in tests, explicit constructor dependencies, hidden dependencies via locator, input validation before infrastructure calls, value object direct instantiation.
Available Case Studies: Payment processor using a payment gateway interface (Processor vs BadProcessor contrast).
Warnings: 
- The demonstration does not use an automated DI framework, relying on manual injection.
- The 12-parameter constructor threshold and the specific value object list (DateTime, Money, Address) are lab-specific heuristics, not industry standards.
- PSR-11 uses "SHOULD NOT" (RFC 2119) for passing containers into objects — a strong recommendation, not a strict prohibition.
- Container lifecycle management (singleton/scoped/transient) and scope validation are discussed in research but NOT demonstrated in this lab's implementation.
- No empirical evidence exists for defect reduction or performance overhead claims.
