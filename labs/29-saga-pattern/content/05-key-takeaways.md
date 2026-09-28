# Key Takeaways

1. **Alternatives to 2PC:** Saga replaces Two-Phase Commit with a sequence of local transactions, avoiding blocking locks across microservices.
2. **Eventual Consistency:** Sagas guarantee eventual consistency rather than strict ACID consistency.
3. **Compensating Transactions:** Rollbacks are achieved via business-level corrective actions (compensations) rather than automatic database rollbacks.
4. **LIFO Execution:** When a step fails, compensation actions execute in reverse order (Last-In, First-Out) of successfully completed steps.
5. **Two Coordination Models:** Choose Orchestration for centralized control and visibility, or Choreography for decoupled event-driven collaboration.
6. **Lack of Isolation:** Sagas do not provide ACID isolation, exposing systems to data anomalies like dirty reads and lost updates.
7. **Semantic Locks:** Use application-level semantic locks to protect intermediate states during saga execution.
8. **Idempotency Requirement:** Retryable and compensable operations must be idempotent to handle transient network failures safely.
9. **Operational Oversight:** Permanent compensation failures require dead-letter queues, monitoring, and manual reconciliation.
