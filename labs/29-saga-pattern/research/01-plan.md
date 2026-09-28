# Research Plan: Saga Pattern — Mengelola Transaksi Terdistribusi Tanpa 2PC

## Research Topic

Saga Pattern — Transaksi terdistribusi alternatif untuk microservices tanpa 2PC

## Objective

Investigasi dan verifikasi klaim teknis dalam lab:
1. 2PC (Two-Phase Commit) terlalu lambat/tidak feasible untuk microservices berbasis database-per-service
2. Saga memecah transaksi terdistribusi menjadi series local transactions
3. Dua pendekatan utama: Choreography (event-driven) dan Orchestration (centralized coordinator)
4. Compensating transactions sebagai kunci rollback saga
5. Perbedaan semantics antara ACID rollback dan compensating transaction
6. Saga menjamin eventual consistency, bukan strict ACID

## Research Questions

### RQ1: Mengapa 2PC Tidak Feasible untuk Microservices
- Bagaimana 2PC bekerja dan keterbatasannya?
- Mengapa latency/availability trade-off membuat 2PC buruk untuk microservices?

### RQ2: Definisi dan Semantics Saga
- Apa definisi formal saga?
- Bagaimana saga memastikan atomicity tanpa isolasi?

### RQ3: Choreography vs Orchestration
- Apa trade-off masing-masing pendekatan?
- Kapan harus pakai choreography vs orchestration?

### RQ4: Compensating Transactions
- Bagaimana compensating transaction berbeda dari rollback?
- Bagaimana idempotency dan ordering semantics?

### RQ5: Anomali Data dan Countermeasures
- Bagaimana sagas menangani anomali isolasi (dirty reads, fuzzy reads)?
- Apa countermeasures yang direkomendasikan?

## Search Strategy

1. Microsoft Azure Architecture Center (official docs)
2. Chris Richardson / Microservices.io (microservices patterns authority)
3. Chris Richardson book references
4. Eventuate.io (Chris Richardson framework)
5. Martin Fowler / InfoQ references

## Expected Primary Sources

1. Microsoft Azure Architecture Center - Saga Design Pattern
2. Chris Richardson / Microservices.io - Saga Pattern
3. Chris Richardson - Microservices Patterns (book reference)

## Risks / Unknowns

- Beberapa URL mungkin berubah (404)
- Akademis paper (Garcia-Molina & Salem, 1987) mungkin behind paywall
- Tidak semua klaim lab bisa diverifikasi dari free sources
