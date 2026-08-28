# Integration Tests for Idempotency Lab

## Running Tests

```bash
# Run all tests
go test ./... -v

# Run with race detector
go test -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Run HTTP handler tests only
go test -v ./tests/...
```

## Test Coverage

| Test | Description |
|------|-------------|
| `TestIdempotentRetry` | Same idempotency key returns cached response |
| `TestConcurrentIdempotentPayment` | 10 concurrent requests, only 1 gateway call |
| `TestIdempotencyKeyConflict` | Different payload with same key returns error |
| `TestDifferentIdempotencyKeys` | Different keys create different payments |
| `TestHTTPHandlerWithIdempotency` | HTTP handler idempotency |
| `TestHTTPConflictOnDifferentPayload` | HTTP 409 on payload mismatch |
| `TestConcurrentPaymentsSameKey` | 20 concurrent requests, idempotent |
| `TestDifferentKeysDifferentPayments` | Multiple distinct payments |

## Expected Results

- **UnSafe tests**: Demonstrate duplicate charges (BUG)
- **Safe tests**: Verify no duplicates, idempotency works
- **Benchmarks**: Measure throughput with idempotency overhead

## Production Checklist

- [ ] Monitor idempotency hit rate (high = client retry issues)
- [ ] Alert on conflict errors (client bug)
- [ ] Alert on gateway call count (should match unique keys)
- [ ] Trace idempotency key propagation
- [ ] Log idempotency key for correlation