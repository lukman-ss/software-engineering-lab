# Revision Plan

Target Lab: labs/21-outbox-pattern
Previous Audit Status: APPROVED_WITH_WARNINGS

## Blocking Issues
None.

## Non-Blocking Issues
1. Source 2 title is mislabeled as "Transactional outbox" instead of "Pattern: Transaction log tailing" in research/02-sources.md:11.
2. Finding 5 overgeneralizes payload design by prescribing minimal identifiers only, omitting Event-Carried State Transfer (fat payload) trade-offs in research/05-report.md:60-72.
3. Operational thresholds (2 seconds vs 47 minutes) cited as empirical guidance rather than illustrative SLA examples in research/05-report.md:73-84.

## Files To Modify
- research/02-sources.md
- research/05-report.md
- research/03-evidence.md (for Gap 2, which is about the weak source for operational metrics)

## Verification Plan
- Source verification: Check that the corrected title matches the source.
- For payload design: Ensure we present both thin and fat event trade-offs.
- For operational metrics: Clarify that the numbers are illustrative examples.
- Check that the changes are consistent across the research files.