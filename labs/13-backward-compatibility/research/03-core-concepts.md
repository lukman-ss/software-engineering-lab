# Evidence

## Evidence 1: Definisi Backward Compatibility

**Claim**: Backward compatibility adalah kemampuan sistem baru untuk menerima input, request, atau struktur data dari versi lama tanpa crash, error logic, atau kerusakan data.

**Evidence**:
- Google AIP-180: "APIs are fundamentally contracts with users, and users often write code against APIs that is then launched into a production service with the expectation that it continues to work."
- Google AIP-180 mendefinisikan tiga lapisan kompatibilitas: Source compatibility ("Code mesti compile dan run dengan versi baru"), Wire compatibility ("Code harus bisa berkomunikasi dengan server versi baru"), Semantic compatibility ("Code harus tetap menerima apa yang wajar diharapkan").
- Stripe: "Old clients **must** be able to work against newer servers (with the same major version number)."

**Source**: Google AIP-180, Stripe API Versioning  
**URL**: https://google.aip.dev/180 | https://stripe.com/blog/api-versioning  
**Tanggal Publikasi**: AIP-180 (2019-07-23), Stripe (2017-08-05)  
**Confidence**: HIGH  
**Corroborated By**: Martin Fowler Feature Toggle menyebutkan backward compatibility sebagai dasar semua pertimbangan versioning.

## Evidence 2: Backward vs Forward Compatibility

**Claim**: Backward compatibility = new system menerima input dari old system; Forward compatibility = old system menerima input dari new system.

**Evidence**:
- Martin Fowler, Parallel Change: "When implementing BranchByAbstraction, parallel change is a good way to introduce the abstraction layer ... The downside of using parallel change is that during the migrate phase the supplier has to support two different versions, and clients could get confused about which version is new versus old."
- Confluent Schema Registry: BACKWARD = "consumer using schema X can process data produced with schema X or X-1" | FORWARD = "data produced using schema X can be read by consumers with schema X or X-1."
- Postel's Law / Robustness Principle: "Be conservative in what you send, be liberal in what you accept" — ini prinsip forward compatibility.

**Source**: Martin Fowler, Confluent Schema Registry  
**URL**: https://martinfowler.com/bliki/ParallelChange.html | https://docs.confluent.io/platform/current/schema-registry/fundamentals/schema-evolution.html  
**Tanggal Publikasi**: Fowler (2014-05-13), Confluent (2026-06)  
**Confidence**: HIGH

## Evidence 3: Apa yang Membuat Breaking Change

**Claim**: Breaking change terjadi ketika ada modifikasi pada kontrak ekspektasi yang dilanggar oleh consumer.

**Evidence**:
- GitHub API Versions dokumentasi: Breaking changes meliputi "Removing or renaming a parameter, Removing or renaming a response field, Adding a new required parameter, Making a previously optional parameter required, Changing the type of a parameter or response field, Changing authentication or authorization requirements."
- Google AIP-180: "Existing fields and messages **must not** have their type changed, even if the new type is wire-compatible" dan "Renaming a component is semantically equivalent to 'remove and add'."
- Confluent (Avro): Menambah required field = tidak backward compatible; menghapus field wajib = tidak forward compatible.

**Source**: GitHub REST API, Google AIP-180, Confluent Schema Registry  
**URL**: https://docs.github.com/en/rest/about-the-rest-api/api-versions | https://google.aip.dev/180 | https://docs.confluent.io/platform/current/schema-registry/fundamentals/schema-evolution.html  
**Tanggal Publikasi**: GitHub (2026-03-10), Google (2019), Confluent (2026-06)  
**Confidence**: HIGH

## Evidence 4: Expand -> Migrate -> Contract Pattern

**Claim**: Parallel Change pattern membagi breaking change menjadi 3 fase non-breaking: Expand (dukung keduanya), Migrate (pindahkan traffic), Contract (hapus lama).

**Evidence**:
- Martin Fowler, Parallel Change: "In the *expand* phase you augment the interface to support both the old and the new versions" — contoh Grid menambahkan `Map<Coordinate, Cell>` bersamaan dengan `Cell[][]`. "During the *migrate* phase you update all clients using the old version to the new version." "Once all usages have been migrated ... perform the *contract* phase to remove the old version."
- Martin Fowler: "This pattern is particularly useful when practicing ContinuousDelivery because it allows your code to be released in any of these three phases."
- Martin Fowler: "A FeatureFlag can be used to control which version of the interface is used. A feature toggle on the client side allows it to be forward-compatible with the new version of the supplier."

**Source**: Martin Fowler  
**URL**: https://martinfowler.com/bliki/ParallelChange.html  
**Tanggal Publikasi**: 2014-05-13  
**Confidence**: HIGH

## Evidence 5: Zero-Downtime Database Migration (PostgreSQL-Specific)

**Claim**: Schema migration dapat dilakukan tanpa downtime dengan prinsip additive change, indexed CONCURRENTLY, dan NOT VALID constraints.

**Evidence**:
- PostgreSQL ALTER TABLE dokumentasi: "Adding a column with a constant default value does not require each row of the table to be updated when the ALTER TABLE statement is executed."
- PostgreSQL: "CREATE INDEX CONCURRENTLY" menghindari lock panjang.
- PostgreSQL ALTER TABLE NOTA: "NOT VALID option ... reduce the impact of adding a constraint on concurrent updates. With NOT VALID, the ADD CONSTRAINT command does not scan the table and can be committed immediately."
- PostgreSQL: "Adding a column with a volatile DEFAULT ... or a stored generated column ... will cause the entire table and its indexes to be rewritten."

**Implementation-specific caveat**: `CONCURRENTLY`, non-rewriting defaults, and `NOT VALID` are PostgreSQL-specific behaviors, not a universal database rule. Other database engines handle locking and table rewrites differently (common MySQL tooling examples: `pt-online-schema-change`, `gh-ost` — illustrative, not formally cited here).

**Source**: PostgreSQL Global Development Group  
**URL**: https://www.postgresql.org/docs/current/sql-altertable.html  
**Tanggal Publikasi**: September 2026 (PostgreSQL 18)  
**Confidence**: HIGH

## Evidence 6: Dual Read (Fallback Read) — Kapan Dibutuhkan

**Claim**: Dual read diperlukan ketika backfill belum selesai; aplikasi mencoba baca struktur baru, jika kosong/null fallback ke struktur lama.

**Evidence**:
- Martin Fowler Parallel Change: Pada migrate phase, "clients can continue to consume the old version" sementara client baru memakai new version — implisit membutuhkan dual read sampai migrasi penuh.
- Stripe API dokumentasi: "Old clients **must** be able to work against newer servers" — ini berarti server harus dapat fallback ke format lama hingga semua klien bermigrasi.

**Source**: Martin Fowler  
**URL**: https://martinfowler.com/bliki/ParallelChange.html  
**Tanggal Publikasi**: 2014-05-13  
**Confidence**: LOW (inferensi dari pattern; Fowler tidak memberikan contoh kode fallback read eksplisit)

## Evidence 7: Risiko Dual Write

**Claim**: Dual write menambah latency, kompleksitas transaksi, dan berisiko data drift jika write tidak atomic.

**Evidence**:
- Microservices.io, Transactional Outbox: "But without using 2PC, sending a message in the middle of a transaction is not reliable. There's no guarantee that the transaction will commit. Similarly, if a service sends a message after committing the transaction there's no guarantee that it won't crash before sending the message."
- Microservices.io, Transactional Outbox: "The Message relay might publish a message more than once. It might, for example, crash after publishing a message but before recording the fact that it has done so. When it restarts, it will then publish the message again. As a result, a message consumer must be idempotent."
- Penerapan ke dual-write schema migration: jika write ke dua tabel/storage terjadi tanpa transaksi lokal tunggal, kegagalan pada write kedua menyebabkan partial update dan drift antara struktur lama dan baru; kegagalan 2PC/desync antar storage tanpa outbox/transaction log tailing menciptakan kondisi yang sama.

**Source**: Microservices.io (Chris Richardson), Transactional Outbox Pattern  
**URL**: https://microservices.io/patterns/data/transactional-outbox.html  
**Tanggal Publikasi**: Diakses 2026-09-26  
**Confidence**: HIGH

## Evidence 8: Backfill Aman

**Claim**: Backfill umum dilakukan secara batch, throttled, idempotent, dan resumable untuk menghindari lock dan memungkinkan recovery.

**Evidence**:
- PostgreSQL: full table scan untuk constraint baru "can take a long time" — implication: batch processing diperlukan untuk menghindari long-running transactions dan lock.
- Prinsip engineering umum: checkpoint-based processing untuk job yang dapat dihentikan; idempotency dan throttling adalah best practice umum bukan spesifik PostgreSQL.

**Source**: PostgreSQL Global Development Group (implikasi performa) + prinsip engineering standar  
**URL**: https://www.postgresql.org/docs/current/sql-altertable.html  
**Tanggal Publikasi**: September 2026  
**Confidence**: LOW (sumber hanya menyatakan performa scan, tidak secara eksplisit mendefinisikan batching/throttle/idempotency/resumability)

## Evidence 9: Compatibility Saat Rolling Deployment

**Claim**: Database schema harus kompatibel dengan version N dan N+1 secara bersamaan selama rolling deployment.

**Evidence**:
- Martin Fowler, Blue Green Deployment: "The trick is to separate the deployment of schema changes from application upgrades. So first apply a database refactoring to change the schema to support both the new and old version of the application, deploy that, check everything is working fine so you have a rollback point, then deploy the new version of the application."
- Google AIP-180: Source compatibility menjadi kunci ketika versi berbeda berjalan serentak.

**Source**: Martin Fowler  
**URL**: https://martinfowler.com/bliki/BlueGreenDeployment.html  
**Tanggal Publikasi**: 2010-03-01  
**Confidence**: HIGH

## Evidence 10: Mendeteksi Consumer Lama Masih Aktif

**Claim**: Observability (metrics, structured logs, deprecation headers) menentukan apakah legacy contract masih dipakai.

**Evidence**:
- GitHub API Versions: "While a version is within its support window but approaching closing down, GitHub includes [Deprecation] and [Sunset] headers in API responses."
- Martin Fowler, Feature Toggle: "Any system using feature flags should expose some way for an operator to discover the current state of the toggle configuration."

**Source**: GitHub REST API documentation, Martin Fowler  
**URL**: https://docs.github.com/en/rest/about-the-rest-api/api-versions  
**Tanggal Publikasi**: GitHub (2026-03-10)  
**Confidence**: HIGH

## Evidence 11: Kapan Legacy Boleh Dihapus

**Claim**: Legacy hanya boleh dihapus saat traffic ke legacy interface/methods/field = 0.

**Evidence**:
- GitHub: "Requests that specify a closing down API version receive a 410 Gone response" hanya setelah sunset date.
- Martin Fowler Parallel Change: "If the contract phase is not executed you might end up in a worse state than you started, therefore you need discipline to finish the transition successfully."

**Source**: GitHub REST API  
**URL**: https://docs.github.com/en/rest/about-the-rest-api/api-versions  
**Tanggal Publikasi**: 2026-03-10  
**Confidence**: HIGH

## Evidence 12: Feature Flag untuk Rollout dan Rollback

**Claim**: Feature flag memisahkan deployment dari aktivasi; memungkinkan canary rollout 1% → 10% → 100% dan instant rollback.

**Evidence**:
- Martin Fowler, Feature Toggle: "Release Toggles are feature flags used to enable trunk-based development for teams practicing Continuous Delivery. They allow in-progress features to be checked into a shared integration branch while still allowing that branch to be deployed to production at any time."
- Martin Fowler, Parallel Change: "During the migrate phase, a FeatureFlag can be used to control which version of the interface is used."
- Martin Fowler, Feature Toggles article: "Canary releasing ... only turning the new feature on for a small percentage of their total userbase - a 'canary' cohort" via kohort berdasarkan user ID.

**Source**: Martin Fowler  
**URL**: https://martinfowler.com/bliki/FeatureFlag.html | https://martinfowler.com/articles/feature-toggles.html  
**Tanggal Publikasi**: Feature Flag (2010-10-29), Feature Toggles (2017-10-09)  
**Confidence**: HIGH

## Evidence 13: Idempotent dan Resumable Migration

**Claim**: Migration harus idempotent dan resumable untuk crash recovery.

**Evidence**:
- PostgreSQL ALTER TABLE: "NOT VALID" memungkinkan menambah constraint tanpa scan penuh, kemudian memvalidasi nanti via "VALIDATE CONSTRAINT."
- Umum: UPSERT dan query `WHERE new_field IS NULL` adalah pola idempotent.

**Source**: PostgreSQL  
**URL**: https://www.postgresql.org/docs/current/sql-altertable.html  
**Tanggal Publikasi**: September 2026  
**Confidence**: MEDIUM (prinsip umum; UPSERT tidak secara eksplisit didokumentasikan di sumber ini)

## Evidence 14: Rollback dan Database Schema

**Claim**: Jika schema hanya ada expand (non-destructive), rollback aman; jika sudah contract, rollback dapat menyebabkan data loss.

**Evidence**:
- Martin Fowler, Blue Green Deployment: "first apply a database refactoring to change the schema to support both the new and old version of the application, deploy that, check everything is working fine so you have a rollback point, then deploy the new version."
- Martin Fowler Parallel Change: "When you have control over all usages of the interface, following this pattern is still useful because it prevents you from spreading breakage across the entire codebase all at once."

**Source**: Martin Fowler  
**URL**: https://martinfowler.com/bliki/BlueGreenDeployment.html | https://martinfowler.com/bliki/ParallelChange.html  
**Tanggal Publikasi**: 2010-03-01, 2014-05-13  
**Confidence**: HIGH

## Evidence 15: API Contract Compatibility vs Database Schema Compatibility

**Claim**: API compatibility lebih sulit dipertahankan karena kontrol terhadap consumer terbatas (mobile app, partner eksternal).

**Evidence**:
- GitHub: "When a new REST API version is released, the previous API version will be supported for at least 24 more months."
- Google AIP-180: "It is not always clear whether a change is compatible or not."
- Stripe: "maintaining compatibility is important, but ... we expect to eventually start retiring our older API versions."

**Source**: GitHub REST API, Google AIP-180  
**URL**: https://docs.github.com/en/rest/about-the-rest-api/api-versions | https://google.aip.dev/180  
**Tanggal Publikasi**: GitHub (2026-03-10), Google (2019)  
**Confidence**: HIGH