# Source Audit — Lab 26: Contract Testing

## Source 1

Claimed Title: Introduction | Pact Docs
Claimed Publisher: Pact Foundation (docs.pact.io)
URL: https://docs.pact.io/

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Publication date in source documentation states "Last updated on Aug 25, 2026 by Matt Fellows" and "Copyright © 2026 Pact Foundation". This matches live site text verbatim.
- Citation accurately quotes definitions of contract tests, test double equivalence, message queues, and consumer-driven contracts.

Assessment:
PASS

---

## Source 2

Claimed Title: Contract Test
Claimed Publisher: Martin Fowler (martinfowler.com/bliki)
URL: https://martinfowler.com/bliki/ContractTest.html

Reachable:
YES

Source Type:
SECONDARY (Expert Technical Publication / Canonical Industry Bliki)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Martin Fowler's 2011/2018 bliki post accurately supports the role of test doubles, running against external services, and evolutionary communication with provider teams.

Assessment:
PASS

---

## Source 3

Claimed Title: Consumer-Driven Contracts: A Service Evolution Pattern
Claimed Publisher: Martin Fowler / Ian Robinson (martinfowler.com/articles)
URL: https://martinfowler.com/articles/consumerDrivenContracts.html

Reachable:
YES

Source Type:
PRIMARY (Canonical defining publication for Consumer-Driven Contracts pattern)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. 2006 Ian Robinson / Martin Fowler article directly defines provider contracts, consumer contracts, and consumer-driven contracts.

Assessment:
PASS

---

## Source 4

Claimed Title: What is contract testing and why should I try it?
Claimed Publisher: Pactflow (Matt Fellows)
URL: https://pactflow.io/blog/what-is-contract-testing/

Reachable:
YES

Source Type:
SECONDARY (Reputable Vendor Technical Blog)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Page metadata shows update date 2 September 2023 (with original ~2019/2021 graphics). The research accurately noted this update date.
- Content accurately supports problems of integrated E2E tests, fast feedback, test pyramid positioning, and independent execution.

Assessment:
PASS

---

## Source 5

Claimed Title: What is consumer driven contract testing?
Claimed Publisher: Pactflow
URL: https://pactflow.io/what-is-consumer-driven-contract-testing

Reachable:
YES

Source Type:
SECONDARY (Reputable Vendor Technical Explainer)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Undated on page as correctly acknowledged by research agent.
- Accurately details explicit contracts vs implicit contracts, contract verification in provider CI builds, and preventing breaking releases.

Assessment:
PASS

---

## Source 6

Claimed Title: Schema-based contract testing with JSON schemas and Open API (Part 1)
Claimed Publisher: Pactflow (Matt Fellows)
URL: https://pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1

Reachable:
YES

Source Type:
SECONDARY (Reputable Vendor Technical Blog)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Publication date on live site is "Updated 5 January 2023" (headline updated date). Research report cites "May 30, 2023" in footer list while `02-sources.md` correctly says "Updated 5 January 2023". Minor internal typo between research files, but URL and substantive points match exact text.

Assessment:
PASS

---

## Source 7

Claimed Title: Lab 06 — API Versioning: Cara Mengubah API Tanpa Merusak Ribuan Client (internal cross-reference)
Claimed Publisher: Lukman SS / software-engineering-lab
URL: /labs/06-api-versioning/README.md (local repo, no external URL)

Reachable:
YES (Local file exists in workspace)

Source Type:
COMMUNITY / INTERNAL (Workspace-internal pedagogical reference)

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Internal repository reference, not an external industry standard. However, the research correctly explicitly disclosed this as an internal reference for taxonomy alignment between labs.

Assessment:
PASS

---

## Verification of Unsuccessful Fetches Noted in Research

- `https://docs.pact.io/getting_started/what_is_pact`: confirmed 404/redirected in reorganization.
- `https://docs.pact.io/getting_started/testing_scope`: confirmed 404/restructured.
The research agent transparently recorded these failures rather than hallucinating content.
