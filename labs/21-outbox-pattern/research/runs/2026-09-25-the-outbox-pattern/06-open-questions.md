# 06-open-questions.md

- How do high-frequency transactions impact the database performance when using a polling publisher versus a CDC-based tailing approach?
- What are the most effective strategies for pruning or archiving processed records from the outbox table to prevent unrestrained table growth?
- How should systems handle "poison" events that are successfully written to the outbox but consistently fail serialization or broker validation during the relay phase?
- What is the empirical latency difference between a typical scheduled Polling Relay and a CDC Log Tailing Relay (e.g., Debezium) under production load?
- Are there ready-made, battle-tested Outbox implementations natively embedded within modern application frameworks (like Spring, Laravel, or Rails) that abstract the relay mechanism?
