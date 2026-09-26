# Sources

Research date: 2026-09-26. All sources opened and inspected via full fetch, not snippets.

## Source 1

Title: Introduction - Pact Docs
Publisher: Pact Foundation
URL: https://docs.pact.io/
Published: Last updated Aug 25, 2026
Accessed: 2026-09-26
Source Tier: Tier 1 (official tool documentation)
Relevance: Core definition of contract testing, consumer-driven contracts, Pact vs schema testing distinction.

## Source 2

Title: How Pact works
Publisher: Pact Foundation
URL: https://docs.pact.io/getting_started/how_pact_works
Published: Last updated Dec 23, 2024
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Consumer test + provider verification mechanics, provider states, Message Pact for async systems.

## Source 3

Title: When to use Pact
Publisher: Pact Foundation
URL: https://docs.pact.io/getting_started/what_is_pact_good_for
Published: Last updated Apr 13, 2022
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Scope and limitations — when Pact is / is not appropriate (public APIs, pass-through APIs, performance).

## Source 4

Title: Consumer-Driven Contracts: A Service Evolution Pattern
Publisher: martinfowler.com / Ian Robinson (ThoughtWorks)
URL: https://martinfowler.com/articles/consumerDrivenContracts.html
Published: 12 June 2006
Accessed: 2026-09-26
Source Tier: Tier 1 (pattern originator, primary historical evidence)
Relevance: Foundational theory: provider vs consumer vs consumer-driven contracts, schema versioning, breaking changes, just-enough validation.

## Source 5

Title: ContractTest (bliki)
Publisher: martinfowler.com / Martin Fowler
URL: https://martinfowler.com/bliki/ContractTest.html
Published: 12 January 2011 (revised 2018-01-01)
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Canonical contract test definition: test doubles vs real service, cadence, communication aspect.

## Source 6

Title: Contract Tests vs Functional Tests
Publisher: Pact Foundation
URL: https://docs.pact.io/consumer/contract_tests_not_functional_tests
Published: Last updated Mar 2, 2022
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Boundary between contract and functional testing; over-specification anti-pattern (validation rules example).

## Source 7

Title: FAQ - Pact Docs
Publisher: Pact Foundation
URL: https://docs.pact.io/faq
Published: Last updated Oct 21, 2025
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Breaking-change workflow (expand/contract), versioning, E2E replacement guidance, DORA metrics, can-i-deploy, GraphQL/SOAP/Protobuf stance.

## Source 8

Title: CI/CD Setup Guide (Pact Nirvana)
Publisher: Pact Foundation
URL: https://docs.pact.io/pact_nirvana
Published: Last updated Jan 13, 2025
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: CI integration stages (Bronze → Diamond), Pact Broker role, independent deployability goal.

## Source 9

Title: Parallel Change (expand and contract)
Publisher: martinfowler.com / Danilo Sato
URL: https://martinfowler.com/bliki/ParallelChange.html
Published: 13 May 2014
Accessed: 2026-09-26
Source Tier: Tier 1
Relevance: Three-phase pattern (expand → migrate → contract) for safe breaking changes; API evolution alternative to versioning.

## Source 10

Title: spring-cloud-contract (archived repository)
Publisher: spring-attic / Spring (GitHub)
URL: https://github.com/spring-attic/spring-cloud-contract
Published: Archived by owner Jul 7, 2026; maintenance moved to Stubborn.sh
Accessed: 2026-09-26
Source Tier: Tier 1 (official company documentation, now historical)
Relevance: Documents provider-driven CDC alternative in Spring ecosystem; confirms project end-of-life — freshness risk for labs recommending it.

## Source 11

Title: AsyncAPI Initiative — Docs (Concepts, Tutorials)
Publisher: AsyncAPI Initiative (Linux Foundation project)
URL: https://www.asyncapi.com/docs/tutorials and https://www.asyncapi.com/docs/concepts
Published: Living docs (no single date); site self-describes as "Building the future of Event-Driven Architectures"
Accessed: 2026-09-26
Source Tier: Tier 1 (official specification community)
Relevance: Event-driven contract analogue to OpenAPI; supports lab claim that contracts apply to Kafka/RabbitMQ/webhooks. Page bodies fetched were navigation-heavy; spec detail NOT VERIFIED beyond positioning.
