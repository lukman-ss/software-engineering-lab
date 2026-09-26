# Contradictions Resolution Log

## Original Contradictions (from research-audit/04-contradictions.md)

### Contradiction 1 (HIGH) — Little's Law Misattribution

**Statement A**: "Little's Law: backlog = (arrival − processing) × waktu" (Research Question Q3 in `01-plan.md`)

**Statement B**: "Little's Law vs formula backlog kumulatif: risiko salah atribusi — bedakan dengan jelas." (Risks / Unknowns in `01-plan.md`)

**Resolution**: RESOLVED in `02-findings.md` §Claim 3

**Action Taken**:
- Explicitly separated Little's Law ($L = \lambda W$) from deterministic queue buildup ($\Delta Q = (r_{in} - r_{out}) \cdot \Delta t$)
- Provided proper mathematical formulation with variable definitions
- Cited Little (1961) original paper with DOI
- Clarified that the fluid model is an approximation, not Little's Law
- Statement A was **incorrect** — the formula does not represent Little's Law
- Statement B was **correct** — identified the misattribution risk

**Outcome**: Statement A corrected; the original claim has been replaced with precise distinction.

---

### Contradiction 2 (LOW) — Irrelevant RFC Citations

**Statement A**: Plan lists RFC 8305 and RFC 5321 under Expected Primary Sources

**Statement B**: "RFC 8305? (mungkin tidak relevan) — ganti: RFC 5321 (SMTP) tidak relevan; fokus: RFC 6585, RFC 9110."

**Resolution**: RESOLVED in `03-sources.md`

**Action Taken**:
- Removed RFC 8305 (Happy Eyeballs — not relevant to rate limiting)
- Removed RFC 5321 (SMTP — not relevant to HTTP rate limiting)
- Focused on RFC 6585, RFC 9110 as correct primary sources
- This was a self-correction within the original plan that has been formalized

**Outcome**: Irrelevant sources removed; source list is now clean.

---

## Additional Contradictions Found During Research

### None

No additional contradictions were found. The original research plan's risk flags were accurate. The primary issue was the absence of a findings document to execute the self-identified corrections.