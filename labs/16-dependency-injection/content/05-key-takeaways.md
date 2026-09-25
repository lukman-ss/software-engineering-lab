# Key Takeaways

1. **Separation of Configuration from Use:** Object instantiation must be externalized to decouple logic from infrastructure, enabling swap-out without touching business layers.
2. **Fast, Isolated Unit Testing:** Constructor injection of interfaces enables replacing concrete integrations with lightweight mocks.
3. **Constructor Injection is the Standard:** It enforces explicit dependency resolution, preventing objects from existing in invalid or uninitialized states.
4. **Service Locator is an Anti-Pattern:** Injecting the container itself obscures underlying requirements and permanently couples domain logic to the framework's API.
5. **Value Objects Do Not Need DI:** Pure data structures or value types lacking external side-effects (like `Money`) should be instantiated directly, not managed by an IoC container.
