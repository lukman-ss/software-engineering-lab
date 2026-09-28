# 01 - Research Audit Plan

Target Lab: `labs/24-slo-sli-error-budget`
Audit Scope: Research Artifacts Only (Pipeline Override: No code/implementation audit in this stage)
Audit Date: 2026-09-28

---

## 1. Files Reviewed

Inside `labs/24-slo-sli-error-budget/research/`:
- `01-plan.md` — Research scope, questions, tiered search strategy.
- `02-sources.md` — Inventory of 13 sources (10 active, 3 recorded as 404/unavailable).
- `03-evidence.md` — 15 extracted evidence items mapped to sources.
- `04-contradictions.md` — 6 contradiction/discrepancy assessments.
- `05-report.md` — Synthesis report, findings, areas of agreement/disagreement, limitations.
- `06-open-questions.md` — Analysis of open questions, weak evidence, and research gaps.

---

## 2. Claims To Verify

1. **Canonical Definitions**: SLI as quantitative ratio metric; SLO as target/range; SLA as contractual commitment with business consequences; Error Budget as `100% - SLO`.
2. **SLI Representation & Quality**: SLI should be represented as ratio of good events to total events; SLIs must be user-centric (symptom-based), not internal infrastructure metrics (e.g. CPU/RAM).
3. **Distribution vs Mean**: Percentiles (P95/P99) must be used instead of arithmetic mean for latency SLIs to avoid obscuring tail latency.
4. **Availability & Nines**: "Nines" calculation for allowed downtime per day/month/year; diminishing returns and escalating cost per additional nine; setting 100% SLO is an anti-pattern.
5. **Burn Rate & Alerting**: Burn rate formula `burn rate = (consumed budget % / consumed time %)`; Google SRE Workbook multi-window multi-burn-rate alerting parameters (Table 5-8: 14.4x over 1h/5m for 2% budget; 6x over 6h/30m for 5% budget; 1x over 3d/6h for 10% budget).
6. **Error Budget Calculation & Policy**: Datadog error budget remaining formula `100 * (current status - target) / (100 - target)`; Error budget policy used as formal release gating / prioritization tool.
7. **Standards & Case Studies**: OpenSLO as vendor-neutral YAML standard; real-world adoption cases (Evernote, The Home Depot VALET framework).

---

## 3. Code To Execute

None.
*Note: Under PIPELINE OVERRIDE, this audit evaluates research artifacts only. Implementation code and tests are out of scope for this stage.*

---

## 4. Primary Risks

1. **Dead or Hallucinated Links**: Citations pointing to 404s, redirected URLs, or fabricated sources.
2. **Overgeneralization**: Presenting Google-specific SRE conventions (such as 4-week rolling windows, request-based availability, or specific alerting burn rates) as universal engineering truths.
3. **Arithmetic Assumptions**: Discrepancies in downtime calculations between a 30-day month convention and a 365.25/12 (~30.44-day) average month convention.
4. **Vendor Divergence**: Conflating vendor-specific dashboard metrics (e.g., Datadog 2-hour burn rate indicator) with Google SRE multi-window alerting patterns.

---

## 5. Audit Strategy

- Verify each cited source in `02-sources.md` for publisher authenticity, domain validity, and content alignment.
- Verify every claim in `03-evidence.md` against authoritative SRE literature.
- Review contradictions in `04-contradictions.md` to ensure legitimate trade-offs are not swept under the rug.
- Assess `05-report.md` for objective representation, unwarranted conclusions, or missing nuances.
- Check `06-open-questions.md` for thorough gap identification.
