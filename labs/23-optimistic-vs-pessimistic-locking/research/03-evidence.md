# Evidence: Optimistic vs Pessimistic Locking

## Evidence 1: Lost Update Definition

Claim: A lost update (write-write conflict) occurs when two transactions read the same data, both modify it based on the old value, and the second commit overwrites the first, causing the first update to be lost.

Evidence: Wikipedia defines write-write conflict as occurring when "transaction requests to write an entity for which an unclosed transaction has already made a write request." The resulting schedule is not serializable because any serial order produces a different outcome than the actual execution. (Source: Wikipedia, "Write–write conflict")

Corroborated By: PostgreSQL docs confirm this behavior under Read Committed isolation level where an UPDATE will apply to the updated version of the row, potentially overwriting concurrent changes. Oracle docs present Table 10-2 showing a concrete lost update scenario in READ COMMITTED.

Confidence: HIGH

## Evidence 2: Pessimistic Locking Mechanism — PostgreSQL

Claim: PostgreSQL implements pessimistic locking via SELECT ... FOR UPDATE, which acquires a row-level exclusive lock preventing other transactions from modifying, deleting, or locking the same row until the current transaction ends.

Evidence: PostgreSQL 18 docs state: "FOR UPDATE causes the rows retrieved by the SELECT statement to be locked as though for update. This prevents them from being locked, modified or deleted by other transactions until the current transaction ends." Other transactions attempting SELECT FOR UPDATE on the same row will be blocked until the first transaction completes.

Confidence: HIGH
Source: https://www.postgresql.org/docs/current/explicit-locking.html

## Evidence 3: Pessimistic Locking Mechanism — Oracle

Claim: Oracle Database implements row-level locks (TX locks) automatically on any row modified by INSERT, UPDATE, DELETE, MERGE, or SELECT ... FOR UPDATE.

Evidence: Oracle 19c Concepts states: "A row lock, also called a TX lock, is a lock on a single row of table. A transaction acquires a row lock for each row modified by an INSERT, UPDATE, DELETE, MERGE, or SELECT ... FOR UPDATE statement. The row lock exists until the transaction commits or rolls back."

Confidence: HIGH
Source: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html

## Evidence 4: Pessimistic Locking Mechanism — MySQL

Claim: MySQL/InnoDB implements SELECT ... FOR UPDATE and SELECT ... FOR SHARE (LOCK IN SHARE MODE) as locking reads that acquire row-level locks.

Evidence: The MySQL documentation URL was inaccessible (403), so this specific claim about MySQL's locking read syntax was not directly verified from the official source during this research session. The pattern is widely documented across secondary sources and consistent with the SQL standard.

Confidence: MEDIUM (not directly verified from official MySQL docs; corroborated by secondary sources and SQL standard)

## Evidence 5: Optimistic Locking Mechanism

Claim: Optimistic concurrency control works by checking whether data has been modified before committing. If a conflict is detected, the transaction is aborted and restarted.

Evidence: Wikipedia states: "OCC assumes that multiple transactions can frequently complete without interfering with each other. While running, transactions use data resources without acquiring locks on those resources. Before committing, each transaction verifies that no other transaction has modified the data it has read. If the check reveals conflicting modifications, the committing transaction rolls back and can be restarted." This was first proposed by H. T. Kung and John T. Robinson in 1979/1981.

Confidence: HIGH
Source: https://en.wikipedia.org/wiki/Optimistic_concurrency_control

## Evidence 6: Optimistic Locking Implementation — Version Column Pattern

Claim: The typical optimistic locking implementation uses a version column: UPDATE ... SET ..., version = version + 1 WHERE id = ? AND version = ?. If 0 rows affected, conflict detected.

Evidence: EF Core documentation confirms this pattern: "EF Core sends the following SQL to the database: UPDATE [People] SET [FirstName] = @p0 WHERE [PersonId] = @p1 AND [Version] = @p2." If a concurrent update occurred, the UPDATE fails to find any matching rows and reports zero affected rows, throwing DbUpdateConcurrencyException.

Confidence: HIGH
Source: https://learn.microsoft.com/en-us/ef/core/saving/concurrency

## Evidence 7: Database Isolation Levels — PostgreSQL

Claim: PostgreSQL default isolation level is Read Committed. Under Read Committed, each command sees a snapshot as of the command start time. Under Repeatable Read, a transaction sees a snapshot as of the transaction start time, and will be rolled back with a serialization error if it attempts to modify a row changed by another committed transaction after the transaction started.

Evidence: PostgreSQL 18 docs state: "Read Committed is the default isolation level in PostgreSQL." Under Repeatable Read: "if the first updater commits (and actually updated or deleted the row, not just locked it) then the repeatable read transaction will be rolled back with the message ERROR: could not serialize access due to concurrent update."

Confidence: HIGH
Source: https://www.postgresql.org/docs/current/transaction-iso.html

## Evidence 8: Database Isolation Levels — Oracle

Claim: Oracle offers only Read Committed (default) and Serializable isolation levels. Under Serializable, a transaction cannot modify rows changed by another transaction that committed after the serializable transaction began, raising ORA-08177.

Evidence: Oracle 19c Concepts confirms: "Oracle Database offers the read committed (default) and serializable isolation levels." Under Serializable: "The database generates an error when a serializable transaction tries to update or delete data changed by a different transaction that committed after the serializable transaction began: ORA-08177: Cannot serialize access for this transaction."

Confidence: HIGH
Source: https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html

## Evidence 9: Isolation Levels — ANSI SQL Standard

Claim: ANSI SQL defines four isolation levels: Read Uncommitted, Read Committed, Repeatable Read, Serializable. Each level prevents certain phenomena (dirty read, nonrepeatable read, phantom read, serialization anomaly).

Evidence: PostgreSQL docs Table 13.1 and Oracle docs Table 10-1 both confirm the four ANSI levels and the phenomena each prevents. PostgreSQL implements Read Uncommitted as Read Committed internally.

Confidence: HIGH
Source: https://www.postgresql.org/docs/current/transaction-iso.html and https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html

## Evidence 10: MVCC as Foundation

Claim: Both PostgreSQL and Oracle use Multiversion Concurrency Control (MVCC) to provide read consistency without blocking readers.

Evidence: PostgreSQL uses MVCC as its core concurrency control architecture. Oracle uses "multiversioning" where "the ability to simultaneously materialize multiple versions of data" allows "read-consistent queries" and "nonblocking queries." MVCC enables optimistic approaches because readers never block writers.

Confidence: HIGH
Sources: https://www.postgresql.org/docs/current/transaction-iso.html and https://docs.oracle.com/en/database/oracle/oracle-database/19/cncpt/data-concurrency-and-consistency.html

## Evidence 11: Deadlocks with Pessimistic Locking

Claim: Pessimistic locking increases the likelihood of deadlocks, which the database must detect and resolve by aborting one transaction.

Evidence: PostgreSQL docs explicitly describe deadlock scenarios: "if transaction 1 acquires an exclusive lock on table A and then tries to acquire an exclusive lock on table B, while transaction 2 has already exclusive-locked table B and now wants an exclusive lock on table A, then neither one can proceed." PostgreSQL automatically detects deadlocks and resolves them by aborting one transaction.

Confidence: HIGH
Source: https://www.postgresql.org/docs/current/explicit-locking.html

## Evidence 12: Martin Fowler — Optimistic Offline Lock

Claim: Optimistic Offline Lock validates that changes about to be committed by one session don't conflict with changes of another session. It assumes conflict is unlikely.

Evidence: Martin Fowler's pattern catalog: "Optimistic Offline Lock solves this problem by validating that the changes about to be committed by one session don't conflict with the changes of another session. A successful pre-commit validation is, in a sense, obtaining a lock indicating it's okay to go ahead with the changes to the record data."

Confidence: HIGH
Source: https://martinfowler.com/eaaCatalog/optimisticOfflineLock.html

## Evidence 13: Martin Fowler — Pessimistic Offline Lock

Claim: Pessimistic Offline Lock prevents conflicts by forcing a business transaction to acquire a lock on data before it starts using it.

Evidence: Martin Fowler's pattern catalog: "Pessimistic Offline Lock prevents conflicts by avoiding them altogether. It forces a business transaction to acquire a lock on a piece of data before it starts to use it, so that, most of the time, once you begin a business transaction you can be pretty sure you'll complete it without being bounced by concurrency control." He also notes its problem: "If several people access the same data within a business transaction, one of them will commit easily but the others will conflict and fail."

Confidence: HIGH
Source: https://martinfowler.com/eaaCatalog/pessimisticOfflineLock.html

## Evidence 14: Atomic UPDATE as Alternative

Claim: For simple operations like decrementing a counter, atomic database operations (UPDATE ... SET x = x - N WHERE condition) can eliminate race conditions without explicit locking.

Evidence: This is well-supported across database documentation as a fundamental SQL capability. PostgreSQL docs note that "UPDATE accounts SET balance = balance + 100.00 WHERE acctnum = 12345;" acquires a row-level lock automatically. Oracle docs confirm row-level locks are acquired on any UPDATE. The atomicity comes from the database executing the read and write as a single operation.

Confidence: HIGH (principle is universal; specific effectiveness depends on isolation level and WHERE clause design)

## Evidence 15: Concurrency Control Categories — Wikipedia

Claim: The three main categories of concurrency control are: Optimistic (proceed without blocking, validate at commit), Pessimistic (block if conflict may occur), and Semi-optimistic (respond based on violation type).

Evidence: Wikipedia states: "The main categories of concurrency control mechanisms are: Optimistic — Allow transactions to proceed without blocking... Pessimistic — Block an operation of a transaction, if it may cause violation... Semi-optimistic — Responds pessimistically or optimistically depending on the type of violation."

Confidence: HIGH
Source: https://en.wikipedia.org/wiki/Concurrency_control
