# Open Questions: Saga Pattern

## Unanswered Questions

1. **Bagaimana implementasi saga dalam teknologi spesifik?**
   - Java Spring Boot: Axon Framework, Eventuate
   - .NET: MassTransit, NServiceBus
   - Node.js:微services saga framework
   - Coverage: LOW (belum ada evidence untuk stack spesifik)

2. **Bagaimana performance comparison saga vs 2PC?**
   - Throughput (ops/sec)
   - Latency (p95, p99)
   - Throughput under load
   - Coverage: NOT VERIFIED

3. **Bagaimana monitoring dan observability saga?**
   - Distributed tracing
   - Saga state visualization
   - Dead letter queue handling
   - Coverage: LOW (Microsoft menyebutkan "Need for monitoring and tracking sagas" tapi detail kurang)

4. **Bagaimana saga handle partial failure dalam compensating transactions?**
   - Jika kompensasi juga gagal, apa recovery mechanism?
   - Coverage: LOW (Microsoft menyebutkan "Compensating transactions might not always succeed")

5. **Eventual consistency - berapa lama consistency tercapai?**
   - Time to convergence
   - SLA implications
   - Coverage: NOT VERIFIED

## Weak Evidence

1. **Countermeasures implementation:** Microsoft menyebutkan 6 countermeasures tapi detail implementasi kurang dalam free sources

2. **Saga vs TCC (Try-Confirm-Cancel) pattern:** Belum ada perbandingan dalam evidence yang dikumpulkan

3. **Saga in CQRS/Event Sourcing context:** Hubungan dengan event sourcing belum terverifikasi secara mendalam

## Claims Needing Deeper Research

1. **Original Garcia-Molina & Salem 1987 paper:** Perlu akses ke full paper untuk definisi formal saga

2. **Real-world production failure cases:** Belum ada case study spesifik dari production failures terkait saga implementation

3. **Idempotency implementation patterns:** Perlu research lebih dalam tentang idempotency keys, status reconciliation

## Next Research Directions

1. **Framework comparison:** Perbandingan saga frameworks (Axon, Eventuate, MassTransit)
2. **Performance benchmarks:** Throughput/latency comparison saga vs 2PC
3. **Production case studies:** Case study dari company yang menggunakan saga di production
4. **Saga + Event Sourcing:** Deep dive ke hubungan saga dengan event sourcing
5. **TCC Pattern:** Investigasi TCC (Try-Confirm-Cancel) sebagai alternatif saga
