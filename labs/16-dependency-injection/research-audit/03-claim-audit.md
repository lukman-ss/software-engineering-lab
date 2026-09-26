# 03 - Claim Audit

## Claim 1
Claim:
Dependency Injection is an object receiving dependencies from an external injector rather than instantiating them internally; it separates object construction from use and fosters loose coupling.

Location:
`research/03-evidence.md` (Evidence 1), `research/05-report.md` (Finding 1)

Evidence Provided:
Direct quotes from Wikipedia DI article, corroboration from Fowler (2004) and Spring Framework 7.0.9 documentation.

Source:
Source 1 (Martin Fowler), Source 3 (Wikipedia DI), Source 5 (Spring Framework)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Foundational definition across literature.

---

## Claim 2
Claim:
Inversion of Control (IoC) is a broad design principle (Hollywood Principle: "Don't call us, we'll call you"); DI is a specific form of IoC controlling dependency implementations rather than overall application control flow.

Location:
`research/03-evidence.md` (Evidence 2, 12, 13), `research/05-report.md` (Finding 2)

Evidence Provided:
Fowler's 2004 rationale for coining DI ("as a result I think we need a more specific name for this pattern. Inversion of Control is too generic a term...") and Wikipedia IoC documentation.

Source:
Source 1 (Martin Fowler), Source 2 (Wikipedia IoC)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately clarifies historical naming ambiguity.

---

## Claim 3
Claim:
DI involves four roles: services, clients, interfaces, and injectors. The injector must not be the client to avoid circular dependencies.

Location:
`research/03-evidence.md` (Evidence 3), `research/05-report.md` (Finding 1)

Evidence Provided:
Wikipedia DI role definitions and Fowler's assembler description.

Source:
Source 1 (Martin Fowler), Source 3 (Wikipedia DI)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Standard structural breakdown.

---

## Claim 4
Claim:
DI provides benefits of reduced coupling, improved testability via test doubles (stubs/mocks), ease of maintenance, and separation of configuration from use.

Location:
`research/03-evidence.md` (Evidence 4, 5, 11), `research/05-report.md` (Finding 1, 6)

Evidence Provided:
Microsoft .NET DI overview, Fowler (2004), Spring documentation, Laravel documentation.

Source:
Source 1 (Martin Fowler), Source 4 (Laravel), Source 5 (Spring), Source 7 (Microsoft .NET)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Universally corroborated across vendors and academic/canonical articles.

---

## Claim 5
Claim:
DI introduces trade-offs and costs: configuration burden, harder code traceability, reflection overhead/tooling hindrance, upfront development effort, and potential framework lock-in.

Location:
`research/03-evidence.md` (Evidence 6), `research/05-report.md` (Finding 10)

Evidence Provided:
Wikipedia DI disadvantages, corroborated by Fowler's analysis of IoC complexity and debugging hurdles.

Source:
Source 1 (Martin Fowler), Source 3 (Wikipedia DI)

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
LOW

Notes:
Accurately balances architectural trade-offs without dogmatism.

---

## Claim 6
Claim:
There are three historical forms of DI: Constructor Injection (Type 3), Setter Injection (Type 2), and Interface Injection (Type 1). Modern mainstream frameworks focus primarily on Constructor and Setter injection, while Interface injection is largely obsolete.

Location:
`research/03-evidence.md` (Evidence 7), `research/05-report.md` (Finding 3)

Evidence Provided:
Fowler's taxonomy (Avalon, PicoContainer, Spring) contrasted with modern framework docs (Spring 7.0.9, .NET, Laravel) which omit Interface Injection.

Source:
Source 1 (Martin Fowler), Source 5 (Spring), Source 7 (Microsoft .NET)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurate historical tracking and assessment of modern software engineering practices.

---

## Claim 7
Claim:
Developers should prefer Constructor Injection as default (ensures valid, immutable object state at birth) and switch/augment with Setter Injection for optional dependencies, circular dependencies, or inheritance complexity.

Location:
`research/03-evidence.md` (Evidence 8), `research/05-report.md` (Finding 4)

Evidence Provided:
Fowler's recommendation, Kent Beck references, and Microsoft .NET architectural guidance.

Source:
Source 1 (Martin Fowler), Source 7 (Microsoft .NET)

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
LOW

Notes:
Standard best practice reflected in industry documentation.

---

## Claim 8
Claim:
Service Locator and DI both decouple clients from concrete implementations. However, Service Locator creates a hard dependency on the locator API in every consumer. PSR-11 recommends that users `SHOULD NOT` pass a container into an object.

Location:
`research/03-evidence.md` (Evidence 9, 10, 17), `research/05-report.md` (Finding 5, 9)

Evidence Provided:
Fowler comparative analysis and PSR-11 ContainerInterface specification section 1.3 verbatim.

Source:
Source 1 (Martin Fowler), Source 6 (PSR-11)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
PSR-11 RFC 2119 keyword is correctly documented as `SHOULD NOT`.

---

## Claim 9
Claim:
Containers manage service lifetimes (Singleton, Scoped, Transient/Prototype); modern containers (such as .NET) validate scopes to prevent resolving Scoped services from the root provider or injecting Scoped services into Singletons (captive dependencies).

Location:
`research/03-evidence.md` (Evidence 16, 19), `research/05-report.md` (Finding 7)

Evidence Provided:
Microsoft .NET dependency injection documentation on lifetimes and scope validation; Laravel singleton/scoped documentation.

Source:
Source 4 (Laravel), Source 7 (Microsoft .NET)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Matches official vendor documentation on dependency lifetime mechanics.

---

## Claim 10
Claim:
A constructor with 12 or more parameters is an anti-pattern indicating a Single Responsibility Principle (SRP) violation.

Location:
`research/05-report.md` (Finding 11), `research/04-contradictions.md` (Divergence 3)

Evidence Provided:
Topic specification guideline; Fowler qualitatively notes "a lot of parameters... often a sign of an over-busy object".

Source:
Topic spec (internal heuristic), Source 1 (Martin Fowler)

Source Actually Supports Claim:
PARTIAL

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
The research appropriately qualified this claim in `05-report.md`: "the 12-parameter threshold is a lab-specific heuristic, not an industry standard". Fowler describes the issue qualitatively without fixing a numeric constant of 12. Correctly flagged as a heuristic.

---

## Claim 11
Claim:
Value objects (such as `DateTime`, `Money`, `Address`) should not be injected via DI containers, but instantiated directly; DI should be reserved for services and external boundary components (Gateways, Repositories, Notifiers).

Location:
`research/05-report.md` (Finding 12), `research/04-contradictions.md` (Divergence 4)

Evidence Provided:
Topic specification guideline; Fowler architectural distinction between Domain Entities/Value Objects and Services.

Source:
Topic spec (internal heuristic), Source 1 (Martin Fowler)

Source Actually Supports Claim:
PARTIAL

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
The specific list (`DateTime`, `Money`, `Address`) originates from internal lab requirements, while the foundational principle (avoid container management for transient state/value objects) is corroborated by DDD and Fowler. Explicitly marked as a lab heuristic in the report.
