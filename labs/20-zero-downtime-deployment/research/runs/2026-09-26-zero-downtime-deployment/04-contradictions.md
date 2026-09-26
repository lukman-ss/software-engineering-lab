# Research Contradictions: Zero-Downtime Deployment

---

## Contradiction 1: Blue-Green vs Rolling Deployment – Database Schema Changes

**SOURCE A (Martin Fowler - Blue Green Deployment):** "The two environments need to be different but as identical as possible. In some situations they can be different pieces of hardware, or they can be different virtual machines running on the same (or different) hardware. They can also be a single operating environment partitioned into separate zones with separate IP addresses for the two slices." "Databases can often be a challenge with this technique, particularly when you need to change the schema to support a new version of the software. The trick is to separate the deployment of schema changes from application upgrades."

**SOURCE B (Kubernetes Documentation - Update a Deployment Without Downtime):** Demonstrates rolling updates using a single Deployment resource with a rolling update strategy. The deployment object manages the transition from old to new Pods automatically.

**ASSESSMENT:** These are complementary patterns that apply in different contexts. Blue-green deployment requires maintaining two full environments (typically separate clusters, VPCs, or at least separate load balancing namespaces) and is typically used for:
1. Large-scale applications with significant infrastructure
2. Applications requiring instant rollback capability
3. Stateful applications where data migration is complex

Rolling deployment (Kubernetes native) applies changes incrementally to the same environment and is typically used for:
1. Container-orchestrated environments
2. Applications that can tolerate brief capacity reduction
3. Stateless or properly load-balanced services

The "contradiction" is not a factual disagreement but a methodological choice. Both patterns achieve zero-downtime, but blue-green requires double the infrastructure capacity (two full environments) while rolling uses the existing capacity with proper health checking.

**Resolution:** Both patterns are valid and not mutually exclusive. The choice depends on infrastructure constraints, application architecture, and rollback requirements.

---

## Contradiction 2: Health Check Implementation – Active vs Passive

**SOURCE A (NGINX Documentation - HTTP Health Checks):** Provides two approaches: 
- Passive health checks (default with NGINX Open Source): monitor failed transactions, mark server unavailable after `max_fails` failures within `fail_timeout`
- Active health checks (NGINX Plus only): send periodic probe requests, check status code and body content

**SOURCE B (Laravel Documentation - Deployment - The Health Route):** Recommends a simple health check endpoint at `/up` that returns 200 on success, 500 on failure, with optional event-based extension for deeper checks.

**SOURCE C (Kubernetes Documentation - Probes):** Recommends readiness probes for traffic routing ("EndpointSlice controller removes the Pod's IP address from the EndpointSlices of all Services that match the Pod").

**ASSESSMENT:** The primary contradiction is cost/feature availability:
1. NGINX Open Source lacks active health checks — relies on passive monitoring of actual requests/responses
2. Laravel's `/up` endpoint is application-level, not infrastructure-level
3. Kubernetes readiness probes operate at the kubelet level

However, the concepts align: all three detect application state and gate traffic accordingly. The differences are in WHERE the check happens (infrastructure vs application) and HOW sophisticated the check can be (passive observation vs active probing with body content matching).

**Resolution:** This is a technical constraint based on tooling cost (NGINX Plus commercial license), not a conceptual disagreement. For zero-downtime deployments:
- NGINX Open Source + Laravel: use passive health checks + application-level `/up` endpoint
- Kubernetes: use readiness probes with application endpoint
- NGINX Plus + other stacks: can use active health checks for more sophisticated failure detection

---

## Contradiction 3: Pod Termination Timing – DNS Propagation vs Immediate Endpoint Update

**SOURCE A (Kubernetes Documentation - Pod Lifecycle):** States: "when the Pod is deleted, the corresponding endpoint in the EndpointSlice will update its conditions: the endpoint ready condition will be set to false, so load balancers will not use the Pod for regular traffic." This suggests immediate removal from service upon deletion.

**SOURCE B (Production Experience/Implementation Reality):** Load balancers often have connection draining windows that can take longer to expire. AWS ELB, AWS NLB, and cloud load balancers may maintain connections for up to 350 seconds (AWS NLB default) even after endpoint is marked unhealthy. Nginx's default `proxy_read_timeout` is 60 seconds.

**ASSESSMENT:** There is a timing gap between:
1. Kubernetes marking the endpoint as "ready=false" (immediate upon pod deletion initiation)
2. Actual load balancer ceasing to send NEW traffic
3. Existing connections completing or timing out
4. All backends confirming pod termination

This gap varies by:
- Load balancer type (cloud vs. nginx vs. envoy)
- Configuration (drain timeouts, proxy timeouts)
- Connection state (new vs. established)

**Resolution:** Kubernetes endpoint removal is not instantaneous from the client perspective. Proper zero-downtime deployment requires:
1. Readiness probes to prevent NEW traffic routing
2. Appropriate grace periods for connection completion
3. Load balancer drain settings coordinated with terminationGracePeriodSeconds
4. Application-level handling of SIGTERM to complete in-flight requests

---

## Contradiction 4: Database Column Addition – NULL Default vs Constant Default

**SOURCE A (PostgreSQL Documentation - Modifying Tables):** "Adding a column with a constant default value does not require each row of the table to be updated when the ALTER TABLE statement is executed... making the ALTER TABLE very fast even on large tables."

**SOURCE B (Evolutionary Database Design):** Recommends: "Add new columns as nullable with no default, then UPDATE to fill values, then add default and/or NOT NULL constraint as a second migration."

**ASSESSMENT:** Both approaches are valid but have different use cases:
- Constant default (e.g., `'pending'`) is metadata-only if the value is the same for all rows
- Volatile default (e.g., `clock_timestamp()`, `nextval()`) requires per-row update

The Evolutionary Database Design article emphasizes making separate migrations for each step (add column, update data, set constraints) which is safer for complex data migrations. PostgreSQL's optimization for constant defaults is an implementation detail that can be leveraged, but the principle of separating schema changes from data migrations remains.

**Resolution:** PostgreSQL's optimization makes constant defaults fast, but the evolutionary approach of small, separate migrations is still recommended because:
1. Complex data transformations (not just simple value assignment) still need UPDATE
2. Separate migrations allow for intermediate verification
3. Some operations (like setting NOT NULL) may still require table rewrites

---

## Contradiction 5: Queue Worker Graceful Termination – Laravel Horizon vs Plain PHP Workers

**SOURCE A (Laravel Horizon Documentation):** "You may gracefully terminate the Horizon process using the `horizon:terminate` Artisan command. Any jobs that are currently being processed will be completed and then Horizon will stop executing."

**SOURCE B (Laravel Queues Documentation):** "Setting `block_for` to `0` will cause queue workers to block indefinitely until a job is available. This will also prevent signals such as SIGTERM from being handled until the next job has been processed."

**ASSESSMENT:** There's a timing-dependent conflict:
- Laravel's default Redis queue worker blocks indefinitely waiting for jobs when `block_for=0`
- SIGTERM signal (used for graceful termination) cannot be processed during blocking wait
- Horizon's `terminate` command likely uses SIGINT or SIGTERM with specific handling

The documentation suggests `block_for` should NOT be 0 for graceful termination to work properly. This conflicts with the observation that some production configurations set very long or infinite `block_for` values for "always-on" workers.

**Resolution:** This is a configuration dependency:
- For graceful termination: set reasonable `block_for` (e.g., 1-5 seconds) so workers can process signals
- For always-on workers accepting that termination may take longer: use Supervisor's `stopsignal` and `stopwaitsecs` configuration to handle the delay
- Horizon provides a higher-level abstraction that manages worker lifecycle more gracefully than plain `queue:work`

---

## Conclusion on Contradictions

No fundamental conceptual contradictions exist in the evidence. The apparent disagreements fall into three categories:

1. **Pattern Selection**: Blue-green vs rolling deployment – different architectural choices, not disagreements
2. **Technical Constraints**: Tool capabilities (NGINX Open Source vs Plus, Kubernetes vs cloud LB) – cost/feature tradeoffs
3. **Configuration Dependencies**: block_for setting, stopwaitsecs values – operational decisions rather than contradictions

All sources agree on the core principles:
- Start new instances before stopping old ones
- Use health checks to gate traffic routing
- Separate database schema changes from application code deployments
- Plan for rollback scenarios
- Handle graceful shutdown for in-flight requests
