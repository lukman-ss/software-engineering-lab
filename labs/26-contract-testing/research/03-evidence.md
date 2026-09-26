# Evidence

## Evidence 1

Claim: Contract testing memastikan dua (atau lebih) aplikasi sepakat tentang bentuk komunikasi mereka, tanpa menjalankan seluruh sistem end-to-end.
Evidence: *"Contract tests assert that inter-application messages conform to a shared understanding that is documented in a contract."*
Source: Pact Docs — Introduction
URL: https://docs.pact.io/
Confidence: HIGH
Corroborated By: Martin Fowler (Source 2): *"check that all the calls against your test doubles return the same results as a call to the external service would"*
Notes: Definisi ini menjadi ground cover seluruh laporan.

## Evidence 2

Claim: Contract = HTTP method, path, status code, headers, serta request/response body (field, type). Bukan hanya field JSON.
Evidence: Contract-nya *"could meliputi HTTP Method, Path, Status, Content-Type, Response schema (id:integer, name:string, phone:string|null), termasuk behavior error (404 vs 200 dengan data:null)."*
Source: Topic spec (laboratorium) — disampaikan sebagai contoh ilustratif, konsisten dengan definisi di atas.
URL: NOT APPLICABLE (topic spec lokal)
Confidence: MEDIUM — ini contoh ilustratif, bukan kutipan primer. Namun selaras dengan Source 1 & 2.
Notes: Diperlukan untuk analisis tiga perubahan pada lab.

## Evidence 3

Claim: Consumer-Driven Contract (CDC) — consumer mendefinisikan expectation (contract), provider menjalankannya di CI untuk verifikasi sebelum deploy.
Evidence: *"[Pact] is a code-first consumer-driven contract testing tool... The contract is generated during the execution of the automated consumer tests. A major advantage... only parts of the communication that are actually used by the consumer(s) get tested."* Topic spec mengilustrasikan alur: build → run provider tests → apakah contract mobile masih terpenuhi? YES → deploy, NO → block deployment.
Source: Pact Docs (Source 1); alur CI: Topic spec (lokal)
URL: https://docs.pact.io/
Confidence: HIGH
Corroborated By: Martin Fowler (Source 3): pattern Consumer-Driven Contracts — provider berhak atas obligasi yang berasal dari luar boundary-nya
Notes: Bagian "block deployment" adalah sintesis antara Pactflow (Source 5: "breaking changes should not be able to make it into a production release") dan topic spec.

## Evidence 4

Claim: Schema/JSON-spec testing (OpenAPI, JSON Schema) ≠ contract testing. Schema hanya menguji kompatibilitas satu sistem terhadap satu skema pada titik waktu — tidak menjamin kedua sistem saling compatible dan tidak mendokumentasikan interaksi/conversations/evolution.
Evidence: *"1. Schema test — asserts that a single system is compatible with a schema... 2. Contract test — asserts that two systems are able to communicate by agreeing on what interactions can be sent between them and providing concrete examples... Contract testing goes beyond schema testing..."* Plus *"Schemas are abstract, and introduce ambiguity... easy to check if a system is compatible with a schema, but it's very difficult to be sure it fully implements the spec."*
Source: Pactflow — Schema vs Contract (Part 1)
URL: https://pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1
Confidence: HIGH
Corroborated By: Pact Docs Source 1: *"Unlike a schema or specification (eg. OAS), which is a static artefact that describes all possible states of a resource, a Pact contract is enforced by executing a collection of test cases... Pact is, in effect, 'contract by example'."*

## Evidence 5

Claim: Contract test vs E2E — contract test lebih cepat, mudah dipelihara, dapat di-debug, repeatable, scale naik, menemukan bug secara lokal sebelum push. E2E lambat, flaky, hard to maintain, scale buruk, temukan bug terlambat.
Evidence: Pactflow blog mencantumkan table properti lawanan (fast vs slow, easier to maintain vs hard to maintain, easy to debug vs hard to fix, repeatable, scale, uncover bugs locally).
Source: Pactflow — What is contract testing
URL: https://pactflow.io/blog/what-is-contract-testing/
Confidence: HIGH

## Evidence 6

Claim: Contract test ditempatkan di "Service Tests layer" Test Pyramid (di atas unit test, di bawah E2E). Bukan semua E2E diganti — hanya yang perlu diverifikasi integrasi.
Evidence: *"Contract tests fit in the Service Tests layer, as they execute quickly and don't need to integrate to external systems to run."* + rekomendasi: "Keep end-to-end integrated tests to a minimum."
Source: Pactflow — What is contract testing
URL: https://pactflow.io/blog/what-is-contract-testing/
Confidence: HIGH

## Evidence 7

Claim: Provider dapat "drift" terhadap schema/dokumen. Dokumentasi (mis OpenAPI) yang bagus tidak otomatis menjamin implementasi melakukan hal yang sama — hanya schema test tidak cukup untuk consumer.
Evidence: *"A good way to deal with [test double accuracy]... is to continue to run your own tests against the double, but in addition to periodically run a separate set of contract tests."* + *"no chance of implementation drift, as real application code is executed"* (code-based pro) vs *"Schemas are abstract, and introduce ambiguity... Code <-> schema drift"* (schema con).
Source: Martin Fowler (Source 2) + Pactflow Schema vs Contract (Source 6)
URL: https://martinfowler.com/bliki/ContractTest.html ; https://pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1
Confidence: HIGH

## Evidence 8

Claim: Breaking change dapat terjadi pada tipe data sekalipun response tetap JSON valid. Contoh: integer → string, name → full_name, status casing (IN_PROGRESS → in_progress).
Evidence: Topic spec + Lab 06 README (Source 7): *"Primitive type berubah (price:number → price:string)"*, *"Rename field (name → full_name)"*, *"Enum semantics berubah"* semuanya termasuk breaking change pada published contract.
Source: Lab 06 README (local) + Topic spec
URL: /labs/06-api-versioning/README.md
Confidence: HIGH (konsisten tiap sumber)

## Evidence 9

Claim: Additive change (tambah field optional di response) biasanya backward-compatible asalkan consumer toleran terhadap unknown field.
Evidence: *"Tambah field optional di response → consumer toleran unknown field"* + *"json.Unmarshal Go: abaikan unknown field → aman"* + *"Consumer contracts are open and incomplete... only parts actually used by consumer(s) get tested"* — jadi field tidak dipakai consumer tidak dijadikan contract.
Source: Lab 06 README (Source 7) + Pact Docs (Source 1) + Martin Fowler CDC (Source 3)
URL: /labs/06-api-versioning/README.md ; https://docs.pact.io/
Confidence: HIGH (catatan: bergantung pada consumer — ketat vs toleran)

## Evidence 10

Claim: Contract tidak harus berisi seluruh response — hanya yang dibutuhkan consumer. Snapshot/response penuh membuat test rapuh dan gagal karena perubahan tidak relevan.
Evidence: *"Pact is a code-first consumer-driven contract testing tool... only parts of the communication that are actually used by the consumer(s) get tested."* + *"Consumer-driven contracts... see exactly which fields each consumer is interested in, allowing unused fields to be removed and new fields to be added... without impacting a consumer."*
Source: Pact Docs (Source 1) + Pactflow (Source 4)
URL: https://docs.pact.io/
Confidence: HIGH

## Evidence 11

Claim: Contract testing berlaku juga untuk event/pesan (message queues/Kafka/RabbitMQ/webhooks), bukan hanya REST/HTTP.
Evidence: *"for an application that used queues, this would be the message that goes on the queue"* + *"For a system that uses message queues, this would involve checking that the provider generates the expected message."* + *"integer → string di payload amount bisa rusak consumer."*
Source: Pact Docs (Source 1) + Topic spec
URL: https://docs.pact.io/
Confidence: HIGH

## Evidence 12

Claim: Versi (versioning) tidak selalu diperlukan — hanya untuk major/breaking changes. Additive change pertahankan versi yang sama.
Evidence: Lab 06 README decision rule: *"Backward-compatible change → pertahankan version yang sama / Contract-breaking change → gunakan major API version baru."*
Source: Lab 06 README (Source 7)
URL: /labs/06-api-versioning/README.md
Confidence: HIGH

## Evidence 13

Claim: Strategi evolusi API ketika breaking change harus dilakukan: (a) versioning V1→V2 paralel, (b) compatibility/migration layer, (c) deprecation cycle + consumer communication + traffic monitoring + sunset criteria.
Evidence: *"Dual Contracts: GET /api/v1 → customer: string / GET /api/v2 → customer: object"* + *"Domain Model ≠ public API contract. DTO terpisah."* + deprecation lifecycle diagram — V1 active → V2 released → V1 deprecated → consumer communication → migration monitoring → V1 traffic monitoring → Sunset criteria → V1 removed.
Source: Lab 06 README (Source 7)
URL: /labs/06-api-versioning/README.md
Confidence: HIGH
Notes: Consumer inventory sebelum breaking change = syarat.

## Evidence 14

Claim: Provider states / test data adalah complexity trade-off pada code-based contract testing; schema-based justru menghilangkan "problem of test data."
Evidence: Schema-based pros: *"Removes the problem of 'test data' - in Pact, we solve this using provider states..."* Code-based cons: kebutuhan provider states, rumit saat banyak consumer.
Source: Pactflow — Schema vs Contract (Part 1)
URL: https://pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1
Confidence: MEDIUM — ini trade-off, bukan fakta absolut.

## Evidence 15

Claim: Provider-driven contract test (schema test provider-only, mis OpenAPI) tidak cukup untuk mencegah integration bugs — tidak memberi kepastian consumer memanggil provider dengan benar.
Evidence: Pact Docs: *"The term 'contract testing'... sometimes used in the context of a standalone provider application... On its own, however, it does not provide any test based assurance that the consumers are calling the provider in the correct manner, or that the provider can meet all its consumers' expectations, and hence, it is not as effective in preventing integration bugs."*
Source: Pact Docs (Source 1)
URL: https://docs.pact.io/
Confidence: HIGH
