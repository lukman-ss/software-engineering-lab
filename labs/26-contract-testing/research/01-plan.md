# Research Plan — Contract Testing

## Research Topic

**Contract Testing — API Bisa Sama-Sama "Lulus Test", Tapi Integrasi Tetap Rusak**

Topic: Contract testing, consumer-driven contracts (CDC), API integration, breaking changes, schema evolution, backward compatibility.

## Objective

Investigate contract testing as a technique for catching integration-breaking API changes in distributed systems. Determine how CDC (consumer-driven contracts) works, how it differs from schema/OpenAPI testing, how it should be wired into CI/CD, and how it compares to end-to-end testing. Apply findings to the lab scenario (three proposed provider changes: enum-casing, rename field, type change).

## Research Questions

1. Apa definisi teknis contract testing dan CDC (Consumer-Driven Contracts)?
2. Bagaimana Pact menerapkan CDC secara code-first, dan peran Pact Broker/PactFlow?
3. Contract testing vs. schema (OpenAPI/JSON Schema) testing — apa perbedaan kunci dan trade-off?
4. Bagaimana contract test dipasangkan ke CI/CD supaya breaking change diketahui sebelum deploy provider?
5. Contract testing untuk pesan/event (message queues) — apakah konsepnya sama?
6. Additive change vs. breaking change — batasan yang dapat dipercaya?
7. Lab scenario: ketiga perubahan backend (status, nama field, tipe) — masing-masing breaking? Bagaimana evolusi minimal?

## Search Strategy

- Primary: Pact official docs (docs.pact.io)
- Primary: Martin Fowler bliki "Contract Test" + article "Consumer-Driven Contracts"
- Primary/Secondary: Pactflow blog posts on contract testing, CDC, schema vs contract
- Cross-check definitions across Pact + Fowler + Pactflow

## Expected Primary Sources

- docs.pact.io (Introduction, How Pact works, Conceptual overview, Testing scope)
- martinfowler.com/bliki/ContractTest.html
- martinfowler.com/articles/consumerDrivenContracts.html
- pactflow.io/blog/what-is-contract-testing
- pactflow.io/blog/contract-testing-using-json-schemas-and-open-api-part-1

## Risks / Unknowns

- Beberapa sumber Pact terbaru (2025) mencantumkan update 2026 — verifikasi tanggal.
- Batasan "additive = aman" bergantung pada consumer yang tolerant/strict.
- Perbedaan definisi "contract testing" provider-only vs integration contract testing.
- Lab menggunakan Go dan REST; sumber banyak berbahasa Kotlin/JS — harus generalisasi.
