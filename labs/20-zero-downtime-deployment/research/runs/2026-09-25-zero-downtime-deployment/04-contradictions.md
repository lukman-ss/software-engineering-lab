# Contradictions & Disagreements

## Material Contradictions
No material contradictions discovered between primary authoritative sources regarding the fundamental principles of zero-downtime deployment, health probes, graceful shutdown, or database parallel change.

## Nuances & Differing Approaches

### Approach 1: Blue-Green vs. Rolling Deployment
- **Blue-Green Deployment**:
  - Requires maintaining two nearly identical environments simultaneously, temporarily doubling resource/infrastructure footprint.
  - Cutover happens instantly at the load balancer / router level.
  - Rollback is nearly instantaneous (switch traffic back to Blue).
- **Rolling Deployment**:
  - Replaces instances incrementally in place without requiring 2x resource overhead.
  - Rollback is slower because replacement instances must be rolled back step-by-step.
  - Requires both old and new versions to coexist and share production database/queue traffic over a longer deployment window.

### Approach 2: Database Migration Strategy
- One school of thought recommends executing database schema expansions as an automated pre-deployment step within CI/CD pipelines before application deployment.
- Another variation decouples schema changes entirely into a distinct, manual or scheduled release cycle prior to deploying the code that relies on it.
- In both cases, there is consensus that breaking changes (dropping old columns, strict non-null additions without defaults) must be deferred to a post-deployment "contract" phase.