# Claim Audit

## Claim 1
Claim: NGINX supports graceful configuration reload using `HUP` signal, enabling zero-downtime worker replacement.

Location:
03-evidence.md Evidence 1, 05-report.md Finding 1

Evidence Provided:
When `nginx -s reload` (HUP signal) is sent, the master opens new configuration; if successful, new worker processes start and old workers receive graceful shutdown message, close listen sockets, continue serving existing connections until finished, then exit.

Source:
NGINX Control.html (https://nginx.org/en/docs/control.html)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately represents NGINX's well-documented zero-downtime reload behavior.

---

## Claim 2
Claim: Kubernetes Readiness probe determines pod inclusion in Service endpoints.

Location:
03-evidence.md Evidence 4, 05-report.md Finding 1

Evidence Provided:
Probe page defines: Liveness=restart policy; Readiness=governs Service endpoint membership; Startup=allow extended boot time.

Source:
Kubernetes Probes (https://kubernetes.io/docs/concepts/workloads/pods/probes/)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Core mechanism correctly identified for gating traffic.

---

## Claim 3
Claim: Laravel `queue:work` can receive SIGTERM for graceful termination; `queue:restart` command triggers graceful reload.

Location:
03-evidence.md Evidence 5, 05-report.md Finding 2

Evidence Provided:
Queues docs note signal handling; `queue:restart` sends message via cache to workers, which finish current job before exiting.

Source:
Laravel Queues (https://laravel.com/docs/11.x/queues)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately highlights the required behavior to prevent job loss during restarts.

---

## Claim 4
Claim: PostgreSQL `ADD COLUMN` without `NOT NULL` or with volatile default requires full rewrite; without those, it is metadata-only and fast.

Location:
03-evidence.md Evidence 9, 05-report.md Finding 3

Evidence Provided:
ALTER TABLE docs: adding column with IF NOT EXISTS or DEFAULT (non-volatile) is instant; changing type/triggering volatile default or adding constraint requires table rewrite + lock.

Source:
PostgreSQL ALTER TABLE (https://www.postgresql.org/docs/current/sql-altertable.html)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
HIGH

Notes:
Critical architectural claim about database backward-compatibility correctly backed by official PostgreSQL documentation.

---

## Claim 5
Claim: NGINX `upstream` supports `drain` mode for connection draining, but this is a commercial feature.

Location:
03-evidence.md Evidence 12, 05-report.md Finding 4 (and limitations)

Evidence Provided:
Upstream module docs: `drain` marks server as draining, only bound requests proxied. OSS health-check is open-source but `drain` parameter is commercial-only.

Source:
NGINX upstream_module (https://nginx.org/en/docs/http/ngx_http_upstream_module.html)

Source Actually Supports Claim:
YES

Classification:
IMPLEMENTATION-SPECIFIC

Severity:
MEDIUM

Notes:
Claim is accurately scoped and acknowledges the difference between NGINX OSS and Plus.

---

## Claim 6
Claim: Redis supports in-place binary upgrade / rolling upgrade without data loss.

Location:
03-evidence.md Evidence 14

Evidence Provided:
Attempted fetch returned HTTP 404; no authoritative source found.

Source:
redis.io (Attempted)

Source Actually Supports Claim:
NO

Classification:
HYPOTHESIS

Severity:
MEDIUM

Notes:
Researcher correctly marked this claim as NOT VERIFIED in the research phase due to the missing source.
