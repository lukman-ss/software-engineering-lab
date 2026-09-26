# Open Questions

## Unanswered Questions
1. **Empirical impact**: No quantitative studies found measuring DI's effect on defect rate, change lead time, or maintainability. Is the "loosely coupled → easier to evolve" claim backed by data or only by expert consensus?
2. **Performance overhead**: What is the measurable cost of container resolution (reflection, proxy creation) vs. manual construction, especially in cold-start / serverless contexts?
3. **Circular dependencies**: How do containers detect/resolve circular graphs? What are best practices to avoid them vs. using property injection or Provider/Lazy patterns?
4. **Compile-time vs runtime DI**: Dagger (Java), Wire (Go), compile-time containers — trade-offs in safety, binary size, startup time not covered in fetched sources.
5. **Scope validation edge cases**: .NET's scope validation catches scoped→singleton leaks. Do Laravel/Spring have equivalent guards? How common are lifetime-mismatch bugs in production?
6. **PPOB-specific**: What does idiomatic DI look like for PHP PPOB providers (Digiflazz, etc.) when providers have different auth, retry, and webhook semantics behind one interface?

## Weak Evidence (Needs Stronger Sources)
- **NestJS/TypeScript DI** (Source 7): webfetch returned redirect only; `@Injectable()`/`@Module()` claims are MEDIUM confidence.
- **Value-object vs service boundary** (Evidence 10/11): Checklist "DateTime/Money no DI; DB/Gateway yes DI" comes from topic spec, not a primary source.
- **"12 constructor parameters" heuristic**: Specific threshold from spec, not from Fowler or framework docs. Fowler says "a lot" qualitatively.
- **Interface Injection current relevance**: Described historically by Fowler 2004; no modern framework docs still document it — confidence that it is obsolete is inferred, not explicitly stated by a source.

## Claims Needing Deeper Research
- Liskov Substitution + interface-based DI: need primary source linking DIP to DI container practice (e.g., Clean Architecture, SOLID original texts).
- Test pyramid implication of DI (more unit tests, fewer integration tests) — asserted but not sourced.
- Configuration: code vs XML vs annotations/attributes — Fowler 2004 predates annotation-based Java config and PHP 8 attributes; modern guidance missing.
- Go idiom: explicit wire-up without container is often preferred; no Go primary source fetched.

## Possible Next Research Directions
1. Fetch Mark Seemann "Dependency Injection in .NET" (2nd ed.) and Robert C. Martin "Clean Architecture" for DIP/LSP grounding.
2. Compare compile-time DI (Dagger, Wire) vs reflection-based containers with benchmarks.
3. Add Go source: `google/wire` docs + Go community DI debate.
4. Empirical study search: IEEE/ACM papers on coupling metrics before/after DI adoption.
5. Laravel deep-dive: attributes `#[Bind]`, `#[Singleton]`, contextual `#[Storage]` — verify against 12.x source code.
6. Spring deep-dive: `@Autowired`, `@Primary`, `@Qualifier` resolution rules and failure modes.
