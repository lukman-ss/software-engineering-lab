# Expand → Migrate → Contract

## Mekanisme
Evidence:
1. **Expand**: Tambahkan schema/endpoint baru bersama dengan yang lama. Sistem mendukung dua versi.
2. **Migrate**: Update client dan data secara bertahap untuk memindahkan traffik dari legacy ke versi baru.
3. **Contract**: Hapus komponen legacy setelah traffic/usage mencapai angka 0.

Source: Parallel Change (Martin Fowler)
URL: https://martinfowler.com/bliki/ParallelChange.html
Confidence: HIGH
