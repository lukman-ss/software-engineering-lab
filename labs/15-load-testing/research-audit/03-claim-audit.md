# Claim Audit

## Claim 1: Load testing has four standardized types (Load, Stress, Spike, Soak)
- Location: `research/05-report.md:Finding 1`
- Evidence Provided: k6 test type documentation, Google SRE Book, Azure Well-Architected Framework.
- Source: Sources 1, 2, 3, 4, 10, 15
- Source Actually Supports Claim: YES
- Classification: FACT
- Severity: LOW (Strongly supported across multiple authoritative sources)
- Notes: Highly consistent definitions found across cloud, tool, and engineering literature.

## Claim 2: Percentiles (P95, P99) are essential over averages to catch tail latency
- Location: `research/05-report.md:Finding 3`
- Evidence Provided: Google SRE Monitoring chapter notes tail latency issues (e.g., 1% of requests taking 50x median). k6 docs demonstrate latency histograms.
- Source: Sources 5, 6, 12
- Source Actually Supports Claim: YES
- Classification: FACT
- Severity: LOW (Standard engineering principle)
- Notes: Supported both mathematically and empirically.

## Claim 3: JMeter supports multiple protocols (HTTP, JDBC, JMS, SOAP, FTP)
- Location: `research/05-report.md:Finding 5`
- Evidence Provided: Secondary confirmation via Azure Load Testing documentation; JMeter official docs timed out.
- Source: Sources 16, 22
- Source Actually Supports Claim: YES (via authoritative secondary source)
- Classification: FACT
- Severity: LOW (Well-known characteristic, though direct docs timed out)
- Notes: Validated via Microsoft documentation for Azure Load Testing integration.

## Claim 4: Real external API calls should be used in load tests
- Location: `research/05-report.md:Finding 7`
- Evidence Provided: Azure Well-Architected Framework PE:06 notes mock hiding real performance issues.
- Source: Source 15
- Source Actually Supports Claim: PARTIAL (Context-dependent)
- Classification: INTERPRETATION
- Severity: MEDIUM
- Notes: The research agent properly qualified this: Real calls are for sandbox/controlled environments, while high-volume stress tests require mocks to prevent rate-limit violations, bans, and excessive cost.

## Claim 5: Gatling has an async/non-blocking architecture
- Location: `research/05-report.md:Finding 5`
- Evidence Provided: Gatling Open Source landing page (high level). Primary architecture docs 403 Forbidden.
- Source: Source 25
- Source Actually Supports Claim: PARTIAL
- Classification: HYPOTHESIS / UNVERIFIED DETAIL
- Severity: MEDIUM
- Notes: High-level claims exist, but detailed architectural evidence is lacking due to inaccessible docs. The Research Report appropriately downgraded confidence to MEDIUM.
