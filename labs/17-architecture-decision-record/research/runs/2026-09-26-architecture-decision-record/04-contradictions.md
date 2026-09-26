# Contradictions

No material contradictions on core principles. Three minor discrepancies and one practice-level tension.

## Minor Discrepancy 1 — Status "Rejected"

SOURCE A (Nygard 2011):
Explicit states listed are `Proposed`, `Accepted`, `Deprecated`, `Superseded`. No formal `Rejected` — implied via unaccepted pull requests or informal notes.

SOURCE B (AWS Prescriptive Guidance, inspected 2026-09-26):
Adds explicit `Rejected` state: "The team can also decide to reject the ADR. In this case, the ADR owner adds a reason for the rejection to prevent future discussions on the same topic. The owner changes the ADR state to Rejected."

ASSESSMENT:
Not contradictory in substance — both preserve rejected proposals to avoid re-litigation. AWS formalizes what Nygard leaves informal. For the lab, support the union set: `Proposed | Accepted | Deprecated | Superseded | Rejected` (validator requirement in topic spec lists four, but `Rejected` is a compatible extension).

## Minor Discrepancy 2 — Template Granularity

SOURCE A (Nygard):
Anatomy limited to Title, Context, Decision, Status, Consequences (1-2 pages).

SOURCE B (MADR / AWS):
Elaborate anatomy: Context and Problem Statement, Decision Drivers, Considered Options, Decision Outcome, Consequences (Good/Bad), Confirmation, Pros and Cons of the Options, More Information (MADR full template). AWS adds "At a minimum ... context, decision, consequences" but encourages richer templates per team.

ASSESSMENT:
Both agree on rationale + context + consequences core. Difference is granularity, not principle. Modern templates decompose Nygard's `Consequences` and `Context` into drivers/options/pros-cons to force explicit alternatives and trade-offs — exactly what the SaaS ERP scenario requires (Laravel Modular Monolith vs Go Microservices).

## Practice-Level Tension — Immutability: Strict vs Living Document

SOURCE A (Nygard + AWS + Azure Well-Architected):
Strict append-only immutability: "When the team accepts an ADR, it becomes immutable." / "The ADR serves as an append-only log. Don't go back and edit accepted records."

SOURCE B (Joel Parker Henderson community repository, inspected 2026-09-26, Tier 2):
"In theory, immutability is ideal. In practice, mutability has worked better for our teams. We insert the new info the existing ADR, with a date stamp, and a note that the info arrived after the decision. This kind of approach leads to a 'living document' that we all can update."

ASSESSMENT:
Community pragmatism vs canonical prescription. Tier 1 sources are mutually consistent on immutability; Tier 2 reports lived-experience deviation. For the lab, follow Tier 1 (immutable + new superseding ADR) because it preserves auditable decision history, which the topic spec explicitly requires ("ADR lama tidak boleh dihapus", "Superseded by ADR-XXX"). The community pattern is worth noting as an observed variant, not as normative guidance.

## Minor Discrepancy 3 — Monolith-First: Consensus vs Dissent

SOURCE A (Martin Fowler, Monolith First, 2015-06-03):
"You shouldn't start a new project with microservices, even if you're sure your application will be big enough to make it worthwhile" due to MicroservicePremium and premature boundary fixation.

SOURCE B (Same Fowler article, counter-argument section):
"While the bulk of my contacts lean toward the monolith-first approach, it is by no means unanimous... The counter argument says that starting with microservices allows you to get used to the rhythm of developing in a microservice environment... especially viable for system replacements where you have... stable-enough boundaries early."

ASSESSMENT:
Not a factual contradiction — both views acknowledge trade-offs. The article's dominant recommendation (monolith-first for greenfield with uncertain boundaries, small team, no platform capability) directly aligns with the lab's SaaS ERP constraints (5 engineers, 3-month deadline, still-evolving requirements, no SRE team, Laravel familiarity). The dissent carves out a narrow exception (experienced microservices team, stable boundaries, system replacement) that does not apply to the lab's primary scenario.

## Minor Discrepancy 4 — Confidence Logging

SOURCE A (Microsoft Azure Well-Architected, 2026-04-10):
"Record the confidence level of the decision. Sometimes an architecturally significant decision is made with relatively low confidence. Documenting that low confidence status could prove useful for future reconsideration decisions."

SOURCE B (Nygard / AWS):
No explicit "confidence" field; confidence is implicit in rationale and consequences.

ASSESSMENT:
Additive, not conflicting. Lab can treat confidence as optional metadata (mirrors MADR's optional front-matter) without violating Nygard/AWS минима.
