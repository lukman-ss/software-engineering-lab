# Open Questions

## Original Open Questions (from research/01-plan.md)

### 1. Numeric Recommendations Are Illustrative

**Original Status**: Identified in plan Risks/Unknowns
**Resolution**: ADDRESSED in `02-findings.md`

Plan claimed: "Beberapa klaim lab adalah ilustrasi numerik internal (mis. 500 req/detik, 200 user per IP)"

**Action**: Clarified in findings and report that numeric values (500 req/s, 200 users/IP) are hypothetical illustrations, not real-world statistics. No specific timeout/retry/threshold values are presented as recommendations without evidence.

**Outcome**: Resolved — illustrative values now explicitly labeled.

---

### 2. Queue Age vs Queue Depth Evidence

**Original Status**: Flagged as "mungkin bersifat praktik (interpretasi)"
**Resolution**: RESOLVED with vendor documentation

**Evidence Found**:
- AWS SQS CloudWatch: `ApproximateAgeOfOldestMessage` — official AWS metric
- Kafka documentation: Consumer Lag monitoring — official Apache documentation
- Google SRE Book: queue saturation leads to latency increase (§22.10)

**Action**: Added AWS SQS and Kafka documentation references in `02-findings.md` §Claim 7 and `03-sources.md` Source #8, #9.

**Outcome**: Resolved — queue age claim now supported by primary vendor sources.

---

### 3. Little's Law Attribution

**Original Status**: Flagged as "risiko salah atribusi"
**Resolution**: RESOLVED in `02-findings.md` §Claim 3

**Action**:
- Explicitly separated Little's Law ($L = \lambda W$) from deterministic queue buildup ($\Delta Q = (r_{in} - r_{out}) \cdot \Delta t$)
- Cited Little (1961) with proper formula and requirements
- Noted misattribution risk and corrected original plan statement

**Outcome**: Resolved — mathematical distinction now documented.

---

### 4. Vendor Blog Volatility

**Original Status**: "Konten yang berubah (blog vendor) — catat tanggal akses"
**Resolution**: ADDRESSED in `03-sources.md`

**Action**: All sources verified as of 2026-09-26. AWS blog still accessible; Google SRE Book permanent online; RFC documents stable.

**Outcome**: Resolved — access date documented.

---

## New Open Questions (Post-Revision)

### 5. Empirically-Validated Threshold Recommendations

**Status**: OPEN — Research Gap

**Problem**: No authoritative source found for universal "best practice" timeout/retry/threshold values (e.g., "timeout = 3 seconds", "retry = 5 times", "threshold = 50%").

**Why Not Addressed**:
- Service types vary dramatically (database vs cache vs external API)
- Optimal values depend on traffic patterns, network latency, backend capacity
- No single vendor publishes universal recommendations
- SRE books recommend service-specific load testing

**Requirement for Resolution**:
- Find peer-reviewed study or industry survey comparing thresholds across service types
- Or: present as service-specific configuration guidance with load testing recommendation

**Can Be Approved Without Fix**: YES (non-blocking)

---

### 6. Cross-Cloud Rate Limiting Comparison

**Status**: OPEN — Enhancement Opportunity

**Problem**: Research primarily cites AWS and Google documentation; Azure, Cloudflare, Twilio implementations not deeply analyzed.

**Why Not Addressed**:
- Pipeline override limited scope to research revision
- Primary blocking issues (findings document, Little's Law) took priority
- Cross-cloud comparison is enhancement beyond original scope

**Requirement for Resolution**:
- Systematic review of rate limiting implementations across AWS API Gateway, Azure API Management, Cloudflare, Google Cloud Endpoints
- Comparative analysis of algorithm choices (token bucket vs sliding window), defaults, observability metrics

**Can Be Approved Without Fix**: YES (non-blocking)

---

### 7. Real-World Cascade Failure Case Studies

**Status**: OPEN — Enhancement Opportunity

**Problem**: Findings rely on theoretical failure patterns from Google SRE Book; specific production incidents with published postmortems not documented.

**Why Not Addressed**:
- Time constraints (pipeline override, research revision scope)
- Available sources cover theoretical patterns adequately for current audit resolution

**Requirement for Resolution**:
- Identify 2-3 documented production failures where rate limiting/backpressure mitigated or failed to prevent cascade
- Examples: AWS outage postmortems, Stripe incident reports, Cloudflare postmortems

**Can Be Approved Without Fix**: YES (non-blocking)

---

## Summary

| # | Question | Status | Blocking? |
|---|----------|--------|-----------|
| 1 | Numeric recommendations illustrative | RESOLVED | — |
| 2 | Queue age vs depth evidence | RESOLVED | — |
| 3 | Little's Law attribution | RESOLVED | YES (was blocking) |
| 4 | Vendor blog volatility | RESOLVED | — |
| 5 | Empirically-validated thresholds | OPEN | NO |
| 6 | Cross-cloud comparison | OPEN | NO |
| 7 | Real-world case studies | OPEN | NO |

**Blocking Open Questions**: 0 (all blocking issues resolved)
**Non-Blocking Open Questions**: 3 (enhancement opportunities)