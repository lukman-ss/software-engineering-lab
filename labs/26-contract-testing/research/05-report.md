# Research Report

## Research Question

How does contract testing prevent integration failures in distributed systems when individual service tests pass, and what patterns, tools, and anti-patterns define effective implementation?

## Executive Summary

Contract testing solves the fundamental problem of distributed systems where individual services have passing unit and integration tests but fail to interoperate due to broken contracts (e.g., field renaming, type changes, or altered error semantics). Consumer-Driven Contract Testing (CDC), pioneered by Martin Fowler and Ian Robinson (2006) and implemented by Pact, shifts contract definition from providers to consumers: consumers define minimal expectations as executable tests, generating a "pact" file that providers verify in CI before deployment. This detects breaking changes pre-production. Unlike schema testing or functional tests, contract tests validate only the messages actually used (request/response, field types, status codes) and avoid over-specification that blocks safe provider evolution. The expand/contract pattern enables safe breaking changes, and Pact Broker integrates contracts into CI/CD pipelines for independent deployability. Event-driven systems (Kafka, webhooks) use Message Pact with identical principles.

## Findings

### Finding 1: Contract testing validates shared message understanding between services, not internal behavior.

**Evidence:** Contract testing is "a technique for testing an integration point by checking each application in isolation to ensure the messages it sends or receives conform to a shared understanding that is documented in a 'contract'" (Source 1). For HTTP, this is request/response; for queues, messages on the queue (Source 2). A contract test uses a test double; a failure indicates the test double no longer matches the real service, requiring updates and possibly a conversation with the service owners (Source 5). 

**Sources:** Source 1 (docs.pact.io), Source 2 (How Pact works), Source 5 (ContractTest bliki)
**Confidence:** HIGH

### Finding 2: Consumer-Driven Contracts shift contract definition to consumers, ensuring only used fields are tested.

**Evidence:** "Only parts of the communication that are actually used by the consumer(s) get tested. This in turn means that any provider behaviour not used by current consumers is free to change without breaking tests" (Source 4). Pact generates contract files during consumer test execution; verification compares provider responses to the consumer's minimal expected response (Source 2). Provider-driven contracts test the entire schema, blocking additive changes; CDC tests only what consumers actually use, enabling provider evolution (Source 4, Source 6).

**Sources:** Source 4 (Consumer-Driven Contracts), Source 2 (How Pact works), Source 6 (Contract vs Functional Tests)
**Confidence:** HIGH

### Finding 3: Unit and integration tests can pass while contracts break; contract testing detects this pre-deployment.

**Evidence:** The lab's Order Service/Customer Service example demonstrates that refactoring `name`→`full_name` causes individual service tests to pass but integration to fail in production (Lab text). Contract testing catches this before deployment: provider verification fails if the response does not match the consumer's expected contract (Source 2). As Martin Fowler states, unit tests verify Function A → Output A; integration tests may test database dependencies; but distributed systems have a boundary where Consumer expectation must equal Provider behavior (Source 4). Contract testing focuses on this boundary.

**Sources:** Lab text (Kasus Nyata), Source 2, Source 4
**Confidence:** HIGH

### Finding 4: A contract includes HTTP semantics (method, path, status, headers), field names/types, and error behavior—not just JSON schema.

**Evidence:** Contracts cover "HTTP Method, Path, Status, Content-Type, Response (id: integer, name: string, phone: string | null) Termasuk behavior error: Customer tidak_found ↓ 404" (Source 4). Pact tests verify status codes, headers, and response bodies; providers must return at least the minimal expected response (Source 2). Semantic meaning matters: changing `total` from integer 450000 to string "Rp450.000" breaks consumers doing arithmetic, even though JSON is valid (Lab text "Contoh Kasus Frontend"). Contract tests catch type changes; JSON schema alone may not.

**Sources:** Source 4, Source 2, Lab text
**Confidence:** HIGH

### Finding 5: Additive changes (new fields) are typically non-breaking; renaming/removing fields or changing types without migration is breaking.

**Evidence:** Adding email to a response `{id, name}` → `{id, name, email}` is additive and usually safe (Lab text "Perubahan Additive Biasanya Lebih Aman"). Renaming `name`→`full_name`, removing `InStock`, or changing `total` from integer to string are breaking changes because consumers depend on the exact structure (Source 4). Martin Fowler CDC states Senior Engineers distinguish additive vs breaking changes before merge. Contract tests fail on breaking changes but pass on additive ones if consumers don't require the new field.

**Sources:** Source 4, Lab text
**Confidence:** HIGH

### Finding 6: Contract tests should focus on message format and error handling, not provider validation rules (anti-pattern).

**Evidence:** Testing validation rules in contracts (e.g., "username max 20 chars", "letters only") creates over-specification that blocks safe provider evolution (Source 6). If provider loosens validation (increases max to 50, allows numbers), contracts fail despite no consumer impact. Recommended: test error responses exist (400 Bad Request) with any error message, not specific validation logic (Source 6). Contract tests should catch: consumer bugs, consumer misunderstanding of endpoints/payload, and provider breaking changes on endpoints/payload—not provider business logic (Source 6).

**Sources:** Source 6 (Contract Tests vs Functional Tests)
**Confidence:** HIGH

### Finding 7: The expand/contract pattern enables safe breaking changes via three phases.

**Evidence:** To make a breaking change (e.g., rename field): 
1. Expand: Add new field/endpoint alongside old; deploy provider.
2. Migrate: Update consumers to use new field; deploy consumers.
3. Contract: Remove old field/endpoint; deploy provider. 
At each step, contract tests remain green if consumers are updated (Source 9, Source 7). This pattern is "particularly useful when practicing Continuous Delivery" and avoids breakage across the entire codebase (Source 9). Pact FAQ explicitly recommends this approach for breaking changes (Source 7).

**Sources:** Source 9 (Parallel Change), Source 7 (Pact FAQ on breaking changes)
**Confidence:** HIGH

### Finding 8: Pact Broker enables CI/CD integration for independent deployability.

**Evidence:** Pact Broker is a "permanently running, externally hosted service with an API and UI that allows you exchange the pacts and verification results" (Source 7). CI/CD integration progresses through levels: 
- Bronze: Manual test + mock service
- Silver: Manual Pact Broker exchange
- Gold: PR pipeline verification
- Platinum: PR pipeline + can-i-deploy with branch tag
- Diamond: Deploy pipeline verification 
This enables teams to "independently deploy any application with the confidence that it will work correctly with the other applications in its environment" (Source 8). Provider verification results can be published back to the broker; consumers check `can-i-deploy` before release (Source 7).

**Sources:** Source 8 (CI/CD Setup Guide), Source 7 (FAQ on Broker, can-i-deploy)
**Confidence:** HIGH

### Finding 9: Contract testing applies to event-driven systems via Message Pact.

**Evidence:** Message Pact supports asynchronous integrations: "Message queues such as ActiveMQ, RabbitMQ, SNS, SQS, Kafka and Kinesis are common... Pact supports messages by abstracting away the protocol and specific queuing technology" (Source 2). Consumer-side: tests handling a message payload (e.g., AWS SNS `id`, `type`, `name`, `version`, `event`). Provider-side: tests producing the correct message structure. Adapter/Port separation isolates protocol-specific code from domain logic (Source 2). This validates lab's claim: "Contract Testing Tidak Hanya untuk REST... sangat relevan pada Kafka, RabbitMQ, Redis Streams, Webhook, Event Bus."

**Sources:** Source 2 (How Pact works, Non-HTTP testing section)
**Confidence:** HIGH

### Finding 10: Contract tests reduce but do not eliminate end-to-end tests; the test pyramid shifts focus.

**Evidence:** Contract tests replace "a certain class of system integration test" (e.g., validating API usage/response) but not tests for "core business logic of your services" (Source 7). The FAQ shows a test pyramid shifting from many E2E tests to fewer, targeted E2E tests after contract test adoption (Source 7). The lab recommends: "Banyak Unit Tests → Contract Tests → Beberapa Integration Tests → Sedikit Critical E2E Tests." Contract tests provide fast feedback; E2E tests validate critical user journeys in production-like environments.

**Sources:** Source 7 (Pact FAQ on E2E tests), Source 8 (CI guide)
**Confidence:** HIGH

## Areas of Agreement

All sources agree on:
- Contract testing's purpose: detecting pre-production integration failures from broken service contracts.
- Consumer-Driven Contracts as superior to provider-driven contracts for enabling evolution.
- Pact's mechanism: consumer tests generate pact file; provider verifies against real service.
- The expand/contract pattern for safe breaking changes.
- Contract tests validate messages (request/response), not provider internal behavior or side effects.
- Over-specification in contracts (testing validation rules) is an anti-pattern.
- Event-driven systems require analogous message contract testing.

## Areas of Disagreement

No substantive disagreements exist between sources. Minor nuances:
- Pact documentation emphasizes its applicability where consumer/provider teams collaborate and control data (Source 3, 7); the lab and Martin Fowler CDC present the technique more universally. These are contextual, not contradictory.
- The lab's exercise analysis (determining breaking changes for three specific changes) is a practical application of principles universally agreed upon.

## Limitations

1. Contract testing requires consumer/provider team collaboration and shared CI/CD pipeline access (Source 7). It is less suited for public APIs where consumers are unknown.
2. Contract tests do not validate provider business logic or data correctness; providers must maintain their own functional and unit tests (Source 6, 7).
3. Test maintenance overhead exists for each consumer-provider pair; managing many consumers can strain provider teams (Source 3, 7).
4. Contract tests alone cannot detect all integration issues (e.g., performance, downtime, complex workflows); targeted E2E or synthetic monitoring complements them (Source 7, 8).
5. The lab's specific exercise (three breaking changes) assumes consumers strictly depend on exact field names/types; real-world tolerance (e.g., lenient JSON parsing) may vary but does not invalidate the principle.

## Conclusion

Contract testing prevents silent integration failures by executable validation of the message contract between consumers and providers. Originating in the Consumer-Driven Contract pattern (Fowler/Robinson, 2006) and implemented by Pact, it shifts contract definition to consumers, ensuring only actually-used fields are tested. This detects breaking changes (field renames, type removals, semantic changes) before deployment via provider verification in CI. Contracts encompass HTTP semantics, field types, status codes, and error behavior—not just JSON schema. The expand/contract pattern enables safe evolution, and Pact Broker integrates contracts into CI/CD for independent deployability. While contract testing reduces the need for brittle end-to-end tests targeting API contracts, it complements (does not replace) tests for core business logic and critical user journeys. Effective implementation avoids over-specification, focuses on minimal consumer needs, and applies equally to REST and event-driven systems via Message Pact.