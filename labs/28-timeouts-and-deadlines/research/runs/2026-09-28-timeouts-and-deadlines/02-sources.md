# Research Sources: Timeouts and Deadlines

## Source 1

Title: Site Reliability Engineering: Chapter 22 - Addressing Cascading Failures
Publisher: Google SRE (Mike Ulrich)
URL: https://sre.google/sre-book/addressing-cascading-failures/
Published: 2016
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Engineering Standard / Primary Reference)
Relevance: Directly addresses the root causes of cascading failure from unconstrained downstream latency, thread/worker pool starvation, missing deadlines, naive retries, deadline propagation across RPC stacks, and bimodal latency distributions.

## Source 2

Title: Exponential Backoff And Jitter
Publisher: AWS Architecture Blog (Marc Brooker, VP & Distinguished Engineer)
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/
Published: 2015-03-04 (Updated 2023-05)
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Cloud Provider Architecture Documentation)
Relevance: Formal mathematical and empirical evaluation of exponential backoff vs un-jittered vs Full Jitter, Equal Jitter, and Decorrelated Jitter algorithms in preventing client synchronization spikes and retry storms.

## Source 3

Title: gRPC Deadlines Guide & Core Concepts
Publisher: gRPC Authors / Cloud Native Computing Foundation (CNCF)
URL: https://grpc.io/docs/guides/deadlines/
Published: 2025-07-07
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Framework Specification / Documentation)
Relevance: Explains client deadline configuration, automatic server cancellation, deadline propagation across multi-tier service calls, and translating absolute deadlines to relative timeouts to eliminate clock skew issues.

## Source 4

Title: PostgreSQL 18 Documentation: Chapter 19.11 - Client Connection Defaults & Statement Behavior
Publisher: The PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/runtime-config-client.html
Published: 2026-09-24
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Database Engine Documentation)
Relevance: Authoritative reference for `statement_timeout`, `transaction_timeout`, `lock_timeout`, `idle_in_transaction_session_timeout`, and `idle_session_timeout` guarding database connection pools, memory, and lock starvation.

## Source 5

Title: Stripe API Documentation: Error Handling, Retries & Idempotent Requests
Publisher: Stripe Developer Documentation
URL: https://docs.stripe.com/error-handling.md?lang=go
Published: 2026
Accessed: 2026-09-28
Source Tier: Tier 1 (Official Payment Provider API Specification)
Relevance: Defines the rule that a network timeout does not imply transaction failure (state ambiguity), how idempotency keys (`Idempotency-Key` header) safely allow retrying mutation requests, and how to reconcile status using webhooks or retrieve endpoints.
