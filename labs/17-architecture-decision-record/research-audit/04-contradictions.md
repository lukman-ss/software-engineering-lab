# Contradictions Audit: Architecture Decision Record Research

## Contradiction 1

Statement A:
"For an early-stage SaaS ERP with a small team (5 engineers) and a 3-month deadline, a Modular Monolith delivers necessary domain boundaries without the operational, networking, and distributed transaction complexity of microservices." (Supported by Cognitect Blog, AWS Prescriptive Guidance)

Location:
`05-report.md` (Finding 4)

Statement B:
The primary sources (Nygard 2011, AWS Prescriptive Guidance, adr.github.io) are strictly process and governance definitions for Architectural Decision Records, not comparative architectural analyses of monoliths versus microservices.

Location:
`02-sources.md`

Type:
INTERNAL

Impact:
The research appropriately frames the ADR methodology within an applied case study (Monolith vs. Microservices). However, it attributes the architectural conclusions of that case study directly to the process guidelines. This makes it appear that AWS Prescriptive Guidance or Nygard provided the benchmark data for the 5-engineer ERP SaaS scenario.

Assessment:
This is a medium-impact internal contradiction regarding source attribution. The architectural trade-off logic is sound, but its source attribution is flawed.

---

## Contradiction 2

Statement A:
Nygard limits status lifecycle to Proposed, Accepted, Deprecated, and Superseded.

Location:
`04-contradictions.md`

Statement B:
AWS Prescriptive Guidance formally incorporates a "Rejected" state into the ADR lifecycle.

Location:
`04-contradictions.md`

Type:
SOURCE_CONFLICT

Impact:
Minimal.

Assessment:
The research correctly identifies and records this minor structural discrepancy between foundational (2011) and modern cloud enterprise (2023) ADR templates.
