# Open Questions: Saga Pattern Research

## Unanswered Questions

1. **Compensation-of-compensation protocol.** Jika Refund Payment gagal setelah Inventory reserve gagal, apa protokol recovery standar? Sumber menyebutkan retry + monitoring + human-in-the-loop, tetapi tidak ada spesifikasi kanonik.

2. **Quantitative threshold choreography vs orchestration.** Berapa jumlah service / branching factor di mana choreography menjadi unmaintainable? Semua sumber kualitatif.

3. **Garcia-Molina 1987 verbatim definition.** PDF scan tidak berhasil di-extract sebagai teks. Definisi original ("A saga is a long-lived transaction that can be written as a sequence of transactions that can be interleaved with other transactions") dikutip sekunder, belum diverifikasi word-for-word dari PDF.

4. **TCC (Try-Confirm-Cancel) vs Saga.** Hubungan TCC (Alibaba Seata, dsb.) dengan saga tidak diteliti di sesi ini. TCC sering disebut sebagai varian 2-phase tanpa lock lama; overlap dengan pivot/retryable Microsoft belum diklarifikasi.

5. **Camunda / Zeebe / Spring Statemachine / Axon Framework.** Implementasi nyata di luar AWS/Temporal/Eventuate tidak dibuka. Relevan jika lab ingin menunjukkan opsi tooling.

6. **Exactly-once end-to-end.** Klaim "exactly-once" di Step Functions Standard dan Temporal workflow level vs at-least-once di participant. Belum ada bukti empiris end-to-end exactly-once tanpa idempotency di participant.

## Weak Evidence

- **Pivot / retryable / compensable taxonomy:** hanya Microsoft Architecture Center yang merinci. Perlu cross-check ke *Microservices Patterns* Chapter 4 (Manning, paywalled) atau Garcia-Molina original.
- **Six isolation countermeasures:** enumerasi lengkap hanya di Microsoft; Richardson merujuk buku tanpa daftar publik.
- **Performance numbers 2PC vs Saga:** tidak ada dataset primer yang dibuka. Klaim "2PC terlalu lambat" diterima sebagai consensus engineering, bukan hasil benchmark sesi ini.

## Claims Needing Deeper Research

- Semantic lock implementation details (semaphore di application level, timeout, deadlock antar-saga).
- Ordering guarantees ketika choreography saga punya "happens-before" (Temporal menyatakan choreography buruk untuk ordered steps; perlu contoh kegagalan konkret).
- Outbox vs Event Sourcing sebagai backing store saga: kapan pilih yang mana.
- Timeout saga vs timeout langkah: bagaimana orchestrator memutuskan saga "stuck" vs langkah lambat.
- Saga dan GDPR/right-to-be-forgotten: event log retention vs compensation history.

## Possible Next Research Directions

1. Buka *Microservices Patterns* Chapter 4 (Richardson) untuk countermeasures dan Create Order Saga walkthrough lengkap.
2. Baca Seata TCC documentation untuk membedakan TCC vs orchestration saga.
3. Inspect Eventuate Tram Sagas example repo (github.com/eventuate-tram/eventuate-tram-sagas-examples-customers-and-orders) sebagai referensi implementasi lab.
4. Bandingkan Camunda BPMN compensation events vs code-first Temporal Saga class.
5. Kaji paper follow-up Garcia-Molina (sagas recovery, nested sagas) jika lab ingin historical depth.
6. Untuk lab exercise: rancang state machine Order (PENDING → PAID → RESERVED → CONFIRMED | COMPENSATING → CANCELLED) dan tabel idempotency key per participant.
