# Research Plan

## Research Topic
SLO, SLI & Error Budget — Mengubah "Sistem Harus Stabil" Menjadi Angka yang Bisa Diputuskan

## Objective
To understand the concepts of SLI, SLO, and Error Budget, their interrelation, best practices, and common pitfalls in the context of reliability engineering for senior software engineers.

## Research Questions
1. What is SLI (Service Level Indicator) and what are common examples in practice?
2. What is SLO (Service Level Objective) and how is it properly set?
3. What is Error Budget and how is it calculated from SLO?
4. How do SLI, SLO, and Error Budget relate to each other?
5. What are best practices for setting SLOs that align with business needs?
6. What are common pitfalls in implementing SLO/SLI/Error Budget?
7. How is burn rate used in relation to error budget for alerting?
6. How should different criticality endpoints have different SLOs?
7. Why are infrastructure metrics (CPU, RAM) poor SLOs compared to user-facing metrics?

## Search Strategy
- Search authoritative sources: Google SRE Book, SLO Book by Alex Hidalgo, cloud provider documentation (GCP, AWS, Azure)
- Search for practical implementations and case studies
- Search for burn rate alerting methodology
- Focus on primary sources (official docs, books) and reputable technical publications

## Expected Primary Sources
- Google SRE Book: https://sre.google/sre-book/table-of-contents/
- Google SRE Workbook: https://sre.google/workbook/table-of-contents/
- SLO Book by Alex Hidalgo: https://slobook.com/
- GCP SLO Monitoring: https://cloud.google.com/monitoring/slos
- AWS Builders Library: https://aws.amazon.com/builders-library/monitoring-service-level-objectives/
- Azure Monitor SLO: https://learn.microsoft.com/en-us/azure/azure-monitor/app/service-level-objectives

## Risks / Unknowns
- Rapid evolution of observability tools might make some sources outdated
- Differing interpretations of SLI/SLO across organizations
- Language barrier: primary sources in English, target lab in Bahasa Indonesia