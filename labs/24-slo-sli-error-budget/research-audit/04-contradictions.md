# 04 Contradictions

## Contradiction 1: Initial Target Setting Strategy
Statement A: "Don't pick a target based on current performance..." (Google SRE Book Ch.4)  
Location: `labs/24-slo-sli-error-budget/research/04-contradictions.md:5-7`  
Statement B: "...your current performance can be a good place to start if you don't have any other information..." (Google SRE Workbook Ch.2)  
Location: `labs/24-slo-sli-error-budget/research/04-contradictions.md:8-10`  
Type: SOURCE_CONFLICT  
Impact: LOW  
Assessment: Documented evolution of SRE guidance from 2016 to 2018. Research report correctly identifies this as contextual advice requiring iteration rather than a factual flaw.

## Contradiction 2: Time Window Alignment Strategy
Statement A: Google recommends a 4-week rolling window to reflect user experience.  
Location: `labs/24-slo-sli-error-budget/research/04-contradictions.md:16-18`  
Statement B: Evernote chose a fixed calendar month window to align with business reviews.  
Location: `labs/24-slo-sli-error-budget/research/04-contradictions.md:19-21`  
Type: INTERNAL  
Impact: LOW  
Assessment: Valid trade-off between user-centric monitoring and business reporting cycles. Properly analyzed in research.

## Contradiction 3: Downtime Per Month Calculation Assumptions
Statement A: Google SRE Book Appendix A specifies 7.2 hours/month for 99% and 4.32 minutes/month for 99.99% based on a 30-day month convention.  
Location: `labs/24-slo-sli-error-budget/research/04-contradictions.md:52-54`  
Statement B: Lab specification lists 7 hours 18 minutes/month for 99% and 4 minutes 23 seconds/month for 99.99% based on an average 30.44-day month (365.25 / 12).  
Location: `labs/24-slo-sli-error-budget/research/04-contradictions.md:55-57`  
Type: INTERNAL  
Impact: LOW  
Assessment: Discrepancy explained by underlying month length assumption. Both calculations are mathematically sound under their stated conditions.

## Contradiction 4: Burn Rate Threshold Implementations Across Tooling
Statement A: Google SRE Workbook recommends multi-window multi-burn-rate (e.g. 14.4x over 1h/5m, 6x over 6h/30m).  
Location: `labs/24-slo-sli-error-budget/research/04-contradictions.md:41-43`  
Statement B: Datadog documentation uses a single rolling 2-hour window indicator (critical >6, warning 1-6).  
Location: `labs/24-slo-sli-error-budget/research/04-contradictions.md:44-46`  
Type: SOURCE_CONFLICT  
Impact: MEDIUM  
Assessment: Highlights vendor-specific simplifications vs canonical Google recommendations. Burn rate math core remains identical.
