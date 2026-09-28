# Open Questions

## Unanswered Questions

1. **Optimal Time Window Selection**
   - Rolling window (4 weeks) vs. calendar window (monthly/quarterly): No data-driven consensus. Google recommends rolling; Evernote chose calendar; Home Depot undecided. Organizations must weigh user experience alignment vs business reporting needs.

2. **Burn Rate Threshold Customization**
   - Google's Table 5-8 provides starting parameters (2%/1h, 5%/6h, 10%/3d), but optimal values depend on traffic volume, on-call cadence, and incident response maturity. No authoritative guidance exists for tuning these thresholds.

3. **OpenSLO Ecosystem Adoption**
   - Primary evidence (OpenSLO spec website) exists, but independent industry adoption rates, real-world usage case studies, and interoperability with other tools (Prometheus, Datadog, Grafana) are not documented.

4. **Error Budget Policy Enforcement at Scale**
   - Google and Home Depot describe error budget policies, but specific automation for halting deployments (e.g., GitOps integration) is not standardized. No reference implementation exists for automated policy enforcement.

## Weak Evidence

1. **Low-Traffic Service Alerting**
   - Google SRE Workbook (p. 119-120) provides recommendations (artificial traffic, service combination, lower SLO), but no empirical data on false positive rates or intervention cost-benefit analysis for these techniques.

2. **Burn Rate Reset Time Metrics**
   - Google documents reset times (e.g., 58 minutes for single-burn-rate alert), but provides no field data on operator response times vs. budget recovery rates for multi-burn-rate configurations.

3. **VALET Framework Scalability**
   - Home Depot reports scaling to 800 services, but latency of updates (e.g., 3 days later still counts toward limit) and correction limit (100 one-time corrections per SLO) are documented but not validated against large-scale incidents.

## Claims Needing Deeper Research

1. **Do multi-dimensional SLOs (availability + latency) cause SLO sprawl?**
   - SRE Book advocates multiple targets; SRE Workbook implicitly suggests this is manageable. Quantitative study of SLO count vs team cognitive load is absent.

2. **Business Impact vs. Technical Degradation**
   - Evernote mentions "user feedback" as input for SLO targets; no methodology exists for translating business impact (e.g., payment webhook failure) to numeric reliability requirements.

3. **SLOs for Event-Driven vs. Request-Driven Systems**
   - "Freshness" SLIs for pipelines are documented, but no systematic comparison between event-driven architectures and request-response architectures for SLO design.

## Possible Next Research Directions

1. **Cross-vendor SLO implementation comparison**: Conduct side-by-side analysis of Google SRE Workbook recommendations vs Datadog's burn rate indicator vs Grafana's SLO plugin. Verify if core algorithms are identical.

2. **Empirical study of SLO alerting effectiveness**: Survey SRE teams on false positive/negative rates for multi-window multi-burn-rate alerts at different thresholds.

3. **SLO target selection methodology**: Develop a business-impact framework (cost of downtime × frequency × revenue impact) to derive SLO targets mathematically rather than qualitatively.

4. **Indonesian SRE community adoption**: Investigate whether Indonesian companies have adopted SLO practices (e.g., Gojek, Tokopedia, Bukalapak) and publish their SLO frameworks.