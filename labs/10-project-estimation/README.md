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

Feature besar dipecah menjadi task-task kecil yang konkret dan dapat diukur, bukan dihitung per halaman atau dianggap satu task atomik.

Misalnya, **"Payment Gateway"** bukan satu task atomik ("Payment Gateway: 5 hari"). Seorang Senior Engineer memecahnya menjadi:

- Memahami dokumentasi API vendor
- Authentication mechanism (API key, OAuth, signature)
- Create transaction / payment request flow
- Payment callback & webhook handling
- Error handling & edge cases (expired, duplicate payment)
- Retry behavior & idempotency
- Testing di environment sandbox

### Breakdown Fitur Aplikasi Booking Servis

| No | Fitur | Task Konkret |
|----|-------|--------------|
| 1 | Login | Auth endpoint, password hashing/token, validasi input, unit test |
| 2 | Booking Online | DB schema booking, API create & list booking, business validation (slot availability), unit test |
| 3 | Pilih Cabang | Master data cabang API, selector UI/endpoint, caching |
| 4 | Pilih Mekanik | Relasi cabang-mekanik, availability check, filter & list endpoint |
| 5 | Payment Gateway | **Spike eksplorasi**, create transaction API, webhook receiver, retry/idempotency, sandbox testing |
| 6 | WhatsApp Notification | Template message, vendor API client, queue/async worker, error retry |
| 7 | Dashboard Admin | Metrics query/aggregator, summary endpoint, UI dashboard |
| 8 | Laporan Excel | Query data filter, Excel generator stream/library, download endpoint |

---

## 4. Known vs Unknown

Membedakan apa yang sudah kita pahami dengan pasti (Known) dan apa yang belum jelas (Unknown):

| Kategori | Karakteristik | Contoh di Kasus Ini | Perlakuan |
|----------|---------------|---------------------|-----------|
| **Known** | Pola & tech stack sudah familiar, tim pernah mengerjakan | Login, Dashboard sederhana, Export Excel | Langsung buat estimasi range (Min/Most Likely/Max) |
| **Unknown** | Vendor API baru, belum jelas flow/webhook/auth-nya | API Payment Gateway baru, behavior webhook vendor | **Wajib Spike** terlebih dahulu |

> **PENTING:**
> - `unknown != 0 hari` (mengabaikan unknown akan meledak di tengah proyek)
> - `unknown != langsung ditebak` (menebak tanpa dasar menghasilkan angka palsu)
> 
> Unknown harus diubah menjadi Known melalui **Spike** sebelum dapat diestimasi dengan akurat.

---

## 5. Risk & Dependency

Setiap task memiliki tingkat risiko teknis dan ketergantungan eksternal:

| Risk Level | Kriteria | Contoh |
|------------|----------|--------|
| **Low** | Well-known, proven tech, minim dependensi | CRUD cabang, login standar |
| **Medium** | Kompleksitas menengah atau ada integrasi internal | Availability slot booking, aggregasi dashboard |
| **High** | Core engine, perubahan skema besar, high cost of failure | Transaksi settlement, logic penugasan mekanik |
| **Unknown** | Belum pernah digunakan, dependensi eksternal pihak ketiga | Vendor Payment Gateway, WhatsApp API provider |

### External Dependency
- **Payment Gateway Vendor**: butuh akun sandbox, API credentials, webhook URL publik.
- **WhatsApp Provider**: butuh registrasi template WhatsApp, API token, quota pesan.
- **Dependency Tim Internal**: desain UI/UX final, requirement spesifik dari product owner.

### Technical & Project Risks
- Dokumentasi payment gateway tidak sesuai dengan kondisi aktual di sandbox
- Flow webhook membutuhkan perubahan arsitektur (misal: butuh reverse proxy / webhook queue)
- Requirement booking dan alokasi slot mekanik berubah di tengah jalan
- Keterlambatan approval atau penerbitan API credential dari vendor pihak ketiga

---

## 6. Spike

Spike adalah eksplorasi teknis yang **di-timebox** (dibatasi waktu secara ketat) untuk menjawab ketidakpastian spesifik.

### Contoh Spike: Payment Gateway
- **Timebox:** Maksimal 1 hari (atau 1–2 hari)
- **Tujuan Spike:**
  1. Cek authentication mechanism (signature / token).
  2. Coba flow create transaction di sandbox.
  3. Pahami payload & verifikasi signature pada webhook/callback.
  4. Cek ketersediaan dan reliabilitas sandbox vendor.
  5. Identifikasi error handling & retry behavior saat payment timeout/gagal.
- **Output Spike:**
  Informasi yang cukup untuk mengubah **Unknown** menjadi **Estimatable Work** (bukan membangun integrasi production siap rilis). Jangan membuat mock infrastructure atau abstraksi yang tidak diperlukan.

```
Payment Gateway (Unknown)
   ↓
Spike 1 hari (eksplorasi timeboxed)
   ↓
Task konkret teridentifikasi:
- Auth & create transaction: 1–2 hari
- Webhook receiver & signature: 1–2 hari
- Error retry & sandbox testing: 1–2 hari
Total Implementation Effort: 3–6 hari (Known)
```

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

## 8. Effort vs Duration (Engineer-days vs Calendar Days)

Bedakan **Effort** (engineer-days, kerja sekumulatif sampai selesai) dengan **Duration** (working days pada kalender, tergantung availability).

Contoh perhitungan:

```
Implementation Effort: 18–24 engineer-days
Spike: 2.5 hari (termasuk dalam effort)
Base Effort: 20.5–26.5 engineer-days
Contingency: 15%
Final Effort: 22–30 engineer-days

1 engineer @ 70% availability (realistis)
↓
Calendar Duration = Final Effort / (1 × 0.7)
↓
Duration: 31–43 working days (~6–9 weeks)
```

> **Jangan menyamakan effort (jumlah jam kerja) dengan duration (hari kalender).**

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

Kerjakan estimasi untuk case study "Aplikasi Booking Servis". Jangan melihat solusi di atas sebelum Anda mengerjakannya.

1. **Breakdown** requirement (misal "Payment Gateway" dipecah: auth, create transaction, webhook, sandbox testing). Jangan gunakan jumlah halaman.
2. Tandai **Known vs Unknown**.
3. Identifikasi **technical risk** dan **external dependency**.
4. Tentukan task mana butuh **Spike** (timeboxed, misal 1-2 hari) + tujuan spike-nya (misal: "mempelajari behavior webhook").
5. Buat **estimation range** (Min/Most Likely/Max) untuk tiap task.
6. Hitung **Effort vs Calendar Duration** (akunkan availability engineer).
7. Tambahkan **contingency** secara masuk akal (buffer risiko, bukan rumus wajib).
8. Tuliskan **minimal 3 assumptions** yang mendasari estimasi.
9. Tuliskan **minimal 3 risks** yang bisa membuat estimasi berubah.
10. Susun **Final Stakeholder Estimate** yang merangkum hasil.

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
