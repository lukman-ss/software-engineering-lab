# Source Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

---

## Source 1

Claimed Title: RFC 6585 - Additional HTTP Status Codes (Section 4: 429 Too Many Requests)  
Claimed Publisher: IETF  
URL: https://www.rfc-editor.org/rfc/rfc6585.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None.

Assessment: PASS

---

## Source 2

Claimed Title: RFC 9110 - HTTP Semantics (Section 10.2.3 / 15.5.20)  
Claimed Publisher: IETF  
URL: https://www.rfc-editor.org/rfc/rfc9110.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None.

Assessment: PASS

---

## Source 3

Claimed Title: RFC 2697 - A Single Rate Three Color Marker (srTCM)  
Claimed Publisher: IETF  
URL: https://www.rfc-editor.org/rfc/rfc2697.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: RFC 2697 formalizes IP packet token bucket metering for DiffServ; application-level rate limiting adapts this principle to requests rather than byte counts. The research explicitly clarifies this scope.

Assessment: PASS

---

## Source 4

Claimed Title: RFC 6598 - IANA-Reserved IPv4 Prefix for Shared Address Space (100.64.0.0/10)  
Claimed Publisher: IETF  
URL: https://www.rfc-editor.org/rfc/rfc6598.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Confirms shared address space deployed in CGNAT contexts.

Assessment: PASS

---

## Source 5

Claimed Title: A Proof for the Queuing Formula: L = λW  
Claimed Publisher: J.D.C. Little, Operations Research (1961)  
URL: https://doi.org/10.1287/opre.9.3.383  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Canonical formal proof for Little's Law in steady-state queuing systems.

Assessment: PASS

---

## Source 6

Claimed Title: Exponential Backoff And Jitter  
Claimed Publisher: AWS Architecture Blog (Marc Brooker, 2015)  
URL: https://aws.amazon.com/blogs/architecture/exponential-backoff-and-jitter/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: Engineering blog, but canonical industry reference for Full Jitter formula and simulations.

Assessment: PASS

---

## Source 7

Claimed Title: Google Site Reliability Engineering / SRE Workbook (Addressing Cascading Failures)  
Claimed Publisher: Google / O'Reilly  
URL: https://sre.google/sre-book/addressing-cascading-failures/  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Authoritative production engineering reference on queue management (§22.10), server overload (§22.2), and retries (§22.8).

Assessment: PASS

---

## Source 8

Claimed Title: SQS CloudWatch Metrics  
Claimed Publisher: AWS Documentation  
URL: https://docs.aws.amazon.com/AWSSimpleQueueService/latest/SQSDeveloperGuide/sqs-cloudwatch-metrics.html  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Directly defines `ApproximateAgeOfOldestMessage` and `ApproximateNumberOfMessagesVisible`.

Assessment: PASS

---

## Source 9

Claimed Title: Kafka Monitoring  
Claimed Publisher: Apache Kafka Documentation  
URL: https://kafka.apache.org/documentation/#monitoring  

Reachable: YES  
Source Type: PRIMARY  
Relevant: YES  
Supports Claimed Topic: YES  

Problems: None. Official documentation for Consumer Lag and message queue latency metrics.

Assessment: PASS

---

## Source 10

Claimed Title: Netflix/Cloudflare/Stripe Engineering Articles  
Claimed Publisher: Various  
URL: N/A  

Reachable: NO (No specific URL provided)  
Source Type: SECONDARY / COMMUNITY  
Relevant: PARTIAL  
Supports Claimed Topic: PARTIAL  

Problems: Listed as secondary/illustrative without direct canonical links. However, research explicitly notes this as non-primary and relies on Sources 1-9 for all factual claims.

Assessment: WARNING
