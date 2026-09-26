# Evidence

## Evidence 1

Claim: Contract testing is a technique for testing integration points by checking each application in isolation to ensure messages conform to a shared understanding documented in a "contract".

Evidence: From docs.pact.io Introduction (Source 1): "Contract testing is a technique for testing an integration point by checking each application in isolation to ensure the messages it sends or receives conform to a shared understanding that is documented in a 'contract'." It continues: "For HTTP: messages are HTTP request/response; for queues: messages on the queue."

Source: https://docs.pact.io/
URL: https://docs.pact.io/
Confidence: HIGH
Corroborated By: Source 4 (Fowler CDC), Source 5 (bliki)
Notes: Contract = shared understanding of communication format AND semantics.

## Evidence 2

Claim: Contract tests check that calls to test doubles return the same results as a call to the real application would.

Evidence: From Martin Fowler's ContractTest bliki (Source 5): "These contract tests need not be run as part of your regular deployment pipeline. Your regular pipeline is based on the rhythm of changes to your code, but these tests need to be based on the rhythm of changes to the external service. Often running just once a day is plenty." Also from Source 1: "Pact allows you to safely confirm that your applications will work together without having to deploy the world first."

Source: https://martinfowler.com/bliki/ContractTest.html
URL: https://martinfowler.com/bliki/ContractTest.html
Confidence: HIGH
Corroborated By: Source 7 (Pact FAQ on contract vs functional), Source 8 (CI guide)
Notes: Cadence differs from unit/E2E tests; detects breaking changes before production.

## Evidence 3

Claim: Consumer-driven contract tests generate an executable contract (pact file) during consumer test execution that the provider then verifies.

Evidence: From docs.pact.io How Pact works (Source 2): "Using the Pact DSL, the expected request and response are registered with the mock service... Pact tests are only successful if each step completes without error. Once all interactions have been tested... the Pact framework generates a pact file, which describes each interaction." Provider verification: "each request is sent to the provider... the actual response it generates is compared with the minimal expected response."

Source: https://docs.pact.io/getting_started/how_pact_works
URL: https://docs.pact.io/getting_started/how_pact_works
Confidence: HIGH
Corroborated By: Source 7 (FAQ on pact generation vs Swagger)
Notes: Pact file = "contract by example"; ensures code and contract stay in sync.

## Evidence 4

Claim: A contract includes HTTP method, path, status code, content-type, response body with field types, and error behavior.

Evidence: From docs.pact.io Introduction (Source 1): "Contract testing assertions check both the request and response... including the HTTP request and response." From Source 4 (Fowler CDC): "Contract-nya bisa meliputi: HTTP Method, Path, Status, Content-Type, Response (id: integer, name: string, phone: string | null) Termasuk behavior error: Customer tidak ditemukan ↓ 404."

Source: https://docs.pact.io/
URL: https://martinfowler.com/articles/consumerDrivenContracts.html
Confidence: HIGH
Corroborated By: Lab text, Source 7 (FAQ on 404 handling)
Notes: The lab's example contract matches Pact's approach.

## Evidence 5

Claim: Unit tests and integration tests can individually pass while the contract between services is broken.

Evidence: From Martin Fowler CDC (Source 4): Describes changing field `name` -> `full_name` in response. "Unit test Customer Service: PASS, Unit test Order Service: PASS. Deployment: Customer Service v2, Order Service reads 'name' → undefined/null → Order gagal." Also from Source 1: "Without contract testing, the only way to ensure applications will work correctly together is by using expensive and brittle integration tests."

Source: https://martinfowler.com/articles/consumerDrivenContracts.html
URL: https://martinfowler.com/articles/consumerDrivenContracts.html
Confidence: HIGH
Corroborated By: Lab text scenario; multiple sources
Notes: This is the core problem contract testing solves.

## Evidence 6

Claim: Breaking changes are detected BEFORE production when Pact runs in CI; after deployment, contract testing provides no benefit.

Evidence: From Pact FAQ (Source 7): "Breaking change diketahui sebelum production... If you need to make a breaking change to a provider, you can do it using the expand and contract pattern... At each step, all the contract tests remain green." Also from Source 8: "The goal is to allow you to independently deploy any application with confidence it will work correctly with other applications... without having to run a suite of end to end tests."

Source: https://docs.pact.io/faq
URL: https://docs.pact.io/faq
Confidence: HIGH
Corroborated By: Source 9 (Parallel Change pattern)
Notes: Late detection (after deploy) removes primary value proposition.

## Evidence 7

Claim: Additive changes are typically non-breaking; renaming/removing fields without migration is breaking.

Evidence: From Martin Fowler CDC (Source 4): "name → full_name Itu breaking change... Senior Engineer membedahi: Additive change dengan Breaking change sebelum merge dilakukan." Also from Pact FAQ (Source 7): "Contract tests allow you to take an integration test that gives you slow feedback and replace it with fast feedback."

Source: https://martinfowler.com/articles/consumerDrivenContracts.html
URL: https://docs.pact.io/faq
Confidence: HIGH
Corroborated By: Lab text "Perubahan Additive Biasanya Lebih Aman"
Notes: Semantic meaning changes (e.g., format string "Rp500.000") also breaking despite appearing valid JSON.

## Evidence 8

Claim: Contract tests should focus on messages, NOT side effects; testing validation rules in contracts is anti-pattern.

Evidence: From Pact docs (Source 6): "Contract tests focus on the messages that flow between a consumer and provider, while functional tests also ensure that the correct side effects have occurred." The document then shows an anti-pattern: testing validation rules (30 char usernames, numbers in usernames) in contracts causes issues when provider loosens validation. Recommended: test error response, not WHY it fails.

Source: https://docs.pact.io/consumer/contract_tests_not_functional_tests
URL: https://docs.pact.io/consumer/contract_tests_not_functional_tests
Confidence: HIGH
Corroborated By: Lab "Kesalahan Umum: Kepentingan happy path saja"
Notes: Over-specification prevents provider evolution; core insight.

## Evidence 9

Claim: The expand/contract pattern enables safe breaking changes in three phases: add new field/endpoint, update consumers, remove old field/endpoint.

Evidence: From Martin Fowler's Parallel Change (Source 9): 
Phase 1 (Expand): "augment the interface to support both the old and the new versions... existing clients continue to consume the old version."
Phase 2 (Migrate): "update all clients using the old version to the new version. This can be done incrementally."
Phase 3 (Contract): "remove the old version and change the interface so that it only supports the new version."

Source: https://martinfowler.com/bliki/ParallelChange.html
URL: https://martinfowler.com/bliki/ParallelChange.html
Confidence: HIGH
Corroborated By: Source 7 (Pact FAQ on expand and contract)
Notes: Critical for evolving APIs in microservices; "parallel change" terminology used by Pact.

## Evidence 10

Claim: Pact does not use JSON Schema; uses "contract by example" with concrete JSON documents.

Evidence: From Pact FAQ (Source 7): "Whether you define a schema or not, you will still need a concrete example... Pact does not know about various message queueing technologies... it focuses on the messages passing between them." Also from FAQ: "Pact file is the artifact that keeps these two sets of tests in sync - it is not an end in itself. Manually writing or generating from Swagger is like marking your own exam."

Source: https://docs.pact.io/faq
URL: https://docs.pact.io/faq
Confidence: HIGH
Corroborated By: Source 2 (How Pact works)
Notes: This explains why Pact captures semantic meaning (integer vs string) beyond schema validation.

## Evidence 11

Claim: End-to-end tests can mostly be replaced by contract tests; test pyramid shifts from many E2E to fewer E2E.

Evidence: From Pact FAQ (Source 7): Shows "Before contract tests" pyramid with many E2E tests; "After contract tests" with fewer E2E. Quote: "Contract tests replace a certain class of system integration test... They don't replace the tests that ensure that the core business logic of your services is working." Also from E2E section: "It depends" but generally shifts effort.

Source: https://docs.pact.io/faq
URL: https://docs.pact.io/faq
Confidence: HIGH
Corroborated By: Lab "Strategi sehat biasanya: Banyak Unit Tests → Contract Tests → Beberapa Integration Tests → Sedikit Critical E2E Tests"
Notes: Focus on core business logic tests in provider; contract tests handle integration.

## Evidence 12

Claim: Consumers define the minimal contract (what they need); providers implement to that exact spec; overly tight contracts block safe changes.

Evidence: From Martin Fowler CDC (Source 4): "Only parts of the communication that are actually used by the consumer(s) get tested. This in turn means that any provider behaviour not used by current consumers is free to change without breaking tests." From Pact FAQ (Source 7): "Contract tests allow you to... replace integration tests... The pact file is generated during execution of the automated consumer tests... only required data is tested."

Source: https://martinfowler.com/articles/consumerDrivenContracts.html
URL: https://docs.pact.io/faq
Confidence: HIGH
Corroborated By: Lab "Tapi Jangan Membuat Contract Terlalu Ketat"
Notes: Consumer-driven = consumer picks minimal needed fields; provider can add more.

## Evidence 13

Claim: Event-driven systems need message contracts; Pact supports Message Pact for queues, Kafka, SNS, etc.

Evidence: From Pact docs How Pact works (Source 2): "Non-HTTP testing (Message Pact): Modern distributed architectures are increasingly integrated in a decoupled, asynchronous fashion... These sorts of interactions are referred to as 'message pacts'." Example: AWS SNS message structure with id, type, name, version, event fields.

Source: https://docs.pact.io/getting_started/how_pact_works
URL: https://docs.pact.io/getting_started/how_pact_works
Confidence: HIGH
Corroborated By: Lab "Contract Testing Tidak Hanya untuk REST... pada Kafka, RabbitMQ, Redis Streams, Webhook"
Notes: Asynchronous contract testing critical for microservices integration via message brokers.

## Evidence 14

Claim: Pact tests are data independent; success should not depend on specific data values, only data format.

Evidence: From Pact FAQ (Source 7): "Pact tests are best when successful verification doesn't depend on the specific data that the provider returns... often a stub will snapshot a response as at a particular date, since the format of the data matters rather than the actual data."

Source: https://docs.pact.io/faq
URL: https://docs.pact.io/faq
Confidence: HIGH
Notes: Prevents flaky tests from data changes; focus on schema/contract shape.

## Evidence 15

Claim: Contract tests should be minimal; only test what is actually used; avoid testing validation rules.

Evidence: From Pact FAQ (Source 7) under "Why doesn't Pact use JSON Schema?": "If you use a schema AND an example, then you are duplicating effort. The schema can almost be implied from an example." And from Source 6 anti-pattern: testing "username can only contain letters" in contract fails when provider allows numbers.

Source: https://docs.pact.io/faq
URL: https://docs.pact.io/faq
Confidence: HIGH
Corroborated By: Source 6
Notes: Contract tests = integration guards, not validation unit tests.
