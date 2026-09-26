# 01-plan.md

Research Topic
The Outbox Pattern — Database Sudah Commit, Tapi Event/Queue Gagal Dikirim

Objective
Investigate the Transactional Outbox pattern, understand how it solves the dual-write problem, cross-check claims about atomicity, at-least-once delivery, and idempotency, and produce a structured research report.

Research Questions
1. What is the dual-write problem and why do database transactions fail to solve it when dealing with external message brokers?
2. How does the Transactional Outbox pattern work to solve the dual-write problem?
3. What are the common methods for relaying messages from the outbox to the message broker?
4. Why is at-least-once delivery guaranteed, and why does it require consumer idempotency?
5. What are the operational considerations, such as payload size, outbox cleanup, and monitoring?

Search Strategy
Target official documentation, well-known software architecture resources (Microservices.io), cloud provider patterns (AWS), and CDC tools documentation (Debezium).

Expected Primary Sources
- Chris Richardson's Microservices.io
- AWS Architecture Center
- Debezium documentation

Risks / Unknowns
- Specific implementations may vary by programming language, framework, or database (e.g., polling vs. CDC).
- Evidence for some operational claims (like payload sizing) might be sparse in theoretical overviews.
