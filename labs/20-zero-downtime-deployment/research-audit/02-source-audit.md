# Source Audit: Zero-Downtime Deployment

**Target Lab:** `labs/20-zero-downtime-deployment`  
**Research Run:** `2026-09-26-zero-downtime-deployment`

---

## Source 1

Claimed Title: Deployments | Kubernetes  
Claimed Publisher: The Kubernetes Authors  
URL: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None.  
Assessment: PASS  

---

## Source 2

Claimed Title: Pod Lifecycle | Kubernetes  
Claimed Publisher: The Kubernetes Authors  
URL: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None.  
Assessment: PASS  

---

## Source 3

Claimed Title: Liveness, Readiness, and Startup Probes | Kubernetes  
Claimed Publisher: The Kubernetes Authors  
URL: https://kubernetes.io/docs/concepts/workloads/pods/probes/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Canonical Kubernetes probe reference.  
Assessment: PASS  

---

## Source 4

Claimed Title: Blue Green Deployment  
Claimed Publisher: Martin Fowler  
URL: https://martinfowler.com/bliki/BlueGreenDeployment.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Minor date inconsistency in sources.md (published 2010, accessed 2026). Content is authentic.  
Assessment: PASS  

---

## Source 5

Claimed Title: HTTP Health Checks | NGINX Documentation  
Claimed Publisher: F5 NGINX  
URL: https://docs.nginx.com/nginx/admin-guide/load-balancer/http-health-check/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Distinguishes OSS passive checks vs Plus active checks.  
Assessment: PASS  

---

## Source 6

Claimed Title: Deployment | Laravel 11.x  
Claimed Publisher: Laravel  
URL: https://laravel.com/docs/11.x/deployment  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Accurately covers `/up` route and optimization commands.  
Assessment: PASS  

---

## Source 7

Claimed Title: Queues | Laravel 11.x  
Claimed Publisher: Laravel  
URL: https://laravel.com/docs/11.x/queues  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Accurately documents worker timeouts and `retry_after`.  
Assessment: PASS  

---

## Source 8

Claimed Title: Laravel Horizon  
Claimed Publisher: Laravel  
URL: https://laravel.com/docs/11.x/horizon  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Documents `horizon:terminate` and `stopwaitsecs`.  
Assessment: PASS  

---

## Source 9

Claimed Title: Evolutionary Database Design  
Claimed Publisher: ThoughtWorks / Pramod Sadalage & Martin Fowler  
URL: https://martinfowler.com/articles/evodb.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Publisher listed as ThoughtWorks; article hosted on Martin Fowler's site co-authored with Sadalage. Minor attribution nuance, content completely valid.  
Assessment: PASS  

---

## Source 10

Claimed Title: 5.7. Modifying Tables | PostgreSQL Documentation  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/ddl-alter.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Documents fast ADD COLUMN with constant default since PG 11.  
Assessment: PASS  

---

## Source 11

Claimed Title: Update a Deployment Without Downtime | Kubernetes  
Claimed Publisher: The Kubernetes Authors  
URL: https://kubernetes.io/docs/tasks/run-application/update-deployment-rolling/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None.  
Assessment: PASS  

---

## Source 12

Claimed Title: docker container stop | Docker Documentation  
Claimed Publisher: Docker  
URL: https://docs.docker.com/reference/cli/docker/container/stop/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Accurately documents SIGTERM -> grace period -> SIGKILL.  
Assessment: PASS  

---

## Source 13

Claimed Title: Server Configuration | PostgreSQL Documentation  
Claimed Publisher: PostgreSQL Global Development Group  
URL: https://www.postgresql.org/docs/current/runtime-config.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Broad source for lock timeouts and runtime config; appropriately used.  
Assessment: PASS  

---

## Source 14

Claimed Title: Pod Termination | Kubernetes  
Claimed Publisher: The Kubernetes Authors  
URL: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-termination  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Anchor link inside Source 2 (Pod Lifecycle). Canonical content.  
Assessment: PASS  
