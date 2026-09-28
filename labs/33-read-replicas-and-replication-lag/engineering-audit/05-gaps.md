# Engineering Gap Analysis

Target Lab: `labs/33-read-replicas-and-replication-lag`

## Summary of Gaps

| Gap ID | Gap Type | Severity | Description | Impact |
|---|---|---|---|---|
| GAP-01 | UNHANDLED_ERROR | LOW | In `Cluster.Write` async mode, WAL entries are dispatched using a non-blocking `select` with empty `default:`. If channel buffer (1024) fills up, entries would drop silently. | Minimal for in-memory lab tests; production systems require backpressure management. |
| GAP-02 | MISSING_EDGE_CASE | LOW | In `Node.WaitForLSN`, timeout on `ctx` does not cancel the background goroutine waiting on `sync.Cond.Wait()`; it waits until the next broadcast. | Minimal memory footprint in simulation; all tests and demos broadcast on write. |

## Unresolved High/Critical Gaps
- None.
