# Content Brief

Topic: Saga Pattern — Mengelola Transaksi Terdistribusi Tanpa 2PC
Target Reader: Software Engineers, Distributed Systems Architects, Backend Developers
Problem: Traditional ACID transactions and Two-Phase Commit (2PC) do not scale or function across multiple independently managed databases in microservices architectures due to latency, lock contention, and availability trade-offs.
Core Mental Model: A saga is a sequence of local transactions where each step updates data in a single service and triggers the next step. If a step fails, compensating transactions execute in reverse order (LIFO) to ensure eventual consistency.
Approved Research Status: APPROVED
Approved Engineering Status: APPROVED
Main Concepts:
- Local Transactions & Saga Atomicity
- Compensating Transactions (LIFO Rollback)
- Orchestration vs. Choreography
- Idempotency & Retryable Transactions
- Semantic Locking Countermeasure for Data Anomalies
Verified Behaviors:
- Sequential orchestration happy path execution
- LIFO compensation rollback on step failure
- Payment idempotency on duplicate retry
- Semantic lock prevention against concurrent pending modifications
- Choreography event-driven message dispatching
Available Case Studies: E-commerce order checkout workflow (Order -> Payment -> Inventory -> Approval)
Warnings:
- Sagas provide eventual consistency, not strict ACID isolation.
- Compensating transactions are business-level corrective actions and may occasionally fail, requiring operational recovery mechanisms (DLQs, manual reconciliation).
