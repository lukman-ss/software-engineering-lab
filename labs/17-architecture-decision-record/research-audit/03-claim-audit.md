# Claim Audit

## Claim 1

Claim:
ADRs must contain Title, Context, Decision, Status, and Consequences sections to be effective.

Location:
`research/03-evidence.md` (Evidence 1) and `research/05-report.md` (Finding 1)

Evidence Provided:
Direct quote from Michael Nygard (2011) defining the five core sections. Corroborated by MADR and AWS Prescriptive Guidance.

Source:
Source 1 (Nygard 2011), Source 3 (MADR), Source 5 (AWS)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
While "must" could sound overly prescriptive, multiple independent sources converge on these exact section semantics as the canonical ADR structure.

---

## Claim 2

Claim:
ADRs prevent teams from repeatedly re-discussing the same architecture decisions by preserving decision rationale.

Location:
`research/03-evidence.md` (Evidence 2) and `research/05-report.md` (Finding 2)

Evidence Provided:
Quotation from AWS Prescriptive Guidance citing "the same topic being discussed multiple times" as an anti-pattern caused by missing justification, and Nygard on future team members understanding past decisions.

Source:
Source 5 (AWS), Source 1 (Nygard)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately reflects the core rationale and observed benefit highlighted by practitioners and official guidance.

---

## Claim 3

Claim:
ADRs should document only architecturally significant decisions, not all technical choices.

Location:
`research/03-evidence.md` (Evidence 3) and `research/05-report.md` (Finding 6)

Evidence Provided:
Direct quotation from Nygard on decisions that affect structure, non-functional characteristics, dependencies, interfaces, or construction techniques, corroborated by Zimmermann's guidance to avoid ADR logs >100 entries.

Source:
Source 1 (Nygard 2011), Source 4 (Zimmermann 2020)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Clear boundary preventing ADR bloat; well-supported by primary and secondary sources.

---

## Claim 4

Claim:
ADRs have lifecycle states (Proposed, Accepted, Deprecated, Superseded) and old ADRs are preserved even when superseded.

Location:
`research/03-evidence.md` (Evidence 4) and `research/05-report.md` (Finding 5)

Evidence Provided:
Nygard quotation on marking old ADRs as superseded and keeping them around; corroborated by MADR specifications.

Source:
Source 1 (Nygard 2011), Source 3 (MADR)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately preserves the append-only/immutable history nature of ADR logs.

---

## Claim 5

Claim:
Microservices architecture introduces significant operational complexity premium compared to monoliths.

Location:
`research/03-evidence.md` (Evidence 5) and `research/05-report.md` (Finding 4)

Evidence Provided:
Quotation from Martin Fowler detailing automated deployment, monitoring, dealing with failure, and eventual consistency costs.

Source:
Source 7 (Fowler 2015 "Microservice Premium")

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Well-known architectural principle clearly documented and supported by cited source.

---

## Claim 6

Claim:
Most successful microservice implementations started as monoliths that later decomposed, not greenfield microservices.

Location:
`research/03-evidence.md` (Evidence 6) and `research/05-report.md` (Finding 3)

Evidence Provided:
Quotation from Martin Fowler stating that almost all successful stories started with a monolith.

Source:
Source 6 (Fowler 2015 "Monolith First")

Source Actually Supports Claim:
PARTIAL

Classification:
INTERPRETATION

Severity:
MEDIUM

Notes:
Fowler explicitly states this is based on anecdotal observations from his network and that dissenting opinions exist. The research report accurately notes this nuance ("Confidence: MEDIUM", acknowledging anecdotal evidence).

---

## Claim 7

Claim:
Early-stage projects with uncertain requirements should prioritize development speed over architectural sophistication.

Location:
`research/03-evidence.md` (Evidence 7) and `research/05-report.md` (Finding 3)

Evidence Provided:
Fowler quote regarding prioritizing speed/cycle time during the initial phase and avoiding microservice overhead.

Source:
Source 6 (Fowler 2015)

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
LOW

Notes:
Consistent with agile principles and well-grounded in cited source text.

---

## Claim 8

Claim:
Multiple ADR templates exist with varying sections, but core elements remain consistent across templates.

Location:
`research/03-evidence.md` (Evidence 8)

Evidence Provided:
WICSA 2015 paper comparing 7 ADR templates and Zimmermann (2020) detailing Y-statements vs Nygard vs arc42.

Source:
Source 4 (Zimmermann 2020), Source 10 (WICSA 2015)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Accurately captures structural diversity while identifying underlying consensus.

---

## Claim 9

Claim:
ADRs serve as "architecture git history" - recording why architecture changed, complementing source control that records how code changed.

Location:
`research/03-evidence.md` (Evidence 9) and `research/05-report.md` (Executive Summary)

Evidence Provided:
Nygard quotation on motivation behind decisions being visible to everyone so nobody wonders "What were they thinking?".

Source:
Source 1 (Nygard 2011), Source 5 (AWS)

Source Actually Supports Claim:
YES

Classification:
INTERPRETATION

Severity:
LOW

Notes:
The phrase "architecture git history" is an apt metaphor synthesising the concepts from Nygard and AWS. The research notes explicitly label this as an interpretation/analogy.

---

## Claim 10

Claim:
Common ADR anti-patterns include: no context, fake alternatives, documentation dumps, missing consequences, no review conditions.

Location:
`research/03-evidence.md` (Evidence 10)

Evidence Provided:
Citations from AWS Prescriptive Guidance and Zimmermann's analysis of bad justifications and anti-patterns.

Source:
Source 5 (AWS), Source 4 (Zimmermann)

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes:
Directly derived from the referenced sources.
