# 04 — Contradictions & Divergences

## Contradiction 1
**Source A (Microsoft .NET)**: Presents DI as a solution to hard-coded `new`, emphasizing that "the app should use a mock or stub... which isn't possible with this approach" (contrasting hard-coded `new MessageWriter = new()`).
**Source B (Fowler 2004, Section "Using a Service Locator")**: States "both [DI and Service Locator] are very amenable to stubbing" if well-designed, and DI is not uniquely about testing.
**Assessment**: The difference is emphasis, not factual contradiction. Microsoft frames DI primarily as a testing solution. Fowler notes testing is the "first benefit noticed" but Service Locator can be equally testable. The core capability (substitution via abstraction) is shared.

## Divergence 1
**Source A (Fowler 2004, Wikipedia DI article)**: Lists three forms — Constructor Injection, Setter Injection, Interface Injection.
**Source B (Spring 7.0.9 docs, Microsoft .NET docs, Laravel docs)**: Only document Constructor and Setter injection as the practical two. Interface injection is absent from modern framework documentation.
**Assessment**: Evolution of practice — Interface injection historically used (Avalon framework), but modern practice favors constructor/setter due to simplicity and explicit dependencies. Not a contradiction; an observed shift.

## Divergence 2
**Source A (Fowler 2004)**: Advocates programmatic builders over XML configuration: "people are over-eager to define configuration files... a programming language makes a straightforward and powerful configuration mechanism."
**Source B (Spring 2006 docs, modern Laravel, .NET)**: Frameworks offer XML, annotations, attributes, and PHP 8 attributes as configuration mechanisms, blurring the "code vs file" distinction.
**Assessment**: Technical evolution. Early 2000s debate (XML vs code) evolved into hybrid approaches (annotations for simple, builders/config for complex). Not a contradiction; a shift in available options.

## Divergence 3
**Source A (Topic specification)**: States "Kalau constructor berisi 12 parameter... Biasanya ada masalah desain" (12+ parameters indicate design problem).
**Source B (Fowler 2004)**: Qualitatively says "If you have a lot of constructor parameters things can look messy... often a sign of an over-busy object."
**Assessment**: Same principle expressed differently. Topic spec provides a heuristic threshold (12). Fowler describes the problem qualitatively. The difference is specificity, not contradiction on the underlying design signal.

## Divergence 4
**Source A (Topic specification)**: Provides a definitive list of when NOT to use DI: "DateTime, Money, Address" as examples.
**Source B (Fowler 2004)**: Discusses "value objects which represent entities in the program's domain" but provides no concrete list.
**Source C (Wikipedia DI)**: States DI "reduces boilerplate code, since all dependency creation is handled by a singular component" implying value objects might still be created manually.
**Assessment**: Topic spec provides practical heuristics without independent cross-check. The concept exists (value objects vs services), but the specific examples/lines are internal guidance rather than external standard. Weak evidence.

## Divergence 5
**Source A (PSR-11)**: Uses "SHOULD NOT" (RFC 2119) for "Users SHOULD NOT pass a container into an object."
**Source B (Java Service Locator implementations, some legacy codebases)**: Service Locator used extensively even with DI containers available.
**Assessment**: PSR-11 is a PHP standard; "SHOULD NOT" is a strong recommendation (not a strict prohibition) for compliant implementations. Legacy/existing Java usage differs. Not a contradiction; PSR-11 explicitly chose to discourage the pattern for interoperability and testability reasons. The difference is in the standard's intent vs actual practice in some ecosystems.

## No Material Contradictions Discovered For:
- DI definition (all sources agree)
- Three DI forms (historically recognized; modern docs omit interface)
- DI benefits (coupling reduction, testability)
- Service lifetimes (singleton/scoped/transient across all containers)
- IoC is broader than DI (Fowler, Wikipedia)