# Research Gap Analysis: Zero-Downtime Deployment

## Gap 1

Type: IMPLEMENTATION_GAP  
Severity: MEDIUM  
Location: `05-report.md:Limitations`, `06-open-questions.md:Question 1 & 2`  
Problem: The behavior of PHP-FPM under SIGTERM and request draining with NGINX in traditional PHP architectures is left as an open question without definitive configuration parameters (e.g. `process_control_timeout`).  
Required Revision: Future implementation notes should explicitly define PHP-FPM `process_control_timeout` and NGINX `proxy_next_upstream` settings for graceful draining.  
Can Be Approved Without Fix: YES (Documented honestly in `06-open-questions.md`).

---

## Gap 2

Type: WEAK_SOURCE  
Severity: LOW  
Location: `06-open-questions.md:Weak Evidence Area 1`  
Problem: Lack of empirical load test metrics showing exact error rate deltas between zero-downtime rolling deployment and standard recreate deployments under heavy concurrent traffic.  
Required Revision: Empirical load benchmark should be collected during lab execution/demo phases.  
Can Be Approved Without Fix: YES (Research phase appropriately identifies the need for empirical validation in lab).

---

## Gap 3

Type: SCOPE_ERROR  
Severity: LOW  
Location: `05-report.md:Limitations`  
Problem: Multi-region distributed database replication synchronization during DDL migration is out of scope.  
Required Revision: Explicitly bounds the research scope to single-cluster / single primary database topology.  
Can Be Approved Without Fix: YES (Sufficiently declared in limitations).
