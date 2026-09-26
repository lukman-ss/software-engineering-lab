# Open Questions dan Research Frontiers

## Unanswered Questions

### 1. Transaction Consistency pada Dual Write

**Pertanyaan**: Bagaimana tim mengelola dual-write consistency antar microservice secara aman tanpa menggunakan two-phase commit (2PC) atau mengorbankan throughput?

**Status**: **UNVERIFIED**
- Prinsip umum: outbox pattern, eventual consistency, atau saga pattern
- Tidak ada sumber Tier 1 yang mendeskripsikan implementasi spesifik untuk microservice

**Posibilitas Next Research**:
- Menguji performance outbox pattern vs dual write synchronous
- Studi kompromi antara consistency dan latency

### 2. Tooling untuk Blocking Destructive Operations

**Pertanyaan**: Apakah migration tools (Flyway, Liquibase, Prisma Migrate) secara native mendukung otomatis memvalidasi bahwa operasi destruktif tidak dijalankan pada fase expand?

**Status**: **PARTIALLY VERIFIED**
- Flyway: Mendukung callbacks dan pelacakan baseline migrations (tidak auto-block destructive)
- Liquibase: Mendukung context filtering (ops=production vs ops=dev) untuk memisahkan destructive changes
- Prisma: Tidak ada dokumentasi resmi tentang blocking destructive operations

**Sumber yang Diperlukan**:
- Dokumentasi resmi flyway.liquibase.org tentang peluncuran migrasi
- PR atau issue di GitHub prisma/migrate yang membahas safe migrations

### 3. Consumer Lag pada API Versioning

**Pertanyaan**: Bagaimana API provider secara aman melakukan contraction bila sebagian kecil kritikal B2B consumer menolak migrasi off deprecated version selama bertahun-tahun?

**Status**: **UNVERIFIED**
- GitHub menyatakan: "Requests that specify a closing down API version receive a 410 Gone response" setelah sunset
- Praktik nyata: Stripe menggunakan rolling versions (tanggal) sehingga setiap request dapat memiliki versi berbeda
- Tidak ada sumber yang mendokumentasikan strategi untuk "uncooperative clients"

**Pertimbangan**:
- Payment gateway seperti Stripe tidak akan pernah memblokir pembayaran karena tidak update API
- Solusi: feature flag pada server-side untuk tetap handle old contract, atau sunset berbayar

### 4. Quantitative Overhead Transformasi

**Pertanyaan**: Overload CPU/Memory dari transformation modules (seperti di Stripe) yang berjalan untuk setiap request - berapa banyak yang dapat diukur?

**Status**: **LOW CONFIDENCE**
- Stripe hanya deskripsikan konseptual: "version change modules" memproses tiap response
- Tidak ada studi kuantitatif tentang overheadnya di PostgreSQL/AWS/other implementations

**Resource yang Diperlukan**:
- Benchmark internal Stripe atau penggunaannya
- Implementasi open-source seperti Prisma atau Hasura dengan metrics

---

## Weak Evidence

### Dual Write & Fallback Read Overhead

**Claim**: Dual write menambah latency pada request path.

**Evidence**: Prinsip konvensional, tidak ada benchmark nyata.

**Sumber**: Umum di dokumentasi engineering; tidak ada sumber Tier 1 yang menyebutkan angka overhead spesifik.

### Batch Backfill vs Lazy Backfill Threshold

**Claim**: Lazy backfill lebih baik untuk dataset kecil/kurang penggunaan intensif, batch backfill untuk dataset besar.

**Evidence**: Conthokan di dokumentasi database; tidak ada studi komparatif kuantitatif.

**Sumber**: Implisit dari PostgreSQL performance guidelines.

---

## Claims Needing Deeper Research

### CDC (Change Data Capture) untuk Migration

**Hipotesis**: Menggunakan CDC (Debezium, logical replication) untuk dual-write/zerodowntime migration berbeda dari mengubah kode aplikasi.

**Potensi Benefit**:
- Mengurangi complexity di application layer
- Memisahkan concern migration dari business logic
- Potentially better consistency guarantees

**Research Direction**:
- Bandingkan CDC-based migration vs application-level dual write
- Studi failure modes pada CDC-based approach
- Benchmark throughput dibandingkan synchronous dual write

### GraphQL Schema Evolution vs REST Versioning

**Hipotesis**: GraphQL memungkinkan evolusi schema secara natural melalui field addition/optional fields, berbeda dengan REST yang memerlukan versioning.

**Research Direction**:
- Analisis Field Masking di GraphQL (Apollo Federation) sebagai mekanisme Expand/Migrate/Contract
- Bandingkan client update frequency GraphQL vs REST
- Studi kegagalan: GraphQL introspection cache versioning

---

## Next Research Directions

1. **Change Data Capture (CDC) Pattern**
   - Fokus: Debezium, PostgreSQL logical replication, Google Cloud Datastream
   - Tujuan: Menilai CDC sebagai alternatif untuk dual-write application-level

2. **GraphQL Schema Evolution**
   - Fokus: Apollo Federation, Hasura, Prisma GraphQL API
   - Tujuan: Bagaimana Field Addition/Removal ditangani secara native di GraphQL

3. **Database-Specific Zero-Downtime Strategies**
   - PostgreSQL: Logical replication, pglogical
   - MySQL:gh-ost, pt-online-schema-change, gh-pro
   - CockroachDB: Online schema changes with distributed transactions
   - Hasura: Metrika untuk GraphQL schema migration

4. **Large-Scale Provider Case Studies**
   - Netflix: Zuul API gateway versioning strategi
   - Amazon AWS: API versioning patterns (Service evolution)
   - Facebook: GraphQL schema versioning

5. **Automated Detection of Consumer Lag**
   - Instrumentasi otomatis untuk deteksi penggunaan field legacy
   - Machine learning untuk memprediksi migration deadline based on consumer behavior

---

## Research Questions yang Belum Terjawab secara Komprehensif

| # | Pertanyaan | Current State | Priority |
|---|------------|---------------|----------|
| 1 | Apakah ada tool open-source yang secara otomatis menggenerasi migration untuk Expand/Migrate/Contract sequence? | NOT VERIFIED | MEDIUM |
| 2 | Bagaimana CDC (Debezium) dapat digunakan untuk menggantikan application-level dual write? | NOT VERIFIED | HIGH |
| 3 | Benchmark overhead request/response transformation di Stripe atau platform serupa | NOT VERIFIED | MEDIUM |
| 4 | Strategi untuk "uncooperative clients" yang tidak update selama bertahun-tahun | PARTIALLY VERIFIED | HIGH |
| 5 | Perbedaan performansi migrasi pada PostgreSQL vs MySQL vs CockroachDB | NOT VERIFIED | MEDIUM |

---

## Sumber yang Diperlukan untuk Verifikasi

1. **Flyway/Liquibase Documentation** — tentang rollback strategies dan destructive change prevention
2. **Debezium Documentation** — CDC untuk zero-downtime migrations
3. **Stripe Engineering Internal Metrics** — overhead transformasi api version module
4. **Enterprise Case Studies** — Netflix, AWS, Uber tentang migrasi skala produksi