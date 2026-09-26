# Contradictions Audit

Target Lab: `labs/25-rate-limiting-and-backpressure`

---

## Contradiction Review

### 1. Little's Law vs Queue Fluid Dynamics

- **Statement A**: Plan originally phrased Q3 around `backlog = (arrival - processing) * time` under Little's Law discussion.
- **Statement B**: `02-findings.md` and `05-report.md` rigorously distinguish Little's Law ($L = \lambda W$) from deterministic fluid buildup ($\Delta Q = (r_{in} - r_{out})\Delta t$).
- **Type**: INTERNAL (Resolved)
- **Impact**: Critical clarification preventing mathematical and conceptual errors in downstream lab code and documentation.
- **Assessment**: RESOLVED. No lingering contradiction in active findings or technical report.

### 2. Irrelevant RFC Citations

- **Statement A**: Plan mentioned RFC 8305 and RFC 5321 as potential primary sources.
- **Statement B**: `03-sources.md` explicitly filtered out RFC 8305 (Happy Eyeballs) and RFC 5321 (SMTP), focusing strictly on RFC 6585, RFC 9110, RFC 2697, and RFC 6598.
- **Type**: SOURCE_CONFLICT (Resolved)
- **Impact**: Eliminated noise and irrelevant standards from the research baseline.
- **Assessment**: RESOLVED.

---

## Overall Assessment

No material contradictions found across `01-plan.md`, `02-findings.md`, `03-sources.md`, `04-contradictions.md`, `05-report.md`, and `06-open-questions.md`.
