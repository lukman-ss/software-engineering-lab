# Changes Made

## Revision 2 (research-only pass, 2026-09-26)

Audit Issue: HIGH/MEDIUM — Overgeneralization: claim that third-party APIs must always be tested with real calls lacks guardrails for sandbox/staging vs high-volume stress tests

Files Changed:
- research/05-report.md (Finding 7)
- research/03-evidence.md (Evidence 17)
- research/04-contradictions.md (Contradiction 2)
- research/02-sources.md (Source 24 — ISO, Source 23 — JMeter component ref)
- research/06-open-questions.md (weak evidence note)

Action:
- Narrowed claim: real calls appropriate in controlled sandbox/staging environments; mocks with artificial latency for high-volume stress tests (Finding 7 now qualified)
- Added qualification to Evidence 17 matching revised Finding 7
- Updated Contradiction 2 assessment to reflect resolved qualification
- Marked Source 24 (ISO/IEC 25010) as DISCLAIMED in sources.md
- Marked Source 23 (JMeter component reference) as DISCLAIMED in sources.md
- Updated open-questions weak evidence to note ISO disclaimer

Verification:
- All edits reviewed against audit/03-claim-audit Claim 6 (PARTIAL, INTERPRETATION)
- Audit required revision 1 now addressed
- Audit required revision 2 now addressed (Sources 23 and 24 disclaimed)

Status: RESOLVED