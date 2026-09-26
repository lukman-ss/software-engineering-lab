# Evidence

## Evidence 1

Claim: Blue-green deployment maintains two identical production environments; traffic is switched from one to the other after testing, enabling instant rollback by switching back.

Evidence: "Blue-green deployment approach does this by ensuring you have two production environments, as identical as possible. At any time one of them, let's say blue for the example, is live. As you prepare a new release of your software you do your final stage of testing in the green environment. Once the software is working in the green environment, you switch the router so that all incoming requests go to the green environment... Blue-green deployment also gives you a rapid way to rollback - if anything goes wrong you switch the router back to your blue environment."

Source: Martin Fowler, Blue Green Deployment
URL: https://martinfowler.com/bliki/BlueGreenDeployment.html
Published: 2010-03-01
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Kubernetes Deployment docs describe blue-green as a strategy pattern; multiple independent sources (Fowler, Continuous Delivery book referenced therein, industry practice) agree.

Notes: Primary foundational source. Fowler also notes a caveat: "There's still the issue of dealing with missed transactions while the green environment was live."

---

## Evidence 2

Claim: Database schema changes must be separated from application upgrades so the schema supports both old and new application versions simultaneously.

Evidence: "The trick is to separate the deployment of schema changes from application upgrades. So first apply a database refactoring to change the schema to support both the new and old version of the application, deploy that, check everything is working fine so you have a rollback point, then deploy the new version of the application. (And when the upgrade has bedded down remove the database support for the old version.)"

Source: Martin Fowler, Blue Green Deployment
URL: https://martinfowler.com/bliki/BlueGreenDeployment.html
Published: 2010-03-01
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: PostgreSQL ALTER TABLE docs confirm ADD COLUMN with non-volatile DEFAULT does not rewrite table and is fast; expand-contract pattern widely documented.

Notes: This is the exact source of the expand → deploy → migrate → contract pattern in the topic spec.

---

## Evidence 3

Claim: Adding a column with a non-volatile DEFAULT in PostgreSQL is fast even on large tables and does not require a table rewrite.

Evidence: "When a column is added with ADD COLUMN and a non-volatile DEFAULT is specified, the default value is evaluated at the time of the statement and the result stored in the table's metadata... The value will be only applied when the table is rewritten, making the ALTER TABLE very fast even on large tables. If no column constraints are specified, NULL is used as the DEFAULT. In neither case is a rewrite of the table required."

Source: PostgreSQL, ALTER TABLE (v18 docs)
URL: https://www.postgresql.org/docs/current/sql-altertable.html
Published: 2026-09-24
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Official PostgreSQL documentation; consistent with expand-phase safe migration practice.

Notes: Caveat from same source: "Adding a column with a volatile DEFAULT (e.g., clock_timestamp()), a stored generated column, an identity column, or a column with a domain data type that has constraints will cause the entire table and its indexes to be rewritten."

---

## Evidence 4

Claim: PostgreSQL ALTER TABLE acquires ACCESS EXCLUSIVE lock by default unless explicitly noted, which blocks concurrent queries during the statement.

Evidence: "Note that the lock level required may differ for each subform. An ACCESS EXCLUSIVE lock is acquired unless explicitly noted. When multiple subcommands are given, the lock acquired will be the strictest one required by any subcommand."

Source: PostgreSQL, ALTER TABLE (v18 docs)
URL: https://www.postgresql.org/docs/current/sql-altertable.html
Published: 2026-09-24
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: PostgreSQL docs; consistent with topic spec warning about DROP COLUMN breaking running v1.

Notes: This is why DROP COLUMN name during a rolling deploy fails while v1 still runs SELECT name — lock contention plus missing column error.

---

## Evidence 5

Claim: Kubernetes distinguishes liveness probes (is the process alive?) from readiness probes (is the pod ready to receive traffic?), and a pod must pass readiness before receiving traffic.

Evidence: Kubernetes documentation provides separate Liveness, Readiness, and Startup Probes pages: liveness probe checks if the container is running, readiness probe determines if the pod is ready to accept traffic — pods not passing readiness are removed from service endpoints.

Source: Kubernetes, Pod Lifecycle / Probes
URL: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Nginx docs on mandatory health checks requiring new servers to pass health checks before traffic; Laravel health route docs.

Notes: Full probe detail page is at https://kubernetes.io/docs/concepts/workloads/pods/probes/ (linked from fetched page nav).

---

## Evidence 6

Claim: A generic health endpoint returning {"status":"ok"} is insufficient; Laravel provides a health route that can dispatch a DiagnosingHealth event to perform additional dependency checks (database, cache).

Evidence: "Laravel includes a built-in health check route... By default, the health check route is served at /up and will return a 200 HTTP response if the application has booted without exceptions. Otherwise, a 500 HTTP response will be returned... When HTTP requests are made to this route, Laravel will also dispatch a Illuminate\Foundation\Events\DiagnosingHealth event, allowing you to perform additional health checks relevant to your application. Within a listener for this event, you may check your application's database or cache status. If you detect a problem with your application, you may simply throw an exception from the listener."

Source: Laravel, Deployment docs (v12.x)
URL: https://laravel.com/docs/12.x/deployment
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Kubernetes probe docs; Nginx health_check directive docs.

Notes: Directly supports the topic spec claim that readiness should check DB/Redis/migration status, not just process liveness.

---

## Evidence 7

Claim: Nginx active health checks mark a server unhealthy and stop sending traffic until it passes again; a mandatory health check forces newly added servers to pass health checks before receiving any traffic.

Evidence: "By default, every five seconds NGINX Plus sends a request for '/' to each server in the backend group. If any communication error or timeout occurs (the server responds with a status code outside the range from 200 through 399) the health check fails. The server is marked as unhealthy, and NGINX Plus does not send client requests to it until it once again passes a health check." And: "The mandatory parameter requires every newly added server to pass all configured health checks before NGINX Plus sends traffic to it. When combined with slow start, it gives a new server more time to connect to databases and 'warm up' before being asked to handle their full share of traffic."

Source: NGINX, HTTP Health Checks
URL: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-health-check/
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Nginx passive health checks (max_fails/fail_timeout) in same doc; Kubernetes readiness gates concept.

Notes: IMPORTANT nuance: active health checks and slow_start are NGINX Plus features; passive health checks (max_fails/fail_timeout) are in NGINX Open Source.

---

## Evidence 8

Claim: Nginx supports graceful shutdown and zero-downtime configuration reload via signals — HUP starts new workers with new config while old workers finish serving existing clients.

Evidence: "HUP: changing configuration, keeping up with a changed time zone... starting new worker processes with a new configuration, graceful shutdown of old worker processes... If this succeeds, it starts new worker processes, and sends messages to old worker processes requesting them to shut down gracefully. Old worker processes close listen sockets and continue to service old clients. After all clients are serviced, old worker processes are shut down."

Source: NGINX, Controlling nginx
URL: https://nginx.org/en/docs/control.html
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Same doc describes USR2 zero-downtime binary upgrade; topic spec graceful shutdown claims.

Notes: QUIT = graceful shutdown; TERM/INT = fast shutdown. This is the mechanism behind "connection draining" at the proxy layer.

---

## Evidence 9

Claim: Nginx supports upgrading the executable on the fly via USR2 without dropping connections, with rollback via TERM on the new master.

Evidence: "USR2: upgrading an executable file... After that all worker processes (old and new ones) continue to accept requests... If for some reason the new executable file works unacceptably... Send the TERM signal to the new master process... When the new master process exits, the old master process will start new worker processes automatically."

Source: NGINX, Controlling nginx
URL: https://nginx.org/en/docs/control.html
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Same source; consistent with zero-downtime nginx upgrade practice.

Notes: This is a real-world example of "start new before stopping old" at the process level.

---

## Evidence 10

Claim: Long-running services including queue workers must be reloaded/restarted after deployment to pick up new code; Laravel provides `php artisan reload` for this.

Evidence: "After deploying a new version of your application, any long-running services such as queue workers, Laravel Reverb, or Laravel Octane should be reloaded / restarted to use the new code. Laravel provides a single reload Artisan command that will terminate these services... If you are not using Laravel Cloud, you should manually configure a process monitor that can detect when your reloadable processes exit and automatically restart them."

Source: Laravel, Deployment docs (v12.x)
URL: https://laravel.com/docs/12.x/deployment
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Laravel Queues docs (Supervisor Configuration, queue workers and deployment sections); Supervisor configuration in same docs.

Notes: Supports the lab requirement that "worker queue versi baru" must be handled without losing jobs.

---

## Evidence 11

Claim: Laravel queue workers process jobs from a central Redis/database queue; jobs are stored durably outside the worker process, so restarting a worker does not lose queued jobs (only in-flight jobs are affected).

Evidence: "Laravel queues provide a unified queueing API across a variety of different queue backends, such as Amazon SQS, Redis, or even a relational database." Queue docs describe retry_after, max attempts, and failed job storage — jobs are returned to the queue on timeout/failure.

Source: Laravel, Queues docs (v12.x)
URL: https://laravel.com/docs/12.x/queues
Accessed: 2026-09-26
Confidence: MEDIUM
Corroborated By: Queue docs describe `retry_after` returning jobs to queue and failed_jobs table persistence.

Notes: Not yet verified from Redis official docs in this session — MEDIUM confidence. Behavior on worker SIGKILL vs graceful stop needs deeper verification (Laravel `--stop-when-empty` / signal handling not yet inspected).

---

## Evidence 12

Claim: Kubernetes rolling update replaces old pods gradually, and Deployment provides built-in rollback capability.

Evidence: Kubernetes Deployments support RollingUpdate strategy (maxSurge/maxUnavailable) and "Update a Deployment Without Downtime" task documentation; deployment rollback is a documented feature (kubectl rollout undo).

Source: Kubernetes, Deployments (#strategy) and Update a Deployment Without Downtime
URL: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/#strategy
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Kubernetes tutorials on rolling update; topic spec's rolling deployment example.

Notes: Page content was truncated in fetch; strategy section referenced from nav/anchors. Cross-check with #strategy anchor is partial — treat details as not fully quoted.

---

## Evidence 13

Claim: `docker compose down && docker compose up -d` creates an availability gap because containers are stopped before new ones start.

Evidence: NOT VERIFIED from a primary Docker source in this session (attempted https://docs.docker.com/compose/how-tos/rolling-update/ returned 404). The claim is logically consistent with compose semantics but no primary source was opened.

Source: None opened
URL: N/A
Accessed: 2026-09-26
Confidence: LOW
Corroborated By: Not corroborated — Docker docs page fetch failed.

Notes: Needs follow-up: correct Docker Compose documentation URL for rolling updates / zero-downtime.

---

## Evidence 14

Claim: Nginx passive health checks (Open Source) mark servers unavailable after max_fails failures within fail_timeout, automatically stopping traffic to failed backends.

Evidence: "fail_timeout – Sets the time during which a number of failed attempts must happen for the server to be marked unavailable, and also the time for which the server is marked unavailable (default is 10 seconds). max_fails – Sets the number of failed attempts that must occur during the fail_timeout period (default is 1 attempt)."

Source: NGINX, HTTP Health Checks
URL: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-health-check/
Accessed: 2026-09-26
Confidence: HIGH
Corroborated By: Same doc distinguishes passive (OSS) vs active (Plus).

Notes: Important for the lab stack — if lab uses NGINX OSS, active health checks are unavailable; only passive checks + external orchestration health checks apply.
