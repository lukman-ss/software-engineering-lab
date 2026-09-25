# Claim Audit

## Claim 1

Claim: Inversion of Control is a broad architectural concept; Dependency Injection is a specific pattern where an assembler/container provides dependencies from the outside, separating configuration from use.

Location:
03-evidence.md (Evidence 1), 05-report.md (Finding 1)

Evidence Provided: Quotes from Martin Fowler ("inversion is about how they lookup a plugin implementation... settled on the name Dependency Injection") and Spring framework docs.

Source:
Martin Fowler, Spring Framework

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Direct textual quotes match the sources.

## Claim 2

Claim: Direct instantiation inside business classes causes tight coupling, scatters configuration, and makes unit testing difficult by preventing substitution with mocks/stubs.

Location:
03-evidence.md (Evidence 2), 05-report.md (Finding 2)

Evidence Provided: Microsoft Learn documentation explicitly discouraging hard-coded dependencies and Fowler's `MovieLister` example.

Source:
Microsoft Learn .NET, Martin Fowler

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Supported by standard literature.

## Claim 3

Claim: Constructor injection ensures objects are in a valid state upon creation and makes dependencies explicit and immutable, making it the preferred injection style over setter/property injection.

Location:
03-evidence.md (Evidence 3), 05-report.md (Finding 3)

Evidence Provided: Fowler's article recommending starting with constructor injection and Microsoft Learn .NET standardizing on it.

Source:
Martin Fowler, Microsoft Learn .NET

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
LOW

Notes: Properly characterized as a strong preference and best practice based on the literature.

## Claim 4

Claim: Using a container as a Service Locator inside classes is an anti-pattern because it hides class dependencies and couples the class to the container interface.

Location:
03-evidence.md (Evidence 4), 05-report.md (Finding 4)

Evidence Provided: PSR-11 Meta Document explicitly discouraging injecting the container for an object to retrieve its own dependencies.

Source:
PSR-11: Container interface Meta Document

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Perfectly aligned with the PSR-11 meta document text.

## Claim 5

Claim: Simple value objects or data structures that lack external side effects and environmental dependencies do not need to be injected via DI.

Location:
03-evidence.md (Evidence 5), 05-report.md (Finding 5)

Evidence Provided: Fowler mentioning that direct instantiation is fine for simple value objects.

Source:
Martin Fowler

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
LOW

Notes: Concept is sound and well-understood.
