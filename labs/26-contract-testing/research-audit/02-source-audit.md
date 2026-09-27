# Source Audit

## Source 1

Claimed Title: Pact Documentation - Introduction
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
- Source publication date claimed in `02-sources.md` as "Aug 25, 2026". Live inspection confirms the page footer state reads "Last updated on Aug 25, 2026 by Matt Fellows".

Assessment:
PASS

---

## Source 2

Claimed Title: Contract Test (bliki)
Claimed Publisher: Martin Fowler
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
- None. Authoritative expert bliki defining integration contract tests and test double alignment.

Assessment:
PASS

---

## Source 3

Claimed Title: Consumer-Driven Contracts: A Service Evolution Pattern
Claimed Publisher: Martin Fowler
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
- Author of article is Ian Robinson (published on martinfowler.com). `02-sources.md` lists publisher as Martin Fowler without explicitly citing Ian Robinson as primary author. Minor attribution precision issue, but content authority is intact.

Assessment:
PASS

---

## Source 4

Claimed Title: Spring Cloud Contract (archived repository)
Claimed Publisher: Spring Attic (archived by Pivotal, now VMware / Broadcom)
URL: https://github.com/spring-attic/spring-cloud-contract

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Repository is archived. It demonstrates alternative CDC tools for JVM, but is no longer actively developed under Spring Attic.

Assessment:
WARNING

---

## Source 5

Claimed Title: Pact Getting Started Guide (5 minute guide)
Claimed Publisher: Pact Foundation
URL: https://docs.pact.io/getting_started/5_minute_getting_started_guide

Reachable:
YES (redirects to https://docs.pact.io/5-minute-getting-started-guide)

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- URL path redirected slightly, but content is available and relevant.

Assessment:
PASS

---

## Source 6

Claimed Title: Pact - How Pact Works
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
- None. Details mock generation, contract JSON emission, and provider verification replay.

Assessment:
PASS

---

## Source 7

Claimed Title: Pact - Contract Tests vs Functional Tests
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
- None. Explicitly details why side-effects and deep business logic validation belong in provider functional tests, not contract tests.

Assessment:
PASS

---

## Source 8

Claimed Title: Pact Broker - Can I Deploy
Claimed Publisher: Pact Foundation
URL: https://docs.pact.io/pact_broker/can_i_deploy

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Covers Pact Matrix verification checks and CI deployment gating.

Assessment:
PASS

---

## Source 9

Claimed Title: Contract Testing Vs Integration Testing
Claimed Publisher: Pactflow / SmartBear
URL: https://pactflow.io/blog/contract-testing-vs-integration-testing/

Reachable:
YES

Source Type:
SECONDARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- Commercial vendor blog (Pactflow/SmartBear). Empirical metrics (e.g. speed/cost benefits) are promotional heuristics rather than peer-reviewed measurements.

Assessment:
WARNING

---

## Source 10

Claimed Title: Google AIP-185 - API Versioning
Claimed Publisher: Google (API Improvement Proposals)
URL: https://google.aip.dev/185

Reachable:
YES

Source Type:
PRIMARY

Relevant:
YES

Supports Claimed Topic:
YES

Problems:
- None. Authoritative guideline on breaking vs non-breaking versioning strategies in enterprise API interfaces.

Assessment:
PASS
