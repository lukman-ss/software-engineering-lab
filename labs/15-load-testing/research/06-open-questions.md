# Open Questions

## Unanswered Questions
1. **Gatling specifics:** Capabilities, architecture, Scala DSL details not verified — docs returned 403. How does Gatling compare quantitatively to k6/Locust/JMeter?
2. **Initial VU selection:** How to determine starting VU count for Booking Bengkel app? No source provided formula linking business traffic (e.g., concurrent bookings) to VU count.
3. **Bottleneck triage decision tree:** No single authoritative source provides step-by-step to distinguish app vs DB vs external API bottleneck when P95 spikes. Synthesized from fragments (Evidence 29, MEDIUM).
4. **JMeter vs Locust vs k6 head-to-head benchmarks:** No independent benchmark comparing same API under same load for throughput, resource overhead, max VUs per generator.
5. **Cost/fidelity tradeoff for prod-like env:** When is scaled-down staging sufficient vs exact prod mirror? Azure says "depends on risk profile" but no quantitative cost-benefit model.
6. **Soak test minimal duration threshold:** k6 lists 3h–72h as typical, but no source defines when 3h is enough vs 24h required.
7. **Real incident case studies:** Netflix/Uber blogs blocked (403/404). No primary case study of HR absensi massal or CMMS WhatsApp bottleneck verified externally.

## Weak Evidence
- **Evidence 15 (Gatling):** MEDIUM (reviser 2026-09-26: vendor page https://gatling.io/open-source/ verified; primary docs still 403)
- **Evidence 14 (JMeter multi-protocol):** MEDIUM-HIGH (reviser 2026-09-26: Azure docs provide authoritative confirmation; direct JMeter component ref still timed out)
- **Evidence 29 (P95 300ms→2.5s investigation steps):** MEDIUM — synthesized best-practice, not single-source authoritative procedure
- **Source 24 (ISO/IEC 25010):** NOT VERIFIED — paywalled, referenced but not opened

## Claims Needing Deeper Research
- Quantitative percentile thresholds: Is P95 <500ms with Error <1% appropriate for Booking Bengkel? Industry SLA benchmarks by domain (HR, CMMS, booking) not found.
- Database connection pool exhaustion vs query slowness: No source details how to distinguish DB bottleneck subtypes during load test.
- Queue behavior under spike: How to test message queue (e.g., WhatsApp notification queue) specifically? Not covered.
- Locust greenlet overhead limits: At what VU count does single Locust process degrade? Needs empirical test.
- Azure Load Testing auto-stop thresholds: Not investigated in detail.

## Possible Next Research Directions
1. Bench same endpoint with k6/Locust/JMeter — measure generator CPU, max RPS, accuracy.
2. Access Gatling docs via alternate mirror; interview Gatling enterprise users.
3. Retrieve Google SRE Workbook chapters on overload/capacity planning for deeper bottleneck methodology.
4. Collect public postmortems where load testing absent — quantify impact (e.g., 504 at 2.5k concurrent users).
5. Research observability integration: OpenTelemetry tracing correlation during load test to pinpoint layer (app/DB/external) automatically.
6. Investigate data seeding strategies for realistic DB state (millions of rows, index warmup).
7. Study chaos + load testing intersection (k6 + xk6-disruptor, Azure Chaos Studio) for resilience validation.
8. Survey senior engineers on Booking Bengkel exercise answers to validate Investigation Plan (question 5).
