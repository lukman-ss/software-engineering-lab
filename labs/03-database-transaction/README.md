# Lab 03 — Database Transaction & External API Dangers

## 1. The Danger of Slow External API Calls Inside Transactions

Executing slow network operations (like HTTP calls to third-party payment gateways, enrichment services, or storage APIs) while holding an active database transaction is an anti-pattern that leads to severe reliability and scalability issues.

### What Happens
1. **Connection Pool Exhaustion**: Each active transaction holds a dedicated connection from the database connection pool (`*sql.DB`). If external calls take 10 seconds, those connections are locked and unavailable for other queries.
2. **Lock Contention & Deadlocks**: Database row/table locks acquired during the transaction are held for the duration of the slow HTTP call. Other transactions trying to access the same rows will block or timeout.
3. **Transaction Timeouts**: Databases often have strict idle-in-transaction limits or statement timeouts. A slow HTTP call can cause the DB to abort the transaction unilaterally.

### Alternative Architecture (Outbox Pattern / Asynchronous Processing)
- **Do not** make HTTP calls inside the DB transaction.
- **Do** record the intent and state changes in the database within a short, fast transaction (e.g., status = `pending_external`).
- **Do** use an asynchronous worker or the **Outbox Pattern** to dispatch the external HTTP call outside the transaction boundary.
- **Do** reconcile the state asynchronously via webhooks or polling.
