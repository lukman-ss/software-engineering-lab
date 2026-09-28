# Research Evidence: Saga Pattern

## Evidence 1: Saga Menggantikan 2PC untuk Database-Per-Service

**Claim:** Traditional database guarantees (ACID) dan 2PC tidak applicable untuk multiple independently managed data stores di microservices architecture

**Evidence:** "Because of these limitations, architectures that rely on interprocess communication, or traditional transaction models like two-phase commit protocol, are often better suited for the Saga pattern." — Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

**Corroborated By:**
- Chris Richardson / Microservices.io: "You have applied the Database per Service pattern... the application cannot simply use a local ACID transaction" (https://microservices.io/patterns/data/saga.html)

---

## Evidence 2: Saga = Sequence of Local Transactions

**Claim:** Saga memecah transaksi terdistribusi menjadi sequence of local transactions, dimana setiap local transaction meng-update database dan memicu transaction berikutnya via event/message

**Evidence:** "The Saga pattern manages transactions by breaking them into a sequence of local transactions. Each local transaction: 1. Completes its work atomically within a single service. 2. Updates the service's database. 3. Initiates the next transaction via an event or message." — Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

**Corroborated By:**
- Chris Richardson / Microservices.io: "A saga is a sequence of local transactions. Each local transaction updates the database and publishes a message or event to trigger the next local transaction in the saga." (https://microservices.io/patterns/data/saga.html)

---

## Evidence 3: Compensating Transactions vs Rollback

**Claim:** Compensating transactions meng-undo perubahan dari local transactions yang sudah completed, TAPI bukan simple rollback - ini adalah business-level corrective action

**Evidence:** "If a local transaction fails, the saga performs a series of compensating transactions to reverse the changes that the preceding local transactions made." — Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

**Corroborated By:**
- Chris Richardson: "lack of automatic rollback - a developer must design compensating transactions that explicitly undo changes made earlier in a saga rather than relying on the automatic rollback feature of ACID transactions" (https://microservices.io/patterns/data/saga.html)

---

## Evidence 4: Dua Pendekatan Implementasi (Choreography vs Orchestration)

**Claim:** Dua pendekatan implementasi saga: Choreography (event-driven, loose coordination) dan Orchestration (centralized coordinator)

**Evidence:**
- "Choreography: In the choreography approach, services exchange events without a centralized controller."
- "Orchestration: In orchestration, a centralized controller, or orchestrator, handles all the transactions and tells the participants which operation to perform based on events."
— Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

**Corroborated By:**
- Chris Richardson: "There are two ways of coordination sagas: Choreography - each local transaction publishes domain events that trigger local transactions in other services. Orchestration - an orchestrator (object) tells the participants what local transactions to execute" (https://microservices.io/patterns/data/saga.html)

---

## Evidence 5: Three Types of Transactions (Compensable, Pivot, Retryable)

**Claim:** Saga terdiri dari tiga jenis transaksi: Compensable (bisa di-undo), Pivot (point of no return), dan Retryable (idempotent, bisa di-retry)

**Evidence:**
- **Compensable transactions:** "can be undone or compensated for by other transactions with the opposite effect"
- **Pivot transactions:** "serve as the point of no return in the saga. After a pivot transaction succeeds, compensable transactions are no longer relevant"
- **Retryable transactions:** "follow the pivot transaction. Retryable transactions are idempotent and help ensure that the saga can reach its final state"
— Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

## Evidence 6: Lack of Isolation (Tidak Ada Isolation Level 'I' dalam ACID)

**Claim:** Sagas tidak memiliki built-in isolation, yang bisa menyebabkan anomali data (lost updates, dirty reads, fuzzy reads)

**Evidence:** "Because each service manages its own data, called participant data, there's no built-in isolation across services. This setup can result in data inconsistencies or durability problems, such as partially applied updates or conflicts between services." — Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

**Corroborated By:**
- Chris Richardson: "Lack of isolation (the 'I' in ACID) - the lack of isolation means that there's risk that the concurrent execution of multiple sagas and transactions can use data anomalies" (https://microservices.io/patterns/data/saga.html)

---

## Evidence 7: Countermeasures untuk Anomali Data

**Claim:** Untuk mengatasi anomali data, tersedia countermeasures: semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency

**Evidence:**
- **Semantic lock:** Use application-level locks when a saga's compensable transaction uses a semaphore to indicate that an update is in progress
- **Commutative updates:** Design updates so that they can be applied in any order while still producing the same result
- **Pessimistic view:** Reorder the sequence of the saga so that data updates occur in retryable transactions to eliminate dirty reads
- **Reread values:** Confirm that data remains unchanged before you make updates
- **Version files:** Maintain a log of all operations performed on a record
— Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

## Evidence 8: Idempotency Requirement

**Claim:** Transaction dalam saga harus idempotent untuk memastikan reliability

**Evidence:** "Handling transient failures and idempotence: The system must handle transient failures effectively and ensure idempotence, when repeating the same operation doesn't alter the outcome" — Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

**Corroborated By:**
- Chris Richardson: "retryable transactions are idempotent" and mentions "Idempotent Consumer pattern" (https://microservices.io/patterns/data/saga.html)

---

## Evidence 9: Limitations of Compensating Transactions

**Claim:** Compensating transactions mungkin tidak selalu berhasil, yang bisa meninggalkan sistem dalam inconsistent state

**Evidence:** "Limitations of compensating transactions: Compensating transactions might not always succeed, which can leave the system in an inconsistent state" — Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

## Evidence 10: Atomicity Guarantee di Level Saga

**Claim:** Saga menjamin atomicity di level saga, bukan di level individual transaction - semua steps harus selesai atau compensations harus dieksekusi

**Evidence:** "A saga is a sequence of local transactions. Each local transaction... Completes its work atomically within a single service... If a local transaction fails, the saga performs a series of compensating transactions to reverse the changes" — Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

## Evidence 11: Anomali Data Types (Lost Updates, Dirty Reads, Fuzzy Reads)

**Claim:** Anomali data yang umum terjadi pada sagas: lost updates, dirty reads, dan fuzzy/nonrepeatable reads

**Evidence:**
- **Lost updates:** "When one saga modifies data without considering changes made by another saga, it results in overwritten or missing updates"
- **Dirty reads:** "When a saga or transaction reads data that another saga has modified, but the modification isn't complete"
- **Fuzzy/nonrepeatable reads:** "When different steps in a saga read inconsistent data because updates occur between the reads"
— Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

## Evidence 12: When NOT to Use Saga

**Claim:** Saga tidak cocok untuk: tightly coupled transactions, cyclic dependencies, atau ketika compensating transactions ada di earlier participants

**Evidence:** "This pattern might not be suitable when: Transactions are tightly coupled. Compensating transactions occur in earlier participants. There are cyclic dependencies." — Microsoft Azure Architecture Center

**Source:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

---

## Evidence 13: Atomically Update State AND Publish Message

**Claim:** Untuk reliable saga, service harus atomically update state DAN publish message/event, tidak bisa pakai traditional distributed transaction yang span database dan message broker

**Evidence:** "In order to be reliable, a service must atomically update its database and publish a message/event. It cannot use the traditional mechanism of a distributed transaction that spans the database and the message broker." — Chris Richardson / Microservices.io

**Source:** https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH
