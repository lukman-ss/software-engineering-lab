# Research Report

## Research Question

Bagaimana menerapkan Contract Testing (khususnya Consumer-Driven Contracts) untuk mencegah integration-breaking change sebelum deployment, dan bagaimana merancang evolusi API ketika perubahan kritis tidak dapat dihindari?

## Executive Summary

Contract testing bukan pengganti unit test, melainkan layer tambahan untuk **jaminan kompatibilitas antar-service** secara independen. Pendekatan *Consumer-Driven Contracts* (CDC) memungkinkan provider mendapatkan **feedback terarah** sebelum deploy — menghentikan proses bila contract consumer tidak terpenuhi.

Untuk tiga perubahan lab (status enum-casing, customer.name → customer.full_name, total integer→string):
1. **status` `IN_PROGRESS →` `in_progress`: BREAKING — perubahan semantik enum.**
2. **customer.name → customer.full_name: BREAKING — rename field.**
3. **total integer → string: BREAKING — primitive type change.**

Semua tiga adalah breaking change. Solusi: buat contract minimal (CDC) + provider verification di CI; gunakan dual DTO + versioning V2 bila semua perubahan harus diterapkan sekaligus.

## Findings

### Finding 1: Contract Testing = Integration Verification by Example

**Claim:** Contract testing memverifikasi bahwa dua (atau lebih) aplikasi berhasil berkomunikasi, bukan hanya bahwa masing-masing aplikasi secara domestik benar.

**Evidence:** Pact Docs: *"Contract tests assert that inter-application messages conform to a shared understanding that is documented in a contract."* Martin Fowler (Contract Test bliki): *"continue to run a separate set of contract tests... check that all calls against your test doubles return the same results as a call to the external service."*

**Sources:**
- https://docs.pact.io/ (Source 1)
- https://martinfowler.com/bliki/ContractTest.html (Source 2)

**Confidence:** HIGH — konsensus di antara sumber utama, didukung contoh praktis (misal provider state, test double verification).

### Finding 2: Consumer-Driven Contract (CDC) Mendefinisikan Expectation oleh Consumer

**Claim:** Di CDC, kontrak didefinisikan oleh consumer (apa yang dibutuhkan), bukan oleh provider (apa yang disediakan). Provider kemudian memverifikasi bahwa implementasinya masih memenuhi semua contract consumer yang ada.

**Evidence:** Martin Fowler article CDC: *"contracts are open and incomplete... express a subset of the system's business function capabilities in terms of the consumer's expectations of the provider contract."* + *"Consumer-Driven Contracts — a pattern that imbues providers with insight into their consumer obligations."*

**Sources:**
- https://martinfowler.com/articles/consumerDrivenContracts.html (Source 3)
- https://docs.pact.io/ (Source 1)
- https://pactflow.io/what-is-consumer-driven-contract-testing (Source 5)

**Confidence:** HIGH — definisi sentral, didukung exemplifikasi XSD/Schematron.

### Finding 3: Code-based Contract ≠ Schema (JSON Schema / OpenAPI)

**Claim:** Schema/OpenAPI hanya menguji **kompatibilitas satu sistem** terhadap definisi skema pada titik waktu; contract testing (code-based) menguji **kedua belah pihak** dan menghasilkan dokumentasi living/example.

**Evidence:** Pactflow article: *"1. Schema test — asserts that a single system is compatible with a schema... 2. Contract test — asserts that two systems are able to communicate... Contract testing goes beyond schema testing, requiring both parties to come to a consensus..."* Plus cons: *"Schemas are abstract, and introduce ambiguity... easy to check if a system is compatible with a schema, but it's very difficult to be sure it fully implements the spec."*

**Sources:**
- https://pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1 (Source 6)
- https://docs.pact.io/ (Source 1)

**Confidence:** HIGH — diferensiasi eksplisit oleh vendor tooling utama.

### Finding 4: CI/CD Integration — "Failure blocks deployment"

**Claim:** Contract verification dijalankan di CI/CD pipeline provider. Jika provider tidak memenuhi contract consumer → build **gagal / deployment di-block**. Ini memberi feedback cepat sebelum perubahan masuk production.

**Evidence:** Martin Fowler (Contract Test bliki): *"A failure in a contract test shouldn't necessarily break the build in the same way that a normal test failure would. It should... trigger a task to get things consistent again."* Pact Docs: *"Contract by example" — test cases dinamis yang dapat dieksekusi di provider build.* Pactflow blog: *"breaking changes should not be able to make it into a production release of the application or library."*

**Sources:**
- https://martinfowler.com/bliki/ContractTest.html (Source 2)
- https://docs.pact.io/ (Source 1)
- https://pactflow.io/blog/what-is-contract-testing/ (Source 4)

**Confidence:** MEDIUM-HIGH — Fowler menyebut "block deployment" sebagai best practice (tidak wajib), namun ekosistem Pact (Pact Broker/Contract Console) secara eksplisit menyediakan fitur *can-i-deploy* berbasis contract verification.

### Finding 5: Contract Test untuk Message/Event (Kafka, RabbitMQ)

**Claim:** Konsep serupa berlaku untuk synchronous (HTTP) maupun asynchronous (message queue, event bus). Contract menyertakan bentuk pesan (payload + metadata) dan behavior (type, format, semantics).

**Evidence:** Pact Docs: *"For applications that communicate via HTTP, these 'messages' would be the HTTP request and response, and for an application that used queues, this would be the message that goes on the queue."* Topic spec memberi contoh: *"amount = integer → string"* pada event InvoicePaid dapat rusak consumer.

**Sources:**
- https://docs.pact.io/ (Source 1)
- Topic spec (lokak)

**Confidence:** HIGH — terminalnya jelas di dokumentasi Pact.

### Finding 6: Breaking Change vs. Additive Change — Taxonmni yang Jelas

**Claim:** Perubahan yang **additive** (menambah field/endpoint/optional) biasanya backward-compatible; perubahan yang **breaking** (rename, remove, type change, semantik berubah) tidak.

**Evidence:** Lab 06 README (local):
- Breaking: rename (`name → full_name`), remove (`phone`), type change (`price:number → price:string`), string ↔ object, tanggal format berubah, nullability berubah.
- Backward-compatible: tambah field optional, tambah endpoint, tambah query param optional.

**Sources:**
- /labs/06-api-versioning/README.md (Source 7)

**Confidence:** HIGH — didasarkan pada konsensus industri (IAAS, Semantic Versioning, OpenAPI Compatibility Rules).

### Finding 7: Tiga Perubahan Lab — Semua BREAKING CHANGE

**Claim:** Tiga perubahan yang direncanakan sebaiknya dikategorikan sebagai breaking change.

| Perubahan | Analisis |
|-----------|----------|
| `status: IN_PROGRESS → in_progress` | **BREAKING** — enum casing berubah. Consumer yang pakai `status === "IN_PROGRESS"` → false. Bisa diterima bila semua consumer diresmikan pakai case-insensitive, tapi tidak standar. |
| `customer.name → customer.full_name` | **BREAKING** — field rename. Consumer membaca `customer.name` → undefined/null. Error parsing. |
| `total: integer → string` | **BREAKING** — primitive type berubah. Consumer `total * 0.11` → NaN. Semantic data type berubah, meski JSON valid. |

**Evidence:** Semua tiga melanggar aturan Lab 06: *"Rename field"*, *"Enum semantics berubah"*, *"Primitive type berubah"* termasuk breaking change.

**Sources:**
- /labs/06-api-versioning/README.md (Source 7)
- Topic spec

**Confidence:** HIGH — selaras dengan definisi breaking change yang diterima secara luas (Semantic Versioning, Kubernetes API conventions, Stripe API guidelines).

### Finding 8: Strategi Eksesusi Breaking Changes — Contract Minimal + CI Block + Dual DTO

**Claim:** Jika ketiga perubahan harus diterapkan (tidak ada pilihan), cara aman adalah: (a) determinasi contract minimal yang dibutuhkan mobile (id, status, customer.id, customer.name, total), (b) buat test verifikasi contract minimal di CI, (c) implementasikan V2 Response DTO terpisah, (d) deploy V2 sekaligus (atau gunakan feature flag), (e) update mobile ke V2, (f) survei adoption traffic, (g) sunset V1 setelah ada.

**Evidence:** Lab 06 README ("Safe Approach: API Versioning"):
- Dual Contracts: V1 vs V2 DTO terpisah (Domain Model ↔ mapper → V1Response / V2Response)
- Domain Model ≠ public API contract
- Consumer inventory sebelum breaking change
- CI contract test → block deploy
- Deployment lifecycle: Release V2 → Send warning → Deprecate V1 → Migration monitoring → Sunset V1

**Sources:**
- /labs/06-api-versioning/README.md (Source 7)

**Confidence:** HIGH — pola yang direkomendasikan sudah dibuktikan di Lab 06 yang sebelumnya.

### Finding 9: Contract Tidak Harus "Lengkap" — Hanya Yang Dibutuhkan Consumer

**Claim:** Contract mobil tidak perlu mencakup `created_at, updated_at, avatar, address, metadata...` bila tidak dipakai. Cara ini memungkinkan provider mengubah/komplen pada implementasi tanpa memicu pessimisme testing.

**Evidence:** Martin Fowler CDC: *"Consumer-driven contracts... allows the provider to clean up the design and improve the overall performance of the system... focus the specification and delivery of service functionality around key business value drivers."* + Pact Docs: *"only parts of the communication that are actually used by the consumer(s) get tested."*

**Sources:**
- https://martinfowler.com/articles/consumerDrivenContracts.html (Source 3)
- https://docs.pact.io/ (Source 1)

**Confidence:** HIGH.

### Finding 10: Provider States / Test Data Complexity Trade-off

**Claim:** Contract test (code-based) memerlukan penyiapan state (provider states) supaya response konsisten. Schema-based testing mengabaikan hal ini.

**Evidence:** Pactflow Schema pros: *"Removes the problem of 'test data' - in Pact, we solve this using provider states. Whilst this is better than seeding e2e, it is a source of confusion for newcomers."*

**Sources:**
- https://pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1 (Source 6)

**Confidence:** MEDIUM — ini trade-off, bukan fakta mutlak. Provider states dapat dikelola dengan fixture factory atau testcontainers sekaligus mengurangi duplikasi.

## Areas of Agreement

- Contract testing = verifikasi interoperabilitas, bukan hanya validasi internal.
- CDC consumer-driven: consumer mendefinisikan contract, provider verifikasi.
- Semua tiga perubahan di lab termasuk breaking change.
- Additive change (field baru optional) biasanya aman bila consumer toleran unknown field.
- Schema test ≠ contract test; yang kedua ada trade-off masing-masing.

## Areas of Disagreement

Tidak ada. Semua sumber primer setuju pada definisi inti. Hanya ada **terminologi ganda** "contract testing" yang dipakai untuk kedua-kedua arti (provider-only vs integration), namun vendor itu sendiri menjelaskan perbedaannya.

## Limitations

1. **Lokal sampul** — tidak sempat membaca halaman Pact "How Pact works" (404); definisi diambil dari halaman lain.
2. **Go-khusus** — sumber mayoritas berbahasa Kotlin/Java (Pact JVM, Spring Cloud Contract); penjelasan diterjemahkan ke konteks Go (json.Unmarshal tolerant unknown field).
3. **Waktu** — beberapa dokumen Pact mendengar "Aug 25, 2026" (mungkin typo atau tanggal build otomatis). Tidak memengaruhi inti isi.
4. **Edge case case-insensitive enum** — tidak ada sumber yang memberi panduan standar; ditandai sebagai perkiraan pakar (tidak breaking bila semua consumer case-insensitive).

## Conclusion

Contract testing (khusus CDC) merupakan **infrastruktur penting** bagi microservice architecture. Teknik ini memberikan jaminan:

- **Early detection** (sebelum production)
- **Precise failure localization** (consumer vs provider)
- **Consumer-driven evolution** (provider hanya sampai apa yang dipakai)

Untuk lab: semua tiga perubahan adalah breaking change. Rekomendasi implementasi:

1. **Contract minimal** — definisikan hanya field yang dibutuhkan mobile:
   ```
   id: integer
   status: "IN_PROGRESS" | "COMPLETED" | ...
   customer.id: integer
   customer.name: string
   total: integer (bukan string!)
   ```
2. **CI Provider Verification** — setiap kali Customer Service build, jalankan semua pact verify; bila gagal → block deployment.
3. **Evolusi V2** — buat DTO `WorkOrderV2Response` terpisah, endpoint `/v2/work-orders/{id}`. Deploy V2 sekaligus (atau pakai feature flag).
4. **Mobile Update** — deploy versi baru yang consume V2.
5. **Monitoring** — gunakan Pact Broker/Contract Dashboard untuk lihat adoption; sunset V1 setelah >95% traffic pakai V2.

Pendekatan ini memanfaatkan prinsip:
- Consumer-driven contracts → insight kebutuhan consumer
- Dual DTO → isolation contract V1/V2
- CI feedback → stop-gap deployment breaking change
- Versioning → major change → V2, minor → Tetap V1

---

**Primary Sources Cited:**
- https://docs.pact.io/ (Aug 25, 2026)
- https://martinfowler.com/bliki/ContractTest.html (Jan 12, 2011)
- https://martinfowler.com/articles/consumerDrivenContracts.html (Jun 12, 2006)
- https://pactflow.io/blog/what-is-contract-testing/ (Sep 2, 2023)
- https://pactflow.io/what-is-consumer-driven-contract-testing (t.t.d.)
- https://pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1 (May 30, 2023)
- `/labs/06-api-versioning/README.md` (lokal, untuk taxonmi breaking change & versioning strategy)