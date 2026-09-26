# Claim Audit — Lab 26: Contract Testing

## Claim 1
Claim: Contract testing asserts that inter-application messages conform to a shared understanding documented in a contract, validating integration without executing full end-to-end environments.
Location: `research/05-report.md: Finding 1`; `research/03-evidence.md: Evidence 1`
Evidence Provided: Pact Docs and Martin Fowler Bliki quotes.
Source: Source 1 (docs.pact.io), Source 2 (martinfowler.com/bliki/ContractTest.html)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Fully verified against live documentation.

---

## Claim 2
Claim: In Consumer-Driven Contracts (CDC), expectations are defined from the perspective of what the consumer needs (open and incomplete subset), and provider verifies adherence.
Location: `research/05-report.md: Finding 2, Finding 9`; `research/03-evidence.md: Evidence 3, 10`
Evidence Provided: Ian Robinson / Martin Fowler (2006), Pact Docs.
Source: Source 3 (martinfowler.com/articles/consumerDrivenContracts.html), Source 1 (docs.pact.io), Source 5 (pactflow.io)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Canonical concept accurately synthesized from original foundational paper.

---

## Claim 3
Claim: Schema/OpenAPI testing tests a single system's adherence to a specification at a point in time, whereas code-based contract testing tests agreement between two communicating systems with concrete examples and execution of real application code.
Location: `research/05-report.md: Finding 3`; `research/03-evidence.md: Evidence 4`
Evidence Provided: Matt Fellows (Pactflow) breakdown of schema vs contract testing.
Source: Source 6 (pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1), Source 1 (docs.pact.io)
Source Actually Supports Claim: YES
Classification: FACT / INTERPRETATION
Severity: LOW
Notes: Accurately captures industry distinction between specification conformance and bidirectional agreement.

---

## Claim 4
Claim: Contract verification can run in provider CI to fail builds or block deployments (`can-i-deploy`) when consumer contracts are broken.
Location: `research/05-report.md: Finding 4`; `research/03-evidence.md: Evidence 3`
Evidence Provided: Pact Docs, Pactflow explainer, Martin Fowler bliki nuance.
Source: Source 1, Source 4, Source 5, Source 2
Source Actually Supports Claim: YES
Classification: FACT / IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: The research report accurately captures the nuance that while Martin Fowler originally suggested contract tests might trigger out-of-band communication rather than immediate build breaks, modern CDC tooling (Pact / Pact Broker) specifically implements automated deployment blocking.

---

## Claim 5
Claim: Contract testing concepts apply equally to asynchronous message-based systems (message queues, event streams like Kafka/RabbitMQ) as to synchronous HTTP APIs.
Location: `research/05-report.md: Finding 5`; `research/03-evidence.md: Evidence 11`
Evidence Provided: Pact Docs quote regarding queues and messages.
Source: Source 1 (docs.pact.io), Source 5 (pactflow.io)
Source Actually Supports Claim: YES
Classification: FACT
Severity: LOW
Notes: Supported directly in Pact documentation.

---

## Claim 6
Claim: All three proposed backend modifications in the lab scenario (`status: IN_PROGRESS -> in_progress`, `customer.name -> customer.full_name`, `total: integer -> string`) constitute breaking changes for consumers expecting the original schema.
Location: `research/05-report.md: Finding 7`; `research/03-evidence.md: Evidence 8`
Evidence Provided: Analysis based on strict JSON deserialization, semantic typing, enum matching, and Lab 06 API versioning taxonomy.
Source: Source 7 (Lab 06), Topic spec
Source Actually Supports Claim: YES
Classification: FACT / INTERPRETATION
Severity: LOW
Notes: Sound analysis. Field rename and primitive type modification are unconditionally breaking in standard JSON parsers/strongly-typed clients; enum casing is breaking unless client-side normalization is explicitly guaranteed across all consumers.

---

## Claim 7
Claim: Additive changes (adding optional fields) are backward-compatible under the condition that consumers are tolerant of unknown fields.
Location: `research/05-report.md: Finding 6`; `research/03-evidence.md: Evidence 9`
Evidence Provided: Robustness Principle, Go `json.Unmarshal` default behavior, Fowler CDC paper.
Source: Source 3, Source 7
Source Actually Supports Claim: YES
Classification: FACT / IMPLEMENTATION-SPECIFIC
Severity: LOW
Notes: The research correctly caveats that backward compatibility of additive changes relies on the consumer parsing mode (tolerant vs strict/fail-on-unknown).

---

## Claim 8
Claim: Managing provider states and test data fixtures represents an operational complexity trade-off in code-based contract testing compared to schema-only testing.
Location: `research/05-report.md: Finding 10`; `research/03-evidence.md: Evidence 14`
Evidence Provided: Pactflow article trade-off analysis.
Source: Source 6 (pactflow.io)
Source Actually Supports Claim: YES
Classification: FACT / INTERPRETATION
Severity: LOW
Notes: Accurately reflected from source text.
