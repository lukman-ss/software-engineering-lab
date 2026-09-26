# Research Gaps

## Gap 1

Type:
MISSING_SOURCE

Severity:
MEDIUM

Location:
03-evidence.md Evidence 14, 06-open-questions.md

Problem:
Official documentation for Redis cluster rolling upgrades was not found (404), leaving the Redis zero-downtime upgrade behavior unverified.

Required Revision:
The research should track down the correct Redis documentation (e.g. from Sentinel or Cluster topologies) or the lab implementation scope must formally exclude Redis version upgrades during zero-downtime application deployments.

Can Be Approved Without Fix:
YES (The researcher explicitly flagged this as NOT VERIFIED and added it to open questions).

---

## Gap 2

Type:
IMPLEMENTATION_GAP

Severity:
LOW

Location:
03-evidence.md Evidence 12, 05-report.md Limitations

Problem:
The NGINX upstream `drain` and active `health_check` parameters are commercial features (NGINX Plus). Open-source users must rely on worker reload (`nginx -s reload`) or external orchestrators (like Kubernetes Service endpoint updates).

Required Revision:
Ensure that any subsequent lab engineering implementation uses the NGINX OSS-compatible `nginx -s reload` pattern, or relies exclusively on the orchestrator (e.g. Docker Compose or K8s) for draining connections.

Can Be Approved Without Fix:
YES (The researcher correctly identified and documented this limitation).
