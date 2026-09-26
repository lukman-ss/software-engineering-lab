# Research Topic
Zero-Downtime Deployment — Deploy Versi Baru Tanpa Membuat User Tahu Ada Deployment

# Objective
To investigate the provided topic, collect reliable evidence on zero-downtime deployment patterns (Rolling, Blue-Green, Expand/Contract migrations, Probes, Queue worker deployments), cross-check important claims, identify any contradictions, and produce a structured research report without creating publication content.

# Research Questions
1. What are the core principles of zero-downtime deployment (such as Rolling and Blue-Green)?
2. How are database migrations managed during zero-downtime deployments to avoid downtime (e.g., Expand and Contract pattern)?
3. How do Liveness and Readiness probes differ, and why are they critical for zero-downtime deployments?
4. What is graceful shutdown / connection draining and how does it prevent dropped requests?
5. How should background queue workers be deployed in a zero-downtime scenario (e.g., in Laravel)?

# Search Strategy
1. Fetch authoritative content from Martin Fowler's bliki regarding Blue-Green Deployment and Parallel Change (Expand and Contract).
2. Fetch Kubernetes official documentation regarding Pod Lifecycle, Liveness vs Readiness Probes, and graceful Pod termination.
3. Fetch Laravel official documentation regarding Queue Workers and deployment strategies.

# Expected Primary Sources
- MartinFowler.com (Martin Fowler, Danilo Sato)
- Kubernetes Official Documentation
- Laravel Official Documentation

# Risks / Unknowns
- Finding exact statistics or metrics for typical deployment durations might require broad assumptions, so they should be avoided unless explicitly stated by the source.
- Specific implementation details vary wildly between tech stacks (e.g., AWS vs Kubernetes), so principles should be kept general but anchored in the requested stack context (Laravel, PostgreSQL, Redis, K8s).