# Research Report

## Research Question
What are the core mechanisms, patterns, and architectural requirements for executing Zero-Downtime Deployments without interrupting HTTP traffic, dropping active background jobs, or breaking shared database states?

## Executive Summary
Zero-downtime deployment (ZDD) ensures continuous availability by starting the newly deployed application version before terminating the old one ("Start new before stopping old"). Implementing ZDD successfully requires addressing the deployment holistically across the load balancer, application lifecycle events, database schema evolution, and background queue workers. The evidence confirms that routing traffic only to "ready" instances, draining active connections via graceful termination, and using the "Expand and Contract" pattern for databases are mandatory prerequisites for achieving true zero-downtime.

## Findings

### Finding 1: Traffic Orchestration via Health Probes
Claim: Deployments require distinguishing between container "Liveness" (is the process running) and "Readiness" (is the app fully initialized and connected to its dependencies) to prevent routing traffic to unprepared nodes.
Evidence: Kubernetes architecture explicitly separates these two concerns. Liveness determines when a container must be restarted (e.g., deadlocks), while Readiness determines when a container should be attached to or removed from the load balancer pool.
Sources: Kubernetes Official Documentation (Pod Lifecycle)
Confidence: HIGH

### Finding 2: Graceful Shutdown (Connection Draining)
Claim: Terminating application instances forcefully drops active, in-flight requests and background jobs. Graceful termination allows the application to finish processing active transactions before exiting.
Evidence: Orchestrators send a SIGTERM signal and initiate a grace period window. The instance is immediately removed from the load balancer to prevent new traffic, but active processes are given time to finish before a SIGKILL is issued. Similarly, Laravel queue workers provide a `queue:restart` command that instructs workers to complete their current job before exiting.
Sources: Kubernetes Official Documentation (Pod Lifecycle), Laravel Official Documentation (Queues)
Confidence: HIGH

### Finding 3: Database and State Backward Compatibility
Claim: Application version 1 (v1) and version 2 (v2) will inevitably run simultaneously during a deployment rollout, necessitating backward-compatible database migrations.
Evidence: The "Parallel Change" or "Expand and Contract" pattern dictates that database migrations must be separated from application rollouts into distinct phases. Expanding (adding columns) occurs first, application code is rolled out while both schemas exist, and contracting (removing old columns) occurs only after all v1 instances are fully terminated.
Sources: Martin Fowler / Danilo Sato (Parallel Change, Blue Green Deployment)
Confidence: HIGH

### Finding 4: Blue-Green Deployment for Safe Rollbacks
Claim: Blue-Green deployment strategies provide the safest rollback mechanism by physically separating the old environment (Blue) from the new environment (Green) behind a router.
Evidence: Martin Fowler asserts that maintaining two identical environments allows for rigorous final testing on the "Green" environment before router switchover, enabling near-instantaneous reversion if issues are detected post-switch.
Sources: Martin Fowler (Blue Green Deployment)
Confidence: HIGH

## Areas of Agreement
- **Coexistence is unavoidable**: All authoritative sources agree that achieving zero-downtime strictly requires application version N and version N-1 to coexist on the infrastructure and query the same database at the same time.
- **Graceful degradation over forced kills**: Process termination must be cooperative (SIGTERM, connection draining, waiting for job completion) rather than forceful (SIGKILL), regardless of the underlying stack (K8s, Laravel, raw VMs).

## Areas of Disagreement
No fundamental disagreements were found. There are simply tradeoffs between strategies (Rolling vs Blue-Green), primarily concerning infrastructure cost overhead vs rollback speed.

## Limitations
- Source evidence concerning queue graceful termination was limited to Laravel specifics per the target lab constraints; other background processing frameworks (e.g., Celery, Sidekiq) behave similarly but use different signalling mechanisms.
- Exact durations for "grace periods" (e.g., K8s default 30 seconds vs AWS typical draining defaults) vary by provider and application workload profile.

## Conclusion
A zero-downtime deployment is an architectural responsibility, not merely an infrastructure pipeline configuration. It mandates designing database migrations for backward compatibility, explicitly defining application readiness logic beyond basic HTTP 200 checks, implementing graceful connection draining, and safely orchestrating long-lived background workers to finish their tasks before yielding to the new version.