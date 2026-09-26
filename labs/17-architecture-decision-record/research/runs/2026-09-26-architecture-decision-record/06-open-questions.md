# Open Questions

## Weak Evidence / Needs Deeper Research

- **Quantified MicroservicePremium for Laravel vs Go SaaS ERP stack.**
  Claim in the report that microservices impose "network failure modes, distributed transactions, observability burden, deployment orchestration" is qualitative (via Fowler). A quantitative baseline (latency overhead per cross-service call, ops-engineer-hour cost per service) would strengthen Finding 5. Priority: LOW. Out of scope per topic spec (excludes Laravel-vs-Go benchmarks), but worth a literature citation if a production case study surfaces.

- **Optimal "Review Trigger" threshold calibration.**
  The report states "reporting deployment cadence diverging from core ERP" is a valid trigger, but the precise measurable threshold (e.g., 10x daily reporting deploys vs weekly core deploys) is not defined by any inspected source. Need production case studies (e.g., Shopify monolith-to-services retrospectives, Netflix decanting). Priority: MEDIUM.

- **Confidence-level field standardization across templates.**
  Azure Well-Architected recommends recording decision confidence levels; Nygard/AWS/MADR do not specify a scale. No inspected source defines a standard 1–5 or 0–100 scale. Priority: LOW. Acceptable as optional metadata (mirrors MADR YAML front-matter).

## Unanswered Questions

- **Cross-repository ADR discovery and ownership.**
  How do organizations with split services across multiple repositories maintain a coherent architecture decision log and ensure cross-repo decisions are visible during code review? Inspected sources (AWS, Azure, MADR, Nygard) address single-repo patterns but do not prescribe multi-repo discovery. Priority: MEDIUM.

- **ADR-to-PR linking automation maturity.**
  adr.github.io references tools (ADR Guard, Decision Guardian) that surface relevant ADRs on PRs, but no inspected source provides empirical effectiveness data (false-positive rate, adoption friction). Priority: LOW for validation tooling scope.

## Possible Next Research Directions

- Inspect a production case study of an Indonesian SaaS ERP migration (monolith → services) to validate review-trigger thresholds against the lab's scenario constraints (small team, limited infra budget).
- Survey open-source Laravel modular monolith scaffolding (e.g., `mlnt Laravel modules`, `spatie` packages) to confirm the lab's "modular monolith" boundary technique is idiomatic, ensuring the lab does not prescribe non-idiomatic structure.
- Examine the `joelparkerhenderson/architecture-decision-record` repository's example ADRs (not just README) to corroborate real-world usage of the Review Triggers and Superseded-by patterns cited in the report.
- Audit a real engineering team's ADR git history to validate the immutability + supersede-by pattern in practice (whether old ADRs survive unchanged when superseded).
