# 04 Contradictions

Target Lab: `labs/24-slo-sli-error-budget`  
Audit Date: 2026-09-26

## Contradiction 1: Monthly Downtime Calculation Baseline (30-day month vs 30.44-day average month)
- **Statement A:** Lab topic / preliminary calculations mention ~7 hours 18 minutes allowed downtime per month for 99% availability (`04-contradictions.md` §1).
- **Statement B:** Google SRE Book Appendix A Table 1-1 gives 7.2 hours/month = 7 hours 12 minutes/month for 99% availability (`03-evidence.md` §Evidence 5).
- **Type:** SOURCE_CONFLICT (Arithmetic baseline variation).
- **Impact:** Negligible (difference is 6 minutes over a month, ~1.4%). Statement A assumes an average Gregorian month length of 365.25 / 12 = 30.44 days (43,830 minutes), whereas Google's table explicitly adopts a nominal 30-day month (43,200 minutes).
- **Assessment:** PASS (Accurately documented and resolved by the research team).

---

## Contradiction 2: 99.99% Downtime Duration (4m 23s vs 4m 19s)
- **Statement A:** Lab preliminary estimate: ~4 minutes 23 seconds/month for 99.99% availability.
- **Statement B:** Google Appendix A: 4.32 minutes = 4 minutes 19.2 seconds/month.
- **Type:** SOURCE_CONFLICT (Arithmetic baseline variation).
- **Impact:** Negligible (difference is 3.8 seconds over a month). Caused by the same 30.44 days vs 30.0 days baseline.
- **Assessment:** PASS (Accurately documented and resolved by the research team).

---

## Contradiction 3: Evaluation Window (30 Calendar Days vs 28 Days / 4 Full Weeks)
- **Statement A:** Common business SLOs often specify a 30-day calendar rolling window.
- **Statement B:** Google SRE Workbook Ch.2 recommends a 4-week rolling window (28 days) to ensure identical weekend/weekday proportions.
- **Type:** INTERNAL (Design recommendation trade-off).
- **Impact:** Minor. 28-day rolling window prevents cyclical day-of-week distortion; 30-day rolling window aligns with calendar billing/monthly reviews.
- **Assessment:** PASS (Both models are valid; the research explicitly highlights the trade-off).

---

## Summary
No material technical contradictions found. All key principles (SLI definitions, SLO thresholds, Error Budget equations, Golden Signals, and multi-burn-rate logic) are consistent across all consulted chapters and research deliverables.
