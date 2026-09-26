# Contradictions Audit — Lab 26: Contract Testing

## Evaluation of Research Contradiction Analysis

The research agent documented potential points of tension across sources in `research/04-contradictions.md`. These were examined during this audit:

### 1. Dual Usage of Term "Contract Testing" (Provider-only vs Integration)
- **Statement A:** In some vendor contexts and literature, "contract testing" refers to provider-only schema validation (e.g., verifying OpenAPI specs against endpoints).
- **Statement B:** In Pact and CDC canonical literature, "contract testing" strictly means integration verification across two communicating parties.
- **Type:** SOURCE_CONFLICT / TERMINOLOGY
- **Impact:** LOW. The research explicitly distinguishes between "provider schema testing" and "integration contract testing" and focuses the lab on the latter.
- **Assessment:** Handled properly and clearly disambiguated.

### 2. Immediate Build Failure vs Out-of-Band Task
- **Statement A (Fowler 2011):** "A failure in a contract test shouldn't necessarily break the build in the same way that a normal test failure would. It should, however, trigger a task to get things consistent again."
- **Statement B (Pact / Pactflow):** Contract verification failures in CI block provider deployments (`can-i-deploy`) to prevent broken contracts in production.
- **Type:** SOURCE_CONFLICT (Historical pattern evolution vs Modern tooling)
- **Impact:** LOW. The research report explicitly notes this evolutionary difference in Finding 4 and Confidence notes.
- **Assessment:** Accurately reconciled.

### 3. Additive Changes Safety vs Strict Deserializers
- **Statement A:** Adding optional fields is backward-compatible.
- **Statement B:** If a consumer configures strict deserialization (e.g., Jackson `FAIL_ON_UNKNOWN_PROPERTIES`), adding fields breaks the client.
- **Type:** IMPLEMENTATION_SPECIFIC_BOUNDARY
- **Impact:** LOW. The research report notes this assumption and records it under Limitations and Open Questions.
- **Assessment:** Accurately qualified.

---

## Direct Document Cross-Check

- `research/01-plan.md` vs `research/05-report.md`: All research questions mapped directly to corresponding findings.
- `research/02-sources.md` vs `research/03-evidence.md`: All evidence citations trace to reachable, evaluated sources.
- `research/05-report.md` vs `research/06-open-questions.md`: Open limitations and boundary conditions are appropriately isolated.

No material internal contradictions found.
