# 04 - Contradictions: Saga Pattern Research

## Contradiction 1: Choreography vs Orchestration Preferred Defaults
Statement A: Temporal (2023-07-13) posits that choreography might look simpler initially, but orchestration is easier to build and maintain when adopted from the start if backed by durable execution.  
Location: `research/04-contradictions.md:5-7`  
Statement B: Microsoft Azure Architecture Center and common industry consensus state that choreography is suitable for simple workflows with few services, while orchestration is better suited for complex workflows.  
Location: `research/04-contradictions.md:8-10`  
Type: SOURCE_CONFLICT (Nuance/Perspective)  
Impact: LOW. Does not invalidate either approach. Temporal's perspective reflects their positioning as an orchestration/durable execution engine, whereas Microsoft reflects general architectural rules of thumb.  
Assessment: Resolved. The research appropriately notes that qualitative complexity thresholds dictate choice, and the 3-step checkout flow in the lab exercise is well-suited for orchestration.

---

## Contradiction 2: Exactly-Once vs At-Least-Once Delivery Guarantees
Statement A: AWS Step Functions Standard workflows guarantee "exactly-once" execution. Temporal guarantees deterministic replay / workflow-level execution.  
Location: `research/04-contradictions.md:70-75`  
Statement B: Message brokers and distributed participants operate under at-least-once delivery; activity executions require participant idempotency regardless of orchestrator guarantees.  
Location: `research/04-contradictions.md:76-78`  
Type: INTERNAL (Layering Distinction)  
Impact: LOW. The tension is reconciled by distinguishing orchestrator state machine execution (workflow level) from RPC/messaging participant invocation (network level).  
Assessment: Resolved. Research correctly insists on idempotency keys at the participant level regardless of workflow engine guarantees.

---

## Contradiction 3: Compensation Failure Protocols
Statement A: Compensating transactions are the mechanism to return the system to consistency.  
Location: `research/03-evidence.md:21-26`  
Statement B: Compensating transactions themselves can fail, leaving the system in an inconsistent state with no standard automatic second-tier rollback protocol.  
Location: `research/04-contradictions.md:34-49`  
Type: INTERNAL (Known Theoretical Limit)  
Impact: LOW to MEDIUM. Represents an inherent limitation of Sagas, not an error in research.  
Assessment: Resolved. Research explicitly highlights that production systems must rely on retries with exponential backoff, dead-letter queues, and human operator intervention.

---

## Material Contradictions Summary
No material, irreconcilable contradictions found. Differences between sources represent legitimate architectural tradeoffs, different layers of abstraction (orchestrator vs participant), and vendor-specific design philosophies.
