# Contradictions Audit: Optimistic vs Pessimistic Locking

## Contradiction 1
Statement A: "Repeatable Read uses Snapshot Isolation... Serializable adds predicate locking (Serializable Snapshot Isolation) on top of Snapshot Isolation, detecting serialization anomalies." (PostgreSQL Documentation 13.2)  
Location: `research/04-contradictions.md:9-11`  
Statement B: "REPEATABLE READ (default): Consistent reads use snapshot from first read. Gap locks / next-key locks prevent phantom rows for locking reads." (MySQL Documentation 15.7.2.1)  
Location: `research/04-contradictions.md:12-14`  
Type: SOURCE_CONFLICT  
Impact: Developers migrating or designing queries across databases may falsely assume identical behavior under the same isolation label.  
Assessment: Accurately identified and resolved. The difference stems from underlying concurrency architectures (MVCC/Snapshot Isolation in PG vs Strict 2PL with gap locks in InnoDB).  

---

## Contradiction 2
Statement A: "Read Uncommitted behaves identically to Read Committed (prevents dirty reads via MVCC)." (PostgreSQL 13.2)  
Location: `research/04-contradictions.md:24-26`  
Statement B: Standard ANSI SQL defines READ UNCOMMITTED as permitting dirty reads.  
Location: `research/04-contradictions.md:27-29`  
Type: SOURCE_CONFLICT  
Impact: Applications attempting to use READ UNCOMMITTED for lock-free read performance gains in PostgreSQL see no effect.  
Assessment: Correctly documented as an engine implementation reality rather than a conceptual error.  

---

## Contradiction 3
Statement A: PostgreSQL re-evaluates the WHERE clause under READ COMMITTED upon encountering updated rows.  
Location: `research/04-contradictions.md:36-38`  
Statement B: MySQL InnoDB uses semi-consistent reads for non-locking reads and updates under READ COMMITTED.  
Location: `research/04-contradictions.md:42-44`  
Type: INTERNAL / ENGINE_VARIATION  
Impact: Clarifies why naive read-modify-write patterns fail under READ COMMITTED across both engines despite different internal evaluation strategies.  
Assessment: Sound analysis; confirms both systems require explicit locking or version guards.  

---

## Overall Assessment
No material unresolved contradictions exist within the research files. Identified discrepancies reflect documented behavioral differences between relational database engines.
