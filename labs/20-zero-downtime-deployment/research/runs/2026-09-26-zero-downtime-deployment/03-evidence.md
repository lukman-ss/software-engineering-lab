# Research Evidence: Zero-Downtime Deployment

---

## Evidence 1

**Claim:** Rolling deployment gradually replaces old Pods with new ones, keeping the application available throughout the process.

**Evidence:** "A rolling update gradually replaces old Pods with new ones, so your application remains available throughout the process." Kubernetes deployments support two update strategy types: RollingUpdate (default) and Recreate (which causes downtime). With RollingUpdate, the parameters `maxUnavailable` (default 25%) and `maxSurge` (default 25%) control how many pods can be unavailable/created simultaneously during an update.

**Source:** Kubernetes Documentation - Update a Deployment Without Downtime

**URL:** https://kubernetes.io/docs/tasks/run-application/update-deployment-rolling/

**Confidence:** HIGH

**Corroborated By:** Kubernetes Deployments documentation (Source 1) confirms the same rolling update mechanism.

**Notes:** The Kubernetes Deployment controller manages this automatically once RollingUpdate strategy is configured. The default behavior ensures at least 75% of pods remain available during updates (25% maxUnavailable).

---

## Evidence 2

**Claim:** Liveness probes determine whether to restart a container, while readiness probes determine whether a container should receive traffic.

**Evidence:** "Liveness probes determine when to restart a container." "Readiness probes determine when a container is ready to accept traffic." "If a container fails its readiness probe, the EndpointSlice controller removes the Pod's IP address from the EndpointSlices of all Services that match the Pod. Readiness probes run on the container during its whole lifecycle."

**Source:** Kubernetes Documentation - Liveness, Readiness, and Startup Probes

**URL:** https://kubernetes.io/docs/concepts/workloads/pods/probes/

**Confidence:** HIGH

**Corroborated By:** NGINX documentation (Source 5) describes similar health-check semantics for traffic routing, with "active health checks" determining which servers receive traffic.

**Notes:** Liveness probe failure triggers container restart (via restartPolicy: Always). Readiness probe failure only removes the pod from the load balancer endpoints. Both probes can use HTTP, TCP, gRPC, or exec mechanisms.

---

## Evidence 3

**Claim:** Blue-green deployment uses two identical production environments, where traffic is switched from old (blue) to new (green) all at once after full deployment and testing.

**Evidence:** "The blue-green deployment approach does this by ensuring you have two production environments, as identical as possible. At any time one of them, let's say blue for the example, is live. As you prepare a new release of your software you do your final stage of testing in the green environment. Once the software is working in the green environment, you switch the router so that all incoming requests go to the green environment - the blue one is now idle."

**Source:** Martin Fowler - Blue Green Deployment

**URL:** https://martinfowler.com/bliki/BlueGreenDeployment.html

**Published:** March 1, 2010

**Confidence:** HIGH

**Corroborated By:** Multiple Kubernetes task tutorials reference blue-green as a standard pattern. AWS documentation references the same pattern.

**Notes:** Rollback is immediate via router switch back to blue. Database changes must be deployed separately from application code according to Fowler.

---

## Evidence 4

**Claim:** Database schema changes that rename or remove columns cause backward compatibility failures when both old and new application versions are running simultaneously.

**Evidence:** "Databases can often be a challenge with this technique [blue-green], particularly when you need to change the schema to support a new version of the software. The trick is to separate the deployment of schema changes from application upgrades." The Evolutionary Database Design article demonstrates a transition phase pattern: "ALTER TABLE customer RENAME to client; CREATE VIEW customer AS SELECT id, first_name, last_name FROM client;" — a view provides backward compatibility while the table is renamed.

**Source:** Martin Fowler - Blue Green Deployment (Source 4), ThoughtWorks - Evolutionary Database Design (Source 9)

**URL:** https://martinfowler.com/bliki/BlueGreenDeployment.html, https://martinfowler.com/articles/evodb.html

**Published:** 2010, 2016

**Confidence:** HIGH

**Corroborated By:** PostgreSQL documentation (Source 10) confirms that column removal (`ALTER TABLE ... DROP COLUMN ... CASCADE`) requires explicit handling of foreign key dependencies.

**Notes:** The expand-deploy-migrate-contract pattern requires adding new columns first, migrating data, deploying application code that uses both old and new schema, then removing old columns in a final step. The transition phase with views is recommended for shared databases.

---

## Evidence 5

**Claim:** A container that is merely "running" does not guarantee the application is ready to serve traffic — health checks must verify database, cache, and dependency availability.

**Evidence:** "container that is merely 'running' doesn't mean the application is ready to accept traffic... The readiness probe checks if your application is actually ready to serve requests." Laravel's health route documentation states: "you may perform additional health checks relevant to your application. Within a listener for [DiagnosingHealth] event, you may check your application's database or cache status."

**Source:** Kubernetes Documentation - Liveness, Readiness, and Startup Probes (Source 3), Laravel Documentation (Source 6)

**URL:** https://kubernetes.io/docs/concepts/workloads/pods/probes/, https://laravel.com/docs/11.x/deployment

**Confidence:** HIGH

**Corroborated By:** NGINX documentation describes checking both status code (200-399) and response body for health check validation.

**Notes:** Kubernetes readiness probes run throughout the container lifecycle. A pod is removed from service endpoints when readiness probe fails, even before the container terminates.

---

## Evidence 6

**Claim:** Pod termination in Kubernetes follows a graceful shutdown flow: SIGTERM sent → grace period (default 30 seconds) → SIGKILL if process doesn't exit.

**Evidence:** "A pod receives a 'grace period' during which it can perform cleanup operations... The kubelet sends a TERM signal to all of the Pod's containers... Then, the kubelet sends a SIGKILL signal to the container, if it has not gone down for a while." Kubernetes documentation also notes: "when the Pod is deleted, the corresponding endpoint in the EndpointSlice will update its conditions: the endpoint ready condition will be set to false."

**Source:** Kubernetes Documentation - Pod Lifecycle (Source 14), Kubernetes Documentation - Pod Termination

**URL:** https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination

**Published:** 2026

**Confidence:** HIGH

**Corroborated By:** Docker documentation confirms the same pattern: "The main process inside the container will receive SIGTERM, and after a grace period, SIGKILL."

**Notes:** The endpoint ready condition is set to false before the pod is fully terminated, which removes it from service load balancing. The terminationGracePeriodSeconds can be configured at pod or probe level (since v1.28). preStop hooks allow custom pre-termination logic.

---

## Evidence 7

**Claim:** NGINX provides both passive and active health checks; active health checks require NGINX Plus (commercial) while passive health checks are available in open source.

**Evidence:** Passive checks use `max_fails` and `fail_timeout` to temporarily remove servers. Active checks use the `health_check` directive to periodically probe servers. The `mandatory` parameter ensures new servers pass health checks before receiving traffic. `slow_start` parameter allows gradual recovery: "A recently recovered server can be easily overwhelmed by connections."

**Source:** NGINX Documentation - HTTP Health Checks

**URL:** https://docs.nginx.com/nginx/admin-guide/load-balancer/http-health-check/

**Published:** 2026

**Accessed:** 2026-09-26

**Confidence:** HIGH

**Corroborated By:** Kubernetes documentation uses similar probe concepts (readiness checks) to achieve the same goal in an orchestrated environment.

**Notes:** Active health checks with `interval`, `fails`, `passes` parameters allow fine-grained control over health detection. The `mandatory` + `persistent` combination preserves server state across config reloads.

---

## Evidence 8

**Claim:** Laravel provides a built-in health endpoint at `/up` (default) that returns 200 if the application booted without exceptions, and allows custom health checks via the `DiagnosingHealth` event.

**Evidence:** "Laravel includes a built-in health check route... By default the health check route is served at /up and will return a 200 HTTP response if the application has booted without exceptions. Otherwise, a 500 HTTP response will be returned." "you may check your application's database or cache status. If you detect a problem with your application, you may simply throw an exception from the listener."

**Source:** Laravel Documentation - Deployment (Source 6)

**URL:** https://laravel.com/docs/11.x/deployment

**Published:** 2026

**Confidence:** HIGH

**Corroborated By:** The health endpoint is designed to be "used to report the status of your application to an uptime monitor, load balancer, or orchestration system such as Kubernetes."

**Notes:** The health route is configured in `bootstrap/app.php` via the `health` parameter on `withRouting()`. Users can customize the URI. The event listener approach allows extending health checks beyond basic boot verification.

---

## Evidence 9

**Claim:** Laravel Horizon provides graceful termination for queue workers via `php artisan horizon:terminate`, with Supervisor's `stopwaitsecs` ensuring long-running jobs complete before shutdown.

**Evidence:** "You may gracefully terminate the Horizon process using the `horizon:terminate` Artisan command. Any jobs that are currently being processed will be completed and then Horizon will stop executing." "During your application's deployment process, you should instruct the Horizon process to terminate so that it will be restarted by your process monitor and receive your code changes." "you should ensure that the value of `stopwaitsecs` is greater than the number of seconds consumed by your longest running job."

**Source:** Laravel Documentation - Horizon (Source 8)

**URL:** https://laravel.com/docs/11.x/horizon

**Published:** 2026

**Confidence:** HIGH

**Corroborated By:** Laravel Queues documentation mentions `retry_after` and `after_commit` options that also affect job lifecycle during deployment.

**Notes:** The `stopwaitsecs=3600` example in the docs is generous. Supervisor's behavior is: send SIGTERM, wait `stopwaitsecs`, then SIGKILL. This is analogous to Kubernetes' terminationGracePeriodSeconds.

---

## Evidence 10

**Claim:** PostgreSQL's `ALTER TABLE ... ADD COLUMN` with a constant default value does not rewrite the table — it is a metadata-only operation that is safe for concurrent access.

**Evidence:** "Adding a column with a constant default value does not require each row of the table to be updated when the ALTER TABLE statement is executed. Instead, the default value will be returned the next time the row is accessed, and applied when the table is rewritten, making the ALTER TABLE very fast even on large tables."

**Source:** PostgreSQL Documentation - Modifying Tables (Source 10)

**URL:** https://www.postgresql.org/docs/current/ddl-alter.html

**Published:** PostgreSQL 18, 2026

**Confidence:** HIGH

**Corroborated By:** Evolutionary Database Design article references PostgreSQL-compatible approaches (ADD COLUMN NULL, UPDATE, then set NOT NULL).

**Notes:** For volatile defaults (e.g., `clock_timestamp()`), each row IS updated, requiring a full table rewrite. This distinction is critical for zero-downtime deployments — only constant/default NULL column additions are safe.

---

## Evidence 11

**Claim:** Docker container stop sends SIGTERM and waits for a grace period (10s Linux, 30s Windows) before sending SIGKILL. The grace period is configurable via `--time` flag or Dockerfile's STOPSIGNAL.

**Evidence:** "The main process inside the container will receive SIGTERM, and after a grace period, SIGKILL. The first signal can be changed with the STOPSIGNAL instruction in the container's Dockerfile, or the --stop-signal option to docker run" "The default timeout can be specified using the --stop-timeout option when creating the container." "The default timeout... is 10 seconds for Linux containers, and 30 seconds for Windows containers."

**Source:** Docker Documentation - docker container stop

**URL:** https://docs.docker.com/reference/cli/docker/container/stop/

**Published:** 2026

**Confidence:** HIGH

**Corroborated By:** Kubernetes documentation confirms the same pattern at the pod level with terminationGracePeriodSeconds.

**Notes:** For Laravel applications, this means the application needs to handle SIGTERM signals to complete in-flight requests. PHP-FPM does not handle SIGTERM gracefully by default, so Nginx + PHP-FPM setups need special handling.

---

## Evidence 12

**Claim:** Database migrations should follow the expand-deploy-migrate-contract pattern: add new columns first (expand), deploy code that uses both (deploy), migrate existing data, then remove old columns (contract).

**Evidence:** From Blue Green Deployment: "The trick is to separate the deployment of schema changes from application upgrades. So first apply a database refactoring to change the schema to support both the new and old version of the application, deploy that, check everything is working fine so you have a rollback point, then deploy the new version of the application. (And when the upgrade has bedded down remove the database support for the old version.)"

**Source:** Martin Fowler - Blue Green Deployment (Source 4), ThoughtWorks - Evolutionary Database Design (Source 9)

**URL:** https://martinfowler.com/bliki/BlueGreenDeployment.html, https://martinfowler.com/articles/evodb.html

**Published:** 2010, 2016

**Confidence:** HIGH

**Corroborated By:** PostgreSQL ALTER TABLE documentation confirms that ADD COLUMN (with NULL default) is fast and reversible, while DROP COLUMN requires CASCADE for dependent objects.

**Notes:** This is the single most important pattern for database zero-downtime. The window between expand and contract is when both versions coexist. The "checking everything is working fine" step is the rollback checkpoint.

---

## Evidence 13

**Claim:** The `maxUnavailable` and `maxSurge` parameters in Kubernetes rolling updates control the trade-off between speed and availability during deployment.

**Evidence:** "For the RollingUpdate strategy, these parameters control how Kubernetes performs the update: maxUnavailable (default 25%) — Maximum number of Pods that can be unavailable during the update; maxSurge (default 25%) — Maximum number of extra Pods that can be created during the update." "Kubernetes calculates percentages from the desired replica count, rounding down for maxUnavailable and rounding up for maxSurge."

**Source:** Kubernetes Documentation - Update a Deployment Without Downtime

**URL:** https://kubernetes.io/docs/tasks/run-application/update-deployment-rolling/

**Published:** 2026

**Confidence:** HIGH

**Corroborated By:** Kubernetes Deployments documentation (Source 1) confirms the same parameters and defaults.

**Notes:** Setting maxUnavailable to 0 and maxSurge to 100% creates a fully parallel deployment. Setting maxUnavailable to 100% creates a recreate-like strategy.
