# Lab 10 — Project Estimation

**Estimasi Proyek Berbasis Breakdown, Risk, dan Uncertainty**

> Estimasi software bukan menebak lama coding. Estimasi adalah proses memecah scope, mengukur effort, dependency, risk, uncertainty, dan menghasilkan range yang dapat dijelaskan.

---

## 1. Problem

Junior engineer sering menganggap estimasi sebagai "berapa hari fitur selesai". Pendekatan ini:

1. **Single-point estimate** — hanya satu angka (misal: "5 hari")
2. **Tanpa risk analysis** — tidak mempertimbangkan ketidakpastian
3. **Tanpa breakdown** — langsung menebak seluruh fitur
4. **Sangat sensitive** terhadap requirement change

Hasil: Estimasi tidak pernah akurat, deadline sering terlewat, tim frustrasi.

---

## 2. Initial Requirement

Studi kasus yang dipakai di seluruh lab ini:

> **Aplikasi Booking Servis** — sistem untuk pelanggan melakukan booking servis kendaraan dengan memilih cabang, memilih mekanik, melakukan pembayaran online, menerima notifikasi WhatsApp, dan bagi admin tersedia dashboard serta laporan Excel.

Fitur yang harus diestimasi:

- Login
- Booking Online
- Pilih Cabang
- Pilih Mekanik
- Payment Gateway
- WhatsApp Notification
- Dashboard Admin
- Laporan Excel

---

## 3. Task Breakdown

Feature besar dipecah menjadi task kecil yang dapat diukur (bukan dihitung per halaman).

| No | Feature | Task kecil (contoh) |
|----|---------|---------------------|
| 1 | Login | Backend/API, Validation, Testing |
| 2 | Booking Online | Database, Backend/API, Validation, Testing |
| 3 | Pilih Cabang | Frontend, Backend/API |
| 4 | Pilih Mekanik | Frontend, Backend/API |
| 5 | Payment Gateway | Backend/API, External API, Webhook, Testing, **Spike** |
| 6 | WhatsApp Notification | Backend/API, External API, **Spike** |
| 7 | Dashboard Admin | Frontend, Backend/API, Testing |
| 8 | Laporan Excel | Backend/API, Database, Export |

Aktivitas lintas task yang biasa muncul:

```
Requirement Analysis
Database Schema
Backend/API
Frontend
Validation
Testing
Code Review
Deployment
UAT
```

---

## 4. Known vs Unknown

| Status | Arti | Perlakuan |
|--------|------|-----------|
| **Known** | Tech & pattern sudah familiar, effort bisa di-range | Langsung estimasi Min/MostLikely/Max |
| **Unknown** | Tech/API/vendor baru, belum jelas effort-nya | Wajib **Spike** dulu sebelum estimasi implementasi |

```go
task := Task{
    Name: "Payment Gateway Vendor API",
    Estimate: EstimateRange{Min: 0, MostLikely: 0, Max: 0},
    Risk: RiskUnknown,
    SpikeDays: 1.5, // Wajib untuk eksplorasi
    Assumptions: []string{"Vendor sandbox available"},
}
```

> **Unknown != 0**. Unknown berarti tidak cukup dipahami untuk dipercaya sebagai estimasi implementasi final.

---

## 5. Risk & Dependency

| Risk Level | Kriteria | Contoh |
|------------|----------|--------|
| **Low** | Well-known, proven tech | Tambah field, CRUD biasa |
| **Medium** | Complexity menengah, ada integrasi ringan | API endpoint baru, query kompleks |
| **High** | Dependency kunci, high cost of failure | Database migration besar, cache strategy |
| **Unknown** | Belum pernah dibuat, teknologi baru | Payment gateway baru, API vendor |

**External dependency** (di studi kasus ini):

- Payment gateway (vendor) → butuh credential, sandbox, dokumentasi.
- WhatsApp provider (vendor) → butuh API key, template, dokumentasi.
- Internal: tim UI/Design, requirement owner, environment Dev/Prod.

**Risiko yang dapat membuat estimasi berubah:**

- Dokumentasi payment gateway tidak sesuai kondisi aktual
- Flow webhook membutuhkan perubahan desain
- Requirement booking berubah
- Dependency eksternal terlambat

---

## 6. Spike

Spike = eksplorasi teknis timeboxed untuk mengubah **unknown** menjadi **estimatable work**.

### Contoh alur:

```
Payment Gateway Vendor API
↓
Unknown Risk (0 effort implementasi)
↓
Spike 1–2 hari
↓
Implementasi Min/MostLikely/Max dapat diestimasi
```

**Spike cukup berupa eksplorasi timeboxed untuk menjawab:**

- authentication mechanism
- create payment flow
- webhook flow
- retry/error behavior
- sandbox availability

**Jangan bangun integrasi payment gateway penuh di spike.** Output spike adalah informasi yang cukup untuk mengubah unknown menjadi estimatable work.

---

## 7. Estimate Range (Effort)

Gunakan range, bukan angka presisi palsu.

```go
type EstimateRange struct {
    Min        float64  // Best case
    MostLikely float64  // Mode
    Max        float64  // Worst case
}
```

### PERT Formula (alat bantu komunikasi, bukan jaminan)

```
Expected = (Min + 4 × MostLikely + Max) / 6
```

```go
func (r EstimateRange) Expected() float64 {
    return (r.Min + 4.0*r.MostLikely + r.Max) / 6.0
}
```

> Ini adalah **model estimasi**, bukan jaminan tanggal selesai.

---

## 8. Effort vs Duration

Bedakan effort (engineer-days) dengan durasi kalender (working days / weeks).

```
Effort: 18–24 engineer-days
+
Spike: 2.5 days
=
Base Effort: 20.5–26.5 engineer-days
+
Contingency: 15%
=
Final Effort: 22–30 engineer-days

1 engineer @ 70% availability
↓
Calendar Duration: 31–43 working days (~6–9 weeks)
```

```go
type DurationRange struct {
    MinDays      float64
    ExpectedDays float64
    MaxDays      float64
}

// Calendar Duration = Effort / (Engineers × Availability)
calendarDays = finalEffort / (engineerCount * availability)
```

> **Jangan menyamakan effort dengan tanggal selesai.**

---

## 9. Contingency

Buffer berbasis profil risiko (contoh, bukan rumus wajib).

| Risk Level | Contingency (contoh) |
|------------|----------------------|
| High / Unknown | 20–25% |
| Medium | 10–15% |
| Low | 5–10% |

`AutoContingency = true` → derive dari risk level tertinggi.
`AutoContingency = false` + `ContingencyRate = 0` → NO buffer (intentional).

---

## 10. Assumptions

Estimasi tanpa assumptions akan punya confidence rendah.

Contoh assumptions yang ditulis di project:

```go
project := Project{
    Assumptions: []string{
        "UI design sudah tersedia",
        "Sandbox payment gateway tersedia",
        "Credential API tersedia",
        "Requirement tidak berubah signifikan",
        "Tidak ada dependency eksternal yang terlambat",
    },
}
```

---

## 11. Final Estimate (Contoh Output)

Contoh final estimate untuk dikomunikasikan ke stakeholder (bukan standar universal):

```text
Task Breakdown
- Login: 1–2 hari
- Booking: 2–3 hari
- Pilih Cabang: 1 hari
- Pilih Mekanik: 1–2 hari
- Payment Gateway: perlu spike
- WhatsApp Notification: 1–2 hari
- Dashboard Admin: 2–3 hari
- Laporan Excel: 1 hari
- Testing/UAT: 2–3 hari

Known Effort:
11–17 hari

Payment Gateway Spike:
1 hari

Contingency:
10–20% (contoh berdasarkan profil risiko)

Estimated Duration:
sekitar 3–4 minggu

Assumptions:
- UI/design sudah tersedia
- Sandbox payment gateway tersedia
- Credential API tersedia
- Requirement tidak berubah signifikan
- Tidak ada dependency eksternal yang terlambat

Risks:
- Dokumentasi payment gateway tidak sesuai kondisi aktual
- Flow webhook membutuhkan perubahan desain
- Requirement booking berubah
- Dependency eksternal terlambat
```

> Angka-angka di atas adalah **contoh** berdasarkan asumsi tertentu. Mereka bukan standar universal.

---

## 12. Failure Scenarios (Validasi Model)

| No | Scenario | Expected |
|----|----------|----------|
| 1 | Invalid range: min > mostLikely | Error |
| 2 | Invalid range: mostLikely > max | Error |
| 3 | Negative estimate | Error |
| 4 | Negative spike days | Error |
| 5 | Availability = 0 | Error |
| 6 | Availability > 1 | Error |
| 7 | Negative contingency | Error |
| 8 | Empty project | Error |
| 9 | Invalid risk level | Error |
| 10 | Unknown risk tanpa SpikeDays > 0 | Error |
| 11 | EngineerCount = 0 | Error |
| 12 | EngineerCount < 0 | Error |
| 13 | Empty task name | Error |
| 14 | Unknown risk dengan SpikeDays = 0 tapi estimasi implementasi non-zero | Allowed (spike sudah selesai) |

---

## Cara Menjalankan Lab

```bash
cd labs/10-project-estimation
go test -v ./...
```

Test menjalankan skenario validasi di atas dan case study Aplikasi Booking Servis untuk mengilustrasikan alur estimasi.

---

## 13. Exercise

Kerjakan estimasi untuk case study "Aplikasi Booking Servis" (atau variasi requirement baru). Jangan melihat solusi di atas sebelum Anda mengerjakannya.

1. **Breakdown** requirement menjadi task kecil (hindari estimasi berbasis jumlah halaman).
2. Tandai **Known vs Unknown** untuk setiap task.
3. Identifikasi **technical risk** dan **external dependency**.
4. Tentukan task mana butuh **Spike** (timeboxed) + tujuan spike-nya.
5. Buat **estimation range** (Min/Most Likely/Max) untuk tiap task.
6. Hitung **Effort vs Calendar Duration** (akunkan availability engineer).
7. Tambahkan **contingency** secara masuk akal (berbasis risiko, bukan rumus wajib).
8. Tuliskan **minimal 3 assumptions** yang mendasari estimasi.
9. Tuliskan **minimal 3 risks** yang bisa membuat estimasi berubah.
10. Susun **Final Stakeholder Estimate** yang siap dipresentasikan.

---

## Senior Engineer Takeaways

1. **Estimasi bukan prediction, tapi range-based analysis** — selalu berikan Min / MostLikely / Max.
2. **Unknown != 0** — kalau belum jelas, spike dulu.
3. **Effort ≠ Calendar Duration** — akuntasi availability engineer.
4. **Risk profile memengaruhi contingency** — contoh: High/Unknown 20–25%, Medium 10–15%, Low 5–10%.
5. **Contingency adalah buffer risiko, bukan rumus wajib** — pilih yang masuk akal untuk konteks Anda.
6. **Range ordering** — Min ≤ Expected ≤ Max.
7. **Validation explicit** — input invalid → error, bukan silent failure.
8. **EngineerCount ≥ 1** — invalid planning input harus error.
9. **TaskName wajib** — setiap task butuh nama untuk dikomunikasi.
10. **Confidence adalah alat komunikasi** — bukan probabilitas statistik.
11. **Deterministic testing** — logika estimasi tidak butuh network/DB.

> **Estimasi ≠ menebak lama coding.**
> **Estimasi = breakdown + effort + uncertainty + risk + dependency + assumptions + communication.**

---

## Navigasi

- **Previous**: [Lab 09 — Code Review](../09-code-review/)
- **Next**: [Lab 11 — Pessimistic Locking](../11-pessimistic-locking/)
