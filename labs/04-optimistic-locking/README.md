# Lab 04 — Optimistic Locking

Optimistic locking assumes that data contention is infrequent. Instead of locking rows preemptively, it reads the record with a `version` (or timestamp), and upon update, verifies that the version has not changed.

## Mechanism
1. Read `(data, version)`
2. Perform computation locally
3. `UPDATE table SET data = newVal, version = version + 1 WHERE id = id AND version = oldVersion`
4. If `RowsAffected == 0`, a concurrent writer modified the record. Conflict!

## Retry Policy
When a conflict occurs, the operation can be retried with jitter/backoff a limited number of times before returning an error to the caller.
