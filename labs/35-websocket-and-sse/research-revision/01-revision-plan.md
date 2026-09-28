# Revision Plan

Target Lab: labs/35-websocket-and-sse

Previous Audit Status: APPROVED_WITH_WARNINGS

## Blocking Issues

1. None.

## Non-Blocking Issues

1. RFC 7540 URL typo fixed.
2. RFC 7540 obsolescence not cited.
3. Weak internal source for scaling claims replaced with Linux kernel docs and Redis docs.
4. Minor editorial typo in URL (multiphase → multipage).
5. Open question on HTTP/3 added clarification.

## Files To Modify

- research/05-report.md
- research/02-sources.md
- research/03-evidence.md
- research/04-contradictions.md
- research/06-open-questions.md

## Verification Plan

- Source verification: check URLs reachable, confirm RFC 9113 obsolescence note.
- Confirm evidence now cites Linux kernel docs, Redis docs.
- Ensure no unsupported claims remain.
- No tests/code to run.
