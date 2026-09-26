# Contradiction Audit — Contract Testing Research

## Contradiction 1: Provider-only schema verification vs Integration Contract Testing
- Statement A: "Contract testing" is sometimes used to describe standalone provider validation against OpenAPI schemas (`04-contradictions.md:Section 1`, citing `docs.pact.io`).
- Statement B: Canonical contract testing requires isolated tests across both consumer and provider (`05-report.md:Finding 1`, citing Fowler & Pact Docs).
- Type: TERMINOLOGY_AMBIGUITY
- Impact: Can cause confusion between schema testing and CDC.
- Assessment: Handled cleanly. The report explicitly documents this distinction and specifies that "integration contract testing" is the standard adopted.

## Contradiction 2: CI failure blocking vs communication trigger
- Statement A: Martin Fowler notes that a contract test failure might not necessarily break the build immediately like a unit test, but trigger reconciliation (`04-contradictions.md`, citing Fowler 2011).
- Statement B: Modern tooling (Pact / Pactflow) emphasizes hard release blocking (`can-i-deploy`) to ensure breaking changes never hit production (`05-report.md:Finding 4`).
- Type: EVOLUTIONARY_PRACTICE
- Impact: Operational difference depending on pipeline tooling.
- Assessment: Handled accurately. Fowler described 2011 external partner integration rhythm; modern microservice CDC uses automated CI gates.

## Contradiction 3: Additive changes safety vs strict consumer validation
- Statement A: Additive changes are backward-compatible (`03-evidence.md:Evidence 9`).
- Statement B: If a consumer enforces strict schema validation (e.g. rejecting unknown properties), an additive change is breaking (`04-contradictions.md:Section 3`).
- Type: INTERNAL_QUALIFICATION
- Impact: Edge-case failure in strict parsers.
- Assessment: Accurately qualified. The report highlights that additive changes rely on consumer tolerance (e.g. Go `json.Unmarshal` default behavior).

No unresolved material contradictions found.
