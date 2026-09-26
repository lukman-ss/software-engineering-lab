# 02 Source Audit

Target Lab: `labs/24-slo-sli-error-budget`  
Audit Date: 2026-09-26

## Source 1
Claimed Title: Chapter 4 - Service Level Objectives  
Claimed Publisher: Google SRE Book (Google Inc.) / O'Reilly Media  
URL: https://sre.google/sre-book/service-level-objectives/  
Reachable: YES (HTTP 200)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Canonical definitions of SLI, SLO, SLA, percentile aggregation, and safety margin.  
Assessment: PASS

## Source 2
Claimed Title: Chapter 6 - Monitoring Distributed Systems  
Claimed Publisher: Google SRE Book (Google Inc.) / O'Reilly Media  
URL: https://sre.google/sre-book/monitoring-distributed-systems/  
Reachable: YES (HTTP 200)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Authoritative source for the Four Golden Signals (Latency, Traffic, Errors, Saturation) and tail latency.  
Assessment: PASS

## Source 3
Claimed Title: Chapter 10 - Practical Alerting from Time-Series Data  
Claimed Publisher: Google SRE Book (Google Inc.) / O'Reilly Media  
URL: https://sre.google/sre-book/practical-alerting/  
Reachable: YES (HTTP 200)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Practical Borgmon time-series alerting, alert thresholds, and duration filtering.  
Assessment: PASS

## Source 4
Claimed Title: Chapter 3 - Embracing Risk  
Claimed Publisher: Google SRE Book (Google Inc.) / O'Reilly Media  
URL: https://sre.google/sre-book/embracing-risk/  
Reachable: YES (HTTP 200)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Explains motivation for error budgets, unreliability risk tolerance, and cost scaling per nine.  
Assessment: PASS

## Source 5
Claimed Title: Appendix A - Availability Table  
Claimed Publisher: Google SRE Book (Google Inc.) / O'Reilly Media  
URL: https://sre.google/sre-book/availability-table/  
Reachable: YES (HTTP 200)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Concrete table of allowed downtime across various time windows.  
Assessment: PASS

## Source 6
Claimed Title: Chapter 5 - Alerting on SLOs  
Claimed Publisher: Google SRE Workbook (Google Inc.) / O'Reilly Media  
URL: https://sre.google/workbook/alerting-on-slos/  
Reachable: YES (HTTP 200)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Definitive guide to burn rate calculation and multiwindow multi-burn-rate alerting with PromQL examples.  
Assessment: PASS

## Source 7
Claimed Title: Chapter 2 - Implementing SLOs  
Claimed Publisher: Google SRE Workbook (Google Inc.) / O'Reilly Media  
URL: https://sre.google/workbook/implementing-slos/  
Reachable: YES (HTTP 200)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Step-by-step implementation guide, SLI specification, 4-week rolling window recommendation, and decision matrix.  
Assessment: PASS

## Source 8
Claimed Title: Appendix B - Example Error Budget Policy  
Claimed Publisher: Google SRE Workbook (Google Inc.) / O'Reilly Media  
URL: https://sre.google/workbook/error-budget-policy/  
Reachable: YES (HTTP 200)  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Note that the "70% outages due to change" statement is an internal Google observation cited as background context rather than an industry-wide empirical standard.  
Assessment: PASS
