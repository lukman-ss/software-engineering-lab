# Contradictions

## Contradiction 1: "Container running" vs "Application ready"

**Source A (Kubernetes):** States that liveness probes detect if a process is alive, and readiness probes determine if a pod is ready to accept traffic. A container can be running but not ready if the app has not finished connecting to its dependencies.

**Source B (Topic spec):** Claims that "proses Node/PHP/Java bisa hidup sementara database belum terkoneksi atau migration belum selesai" — process can be alive while DB is not connected.

**Assessment:** No contradiction. Both agree that process liveness ≠ application readiness. Kubernetes formalizes this with separate probes; the topic spec describes the same problem in Laravel terms.

---

## Contradiction 2: NGINX active health checks availability

**Source A (NGINX Plus docs):** Active health checks, mandatory health checks, and slow_start are explicitly labeled as NGINX Plus features requiring a subscription.

**Source B (NGINX Open Source):** Only passive health checks (max_fails/fail_timeout) are documented for NGINX Open Source.

**Assessment:** Not a contradiction — an important distinction for the lab. If the lab uses NGINX Open Source (free), active health checks are unavailable. The topic spec's claim that "instance baru jangan dimasukkan ke load balancer sebelum readiness berhasil" requires either NGINX Plus, an external orchestrator (Kubernetes), or custom scripts for NGINX OSS.

**Impact on lab design:** The lab must either specify NGINX Plus, use Kubernetes/External readiness gates, or implement a script-based approach for NGINX OSS.

---

## Contradiction 3: Blue-green vs rolling deployment resource cost

**Source A (Martin Fowler):** Blue-green requires "two production environments, as identical as possible" — implying double infrastructure cost.

**Source B (Topic spec):** Rolling deployment increments instances one-by-one, requiring only modest additional capacity at any time.

**Assessment:** Not a contradiction — these are two different trade-offs. Blue-green provides faster rollback but costs more infrastructure; rolling is more resource-efficient but rollback requires health-check-based traffic shifting (slower). Both are valid zero-downtime strategies with different cost profiles.

---

## Contradiction 4: PostgreSQL DROP COLUMN behavior

**Source A (PostgreSQL ALTER TABLE docs):** "DROP COLUMN form does not physically remove the column, but simply makes it invisible to SQL operations." Space is reclaimed over time as existing rows are updated.

**Source B (Topic spec):** "DROP COLUMN name" during v2 deployment would break v1's "SELECT name" because the column is "gone."

**Assessment:** Partial contradiction. PostgreSQL DROP COLUMN is fast and doesn't rewrite the table, BUT it makes the column invisible to all sessions immediately. V1's SELECT name would fail with a column-not-found error. The topic spec's "Boom." is correct about the failure, but the mechanism is not a table rewrite — it's a metadata-level visibility change that affects all queries immediately.

---

## Contradiction 5: Rolling deployment safety claims

**Source A (Kubernetes):** RollingUpdate uses maxSurge and maxUnavailable to control how many pods can be added/removed simultaneously.

**Source B (Topic spec):** Rolling deployment example shows one instance at a time going through health check before the next is updated.

**Assessment:** Not contradictory — Kubernetes allows configuration of maxSurge=1, maxUnavailable=0 to achieve exactly the topic spec's one-at-a-time pattern. The spec describes a specific configuration of the more general Kubernetes mechanism.
