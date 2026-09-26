# Source Audit: Transactional Outbox Pattern

Target Lab: `labs/21-outbox-pattern`
Research Set Under Audit: `research/2026-09-26-outbox-pattern/`
Date: 2026-09-26

---

## Source 1

Claimed Title: Pattern: Transactional outbox  
Claimed Publisher: Microservices.io (Chris Richardson)  
URL: https://microservices.io/patterns/data/transactional-outbox.html  

Reachable: YES (HTTP 200)  
Source Type: PRIMARY / CANONICAL ARCHITECTURAL PATTERN  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Canonical pattern description.

Assessment: PASS

---

## Source 2

Claimed Title: Pattern: Polling publisher  
Claimed Publisher: Microservices.io (Chris Richardson)  
URL: https://microservices.io/patterns/data/polling-publisher.html  

Reachable: YES (HTTP 200)  
Source Type: PRIMARY / CANONICAL ARCHITECTURAL PATTERN  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Directly describes the polling implementation for message relay.

Assessment: PASS

---

## Source 3

Claimed Title: Pattern: Transaction log tailing  
Claimed Publisher: Microservices.io (Chris Richardson)  
URL: https://microservices.io/patterns/data/transaction-log-tailing.html  

Reachable: YES (HTTP 200)  
Source Type: PRIMARY / CANONICAL ARCHITECTURAL PATTERN  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Directly describes CDC / transaction log tailing implementation for message relay.

Assessment: PASS

---

## Source 4

Claimed Title: Outbox Event Router :: Debezium Documentation  
Claimed Publisher: Debezium Community  
URL: https://debezium.io/documentation/reference/stable/transformations/outbox-event-router.html  

Reachable: YES (HTTP 200)  
Source Type: PRIMARY (Official Engineering Documentation)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Official reference for Debezium's outbox SMT transformation, schema defaults, and routing.

Assessment: PASS

---

## Source 5

Claimed Title: Reliable Microservices Data Exchange With the Outbox Pattern  
Claimed Publisher: Debezium Blog (Gunnar Morling)  
URL: https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/  

Reachable: YES (HTTP 200)  
Source Type: PRIMARY / OFFICIAL ENGINEERING BLOG  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Detailed end-to-end walkthrough of Outbox with Postgres WAL, Quarkus, Kafka, and Idempotent Consumer.

Assessment: PASS

---

## Source 6

Claimed Title: Transactions in Apache Kafka  
Claimed Publisher: Confluent Blog (Ap行业/Jason Gustafson et al.)  
URL: https://www.confluent.io/blog/transactions-apache-kafka/  

Reachable: YES (HTTP 200)  
Source Type: PRIMARY (Official Kafka/Confluent Blog)  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Accurately demarcates Kafka's transactional boundaries (Kafka-to-Kafka EOS) from cross-system DB-to-broker transactions.

Assessment: PASS

---

## Source 7

Claimed Title: Event-Driven  
Claimed Publisher: Martin Fowler  
URL: https://martinfowler.com/articles/201701-event-driven.html  

Reachable: YES (HTTP 200)  
Source Type: PRIMARY / CANONICAL REFERENCE  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Provides architectural taxonomy of event styles (Notification, Event-carried State Transfer, Event Sourcing).

Assessment: PASS

---

## Source 8

Claimed Title: Transactional Client  
Claimed Publisher: Enterprise Integration Patterns (Gregor Hohpe & Bobby Woolf)  
URL: https://www.enterpriseintegrationpatterns.com/patterns/messaging/TransactionalClient.html  

Reachable: YES (HTTP 200)  
Source Type: PRIMARY / CANONICAL REFERENCE  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Foundational pattern reference on transactional boundaries in message publishing.

Assessment: PASS

---

## Source 9

Claimed Title: Idempotent Receiver  
Claimed Publisher: Enterprise Integration Patterns (Gregor Hohpe & Bobby Woolf)  
URL: https://www.enterpriseintegrationpatterns.com/patterns/messaging/IdempotentReceiver.html  

Reachable: YES (HTTP 200)  
Source Type: PRIMARY / CANONICAL REFERENCE  
Relevant: YES  
Supports Claimed Topic: YES  

Problems:
- None. Foundational pattern reference on handling duplicate messages.

Assessment: PASS
