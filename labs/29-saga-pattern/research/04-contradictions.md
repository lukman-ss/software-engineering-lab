# Contradictions and Uncertainties: Saga Pattern

## Contradiction 1: Choreography vs Orchestration as Default Choice

**SOURCE A (Temporal, 2023-07-13):**
Choreography is "easier to implement, at least initially" but "orchestration is often easier to build when one uses it from the start." Counterintuitive: "one wants to reach for the light, agile option (choreography) in the early days... but orchestration is often easier to build when one uses it from the start."

**SOURCE B (Microsoft Azure Architecture Center):**
Choreography: "Good for simple workflows that have few services and don't need a coordination logic." Orchestration: "Better suited for complex workflows or when you add new services."

**SOURCE C (Lab specification):**
Presents choreography as "Sederhana untuk workflow kecil" and orchestration as the recommended approach for the e-commerce checkout exercise (3-step workflow).

**ASSESSMENT:**
Not a true contradiction — sources agree on the complexity threshold. Difference is emphasis: Temporal argues orchestration is simpler even for modest workflows if you have a durable execution engine; Microsoft/lab treat 3-service checkout as orchestration-appropriate because failure handling and compensation are explicit. Lab specification is consistent with Microsoft's framing.

**Remaining uncertainty:** No quantitative threshold (e.g., number of services or branching complexity) at which choreography becomes unmaintainable. Community consensus is qualitative.

---

## Contradiction 2: Isolation Countermeasures Completeness

**SOURCE A (Microsoft Azure Architecture Center):**
Enumerates 6 specific countermeasures: semantic lock, commutative updates, pessimistic view, reread values, version files, risk-based concurrency.

**SOURCE B (Microservices.io):**
References "countermeasures, which are design techniques that implement isolation" and points to Chapter 4 of *Microservices Patterns* (Manning) without enumerating them in the public page.

**ASSESSMENT:**
Not a disagreement — Microsoft's list appears to be derived from Richardson's book. Public Microservices.io page does not independently verify the six-item taxonomy. Confidence for the specific six-item list is MEDIUM (one detailed source). This enumeration reflects Microsoft's formalization, not a universal standard across all literature.

---

## Contradiction 3: Compensating Transaction Failure Handling

**SOURCE A (Microsoft):**
"Compensating transactions might not always succeed, which can leave the system in an inconsistent state."

**SOURCE B (Temporal compensating samples):**
On compensation failure, Temporal logs error and continues remaining compensations (or parallel). No automatic recovery of failed compensation — operator intervention implied.

**SOURCE C (AWS Step Functions):**
Catch/Retry fields can retry compensation; after retries exhausted, state machine can fail leaving compensation incomplete.

**ASSESSMENT:**
All sources agree compensation can fail. None provide a complete, general-purpose recovery protocol for "compensation of compensation." This is an acknowledged limitation of the pattern, not a source disagreement.

**Remaining uncertainty:** Production systems typically add: retry with backoff, dead-letter queue, human-in-the-loop, or "semantic lock" remaining until operator resolves. No standard protocol exists.

---

## Contradiction 4: Dual-Write Problem Attribution

**SOURCE A (Debezium/Morling 2019):**
Explicitly names "The Issue of Dual Writes" as writing to database AND Kafka without shared transaction.

**SOURCE B (Microservices.io Resulting Context):**
Describes same problem without using "dual-write" terminology: "a service must atomically update its database and publish a message/event. It cannot use the traditional mechanism of a distributed transaction."

**SOURCE C (Garcia-Molina 1987 via secondary citations):**
Original paper addressed long-lived transactions in a *single* database system, not the dual-write problem of microservices + message brokers. Dual-write is a *modern* problem that sagas inherited when adapted to event-driven microservices.

**ASSESSMENT:**
Terminology difference, not factual contradiction. Dual-write is a prerequisite reliability concern for *choreography* sagas (event publication) and for orchestration sagas that use async command/reply. Lab spec lists Dual-Write Problem as a keyword; this is accurate as a related concern, not as part of the original 1987 definition.

---

## Contradiction 5: Exactly-Once vs At-Least-Once Semantics

**SOURCE A (AWS Step Functions):**
Standard workflows: exactly-once execution. Express workflows: at-least-once execution.

**SOURCE B (Temporal):**
Activities retried automatically; "you, the programmer, need to make sure each Temporal Activity is idempotent." Temporal claims "exactly-once" at workflow level via deterministic replay, but activities are at-least-once.

**SOURCE C (Debezium Outbox):**
Pipeline is at-least-once; consumers must detect duplicates via event UUID.

**ASSESSMENT:**
Apparent tension is actually layered: *workflow/orchestrator* may be exactly-once (Step Functions Standard, Temporal workflow history), but *participant services* always see at-least-once invocation and must be idempotent. Lab exercise should treat Payment debit and refund as requiring idempotency keys regardless of orchestrator choice.

---

## No Material Contradictions Discovered

The core claims of the lab specification (saga as sequence of local transactions + compensation; two coordination styles; when to use/not use) are corroborated by all opened Tier 1 and Tier 2 sources. Remaining differences are emphasis, terminology, and implementation detail — not conflicting facts.
