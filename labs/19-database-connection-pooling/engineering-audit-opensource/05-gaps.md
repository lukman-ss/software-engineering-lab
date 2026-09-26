# Gaps

## Gap 1
WARNING - timing-sensitive assertions in TestDirectConnectionOverhead & TestTotalCreatedPoolReuse.
Notes: Pass currently; nondeterministic under load.

## Gap 2
MISSING_EDGE_CASE - query execution failure path:
mockStmt.Exec never errors; service never tests failed ExecContext recovery.
Low impact (no data loss path).

## Gap 3
UNHANDLED_ERROR - MockConnector.Connect ignores context during connectDelay;
cannot cancel in-progress connect. Not used in failure-path tests.
