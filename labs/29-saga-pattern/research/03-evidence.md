# Evidence Register: Saga Pattern

## Evidence 1

**Claim:** Saga pattern originated in 1987 paper by Garcia-Molina & Salem as mechanism for long-lived database transactions using compensating transactions.

**Evidence:** Temporal blog footnote states "sagas aren't a trend; they've been around in databases since the 80s" linking to Garcia-Molina paper; Microsoft Azure docs reference saga history as extension of long-running transaction concept. Cornell mirror hosts the original PDF at /andru/cs711/2002fa/reading/sagas.pdf.

**Source:** Hector Garcia-Molina & Kenneth Salem — "Sagas" (1987) via Cornell mirror; Temporal Blog 2023-05-24 footnote 2.

**URL:** https://www.cs.cornell.edu/andru/cs711/2002fa/reading/sagas.pdf ; https://temporal.io/blog/saga-pattern-made-easy

**Confidence:** HIGH

**Corroborated By:** Microsoft Azure Architecture Center (references saga historical context); multiple secondary sources citing DOI 10.1145/62224.62226.

**Notes:** PDF is scanned image with LZW compression; text extraction garbled. Historical claim verified via citation chain rather than direct PDF text parsing. DOI and authorship cross-checked via ACM Digital Library references in multiple sources.

---

## Evidence 2

**Claim:** Saga is a sequence of local transactions where each local transaction updates its own database and publishes a message/event to trigger the next; if one fails, compensating transactions undo preceding changes.

**Evidence:** "Implement each business transaction that spans multiple services as a saga. A saga is a sequence of local transactions. Each local transaction updates the database and publishes a message or event to trigger the next local transaction in the saga. If a local transaction fails because it violates a business rule then the saga executes a series of compensating transactions that undo the changes that were made by the preceding local transactions."

**Source:** Chris Richardson — Microservices.io: Pattern Saga

**URL:** https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH

**Corroborated By:** Microsoft Azure Architecture Center — "The Saga pattern manages transactions by breaking them into a sequence of local transactions. Each local transaction completes its work atomically within a single service... If a local transaction fails, the saga performs a series of compensating transactions to reverse the changes." ; AWS Prescriptive Guidance Saga Pattern.

**Notes:** Definition consistent across all Tier 1/2 sources opened.

---

## Evidence 3

**Claim:** Two coordination approaches exist: Choreography (each service publishes domain events triggering next service) and Orchestration (central orchestrator tells participants what to execute).

**Evidence:** Microservices.io: "There are two ways of coordination sagas: Choreography - each local transaction publishes domain events that trigger local transactions in other services; Orchestration - an orchestrator (object) tells the participants what local transactions to execute." Microsoft: "The two typical saga implementation approaches are choreography and orchestration."

**Source:** Microservices.io: Pattern Saga + Microsoft Azure Architecture Center

**URL:** https://microservices.io/patterns/data/saga.html ; https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

**Corroborated By:** Temporal blog "To choreograph or orchestrate your saga" (2023-07-13); AWS Prescriptive Guidance.

**Notes:** Lab specification correctly reflects this dichotomy with matching diagrams.

---

## Evidence 4

**Claim:** Choreography benefits: loose coupling, no coordinator, no SPOF, good for simple workflows; drawbacks: spaghetti events, cyclic dependencies, hard to track/debug, difficult integration testing.

**Evidence:** Microsoft table: Benefits — "Good for simple workflows that have few services and don't need a coordination logic; Doesn't introduce a single point of failure." Drawbacks — "Workflow can be confusing when you add new steps; risk of cyclic dependency; Integration testing is difficult because all services must run." Temporal: "just from looking at each microservice's individual codebase, it's difficult to understand the order that the system should have."

**Source:** Microsoft Azure Architecture Center + Temporal Blog 2023-07-13

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga ; https://temporal.io/blog/to-choreograph-or-orchestrate-your-saga-that-is-the-question

**Confidence:** HIGH

**Corroborated By:** Microservices.io related discussion.

**Notes:** Lab spec claim ("Sederhana untuk workflow kecil, loose coupling" / "Sulit di-track jika workflow kompleks - spaghetti event") verified.

---

## Evidence 5

**Claim:** Orchestration benefits: better for complex workflows, avoids cyclic dependencies, clear separation of responsibilities, easier debugging; drawbacks: requires coordination logic, introduces single point of failure.

**Evidence:** Microsoft table: Benefits — "Better suited for complex workflows or when you add new services; Avoids cyclic dependencies; Clear separation of responsibilities simplifies service logic." Drawbacks — "Other design complexity requires an implementation of a coordination logic; Introduces a point of failure." AWS: "The saga pattern is difficult to debug and its complexity increases with the number of microservices."

**Source:** Microsoft Azure Architecture Center ; AWS Prescriptive Guidance

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga ; https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html

**Confidence:** HIGH

**Corroborated By:** Temporal blog (orchestration centralized control flow easier to understand/debug; SPOF is "glaring Achilles' heel").

**Notes:** Lab spec orchestration diagram and description verified.

---

## Evidence 6

**Claim:** Compensating transactions are semantic undo operations, not automatic rollback; required for each reversible step (e.g., Reserve Stock → Release Stock; Charge Card → Refund).

**Evidence:** Microsoft: "Compensable transactions can be undone or compensated for by other transactions with the opposite effect." Lab spec mapping (Reserve Stock→Release Stock etc.) matches definition. Microservices.io: "a developer must design compensating transactions that explicitly undo changes made earlier in a saga rather than relying on the automatic rollback feature of ACID transactions" and "Lack of isolation (the 'I' in ACID)."

**Source:** Microsoft Azure Architecture Center ; Microservices.io

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga ; https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH

**Corroborated By:** AWS Saga Pattern; Temporal compensating-actions blog linked from main saga post.

**Notes:** Lab spec compensation table verified as correct examples.

---

## Evidence 7

**Claim:** Sagas include compensable, pivot (point of no return), and retryable transaction types.

**Evidence:** Microsoft: "Pivot transactions serve as the point of no return in the saga. After a pivot transaction succeeds, compensable transactions are no longer relevant. All subsequent actions must be completed for the system to achieve a consistent final state." "Retryable transactions follow the pivot transaction... idempotent and help ensure that the saga can reach its final state."

**Source:** Microsoft Azure Architecture Center – satu sumber otoritatif yang merinci taksonomi tiga jenis langkah ini secara eksplisit

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** MEDIUM

**Corroborated By:** Garcia-Molina original paper (pivot concept referenced in secondary literature); Temporal implementation shows analogous LIFO compensation ordering.

**Notes:** Single authoritative source for this specific taxonomy; not directly corroborated by Microservices.io in opened page. Lab spec does not mention pivot — extended concept.

---

## Evidence 8

**Claim:** 2PC is not viable for microservices due to blocking locks, tight coupling, availability reduction; sagas trade ACID isolation for availability and loose coupling with eventual consistency.

**Evidence:** Microservices.io Forces: "2PC is not an option." Microsoft Context: "it can be more complex to achieve ACID compliance across multiple services... architectures that rely on interprocess communication, or traditional transaction models like two-phase commit protocol, are often better suited for the Saga pattern." AWS: "You should consider using this pattern if... There are long-lived transactions and you don't want other microservices to be blocked if one microservice runs for a long time." Microsoft lists sagas sacrifice isolation requiring countermeasures.

**Source:** Microservices.io ; Microsoft Azure Architecture Center ; AWS Prescriptive Guidance

**URL:** https://microservices.io/patterns/data/saga.html ; https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

**Corroborated By:** All three sources converge; lab spec ("2PC terlalu lambat atau tidak memungkinkan") verified.

**Notes:** Lab spec claim that sagas use eventual consistency not strict ACID instant consistency is verified — Microsoft explicitly notes compensating transactions might not always succeed leaving inconsistent state.

---

## Evidence 9

**Claim:** Dual-write problem: atomically updating database and publishing message/event is impossible with distributed transaction across DB and message broker; outbox pattern solves by writing business data + outbox record in single local ACID transaction and relaying via CDC/polling.

**Evidence:** Debezium: "So how can this situation be avoided? The answer is to only modify one of the two resources... The idea of this approach is to have an 'outbox' table in the service's database. When receiving a request... not only an INSERT into the PurchaseOrder table is done, but, as part of the same transaction, also a record representing the event to be sent is inserted into that outbox table." Microservices.io Related Patterns: "In order to be reliable, a service must atomically update its database and publish a message/event... It cannot use the traditional mechanism of a distributed transaction... Instead, it must use one of the patterns listed below: Event sourcing, Transactional Outbox."

**Source:** Gunnar Morling / Debezium Blog 2019-02-19 ; Microservices.io Resulting Context

**URL:** https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/ ; https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH

**Corroborated By:** Both sources describe same mechanism; outbox table schema (id uuid, aggregatetype, aggregateid, type, payload jsonb) documented in Debezium post.

**Notes:** Lab spec implicitly assumes reliable event publication; outbox is prerequisite for production correctness.

---

## Evidence 10

**Claim:** Idempotency mandatory for saga participants; duplicate detection via PROCESSED_MESSAGES table with (subscriberId, messageID) primary key or event UUID header.

**Evidence:** Microservices.io Idempotent Consumer: "After starting the database transaction, the message handler inserts the message's ID into the PROCESSED_MESSAGE table. Since the (subscriberId, messageID) is the PROCESSED_MESSAGE table's primary key the INSERT will fail if the message has been already processed... The message handler can then rollback the transaction and ignore the message." Debezium consumer example: eventId UUID propagated as Kafka header allows efficient duplicate detection. Temporal: compensating activities must handle idempotency key (clientId or clientId+workflowId).

**Source:** Microservices.io Idempotent Consumer Pattern ; Debezium Blog ; Temporal Blog 2023-05-24

**URL:** https://microservices.io/patterns/data/idempotent-consumer.html ; https://debezium.io/blog/2019/02/19/reliable-microservices-data-exchange-with-the-outbox-pattern/ ; https://temporal.io/blog/saga-pattern-made-easy

**Confidence:** HIGH

**Corroborated By:** All three sources agree; at-least-once delivery semantics make idempotency essential.

**Notes:** Directly relevant to lab exercise where Payment refund compensation must be idempotent if retried.

---

## Evidence 11

**Claim:** Sagas lack isolation, causing anomalies: lost updates, dirty reads, fuzzy/nonrepeatable reads.

**Evidence:** Microsoft: "Typical problems include: Lost updates: When one saga modifies data without considering changes made by another saga... Dirty reads: When a saga or transaction reads data that another saga has modified, but the modification isn't complete. Fuzzy, or nonrepeatable, reads: When different steps in a saga read inconsistent data because updates occur between the reads." Microservices.io: "Lack of isolation (the 'I' in ACID) - the lack of isolation means that there's risk that the concurrent execution of multiple sagas and transactions can use data anomalies."

**Source:** Microsoft Azure Architecture Center ; Microservices.io

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga ; https://microservices.io/patterns/data/saga.html

**Confidence:** HIGH

**Corroborated By:** Both sources list identical anomaly types; Microsoft adds 6 countermeasures (semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency).

**Notes:** Lab spec does not explicitly mention isolation anomalies but states eventual consistency — this is the technical explanation. Daftar 6 countermeasures spesifik berasal dari Microsoft; Microservices.io mengonfirmasi konsep countermeasures tanpa enumerasi.

---

## Evidence 12

**Claim:** Countermeasures for isolation anomalies include semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency.

**Evidence:** Microsoft Strategies section lists all six countermeasures verbatim with descriptions.

**Source:** Microsoft Azure Architecture Center – satu-satunya sumber yang menyediakan enumerasi lengkap dari enam countermeasures spesifik ini

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** MEDIUM

**Corroborated By:** Microservices.io chapter 4/section 4.3 reference mentions "countermeasures, which are design techniques that implement isolation. Moreover, careful analysis is needed to select and correctly implement the countermeasures."

**Notes:** Only one source (Microsoft) provides detailed enumeration; second source (Microservices.io) confirms existence of countermeasure concept without enumerating. Lab could mention at least semantic lock concept. Taxonomy this specific is Microsoft's formalization, not universally agreed across all literature.

---

## Evidence 13

**Claim:** Orchestration saga can be implemented via AWS Step Functions state machine with Task/Choice/Retry/Catch states; Standard workflows provide exactly-once execution, Express provides at-least-once.

**Evidence:** AWS: illustration shows order processing with Step Functions; each step has success/failure branches. AWS Docs Welcome page: Standard workflows have exactly-once execution, Express have at-least-once. Temporal alternative uses deterministic replay and avoids orchestrator SPOF.

**Source:** AWS Prescriptive Guidance + AWS Step Functions Docs

**URL:** https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html ; https://docs.aws.amazon.com/step-functions/latest/dg/welcome.html

**Confidence:** HIGH

**Corroborated By:** AWS Prescriptive Guidance example lists Use Cases including saga; Step Functions docs confirm execution semantics.

**Notes:** Lab exercise maps cleanly to Step Functions orchestration: Create Order → Reserve Payment → Reserve Stock → (on failure) Refund Payment.

---

## Evidence 14

**Claim:** Compensating transactions may themselves fail, leaving system in inconsistent state; monitoring and tracking workflow essential.

**Evidence:** Microsoft Problems and considerations: "Limitations of compensating transactions: Compensating transactions might not always succeed, which can leave the system in an inconsistent state." "Need for monitoring and tracking sagas: Monitoring and tracking the workflow of a saga are essential." AWS Important: "The saga pattern is difficult to debug and its complexity increases with the number of microservices."

**Source:** Microsoft Azure Architecture Center ; AWS Prescriptive Guidance

**URL:** https://learn.microsoft.com/en-us/azure/architecture/patterns/saga ; https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html

**Confidence:** HIGH

**Corroborated By:** Both sources identify same risk; Temporal compensating sample code logs compensation failure but does not auto-resolve — requires operational intervention.

**Notes:** Important lab design consideration: orchestrator must handle compensation failure with retry + dead-letter/alert, not silently ignore.

---

## Evidence 15

**Claim:** Appropriate use cases for saga: multi-service/database transactions requiring eventual consistency; NOT suitable for single-database monolith (use local ACID) or when SERIALIZABLE isolation across services required.

**Evidence:** AWS "You should consider using this pattern if: Application needs to maintain data consistency across multiple microservices without tight coupling; There are long-lived transactions." Microsoft When to use: "You need to ensure data consistency in a distributed system without tight coupling; You need to roll back or compensate if one operation fails." Microsoft Not suitable when: "Transactions are tightly coupled." Lab spec: "Jangan gunakan jika masih dalam satu database monolith" and "Membutuhkan isolation level SERIALIZABLE secara instantaneous across services" — matches authoritative guidance.

**Source:** AWS Prescriptive Guidance ; Microsoft Azure Architecture Center

**URL:** https://docs.aws.amazon.com/prescriptive-guidance/latest/modernization-data-persistence/saga-pattern.html ; https://learn.microsoft.com/en-us/azure/architecture/patterns/saga

**Confidence:** HIGH

**Corroborated By:** All sources converge; lab spec correctly captures both positive and negative cases.

**Notes:** Verifies lab spec's "Kapan Memakai Saga?" section is accurate per authoritative sources.

