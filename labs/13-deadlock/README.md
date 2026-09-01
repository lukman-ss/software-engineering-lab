# Lab 13 — Deadlock Prevention & Resolution

A deadlock occurs when two or more transactions permanently block each other by holding locks the other transactions need.

## 1. Deadlock Reproduction
Happens when Transaction 1 locks Account A then tries to lock Account B, while Transaction 2 concurrently locks Account B then tries to lock Account A.

## 2. Prevention via Deterministic Lock Ordering
Always acquire locks in a globally consistent order (e.g., ascending order of account IDs). This breaks the circular wait condition required for deadlocks.

## 3. Transient Error Retry Policy
- **Detect**: Check error codes (Postgres `40P01` Deadlock Detected, `40001` Serialization Failure).
- **Retry**: Use bounded retries with exponential backoff and jitter.
- **Do Not Retry**: Permanent errors like insufficient funds or syntax errors.

---

## Navigasi

- **Previous**: [Lab 12 — Optimistic Locking](../12-optimistic-locking/)
