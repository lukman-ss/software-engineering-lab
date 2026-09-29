# Revision Plan

Target Lab: `labs/38-mutation-testing`

Previous Audit Status: `APPROVED_WITH_WARNINGS`

## Audit Summary

- Critical Issues: 0
- High Issues: 1 (Martin Fowler bliki draft URL returned HTTP 404)
- Medium Issues: 1 (Unsourced parenthetical numeric range `(80%, 85%, 90%)` in Finding 11)
- Low Issues: 2 (Secondary attribution for academic papers, Tool pattern overgeneralization regarding generative tools)

## Blocking Issues

None.

## Non-Blocking Issues

1. **Source 3 URL Unreachable (HIGH)**: `https://martinfowler.com/bliki/MutationTesting.html` returned HTTP 404. Source 3 was already marked REMOVED in `research/02-sources.md` during research phase, but remaining citations in report need to reflect pure backing by primary reachable sources.
2. **Finding 11 Parenthetical Numbers (MEDIUM)**: Parenthetical numeric target values `(80%, 85%, 90%)` in Finding 11 and open questions were unsourced.
3. **Generative Tool Workflow Overgeneralization (LOW)**: Claim 8 stated "All major mutation testing tools follow..." without distinguishing evaluation-focused tools from generative LLM tools (e.g. ACH).
4. **Academic Secondary Attribution (LOW)**: Primary 1978 and 2009 papers cited via Wikipedia references; properly disclosed in sources.

## Files To Modify

- `labs/38-mutation-testing/research/05-report.md`
- `labs/38-mutation-testing/research/03-evidence.md`
- `labs/38-mutation-testing/research/06-open-questions.md`

## Verification Plan

- Verify removal of unsourced numbers from report and open questions.
- Verify clarification of tool execution pattern for generative vs evaluation tools.
- Ensure all claims retain valid, reachable primary/secondary sources.
