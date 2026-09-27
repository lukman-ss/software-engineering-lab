# Contradiction Audit

Target Lab: `labs/21-outbox-pattern`

---

## Internal & Source Consistency Analysis

### 1. Research Report vs Evidence Files
- `research/05-report.md` matches `research/03-evidence.md` across all findings (Dual-write problem, Atomicity guarantees, Relay alternatives, Idempotency, Schema design, Monitoring).

### 2. Research Report vs Primary Sources
- **CDC vs Polling**: `microservices.io` treats Polling Publisher and Transaction Log Tailing as separate options. Debezium docs focus heavily on log tailing via WAL. `05-report.md` correctly synthesizes both as valid implementations of the Message Relay component without presenting CDC as the only way.
- **Idempotency Responsibility**: `microservices.io` places idempotency on the consumer. Debezium blog shows both consumer-side deduplication table (`MessageLog`) and immediate outbox row deletion (`persist` + `remove` in single DB transaction). `05-report.md` accurately notes that both mechanisms can coexist.

### 3. Claimed Semantic Guarantees
- No claims of "magical exactly-once delivery" without idempotent handling. The report explicitly specifies that outbox delivers **at-least-once** semantics.

---

## Conclusion
No material contradictions found.
