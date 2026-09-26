# 02 — Source Audit

Target Lab: `labs/20-zero-downtime-deployment`
Research Run: `research/runs/2026-09-26-zero-downtime-deployment/`
Audit Date: 2026-09-26

---

## Source 1

Claimed Title: Deployment — Laravel 11.x docs
Claimed Publisher: Laravel
URL: https://laravel.com/docs/11.x/deployment

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 2

Claimed Title: Queues — Laravel 11.x docs
Claimed Publisher: Laravel
URL: https://laravel.com/docs/11.x/queues

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 3

Claimed Title: Laravel Octane — Laravel 11.x docs
Claimed Publisher: Laravel
URL: https://laravel.com/docs/11.x/octane

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 4

Claimed Title: Laravel Horizon — Laravel 11.x docs
Claimed Publisher: Laravel
URL: https://laravel.com/docs/11.x/horizon

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 5

Claimed Title: Task Scheduling — Laravel 11.x docs
Claimed Publisher: Laravel
URL: https://laravel.com/docs/11.x/scheduling

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 6

Claimed Title: Blue Green Deployment — bliki
Claimed Publisher: Martin Fowler (martinfowler.com)
URL: https://martinfowler.com/bliki/BlueGreenDeployment.html

Reachable: YES
Source Type: PRIMARY (Canonical pattern definition)
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 7

Claimed Title: Parallel Change — bliki
Claimed Publisher: Danilo Sato / martinfowler.com
URL: https://martinfowler.com/bliki/ParallelChange.html

Reachable: YES
Source Type: PRIMARY (Canonical pattern definition)
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 8

Claimed Title: Pod Lifecycle — Kubernetes docs
Claimed Publisher: Kubernetes
URL: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 9

Claimed Title: Liveness, Readiness, and Startup Probes — Kubernetes docs
Claimed Publisher: Kubernetes
URL: https://kubernetes.io/docs/concepts/workloads/pods/probes/

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 10

Claimed Title: Controlling nginx
Claimed Publisher: NGINX
URL: https://nginx.org/en/docs/control.html

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 11

Claimed Title: Module ngx_http_upstream_module
Claimed Publisher: NGINX
URL: https://nginx.org/en/docs/http/ngx_http_upstream_module.html

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: Correctly identifies `drain` and active health checks as commercial/NGINX Plus features.
Assessment: PASS

---

## Source 12

Claimed Title: Chapter 5. Data Definition — PostgreSQL 18 docs
Claimed Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/ddl.html

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 13

Claimed Title: ALTER TABLE — PostgreSQL docs
Claimed Publisher: PostgreSQL Global Development Group
URL: https://www.postgresql.org/docs/current/sql-altertable.html

Reachable: YES
Source Type: PRIMARY
Relevant: YES
Supports Claimed Topic: YES
Problems: None.
Assessment: PASS

---

## Source 14

Claimed Title: Redis upgrade doc (attempted)
Claimed Publisher: redis.io
URL: https://redis.io/docs/latest/operate/oss_and_stack/management/upgrading/

Reachable: NO (HTTP 404)
Source Type: UNKNOWN
Relevant: PARTIAL
Supports Claimed Topic: NO
Problems:
- The research agent explicitly tagged this source as a failed fetch and marked the corresponding claim as NOT VERIFIED.
- Properly segregated from supporting evidence.
Assessment: WARNING (Recorded as failed fetch / unverified, not fabricated)

---

## Source Audit Summary

- Total Sources Listed: 14
- PASS: 13
- WARNING: 1 (Properly declared failed fetch)
- FAIL: 0
- Source Integrity: PASS
