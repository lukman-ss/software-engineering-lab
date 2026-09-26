# Research Claim Audit: Zero-Downtime Deployment

## Claim 1

Claim: Rolling deployment gradually replaces old Pods with new ones, keeping the application available throughout the process, governed by `maxUnavailable` and `maxSurge` parameters (default 25%).

Location: `05-report.md:Finding 1`, `03-evidence.md:Evidence 1 & 13`

Evidence Provided: Official Kubernetes deployment and rolling update documentation citations.

Source: Source 1 (`kubernetes.io/docs/concepts/workloads/controllers/deployment/`), Source 11 (`kubernetes.io/docs/tasks/run-application/update-deployment-rolling/`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Accurately reflects Kubernetes Deployment controller specification and default percentage rounding rules.

---

## Claim 2

Claim: Liveness probes determine whether to restart a container, while readiness probes determine whether a container should receive traffic.

Location: `05-report.md:Finding 2`, `03-evidence.md:Evidence 2`

Evidence Provided: Official Kubernetes probe documentation quotes.

Source: Source 3 (`kubernetes.io/docs/concepts/workloads/pods/probes/`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Clean separation of concerns between container liveness and network traffic readiness.

---

## Claim 3

Claim: Blue-green deployment uses two identical production environments, switching router traffic all at once from blue to green after verification, allowing instant rollback.

Location: `05-report.md:Finding 3`, `03-evidence.md:Evidence 3`

Evidence Provided: Martin Fowler's canonical article text.

Source: Source 4 (`martinfowler.com/bliki/BlueGreenDeployment.html`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Martin Fowler notes that schema changes must be separated from application deployment for this to function safely.

---

## Claim 4

Claim: Database schema changes that rename or remove columns break backward compatibility when old and new application versions run simultaneously; expand-contract or views are necessary.

Location: `05-report.md:Finding 4`, `03-evidence.md:Evidence 4 & 12`

Evidence Provided: Evolutionary Database Design transition phase and view compatibility examples.

Source: Source 4, Source 9 (`martinfowler.com/articles/evodb.html`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Transition patterns (views, parallel columns, triggers) prevent query failures across versions.

---

## Claim 5

Claim: A running container does not guarantee application readiness; health checks must verify downstream dependencies such as databases and caches.

Location: `05-report.md:Finding 5`, `03-evidence.md:Evidence 5`

Evidence Provided: Kubernetes probe concepts and Laravel `DiagnosingHealth` event documentation.

Source: Source 3, Source 6 (`laravel.com/docs/11.x/deployment`)

Source Actually Supports Claim: YES

Classification: FACT / BEST_PRACTICE

Severity: LOW

Notes: Properly identifies risk of black-hole routing to uninitialized or disconnected instances.

---

## Claim 6

Claim: Pod termination follows SIGTERM → grace period (default 30s) → SIGKILL, while endpoint ready status is updated to false.

Location: `05-report.md:Finding 6`, `03-evidence.md:Evidence 6`

Evidence Provided: Kubernetes Pod lifecycle and termination documentation.

Source: Source 2, Source 14 (`kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Distinguishes between endpoint removal and process signal escalation.

---

## Claim 7

Claim: NGINX provides passive and active health checks; active health checks require NGINX Plus (commercial), while passive checks are available in open source.

Location: `05-report.md:Finding 7`, `03-evidence.md:Evidence 7`

Evidence Provided: NGINX HTTP Health Checks documentation.

Source: Source 5 (`docs.nginx.com/nginx/admin-guide/load-balancer/http-health-check/`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Explicitly calls out commercial boundary of NGINX active probes (`health_check` directive requires NGINX Plus).

---

## Claim 8

Claim: Laravel 11 includes a built-in `/up` health route returning 200 on boot success and dispatches `DiagnosingHealth` for custom checks.

Location: `05-report.md:Finding 8`, `03-evidence.md:Evidence 8`

Evidence Provided: Laravel 11.x deployment documentation.

Source: Source 6 (`laravel.com/docs/11.x/deployment`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Verified against official Laravel 11.x docs.

---

## Claim 9

Claim: Laravel Horizon terminates gracefully via `php artisan horizon:terminate`, and Supervisor `stopwaitsecs` must exceed the longest running job duration.

Location: `05-report.md:Finding 9`, `03-evidence.md:Evidence 9`

Evidence Provided: Laravel Horizon deployment documentation.

Source: Source 8 (`laravel.com/docs/11.x/horizon`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Accurately highlights Supervisor's termination timeout interaction.

---

## Claim 10

Claim: PostgreSQL `ALTER TABLE ... ADD COLUMN` with a constant default value does not rewrite the table in modern PostgreSQL and is a metadata-only operation.

Location: `05-report.md:Finding 10`, `03-evidence.md:Evidence 10`

Evidence Provided: PostgreSQL 18 Documentation (Modifying Tables section 5.7.1).

Source: Source 10 (`www.postgresql.org/docs/current/ddl-alter.html`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: MEDIUM

Notes: True for constant default values. Volatile default expressions (e.g. `clock_timestamp()`, `gen_random_uuid()`) or lock queue contention during DDL acquisition still present table rewrite or lock wait risks. Research acknowledged volatile expressions in Evidence 10.

---

## Claim 11

Claim: Docker container stop issues SIGTERM and waits for a grace period (10s Linux default) before SIGKILL, configurable via `--time` or `STOPSIGNAL`.

Location: `05-report.md:Finding 11`, `03-evidence.md:Evidence 11`

Evidence Provided: Docker CLI documentation for `docker container stop`.

Source: Source 12 (`docs.docker.com/reference/cli/docker/container/stop/`)

Source Actually Supports Claim: YES

Classification: FACT

Severity: LOW

Notes: Verified against official Docker documentation.
