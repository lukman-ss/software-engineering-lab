# Sumber - Backward Compatibility Research

## Source 1: Parallel Change (Expand and Contract Pattern)
**Judul**: Parallel Change  
**Penulis**: Joshua Kerievsky (asli), Danilo Sato, Martin Fowler  
**Publisher**: martinfowler.com (ThoughtWorks)  
**URL**: https://martinfowler.com/bliki/ParallelChange.html  
**Tanggal Publikasi**: 2014-05-13  
**Tanggal Diakses**: 2026-09-25  
**Tier Sumber**: Tier 1 (Dokumentasi Arsitektur Profesional)  
**Relevansi**: Fundamentasi pola Expand -> Migrate -> Contract; contoh kode Grid dengan Coordinate class; grafik visual fase expand/migrate/contract.

## Source 2: Stripe API Versioning
**Judul**: APIs as infrastructure: future-proofing Stripe with versioning  
**Penulis**: Brandur Leach, API Experience  
**Publisher**: Stripe Engineering Blog  
**URL**: https://stripe.com/blog/api-versioning  
**Tanggal Publikasi**: 2017-08-05  
**Tanggal Diakses**: 2026-09-25  
**Tier Sumber**: Tier 1 (Dokumentasi Teknis Asli dari Perusahaan Infrastruktur)  
**Relevansi**: Strategi versioning API tanggal (rolling versions), version change modules, transformation pipelines, API resource versioning, contoh kode ChargeAPIResource dan VersionChanges.

## Source 3: AIP-180 - Google API Backward Compatibility
**Judul**: Backwards compatibility (AIP-180)  
**Publisher**: Google API Improvement Proposals  
**URL**: https://google.aip.dev/180  
**Tanggal Publikasi**: 2019-07-23 (Diperbarui 2025-10-21)  
**Tanggal Diakses**: 2026-09-25  
**Tier Sumber**: Tier 1 (Standar Industri API desain)  
**Relevansi**: Definisi tiga tipe kompatibilitas (source, wire, semantic); panduan mendetail tentang penambahan/menghapus komponen; aturan untuk enum, string length, resource names.

## Source 4: Feature Toggle Patterns
**Judul**: Feature Toggle  
**Penulis**: Martin Fowler, Pete Hodgson  
**Publisher**: martinfowler.com (ThoughtWorks)  
**URL**: https://martinfowler.com/bliki/FeatureFlag.html  
**Tanggal Publikasi**: 2010-10-29 (Artikel 2017)  
**Tanggal Diakses**: 2026-09-25  
**Tier Sumber**: Tier 1 (Literatur Arsitektur Distribusi)  
**Relevansi**: Kategori toggle (Release, Experiment, Ops, Permissioning); implementasi polymorphic substitution; pengelolaan konfigurasi; testing strategies.

## Source 5: PostgreSQL ALTER TABLE Documentation
**Judul**: ALTER TABLE  
**Publisher**: PostgreSQL Global Development Group  
**URL**: https://www.postgresql.org/docs/current/sql-altertable.html  
**Tanggal Publikasi**: September 2026 (PostgreSQL 18)  
**Tanggal Diakses**: 2026-09-25  
**Tier Sumber**: Tier 1 (Dokumentasi Database Resmi)  
**Relevansi**: Operasi ALTER TABLE non-blocking; cara menambahkan kolom dengan default; mencegah table rewrite; DROP COLUMN melembapkan data.

## Source 6: GitHub API Versioning & Deprecation
**Judul**: API Versions  
**Publisher**: GitHub Developer Documentation  
**URL**: https://docs.github.com/en/rest/about-the-rest-api/api-versions  
**Tanggal Publikasi**: 2026-03-10 (versi API terbaru)  
**Tanggal Diakses**: 2026-09-25  
**Tier Sumber**: Tier 1 (Dokumentasi Platform API)  
**Relevansi**: Daftar breaking changes vs additive changes; header Deprecation dan Sunset (RFC 8594); jendela dukungan 24 bulan.

## Source 7: Confluent Schema Registry - Schema Evolution
**Judul**: Schema Evolution and Compatibility for Schema Registry  
**Publisher**: Confluent Documentation  
**URL**: https://docs.confluent.io/platform/current/schema-registry/fundamentals/schema-evolution.html  
**Tanggal Publikasi**: Juni 2026  
**Tanggal Diakses**: 2026-09-25  
**Tier Sumber**: Tier 1 (Dokumentasi Streaming Data)  
**Relevansi**: Definisi BACKWARD, FORWARD, FULL compatibility; aturan evolusi schema Avro/Protobuf/JSON Schema; urutan upgrade klien.

## Source 8: Blue Green Deployment
**Judul**: Blue Green Deployment  
**Penulis**: Martin Fowler  
**Publisher**: martinfowler.com (ThoughtWorks)  
**URL**: https://martinfowler.com/bliki/BlueGreenDeployment.html  
**Tanggal Publikasi**: 1 Maret 2010  
**Tanggal Diakses**: 2026-09-25  
**Tier Sumber**: Tier 1 (Literatur Deployment)  
**Relevansi**: Membandingkan dua environment; rollback cepat; database schema changes terpisah dari application upgrade.

## Source 9: Semantic Versioning 2.0.0
**Judul**: Semantic Versioning 2.0.0  
**Penulis**: Tom Preston-Werner  
**Publisher**: semver.org  
**URL**: https://semver.org/  
**Tanggal Publikasi**: 2013 (versi 2.0.0)  
**Tanggal Diakses**: 2026-09-25  
**Tier Sumber**: Tier 1 (Standar Industri)  
**Relevansi**: Definisi MAJOR (breaking), MINOR (backward compatible), PATCH (bug fix); aturan peningkatan versi.

## Source 10: Transactional Outbox Pattern (Dual-Write Failure Modes)
**Judul**: Pattern: Transactional outbox  
**Penulis**: Chris Richardson  
**Publisher**: Microservices.io  
**URL**: https://microservices.io/patterns/data/transactional-outbox.html  
**Tanggal Publikasi**: Diakses 2026-09-26  
**Tier Sumber**: Tier 1 (Dokumentasi Pola Arsitektur Praktik Industri)  
**Relevansi**: Menggambarkan dual-write problem saat atomisitas database dan message broker tidak dapat dipenuhi oleh 2PC; penjelasan tentang data inconsistency, message ordering, dan keharusan idempotency consumer. Digunakan sebagai sumber utama untuk klaim risiko dual write pada research/03-core-concepts.md Evidence 7 dan research/08-failure-modes.md Failure Mode 2.
