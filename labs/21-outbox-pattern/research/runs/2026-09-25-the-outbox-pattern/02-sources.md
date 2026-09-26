# 02-sources.md

## Source 1

Title: Pattern: Transactional outbox
Publisher: Microservices.io (Chris Richardson)
URL: https://microservices.io/patterns/data/transactional-outbox.html
Published: Unknown
Accessed: 2026-09-25
Source Tier: Tier 2 (Established industry publication by a domain expert)
Relevance: Defines the pattern and its context, problem, forces, solution, drawbacks, and related patterns.

## Source 2

Title: Transactional outbox pattern
Publisher: AWS Prescriptive Guidance
URL: https://docs.aws.amazon.com/prescriptive-guidance/latest/cloud-design-patterns/transactional-outbox.html
Published: Unknown
Accessed: 2026-09-25
Source Tier: Tier 1 (Official company documentation)
Relevance: Details the dual-write problem, rollback edge cases, and provides architecture implementations for both relational databases (polling publisher via SQS) and DynamoDB (CDC via DynamoDB Streams).

## Source 3

Title: Reliable Microservices Data Exchange With the Outbox Pattern
Publisher: Debezium Blog
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/
Published: February 19, 2019
Accessed: 2026-09-25
Source Tier: Tier 1 (Official project blog)
Relevance: Discusses the pattern using log-based Change Data Capture (CDC) via Debezium, transaction log tailing, aggregate event design, duplicate detection in the consumer, and table cleanup.
