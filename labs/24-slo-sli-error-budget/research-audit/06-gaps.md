# 06 Research Gap Analysis

Target Lab: `labs/24-slo-sli-error-budget`  
Audit Date: 2026-09-26

## Gap 1
- **Type:** OVERGENERALIZATION
- **Severity:** MEDIUM
- **Location:** `03-evidence.md` §Evidence 17, `05-report.md` §Limitations
- **Problem:** The figure "Changes represent roughly 70% of outages" originates from Google SRE Workbook Appendix B as internal background context without cited methodology or broader industry replication.
- **Required Revision:** Ensure lab text explicitly clarifies that this figure is an internal historical observation from Google, not an undisputed universal industry statistic.
- **Can Be Approved Without Fix:** YES (The research report already appropriately marked this with LOW confidence and recorded it in §Limitations).

---

## Gap 2
- **Type:** SCOPE_ERROR (Single-vendor source distribution)
- **Severity:** MEDIUM
- **Location:** `02-sources.md`, `05-report.md` §Limitations
- **Problem:** 100% of the cited Tier 1 sources are published by Google (Google SRE Book & Google SRE Workbook). No third-party engineering blogs or cross-cloud literature (e.g. AWS Builders Library, CNCF, Microsoft Azure architecture center) are incorporated.
- **Required Revision:** While Google SRE established the modern terminology and best practices for SLI/SLO/Error Budget, future revisions can incorporate cross-cloud references or case studies (e.g., Alex Hidalgo's SLO Book or AWS Builders' Library).
- **Can Be Approved Without Fix:** YES (Google SRE Book/Workbook is the canonical, universally recognized origin of these principles).

---

## Gap 3
- **Type:** MISSING_CASE
- **Severity:** LOW
- **Location:** `06-open-questions.md` §Unanswered Questions (Item 3)
- **Problem:** Research focus is heavily weighted toward request-driven synchronous HTTP/RPC workloads; guidance for asynchronous background jobs (e.g., Kafka streaming consumers, daily ETL batch pipelines) is only briefly acknowledged.
- **Required Revision:** Provide brief conceptual examples of SLIs for event-driven systems (e.g., consumer lag, pipeline freshness/throughput) during lab implementation or documentation.
- **Can Be Approved Without Fix:** YES (Request-driven APIs serve as the primary teaching vehicle for Lab 24).
