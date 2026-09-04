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

## 3. Naive Estimation (Mindset Junior)

Contoh pendekatan naif yang sering gagal:

```text
Hitung halaman: 8 fitur / 8 halaman
Asumsi: 1 hari per halaman
Estimasi total: 8 hari (pasti)
```

**Kelemahan pendekatan naif:**
- Mengabaikan integrasi eksternal (Payment Gateway & WhatsApp)
- Mengasumsikan 100% waktu habis untuk coding tanpa code review / meeting / testing
- Mengabaikan risiko ketidaksesuaian dokumentasi API pihak ketiga
- Tidak memiliki buffer (contingency) jika terjadi kendala tak terduga

---

## 4. Task Breakdown

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

## 5. Known vs Unknown

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

## 6. Risk & Dependency

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

## 7. Spike

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

## 8. Estimate Range (Uncertainty)

Estimasi harus merepresentasikan **ketidakpastian (uncertainty)** melalui *range*. 

Hindari memberikan angka tunggal yang seolah-olah pasti:
- ❌ Payment Gateway = 3 hari pasti.
- ✅ Payment Gateway = 2–4 hari.
- ✅ Payment Gateway = belum dapat diestimasi (butuh spike 1 hari).

**Penting:** Angka estimasi sangat bergantung pada:
- Requirement yang ada
- Pengalaman engineer yang mengerjakan
- Kondisi codebase existing
- Kesiapan environment
- Proses review & testing
- Dependensi eksternal

```go
type EstimateRange struct {
    Min        float64  // Best case
    MostLikely float64  // Mode
    Max        float64  // Worst case
}
```

### Kalkulasi Sederhana Range (Bukan Jaminan)

```
Expected = (Min + 4 × MostLikely + Max) / 6
```

---

## 8. Effort vs Duration

Bedakan **Effort** (total jam/hari kerja aktual) dengan **Duration** (durasi di kalender).

Misalnya:
- **Effort:** 10–14 engineer-days
- **Duration:** sekitar 3 minggu kalender

Mengapa berbeda? Karena dalam hari kerja selalu ada:
- Code review
- Meeting
- UAT
- Deployment
- Menunggu dependency
- Context switching

Contoh perhitungan Duration dengan asumsi availability engineer 70%:

```
Final Effort: 22–30 engineer-days
1 engineer @ 70% availability (0.7)
↓
Calendar Duration = Final Effort / 0.7
↓
Duration: 31–43 working days (~6–9 minggu kalender)
```

---

## 9. Contingency (Bukan Rumus Wajib)

Contingency adalah buffer yang ditambahkan berdasarkan *risk profile*. Ini harus dikaitkan dengan ketidakpastian yang nyata, **bukan sekadar "tambah 20% tanpa alasan"**.

| Risk Level | Contingency (contoh, bukan aturan wajib) |
|------------|------------------------------------------|
| High / Unknown | 20–25% (karena banyak dependensi API luar) |
| Medium | 10–15% |
| Low | 5–10% |

---

## 10. Assumptions & Risk

Final estimate **wajib** mencantumkan asumsi (kondisi ideal yang diharapkan) dan risiko (apa yang bisa membuat estimasi meleset).

**Contoh Assumptions:**
- Desain UI final sudah tersedia.
- Lingkungan sandbox payment gateway siap dipakai.
- Credential vendor tersedia sejak hari H.
- Dikerjakan oleh 1 engineer tanpa interruption.
- Tidak ada perubahan requirement besar di tengah jalan.
- API eksternal berjalan persis seperti dokumentasi.

**Contoh Risk (Risiko yang memengaruhi estimasi):**
- Dokumentasi payment gateway tidak akurat dengan kondisi sandbox.
- Webhook vendor memiliki *behavior* yang tidak terduga.
- Requirement booking dari stakeholder berubah.

**Note:** Bedakan Risk dan Dependency.
- **Dependency:** Credential payment gateway, API vendor, desain UI.
- **Risk:** Dokumentasi vendor tidak akurat, behavior webhook tidak standar.

---

## 11. Final Estimate Communication

Senior engineer tidak sekadar memberikan satu angka. Mereka mengkomunikasikan **range + assumptions + risks + unknowns**.

Contoh cara mengkomunikasikan hasil estimasi:

> "Berdasarkan requirement saat ini, estimasi waktu penyelesaian (Duration) adalah sekitar **3–4 minggu kalender** (15-23 hari kerja).
> 
> Estimasi ini menggunakan **asumsi** bahwa desain UI sudah final, credential sandbox vendor sudah tersedia, dan tidak ada perubahan requirement besar di tengah jalan.
> 
> **Risiko terbesar** ada di integrasi payment gateway dan webhook-nya. Bagian ini belum bisa kami estimasi (unknown), sehingga **membutuhkan Spike selama 1 hari** terlebih dahulu. Setelah spike selesai, kami bisa memberikan estimasi implementasi pastinya untuk bagian tersebut."

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

Kerjakan estimasi untuk case study "Aplikasi Booking Servis". Jangan melihat contoh acuan sebelum Anda mengerjakannya secara mandiri.

### Instruksi:
1. **Break down** fitur Aplikasi Booking Servis menjadi task-task kecil yang konkret.
2. Tandai mana task **Known** dan mana yang **Unknown**.
3. Identifikasi **external dependency** (vendor credentials, API sandbox, dll).
4. Identifikasi **high-risk task**.
5. Pilih task yang membutuhkan **Spike** (timeboxed), tentukan durasi dan tujuannya.
6. Buat **effort estimate** dalam bentuk range (Min/Most Likely/Max).
7. Hitung dan bedakan antara **effort** (engineer-days) dan **calendar duration** (working days/weeks).
8. Tentukan **contingency** yang wajar berdasarkan risiko yang teridentifikasi.
9. Tuliskan **assumptions** yang mendasari estimasi Anda.
10. Tuliskan **minimal 3 risks** yang dapat memengaruhi estimasi.
11. Buat **final estimate** yang siap dikomunikasikan ke stakeholder.

---

## Reference Example

*Bagian ini dapat digunakan sebagai pembanding setelah Anda selesai mengerjakan Exercise.*

<details>
<summary>Lihat Reference Example (Solusi Acuan)</summary>

```text
Task Breakdown & Estimation Range (Known):
- Login: 1–2 hari (Low Risk)
- Booking Online: 2–4 hari (Medium Risk)
- Pilih Cabang: 1 hari (Low Risk)
- Pilih Mekanik: 1–2 hari (Medium Risk)
- WhatsApp Notification: 1–2 hari (Unknown -> Spike 0.5 hari)
- Dashboard Admin: 2–3 hari (Medium Risk)
- Laporan Excel: 1–2 hari (Medium Risk)
- Testing / UAT: 2–3 hari

Unknown / Spike:
- Payment Gateway API: Spike 1 hari (eksplorasi auth, sandbox, webhook)

Known Implementation Effort: 11–19 engineer-days
Spike Effort: 1.5 engineer-days
Contingency (15% buffer): ~2–3 hari

Total Effort: ~15–23 engineer-days

Calendar Duration (1 engineer @ 70% availability):
15 / 0.7 = ~21 hari kerja (sekitar 3–4 minggu kalender)

Assumptions:
- UI design sudah final
- Credentials & sandbox Payment Gateway sudah siap
- 1 engineer dedicated dengan availability 70%
- API eksternal berjalan sesuai spesifikasi dokumentasi

Risks:
- Dokumentasi payment gateway tidak akurat dengan sandbox
- Behavior webhook vendor memerlukan perubahan arsitektur receiver
- Requirement alokasi slot booking berubah di tengah proyek
```

</details>

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
- **Next**: [Lab 11 — Debugging Sistem Production: Hypothesis-Driven Debugging](../11-debugging-production/)
