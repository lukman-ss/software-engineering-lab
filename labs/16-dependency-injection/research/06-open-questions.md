# 06 — Open Questions

## Unanswered Questions

### Empirical Evidence Gap
Q1. Does DI measurably reduce defect rates, improve regression testing speed, or accelerate onboarding vs hard-coded dependencies in real-world projects?
Evidence status: NOT VERIFIED — no empirical studies found; only authoritative opinion (Fowler's testing benefit emphasis).

Q2. What is quantified performance overhead of container resolution vs direct `new` in high-throughput systems? Does it matter at scale?
Evidence status: NOT VERIFIED — not discussed in framework docs or surveys.

### Language-Specific Nuances
Q3. How does DI interact with Rust's ownership/borrowing model? TypeScript's structural typing? Go's composition-over-inheritance culture?
Evidence status: WEAK — Wikipedia Rust/Go/TypeScript snippets present; no deep dive on trade-offs vs idioms.

Q4. In PHP, how do attributes (PHP 8) for binding (`#[Bind]`, `#[Singleton]`, `#[Scoped]`) compare to explicit service provider bindings in complexity, runtime cost, and tooling support?
Evidence status: MEDIUM — Laravel docs show attributes; no comparative analysis with programmatic binding.

### Design Boundary Questions
Q5. Precise criteria for "value object" vs "service" beyond heuristics. What about a `MailTemplateRenderer` that depends on `Twig` vs a `PaymentGateway` that depends on `HttpClient`? Where is the line?
Evidence status: NOT VERIFIED — Topic spec lists DateTime/Money/Address, but no principled definition (no external dependency, no interface contract needed).

Q6. Edge cases of constructor over-injection: Is 8 params acceptable for a well-encapsulated service? 12 is a heuristic from spec; what does actual codebase analysis show?
Evidence status: MEDIUM — Fowler says "a lot" qualitatively; no data. 12 may be rule-of-thumb for a specific project, not universal.

Q7. When should you use a factory pattern versus DI for complex object graphs? Where does DI stop and factory begin?
Evidence status: WEAK — topic spec exercise mentions factory implicitly but no guidance.

### Framework/Container Design
Q8. How do keyed services (multiple impls with keys) compare to tagged services (collections) for plugin systems? When to prefer one over the other?
Evidence status: MEDIUM — Laravel/Spring/.NET support both; no comparison guidance in docs.

Q9. What is the best practice for circular dependencies when using DI? What should container behavior be (throw error, lazy, scoped proxy)?
Evidence status: WEAK — Wikipedia mentions circular dependency risk but no resolution guidance.

Q10. How does DI interact with service mesh (Envoy/Istio), function-as-a-service (AWS Lambda), or sidecar architectures? Does DI matter at the infrastructure layer?
Evidence status: NOT VERIFIED — no coverage in docs or articles.

## Weak Evidence Items
- Value object list (DateTime/Money/Address) — from spec, not primary source
- "12 parameter" heuristic — spec-specific, not academic or cross-framework validated
- NestJS/TypeScript examples — fetch failed, evidence absent

## Possible Next Research Directions
1. Survey industrial case studies (engineering blogs, conference talks) on DI adoption outcomes
2. Analyze large open-source PHP/.NET/Java repos for dependency counts vs test coverage correlation
3. Benchmark DI container resolution overhead across frameworks (Spring vs Pico vs Laravel vs .NET)
4. Compare DI vs Composition Over Inheritance in Go/TypeScript codebases
5. Deep dive on circular dependency resolution strategies across containers