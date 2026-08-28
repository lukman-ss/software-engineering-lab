# Architecture

## System Overview

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Client    │────▶│    API      │────▶│  PostgreSQL │
└─────────────┘     └─────────────┘     └─────────────┘
                           │
                           ▼
                    ┌─────────────┐
                    │    Redis    │
                    └─────────────┘
                           │
                           ▼
                    ┌─────────────┐     ┌─────────────────┐
                    │   Worker    │────▶│  Outbox Relay   │
                    └─────────────┘     └─────────────────┘
```

## Components

### API Server (`cmd/api`)
- HTTP/1.1 + gRPC (future)
- Request validation, rate limiting
- OpenTelemetry instrumentation
- Health/readiness endpoints

### Workers (`cmd/worker`)
- Outbox event processor
- Scheduled jobs
- Dead letter handling

### Mock Provider (`cmd/mock-provider`)
- Simulates external payment gateway
- Configurable failure modes
- Latency injection

### Domains (`internal/`)
Each domain owns its data and business logic:

| Domain | Responsibility |
|--------|----------------|
| order | Order lifecycle, items, state machine |
| payment | Payment processing, idempotency |
| inventory | Stock reservation, allocation |
| wallet | Balance, ledger, transfers |
| notification | Email, SMS, push notifications |

### Packages (`pkg/`)

| Package | Purpose |
|---------|---------|
| database | Connection pool, migrations, transactions |
| observability | Metrics, tracing, logging setup |
| resilience | Retry, circuit breaker, rate limiter |

## Data Flow

### Order Creation
```
POST /orders
  ▼
Validate request
  ▼
Begin transaction
  ├─ Create order (pending)
  ├─ Reserve inventory (SELECT FOR UPDATE)
  ├─ Create outbox event: OrderCreated
  └─ Commit
  ▼
Return order ID
```

### Payment Processing
```
POST /payments
  ▼
Extract Idempotency-Key
  ▼
Check idempotency store (Redis + DB)
  ├─ If exists → Return cached response
  └─ If not exists → Process payment
       ├─ Call external provider
       ├─ Store result with idempotency key
       └─ Return response
```

## Failure Domains

| Component | Failure Mode | Mitigation |
|-----------|--------------|------------|
| PostgreSQL | Connection loss | Pool, retry, circuit breaker |
| Redis | Unavailable | Fallback to DB, degrade gracefully |
| External API | Timeout, 5xx | Retry with backoff, circuit breaker |
| Network | Partition | Idempotency, eventual consistency |

## Observability

### Metrics
- Request rate, latency, errors (RED)
- Business metrics: orders, payments, inventory
- Resource: CPU, memory, goroutines, DB pool

### Tracing
- HTTP middleware auto-instrumentation
- DB query spans
- External call spans

### Logging
- Structured JSON (zerolog)
- Request ID correlation
- Error stack traces

## Security

- Parameterized queries (no SQL injection)
- Input validation on all boundaries
- Rate limiting per client
- Secrets via environment (never in code)
- TLS in production