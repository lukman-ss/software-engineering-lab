# Engineering Notes

## Design Decisions

### Why Go?
- Standard library covers most production needs
- Built-in concurrency (goroutines, channels)
- Race detector (`go test -race`)
- Fast compile, single binary deployment
- Strong ecosystem for observability

### Why PostgreSQL?
- ACID transactions
- Rich data types (JSONB, UUID, arrays)
- Advisory locks for distributed coordination
- Mature, well-understood
- `pg_locks` for deadlock analysis

### Why Redis?
- Sub-millisecond latency for locks/cache
- Atomic operations (INCR, SETNX, Lua scripts)
- Pub/sub for cache invalidation
- Sorted sets for rate limiting windows

### Why not ORM?
- Explicit SQL is readable and debuggable
- Full control over queries, indexes, plans
- No hidden N+1, lazy loading surprises
- sqlc for type-safe queries (future)

### Why not Framework (Gin, Echo, Chi)?
- `net/http` + `context` is sufficient
- Middleware is just functions
- Less magic, easier to trace
- stdlib `ServeMux` improved in Go 1.22

## Patterns Used

### Error Handling
```go
// Wrap with context
return fmt.Errorf("create order: %w", err)

// Sentinel errors for callers
var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

// Check with errors.Is/As
if errors.Is(err, ErrNotFound) { ... }
```

### Context Propagation
```go
func (s *Service) Do(ctx context.Context, req Request) (Response, error) {
    ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
    defer cancel()
    // All I/O uses ctx
}
```

### Transaction Pattern
```go
func (s *Service) Transfer(ctx context.Context, from, to UUID, amount int) error {
    return s.db.WithTx(ctx, func(tx *sql.Tx) error {
        // All operations in tx
        return nil
    })
}
```

### Metrics Naming
```
http_requests_total{method, path, status}
http_request_duration_seconds{method, path}
db_query_duration_seconds{query, table}
business_orders_total{status}
business_payment_total{status, provider}
```

## Trade-offs Documented

### Idempotency (Lab 01)
| Approach | Pros | Cons |
|----------|------|------|
| DB unique constraint | Simple, durable | Extra index, write latency |
| Redis SETNX | Fast, TTL support | Can lose on failover |
| Hybrid (Redis + DB) | Best of both | Complexity |

**Choice**: Hybrid — Redis for fast path, DB as source of truth

### Caching Strategy (Lab 04)
| Approach | Pros | Cons |
|----------|------|------|
| Cache Aside | Simple to implement, standard | Cache stampede on expiration |
| Cache Aside + Single Flight | Deduplicates concurrent DB queries | Still blocks if DB is slow |
| Cache Aside + Jitter | Prevents expiration clusters | Slightly harder to test |
| Redis Distributed Lock | Mutual exclusion across instances | Network overhead, split-brain risk |

**Choice**: Cache Aside with Single Flight for high-traffic endpoints.

### Optimistic vs Pessimistic Locking (Lab 11/12)
| Approach | Pros | Cons |
|----------|------|------|
| Optimistic | No lock wait, high throughput | Retry on conflict, not for high contention |
| Pessimistic | Strong consistency | Lock wait, deadlock risk |

**Choice**: Optimistic for inventory (low contention), Pessimistic for payments (high value)

### Circuit Breaker (Lab 09)
| Config | Effect |
|--------|--------|
| Failure threshold: 5 | Open after 5 failures |
| Timeout: 30s | Half-open after 30s |
| Success threshold: 2 | Close after 2 successes |

**Tuning**: Start conservative, adjust based on error budget

## What Can Still Fail

Even with all patterns:

1. **Network partition** — Idempotency helps, but cannot guarantee exactly-once without distributed consensus
2. **Clock skew** — Rate limiting windows, TTL expiration affected
3. **GC pauses** — Long STW can trigger false timeouts
4. **Schema migration** — Backward compatibility required
5. **Human error** — Config typos, wrong migration order
6. **Dependency bugs** — Library vulnerabilities, driver bugs

## Future Labs

- [ ] 11. Bulkhead pattern
- [ ] 12. Saga pattern (distributed transactions)
- [ ] 13. Event sourcing
- [ ] 14. CQRS
- [ ] 15. Leader election
- [ ] 16. Distributed tracing deep dive
- [ ] 17. Load testing / chaos engineering
- [ ] 18. Database connection pooling tuning
- [ ] 19. Graceful shutdown
- [ ] 20. Feature flags / dark launches

## References

- [Go Proverbs](https://go-proverbs.github.io/)
- [Effective Go](https://go.dev/doc/effective_go)
- [Database Transactions](https://www.postgresql.org/docs/current/tutorial-transactions.html)
- [Redis Distributed Locks](https://redis.io/docs/latest/develop/use/patterns/distributed-locks/)
- [Circuit Breaker Pattern](https://martinfowler.com/bliki/CircuitBreaker.html)
- [Outbox Pattern](https://microservices.io/patterns/data/transactional-outbox.html)