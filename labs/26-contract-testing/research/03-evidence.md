# Evidence

## Evidence 1
Claim: Contract testing checks that inter-application messages conform to a shared contract, without needing to deploy both applications together. It avoids expensive brittle integration tests.
Evidence: "Contract tests assert that inter-application messages conform to a shared understanding that is documented in a contract. Without contract testing, the only way to ensure that applications will work correctly together is by using expensive and brittle integration tests." / "Contract testing is a technique for testing an integration point by checking each application in isolation to ensure the messages it sends or receives conform to a shared understanding that is documented in a 'contract'."
Source: Pact Foundation — Introduction
URL: https://docs.pact.io/
Confidence: HIGH
Corroborated By: Martin Fowler Contract Test (https://martinfowler.com/bliki/ContractTest.html), Pactflow blog
Notes: Definition consistent across sources. Tier 1.

## Evidence 2
Claim: A Pact contract is code-first, consumer-driven, and enforced by example — a collection of concrete request/response pairs, not a static schema describing all possible states.
Evidence: "Pact is a code-first consumer-driven contract testing tool... The contract is generated during the execution of the automated consumer tests. A major advantage... only parts of the communication that are actually used by the consumer(s) get tested." / "Unlike a schema or specification (eg. OAS), which is a static artefact that describes all possible states of a resource, a Pact contract is enforced by executing a collection of test cases, each of which describes a single concrete request/response pair - Pact is, in effect, 'contract by example'."
Source: Pact Foundation — Introduction
URL: https://docs.pact.io/
Confidence: HIGH
Corroborated By: https://martinfowler.com/articles/consumerDrivenContracts.html
Notes:

## Evidence 3
Claim: Contract tests verify that calls against test doubles return same results as calls to real external service. They complement consumer tests that run against doubles.
Evidence: "In practice, a common way of implementing contract tests (and the way Pact does it) is to check that all the calls to your test doubles return the same results as a call to the real application would." / "A good way to deal with this is to continue to run your own tests against the double, but in addition to periodically run a separate set of contract tests. These check that all the calls against your test doubles return the same results as a call to the external service would."
Source: Martin Fowler — Contract Test bliki; Pact docs Introduction
URL: https://martinfowler.com/bliki/ContractTest.html
Confidence: HIGH
Corroborated By: Pact docs https://docs.pact.io/
Notes: HIGH due to two authoritative Tier 1.

## Evidence 4
Claim: Consumer-Driven Contracts pattern: consumer expectations drive provider contract. Provider contract is singular/non-authoritative derived from union of consumer expectations; consumer contracts are multiple/non-authoritative, open/incomplete.
Evidence: Table: Provider Closed/Complete/Single/Authoritative; Consumer Open/Incomplete/Multiple/Non-authoritative; Consumer-Driven Closed/Complete/Single/Non-authoritative. Text: "When a provider accepts and adopts the reasonable expectations expressed by a consumer, it enters into a consumer contract" / "provider contracts emerge to meet consumer expectations and demands... we call such provider contracts consumer-driven contracts or derived contracts."
Source: Ian Robinson / Martin Fowler — Consumer-Driven Contracts
URL: https://martinfowler.com/articles/consumerDrivenContracts.html
Published: 2006-06-12
Confidence: HIGH
Corroborated By: Pact docs describing consumer-driven
Notes: Foundational pattern definition.

## Evidence 5
Claim: Pact workflow has two phases: Consumer testing (mock provider, generate pact file) and Provider verification (replay requests against real provider, compare minimal expected response). Verification passes if response contains at least described data.
Evidence: "Consumer Pact tests operate on each interaction... does consumer code correctly generate request and handle expected response?" Steps 1-4 diagram. "In provider verification, each request is sent to the provider, and the actual response it generates is compared with the minimal expected response... Provider verification passes if each request generates a response that contains at least the data described in the minimal expected response." Plus provider states for preconditions.
Source: Pact Foundation — How Pact Works
URL: https://docs.pact.io/getting_started/how_pact_works
Confidence: HIGH
Corroborated By: Pact intro
Notes:

## Evidence 6
Claim: Contract testing applies to asynchronous messaging (Kafka, RabbitMQ, SNS/SQS, Kinesis) by abstracting protocol and focusing on message content. Pact supports message pacts with consumer handling / producer producing checks.
Evidence: "Modern distributed architectures are increasingly integrated in a decoupled, asynchronous fashion. Message queues such as ActiveMQ, RabbitMQ, SNS, SQS, Kafka and Kinesis are common... These sorts of interactions are referred to as 'message pacts'." / "Pact supports messages by abstracting away the protocol and specific queuing technology and focusses on the messages passing between them."
Source: Pact Foundation — How Pact Works (Non-HTTP testing)
URL: https://docs.pact.io/getting_started/how_pact_works
Confidence: HIGH
Corroborated By: Topic specification mentions same queues
Notes:

## Evidence 7
Claim: Contract tests should focus on messages, not provider behaviour/side-effects. Should test how validation fails, not why.
Evidence: "Contract tests should focus on the messages (requests and responses) rather than the behaviour. It can be tempting to use contract tests to write general functional tests for the provider. Experience shows this to leads to painful experiences with brittle tests." / Example: don't test each username validation rule (blank, length 20, letters), instead one generic: "When creating a user with an invalid username ... Response is 400 ... error is <any string>" / Table: "Does the provider do the right thing with the request? -> Provider's own functional tests"
Source: Pact Foundation — Contract Tests vs Functional Tests
URL: https://docs.pact.io/consumer/contract_tests_not_functional_tests
Confidence: HIGH
Corroborated By: Pactflow blog about scope
Notes:

## Evidence 8
Claim: Contract tests are fast, isolate single integration, need no dedicated environment, scale linearly, enable independent deployment.
Evidence: From comparison chart: Benefits: "Contract tests are fast. They focus on testing a single integration at a time without the need to deploy. With contract tests, there's no need for dedicated test environments. They provide fast, reliable feedback, that is easier to debug. Tests scale linearly with the number of integrations, instead of exponentially. They allow you to deploy services independently."
Source: Pactflow — Contract Testing Vs Integration Testing (SmartBear/Pactflow)
URL: https://pactflow.io/blog/contract-testing-vs-integration-testing/
Published: Updated 2023-01-04
Confidence: MEDIUM (Tier 2, reputable vendor but secondary)
Corroborated By: Pact docs Introduction ("without having to deploy the world first")
Notes: Consistent with Pact docs claims.

## Evidence 9
Claim: Test pyramid rebalanced with contract testing reduces reliance on E2E: Many Unit -> Contract -> Some Integration -> Few Critical E2E. E2E remain but minimized due to slowness/flakiness/complexity.
Evidence: Visual rebalanced pyramid: "There is more reliance on base and mid-level test techniques... doesn't remove need for E2E tests, but reduces reliance" / "E2E tests of any form can be slow... fragile... challenging to debug... expensive overhead". Chart: E2E most confidence but slowest/most cost.
Source: Pactflow — Contract Testing Vs Integration Testing
URL: https://pactflow.io/blog/contract-testing-vs-integration-testing/
Confidence: MEDIUM
Corroborated By: Widely known pyramid heuristic (Mike Cohn)
Notes: Topic spec matches same pyramid ordering.

## Evidence 10
Claim: Pact Broker provides can-i-deploy gate using matrix of consumer/provider versions and verification results; record-deployment notifies environment. Blocks breaking deployment pre-production.
Evidence: "The Pact way... is to use the Pact 'Matrix' and the can-i-deploy tool. The Matrix is the grid created when you create a table of all the consumer and provider versions that have been tested against each other." / Commands: "pact-broker can-i-deploy --pacticipant Foo --version 23 --to-environment production" exit 0 yes, 1 no / "pact-broker record-deployment --pacticipant Bar --version 56 --environment production"
Source: Pact Foundation — Can I Deploy
URL: https://docs.pact.io/pact_broker/can_i_deploy
Confidence: HIGH
Corroborated By: Pact docs pact_nirvana
Notes: Tier 1 official.

## Evidence 11
Claim: Breaking vs additive changes: adding fields is generally safe; renaming/removing/changing type is breaking. Contract testing distinguishes by minimal expected response.
Evidence: Implied by provider verification "contains at least the data described" -> additive passes. Topic spec example is canonical but supported by: "any provider behaviour not used by current consumers is free to change without breaking tests" (Pact intro) and "what if Service Team increases [username length] to 50... numbers allowed? ... Consumer should be unaffected... These are not breaking changes, but by over-specifying we are stopping Team" (Contract vs Functional). Google AIP guidance: new major version only for incompatible change.
Source: Pact docs Introduction + Contract vs Functional + Google AIP-185
URL: https://docs.pact.io/ ; https://google.aip.dev/185
Confidence: MEDIUM (combination, no single source explicitly states additive vs breaking terminology but behavior described)
Corroborated By: Consumer-Driven Contracts article discussing removing InStock field as breaking if validated
Notes: Need explicit cite for additive safe — verified via minimal response logic.

## Evidence 12
Claim: Overly strict contracts cause false failures; contracts should reflect only what consumer actually needs.
Evidence: "only parts of the communication that are actually used by the consumer(s) get tested. This in turn means that any provider behaviour not used by current consumers is free to change without breaking tests." / Also anti-pattern: "These scenarios are going too far and create an unnecessarily tight contract" + Schematron example validating only CatalogueID/Name/Price not InStock.
Source: Pact Introduction; Consumer-Driven Contracts (Schematron section)
URL: https://docs.pact.io/ ; https://martinfowler.com/articles/consumerDrivenContracts.html
Confidence: HIGH
Corroborated By: Topic spec's "Jangan membuat contract terlalu ketat"
Notes:

## Evidence 13
Claim: Documentation (OAS/Swagger) is not sufficient protection; executable contract enforced by tests is stronger.
Evidence: "The term 'contract testing', or 'provider contract testing', is sometimes used... means: ensuring provider's actual behaviour conforms to its documented contract (e.g., OpenAPI). This type of contract testing helps avoid... but on its own... does not provide assurance that consumers are calling provider correctly or provider can meet all consumers' expectations, hence not as effective in preventing integration bugs." Pact distinguishes integration contract testing vs provider contract testing.
Source: Pact Introduction
URL: https://docs.pact.io/
Confidence: HIGH
Notes:

## Evidence 14
Claim: Google API versioning requires major version in package and URI path; incompatible changes need new major version; ability to run multiple versions concurrently during transition.
Evidence: "All Google API interfaces must provide a major version number, which is encoded at the end of the protobuf package, and included as the first part of the URI path for REST APIs. In the event an API needs to make an incompatible change... Different versions of same API must be able to work at same time within a single client application for a reasonable transition period."
Source: Google AIPs — AIP-185 API Versioning
URL: https://google.aip.dev/185
Published: 2024-10-22
Confidence: HIGH (Tier 1 standards)
Corroborated By: AIP-180 Backwards compatibility
Notes: Relevant to lab exercise about evolving API with breaking changes.
