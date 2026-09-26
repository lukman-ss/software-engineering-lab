# Source Audit

## Source 1

Claimed Title: Introduction - Pact Docs
Claimed Publisher: Pact Foundation
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
- None. Canonical official documentation defining contract testing and consumer-driven contracts.

Assessment:
PASS

---

## Source 2

Claimed Title: How Pact works
Claimed Publisher: Pact Foundation
URL: https://docs.pact.io/getting_started/how_pact_works

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Details the two-sided contract generation and verification workflow, provider states, and Message Pact.

Assessment:
PASS

---

## Source 3

Claimed Title: When to use Pact
Claimed Publisher: Pact Foundation
URL: https://docs.pact.io/getting_started/what_is_pact_good_for

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Documents operational boundaries, non-goals (e.g. public APIs with unknown consumers, functional/load testing).

Assessment:
PASS

---

## Source 4

Claimed Title: Consumer-Driven Contracts: A Service Evolution Pattern
Claimed Publisher: martinfowler.com / Ian Robinson (ThoughtWorks)
URL: https://martinfowler.com/articles/consumerDrivenContracts.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Foundational 2006 pattern paper outlining Provider Contracts, Consumer Contracts, and Consumer-Driven Contracts.

Assessment:
PASS

---

## Source 5

Claimed Title: ContractTest (bliki)
Claimed Publisher: martinfowler.com / Martin Fowler
URL: https://martinfowler.com/bliki/ContractTest.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Defines contract test as verification of test doubles against real external provider behavior.

Assessment:
PASS

---

## Source 6

Claimed Title: Contract Tests vs Functional Tests
Claimed Publisher: Pact Foundation
URL: https://docs.pact.io/consumer/contract_tests_not_functional_tests

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Documents test boundaries and anti-pattern of testing provider validation logic inside consumer contracts.

Assessment:
PASS

---

## Source 7

Claimed Title: FAQ - Pact Docs
Claimed Publisher: Pact Foundation
URL: https://docs.pact.io/faq

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Covers schema vs contract by example, expand/contract integration, and test pyramid transitions.

Assessment:
PASS

---

## Source 8

Claimed Title: CI/CD Setup Guide (Pact Nirvana)
Claimed Publisher: Pact Foundation
URL: https://docs.pact.io/pact_nirvana

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative guide on progressive CI/CD maturity (Bronze to Diamond), Pact Broker, and `can-i-deploy`.

Assessment:
PASS

---

## Source 9

Claimed Title: Parallel Change (expand and contract)
Claimed Publisher: martinfowler.com / Danilo Sato
URL: https://martinfowler.com/bliki/ParallelChange.html

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Primary pattern bliki entry for expand-migrate-contract technique across breaking interface evolutions.

Assessment:
PASS

---

## Source 10

Claimed Title: spring-cloud-contract (archived repository)
Claimed Publisher: spring-attic / Spring (GitHub)
URL: https://github.com/spring-attic/spring-cloud-contract

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
PARTIAL

Problems:
- Project was archived in July 2026. While relevant as historical context for provider-driven CDC in Spring, it no longer represents an actively maintained standard path. The research explicitly documents this limitation.

Assessment:
WARNING

---

## Source 11

Claimed Title: AsyncAPI Initiative — Docs (Concepts, Tutorials)
Claimed Publisher: AsyncAPI Initiative (Linux Foundation project)
URL: https://www.asyncapi.com/docs/tutorials and https://www.asyncapi.com/docs/concepts

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
PARTIAL

Problems:
- Source covers schema and API specification for event-driven systems, but does not specify executable consumer-driven contract test enforcement tooling directly comparable to Pact. The research notes that page bodies fetched were navigation-heavy and spec details remain unverified for test harness mechanics.

Assessment:
WARNING
