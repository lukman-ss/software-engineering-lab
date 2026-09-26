# Gap Analysis

## MISSING_TEST
### G001
Location: tests/outbox_test.go (missing test scenario)
Description: No test verifies relay retry behavior after broker publish failure. The relay logic handles this (message stays pending on broker failure, retried on next poll) but it is not exercised by any test.
Reproduction: Would require a broker that fails on first publish attempt, succeeds on second, and asserts that the message is eventually published and outbox marked processed.
Risk: Relay recovery logic could regress without test coverage. Severity: MEDIUM.

### G002
Location: tests/outbox_test.go:TestTransactionalOutbox_ConcurrentWrites (lines 141-170)
Description: The concurrent writes test reuses the same order ID across all goroutines (overwrite semantics) and contains zero assertions after the write burst. It does not verify that concurrent writes of distinct orders preserve all writes, that no orders are lost, or that broker receives correct number of messages.
Reproduction: Change test to use distinct order IDs (e.g., fmt.Sprintf("o-%d-%d", wID, j)) and assert: total orders in DB == workers*ordersPerWorker, total broker messages == workers*ordersPerWorker.
Risk: Concurrency correctness (counts, data integrity) is not verified by tests. Race detector only checks for data races, not logical correctness. Severity: MEDIUM.

## UNHANDLED_ERROR
### G003
Location: internal/outbox/relay.go:Relay.Stop (lines 39-41)
Description: Relay.Stop() closes stopChan without guarding against multiple calls. If Stop() is called more than once, close(stopChan) will panic ("close of closed channel"). While the tests only call Stop() once via defer, this is a latent defect.
Reproduction: Call relay.Stop(); relay.Stop() in test or code.
Risk: Panic under double-stop. Severity: LOW (not exercised in normal test/demo flow).

## MISSING_EDGE_CASE
### G004
Location: tests/outbox_test.go (missing edge case)
Description: No test for extremely short poll interval (e.g., 1ms) to stress test relay goroutine scheduling and potential starvation or excessive CPU usage.
Reproduction: Set pollInterval to 1*time.Millisecond and run concurrent writes test with assertions.
Risk: Relay could miss ticks or behave unexpectedly under high frequency. Severity: LOW.