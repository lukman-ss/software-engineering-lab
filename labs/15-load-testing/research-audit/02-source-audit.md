# Source Audit: Lab 15 (Load Testing)

## Source 1
Claimed Title: Load test types  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Canonical guide defining test types (smoke, average-load, stress, spike, soak, breakpoint).  
Assessment: PASS  

## Source 2
Claimed Title: Average-load testing  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/load-testing/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Authoritative definition of average-load testing, ramp-up ratios, plateau stages.  
Assessment: PASS  

## Source 3
Claimed Title: Stress testing  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/stress-testing/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Documents testing beyond average capacity without arbitrary fixed increments.  
Assessment: PASS  

## Source 4
Claimed Title: Spike testing  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/spike-testing/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Documents steep surges in traffic and recovery evaluation.  
Assessment: PASS  

## Source 5
Claimed Title: Soak testing  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/soak-testing/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Documents long-duration tests for memory/resource leaks.  
Assessment: PASS  

## Source 6
Claimed Title: Breakpoint testing  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/breakpoint-testing/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Documents capacity limit hunting and warning regarding cloud auto-scaling bill escalation.  
Assessment: PASS  

## Source 7
Claimed Title: Smoke testing  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/smoke-testing/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Duplicate citation with Source 17.  
Assessment: WARNING  

## Source 8
Claimed Title: Ramping arrival rate  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/using-k6/scenarios/executors/ramping-arrival-rate/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None.  
Assessment: PASS  

## Source 9
Claimed Title: Testing for Reliability (Chapter 17)  
Claimed Publisher: Google SRE (Site Reliability Engineering) Book  
URL: https://sre.google/sre-book/testing-reliability/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. High-authority architectural guidance on testing at scale and reliability equivalence.  
Assessment: PASS  

## Source 10
Claimed Title: ISO/IEC 25010  
Claimed Publisher: International Organization for Standardization (ISO) / Wikipedia  
URL: https://en.wikipedia.org/wiki/ISO/IEC_25010  
Reachable: YES  
Source Type: SECONDARY / COMMUNITY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Source is cited as Tier 1 standard, but URL provided is a Wikipedia article. While content accurately summarizes Performance Efficiency (Time behaviour, Resource utilisation, Capacity), secondary Wikipedia references should be clearly designated as Tier 2/Community summary, not Tier 1 primary standard.  
Assessment: WARNING  

## Source 11
Claimed Title: Architecture Strategies for Performance Testing  
Claimed Publisher: Microsoft Azure Well-Architected Framework  
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Comprehensive cloud architecture guidance on performance testing, antipatterns, and iterative tuning.  
Assessment: PASS  

## Source 12
Claimed Title: Performance Efficiency Pillar  
Claimed Publisher: AWS Well-Architected Framework  
URL: https://docs.aws.amazon.com/wellarchitected/latest/performance-efficiency-pillar/welcome.html  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: High-level architectural overview; provides supporting framework context rather than specific testing implementation metrics.  
Assessment: PASS  

## Source 13
Claimed Title: What is Locust?  
Claimed Publisher: Locust Official Documentation  
URL: https://docs.locust.io/en/stable/what-is-locust.html  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Valid primary documentation for Locust tool evaluation.  
Assessment: PASS  

## Source 14
Claimed Title: Your first test  
Claimed Publisher: Locust Official Documentation  
URL: https://docs.locust.io/en/stable/quickstart.html  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None.  
Assessment: PASS  

## Source 15
Claimed Title: Thresholds  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/using-k6/thresholds/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Documents metric assertions, percentiles syntax, and abort criteria.  
Assessment: PASS  

## Source 16
Claimed Title: Introduction to the Spring IoC Container and Beans  
Claimed Publisher: VMware Spring Framework Documentation  
URL: https://docs.spring.io/spring-framework/reference/core/beans/introduction.html  
Reachable: YES  
Source Type: PRIMARY  
Relevant: NO  
Supports Claimed Topic: NO  
Problems: Leftover artifact from Dependency Injection lab (Lab 16). Explicitly marked as excluded in notes, but should not remain in active source catalog for Load Testing.  
Assessment: FAIL  

## Source 17
Claimed Title: Smoke testing (k6)  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/testing-guides/test-types/smoke-testing/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Duplicate of Source 7. Redundant entry.  
Assessment: WARNING  

## Source 18
Claimed Title: Calculate concurrent users for load tests  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/testing-guides/calculate-concurrent-users/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Critical primary reference for user concurrency sizing calculations.  
Assessment: PASS  

## Source 19
Claimed Title: Architecture Strategies for Performance Testing (duplicate of Source 11)  
Claimed Publisher: Microsoft Azure Well-Architected Framework  
URL: https://learn.microsoft.com/en-us/azure/well-architected/performance-efficiency/performance-test  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Explicitly acknowledged as duplicate of Source 11 in text. Redundant entry inflating source count.  
Assessment: WARNING  

## Source 20
Claimed Title: Built-in metrics reference  
Claimed Publisher: Grafana k6 Documentation  
URL: https://grafana.com/docs/k6/latest/using-k6/metrics/reference/  
Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: None. Full breakdown of HTTP request phases (`http_req_waiting`, etc.).  
Assessment: PASS  

## Source 21
Claimed Title: What is HTTP and how does it work  
Claimed Publisher: Mozilla Developer Network (MDN)  
URL: https://developer.mozilla.org/en-US/docs/Web/HTTP/Overview  
Reachable: YES  
Source Type: PRIMARY / SECONDARY  
Relevant: YES  
Supports Claimed Topic: YES  
Problems: Generic background context on HTTP protocol.  
Assessment: PASS  
