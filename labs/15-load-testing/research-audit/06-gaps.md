# Research Gaps

Target Lab: `labs/15-load-testing`  
Audit Scope: Missing evidence, disclaimed sources, and limitations in research artifacts.  

---

## Gap 1
- **Type:** OUTDATED_SOURCE / WEAK_SOURCE
- **Severity:** MEDIUM
- **Location:** `02-sources.md:Source 25`, `05-report.md:Limitations (Item 1)`
- **Problem:** Gatling primary documentation (`https://gatling.io/docs/gatling/guides/concepts/`) returned HTTP 403. While the vendor landing page (`https://gatling.io/open-source/`) confirmed Gatling's open-source edition and non-blocking model, technical claims regarding Scala DSL syntax and engine internals lack direct primary document verification.
- **Required Revision:** If Gatling is implemented as a primary lab engine or code example in future phases, fetch verified documentation or release notes via accessible mirrors.
- **Can Be Approved Without Fix:** YES (Gatling is treated as a secondary tool comparison; k6 and Locust are the primary focus).

---

## Gap 2
- **Type:** MISSING_SOURCE
- **Severity:** LOW
- **Location:** `02-sources.md:Source 24`, `06-open-questions.md:Weak Evidence`
- **Problem:** ISO/IEC 25010:2011 is a paywalled international standard. Its sub-characteristics (time behavior, resource utilization, capacity) were cited based on secondary engineering literature rather than full-text inspection.
- **Required Revision:** The research correctly disclaimed this source and removed it from primary evidence. In the future, verify through open academic citations or ISO publicly available summaries.
- **Can Be Approved Without Fix:** YES (Explicitly disclaimed and not relied upon for core technical conclusions).

---

## Gap 3
- **Type:** WEAK_SOURCE
- **Severity:** LOW
- **Location:** `02-sources.md:Source 22, 23`, `03-evidence.md:Evidence 14`
- **Problem:** Direct HTTP access to Apache JMeter official documentation (`https://jmeter.apache.org/`) timed out during research. JMeter protocol support (HTTP, JDBC, JMS, SOAP) was corroborated via Microsoft Azure Load Testing documentation (`Source 16`).
- **Required Revision:** Corroboration via Azure Load Testing is authoritative for cloud testing capabilities, but direct JMeter manual citations should be captured once connectivity is restored.
- **Can Be Approved Without Fix:** YES (Microsoft Learn is Tier 1 and confirms JMeter capabilities).

---

## Gap 4
- **Type:** UNVERIFIED_CLAIM
- **Severity:** MEDIUM
- **Location:** `03-evidence.md:Evidence 29`, `05-report.md:Limitations (Item 6)`
- **Problem:** The specific bottleneck triage decision tree when P95 jumps from 300ms to 2.5s at 800 VUs is a synthesis of general observability principles rather than an industry-standardized decision algorithm.
- **Required Revision:** The research explicitly notes this as synthesized guidance. Lab implementation and instructional draft should label this triage plan as an engineering heuristic/framework rather than an industry standard.
- **Can Be Approved Without Fix:** YES (Transparently qualified in limitations and open questions).
