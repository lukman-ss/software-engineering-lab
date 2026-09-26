# Revision Plan

Target Lab: labs/13-backward-compatibility
Previous Audit Status: APPROVED_WITH_WARNINGS

## Blocking Issues

None.

## Non-Blocking Issues

1. **Gap 1 (OVERGENERALIZATION, MEDIUM)**: The research claims that a legacy interface should only be removed after metric usage shows zero hits for 30 consecutive days. This is presented as an absolute rule but lacks any citation from the reviewed sources (Stripe, Fowler, Prisma). It is an operational heuristic rather than a universal standard.

## Files To Modify

- research/11-final-research.md (Section Q11)
- research/08-failure-modes.md (Section "4. Failure Mode: Permanent Expand State")

## Verification Plan

- research/11-final-research.md: confirm 30-day period labeled as heuristic
- research/08-failure-modes.md: confirm same labeling
- no source changes required (audit only reviewed Sources 1-3, all PASS)

## Changes Made

- research/11-final-research.md: Added `> [!NOTE] Industry heuristic — NOT VERIFIED` callout; removed "misal: 30 hari berturut-turut" parenthetical as primary statement
- research/08-failure-modes.md: Added "operational heuristic, NOT VERIFIED as universal standard" to 30-day alert description

## Verification Results

| Check | Result |
|---|---|
| 30-day heuristic labeled in 11-final-research.md | PASS |
| 30-day heuristic labeled in 08-failure-modes.md | PASS |
| No source inventory changes required | PASS |
| Only audit-flagged files touched | PASS |