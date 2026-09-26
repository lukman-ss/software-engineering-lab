# Open Questions and Unanswered Areas

## Unanswered Questions

1. **Industry-wide SLO adoption rates**: What percentage of engineering organizations actually implement SLO-based error budgeting in production, and with what success rates?

2. **Optimal SLO time windows**: While quarterly SLOs are standard at Google, what evidence exists for the optimal time window for other organization sizes and service types (e.g., 7-day, 30-day, 90-day)?

3. **Burn rate thresholds standard**: The specific numeric thresholds for burn rate alerting (e.g., Datadog's 1-6 elevated, 6+ critical) appear vendor-specific. Is there an authoritative standard for burn rate thresholds from the original SRE practice?

4. **Error budget policy implementation**: How do organizations without dedicated SRE teams implement error budget-based release controls?

5. **Multi-service error budget allocation**: When an organization has multiple services with different SLOs, how should error budget be allocated across services?

6. **SLO for data correctness**: How to practically measure data correctness as a service level indicator when correctness is often a property of data rather than infrastructure?

## Weak Evidence Areas

1. **Burn rate numeric thresholds**: The burn rate ranges (1-6, 6+) come solely from Datadog's implementation documentation. Google SRE Book discusses error budget consumption but does not provide specific numeric thresholds for "elevated" vs "critical" burn rates.

2. **Error budget remaining formula**: The specific formula (error budget remaining = 100 * (current - target) / (100 - target)) comes from Datadog. No independent source corroborates this specific formula, though it is mathematically sound.

3. **Cost of reliability (100x per nine)**: The claim that "each additional nine costs ~100x more" is stated in the SRE Book but lacks specific empirical data or methodology for calculation.

4. **User preference for lower variance**: The claim that "people prefer a slightly slower system to one with high variance" lacks cited references for the underlying user studies.

## Claims Needing Deeper Research

1. **Availability calculation in topic specification**: The topic specification states "99% → ~7 jam 18 menit/bulan" while Google's authoritative source shows 7.2 hours (7 hours 12 minutes). A deeper verification across cloud providers (AWS, Azure, GCP) would confirm the standard values.

2. **SLO impact on deployment frequency**: While the SRE Book claims error budgets drive release velocity, empirical data on actual deployment frequency changes when SLO-based controls are implemented is limited.

3. **SLO failure rate calculation**: The topic specification provides a specific scenario (Payment Webhook with 15 failures out of 200,000 requests) - this should be verified against the actual SLO error budget math.

## Possible Next Research Directions

1. **Vendor-neutral SLO implementation guides**: Compare AWS CloudWatch SLOs, GCP SLOs, and Azure Monitor to verify consistency with Google SRE principles.

2. **Case studies**: Search for detailed case studies from companies that implemented SLO-based error budgeting (Netflix, Shopify, etc.).

3. **Academic research**: Search academic databases for peer-reviewed research validating or challenging SLO/SLI effectiveness metrics.

4. **Error budget burn rate analysis**: Research the mathematical basis for burn rate alerting thresholds and whether there's a consensus beyond vendor implementations.

5. **Industry surveys**: Look for SRE surveys (e.g., CNCF survey, Splunk SLO survey) that provide quantitative data on SLO adoption and effectiveness.