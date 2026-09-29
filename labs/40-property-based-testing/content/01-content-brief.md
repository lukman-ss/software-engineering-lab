# Content Brief: Property-Based Testing (PBT)

Topic: Property-Based Testing (PBT) vs. Example-Based Testing
Target Reader: Software Engineer, Backend Engineer, QA Automation Engineer yang terbiasa dengan unit test klasik (example-based) dan ingin menangkap bug laten/edge-case secara otomatis.
Problem: Example-based testing (unit test biasa) hanya memverifikasi skenario yang dipikirkan oleh programmer. Boundary values tersembunyi, floating-point precision error, data tidak terurut, serta kombinasi input ekstrim sering lolos ke produksi karena ruang input tak hingga tidak teruji.
Core Mental Model: Alih-alih mencocokkan input spesifik dengan output spesifik `f(x) == y`, programmer mendefinisikan *universal invariant* yang harus selalu berlaku `∀x ∈ Domain: Property(f(x)) == true`. Mesin PBT mengacak ratusan/ribuan input, menemukan pelanggaran, lalu melakukan *shrinking* ke counterexample terkecil.
Approved Research Status: APPROVED (research-audit/07-verdict.md)
Approved Engineering Status: APPROVED (engineering-audit/06-verdict.md)

Main Concepts:
1. Universal Invariants vs. Example Scenarios.
2. Canonical Property Patterns:
   - Roundtrip: `Decode(Encode(x)) == x`
   - Idempotence: `f(f(x)) == f(x)`
   - Equivalence / Oracle: `f_optimized(x) == f_reference(x)`
   - Hard-to-prove, Easy-to-verify.
3. Counterexample Shrinking (Binary chunk removal & element reduction).
4. Biased Random Generation (zero, negative, boundaries).
5. Deterministic Reproducibility via Seeds.

Verified Behaviors:
1. Naive float currency lolos 5/5 example tests tetapi gagal 1000/1000 pada roundtrip PBT akibat floating-point precision loss.
2. Robust int64 cents currency lolos 1000/1000 roundtrip iterations pada rentang negatif, nol, hingga multi-miliar cent.
3. Naive interval merge lolos pada data terurut tetapi gagal oracle check (discrepancy 85/100) saat input acak tidak terurut. Robust merge lolos idempotence dan non-overlapping invariant 1000/1000.
4. Shrinker algoritma mereduksi slice acak 10 elemen yang gagal `[137 18 100 -27 95 107 126 78 -7 79]` dalam 21 langkah menjadi minimal counterexample 1 elemen `[-1]`.

Available Case Studies:
- Moneter / Currency Parser & Formatter: float64 precision leak vs int64 cents roundtrip.
- Interval Scheduling & Merging: naive sorted assumption vs robust sorted-merge idempotence.
- Counterexample Shrinker: step-by-step reduction dari 10 elemen menjadi `[-1]`.

Warnings:
1. PBT bukan pengganti unit test atau integrasi test spesifik domain; keduanya komplementer.
2. Generator testing/quick Go berbasis reflection default 100 iterasi; custom biased generator diperlukan untuk eksplorasi boundary.
3. Shrinking dapat terjebak dalam local minima jika generator/predikat memiliki kondisi multi-variabel kompleks.
4. Implementasi lab bersifat self-contained Go standard library (`testing/quick`), bukan benchmark performa framework pihak ketiga (Hypothesis/fast-check).
