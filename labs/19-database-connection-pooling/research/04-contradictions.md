# Contradictions & Disagreements

## Contradiction 1: Sizing Formula for Modern NVMe SSDs vs Spindle Hard Drives

### SOURCE A:
PostgreSQL Wiki & HikariCP Wiki
"formula: active connections = ((core_count * 2) + effective_spindle_count)... There hasn't been any analysis so far regarding how well the formula works with SSDs."

### SOURCE B:
Modern Cloud & Database Engineering Benchmarks (HikariCP commentary)
"Don't be tricked into thinking, 'SSDs are faster and therefore I can have more threads'. That is exactly 180 degrees backwards. Faster, no seeks, no rotational delays means less blocking and therefore fewer threads [closer to core count] will perform better than more threads."

### ASSESSMENT:
Classical database tuning assumed disk seeks introduced idle I/O wait times where the CPU could switch to other threads. For modern NVMe/SSD storage and high memory buffer cache hit rates, effective spindle count is essentially 0. Therefore, modern pool sizing should approach `core_count` or `core_count * 2` rather than inflated numbers. There is agreement on the underlying physics (fewer threads when blocking is absent), but historical formula implementations can mislead practitioners if spindle count is assumed to be large.

---

## Contradiction 2: Application-Side Pooling vs Middleware/Proxy-Level Pooling

### SOURCE A:
Application Framework View (HikariCP, Java EE)
Pooling should reside directly in the client application to eliminate protocol hops and optimize connection borrowing latency, allowing threads to block directly at the in-memory pool.

### SOURCE B:
Database Operations & Architecture View (PostgreSQL Core / PgBouncer)
"The decision not to include a connection pooler inside the PostgreSQL server itself has been taken deliberately and with good reason: In many cases you will get better performance if the connection pooler is running on a separate machine... having pooling outside the core server maintains flexibility." In horizontally scaled microservices (e.g., hundreds of Kubernetes pods or background workers), individual application pools accumulate and overwhelm database limits unless capped by an intermediary proxy like PgBouncer.

### ASSESSMENT:
They do not fundamentally contradict; they address different deployment topologies:
- For a monolithic or small cluster setup, application-side connection pooling (e.g., HikariCP) is sufficient and low-latency.
- For high-concurrency horizontally scaled microservices / serverless architectures where `instance_count * workers * local_pool_size > database max_connections`, proxy pooling (such as PgBouncer in transaction mode) is mandatory to prevent backend connection exhaustion.
