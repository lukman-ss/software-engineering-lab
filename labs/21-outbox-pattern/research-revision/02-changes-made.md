# Changes Made

## Revision 1

Audit Issue: 
1. Source 2 title mislabeled in research/02-sources.md:11
2. Finding 5 overgeneralizes payload design in research/05-report.md:60-72
3. Finding 6 presents operational metrics as empirical guidance in research/05-report.md:73-84

Files Changed:
- research/02-sources.md
- research/05-report.md
- research/03-evidence.md

Actions:
- Fixed Source 2 title from "Transactional outbox" to "Pattern: Transaction log tailing"
- Rewrote Finding 5 to present both thin-event (minimal identifiers) and fat-event (Event-Carried State Transfer) approaches with their trade-offs
- Rewrote Finding 6 to explicitly label 2s/47min thresholds as illustrative SLO examples, not universal constants
- Updated Evidence 10 similarly to clarify metric values are illustrative

Verification:
- Cross-checked with microservice.io and Debezium sources
- Confirmed Debezium blog example documents fat events with full order state
- Verified all changes maintain consistency with cited sources

Status:
RESOLVED