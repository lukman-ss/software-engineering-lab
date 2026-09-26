# 04 - Contradictions

## Contradiction 1
Statement A:
"The 12-parameter threshold is a lab-specific heuristic, not an industry standard"

Location:
`research/05-report.md` (Finding 11, Limitations)

Statement B:
The target lab instructions/heuristics claim that if a constructor has 12 parameters, "Biasanya ada masalah desain".

Location:
Topic spec references / internal heuristics.

Type:
SOURCE_CONFLICT

Impact:
Minor. The research acknowledges that the target lab uses a specific heuristic (12 parameters) while standard industry literature (Fowler) refers to it qualitatively ("a lot of parameters"). The research does an excellent job of reconciling this by flagging it as a lab-specific threshold rather than an academic truth.

Assessment:
No further action needed. Appropriately handled by the research agent.

---

## Contradiction 2
Statement A:
DI and Service Locator "are very amenable to stubbing" if well-designed.

Location:
Fowler 2004, referenced in `research/04-contradictions.md` (Contradiction 1).

Statement B:
Microsoft .NET documentation presents DI primarily as a testing solution, implying hard-coded or non-DI approaches block stubbing.

Location:
Microsoft Learn Docs, referenced in `research/04-contradictions.md` (Contradiction 1).

Type:
SOURCE_CONFLICT

Impact:
Minor. Microsoft frames DI as the definitive testing enabler, while Fowler takes a broader architectural view that Service Locator also supports stubs. The research agent successfully categorized this as a difference in emphasis rather than a factual blocker.

Assessment:
No further action needed. Appropriately contextualized by the research agent.

---

## Conclusion
No material or internal contradictions found in the finalized research documentation. Previous issues (e.g., misattribution of PSR-11's RFC 2119 keyword) appear to have been fully remediated in prior revisions, and the current research outputs are internally consistent and accurately align with authoritative external literature.
