# Claim Audit

## Claim 1

Claim: Chaos Engineering is the discipline of experimenting on a system in order to build confidence in the system’s capability to withstand turbulent conditions in production.

Location: `research/03-evidence.md:3`, `research/05-report.md:7`

Evidence Provided: Exact verbatim quote from Principles of Chaos Engineering.

Source: https://principlesofchaos.org/

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Verbatim canonical definition accepted across industry.

---

## Claim 2

Claim: Steady state determination must focus on business/systemic output metrics (throughput, error rate, latency) rather than internal component conditions.

Location: `research/03-evidence.md:19`, `research/05-report.md:12-14`

Evidence Provided: Principles of Chaos Engineering excerpt on systemic output vs internal health.

Source: https://principlesofchaos.org/

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Essential concept distinguishing chaos engineering from low-level unit/integration assertions.

---

## Claim 3

Claim: Chaos hypotheses must posit that steady state continues despite fault injection.

Location: `research/03-evidence.md:35-38`, `research/05-report.md:12-15`

Evidence Provided: Principles of Chaos Engineering principle 2 (Hypothesize about steady state).

Source: https://principlesofchaos.org/

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Standard hypothesis structure in chaos experiments.

---

## Claim 4

Claim: Blast radius must be constrained by starting small (canary/staging) and enforcing automated abort mechanisms.

Location: `research/03-evidence.md:51-54`, `research/05-report.md:18-20`

Evidence Provided: AWS Well-Architected Framework Reliability Pillar guidance.

Source: https://docs.aws.amazon.com/wellarchitected/latest/reliability-pillar/chaos-engineering.html

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Widely recognized safety requirement for running chaos in production environments.

---

## Claim 5

Claim: Downstream latency/failures without timeouts and circuit breakers cause thread starvation and cascading failures.

Location: `research/05-report.md:24-26`

Evidence Provided: Google SRE Book Chapter 22 & Netflix architecture patterns.

Source: Google SRE Book & Netflix TechBlog

Source Actually Supports Claim:
YES

Classification:
FACT

Severity:
LOW

Notes: Standard distributed systems failure mode demonstrated repeatedly in industry literature.
