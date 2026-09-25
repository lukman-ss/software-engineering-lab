# Content Brief

Topic: Dependency Injection (DI) and Inversion of Control (IoC)
Target Reader: Software Engineers
Problem: Hard-coded dependencies create tightly coupled code, making isolated testing and swapping implementations difficult.
Core Mental Model: Externalize object creation. Instead of a class building its own dependencies, an external container or assembler provides them (usually via constructors).
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts: Separation of Configuration from Use, Constructor Injection, Service Locator Anti-Pattern, Value Objects.
Verified Behaviors: Dependency mocking in tests, explicit constructor dependencies, hidden dependencies via locator.
Available Case Studies: Payment processor using a payment gateway interface.
Warnings: The `BadProcessor` error paths (invalid amount, gateway failure) are omitted from the test suite to keep the focus on structural dependency differences. The demonstration does not use an automated DI framework, relying on manual injection.
