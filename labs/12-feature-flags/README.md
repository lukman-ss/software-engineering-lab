# Lab 12 — Feature Flags

Lab ini mengajarkan bahwa **Deploy ≠ Release**. Feature Flag adalah mekanisme untuk mengendalikan risiko ketika merilis perubahan ke production, bukan sekadar if/else.

## Problem

Aplikasi bengkel (CMMS) akan merilis fitur baru: `online_booking`.
Masalahnya: Fitur tersebut masih baru dan belum diketahui apakah aman digunakan oleh seluruh user di production.

Tanpa feature flag:
Deploy → 100% User → Bug → Semua User Terdampak → Hotfix / Rollback / Deploy Lagi.

## Why Deploy != Release

- **Deploy**: Memindahkan kode aplikasi ke server production.
- **Release**: Membuat fitur tersedia (visible) untuk end user.

Kode boleh sudah berada di production, tetapi fiturnya tidak harus langsung tersedia untuk seluruh user. Feature flag memisahkan *deployment* dari *release*.

## Architecture / Flow

Gunakan `FeatureService` sebagai abstraction boundary. Business logic jangan tersebar dengan pengecekan flag di banyak tempat.

```go
if featureService.IsEnabled("online_booking", req.UserID) {
    return newBookingFlow(req)
}
return legacyBookingFlow(req)
```

## Running the Lab

Jalankan simulasi rollout dari CLI:

```bash
make lab-12-demo
# atau
go run ./labs/12-feature-flags/cmd/simulate/
```

## Feature OFF

Saat deployment pertama kali, fitur diset `OFF` (`enabled = false`).
Aplikasi tidak crash, melainkan semua request masuk ke fallback (Legacy Flow).

## Internal Rollout

Sebelum fitur dirilis ke publik, fitur diuji di production hanya untuk internal team.
User normal tetap masuk ke legacy flow.

## Percentage Rollout

Mengontrol eksposur fitur ke sebagian kecil user.
Contoh: 10% user menggunakan new flow, 90% user menggunakan legacy flow.

## Deterministic Bucketing

Percentage rollout tidak boleh menggunakan fungsi `random(1, 100)` per request.
Jika random: User yang sama dapat menerima `request 1 -> new`, `request 2 -> old`. Ini pengalaman yang buruk.

Gunakan deterministic bucketing:
`hash(user_id + feature_key) % 100`
Sehingga user yang sama selalu masuk ke bucket yang konsisten.

## 10% → 50% → 100%

Proses rollout dilakukan bertahap berdasarkan sinyal metrik.
Jika error rate normal pada 10%, naikkan ke 50%, lalu 100%.

## Kill Switch

Kondisi: Booking error rate meningkat!
Engineer merubah flag `online_booking = OFF`.
Semua request berikutnya langsung menggunakan legacy flow tanpa perlu redeploy aplikasi.
Digunakan untuk memperkecil blast radius ketika fitur baru menyebabkan incident.

## Metrics

Rollout keputusan harus berdasarkan data, bukan perasaan. Track metrik sederhana:
- `total_requests`
- `new_flow_requests` / `legacy_flow_requests`
- `booking_success` / `booking_failed` (Error Rate)

## Feature Flag vs Authorization

Feature flag bukan authorization system.
- **Authorization**: Apakah user BOLEH melakukan sesuatu? (`authorize(user, "booking:create")`)
- **Feature Flag**: Behavior versi fitur mana yang AKTIF untuk user ini? (`featureService.isEnabled("online_booking", user.id)`)

## Testing ON and OFF

Setiap feature flag menciptakan lebih dari satu execution path. Jalur ON dan jalur OFF (serta fallback) wajib diuji.

## Blast Radius

Blast radius adalah seberapa banyak sistem/user yang terdampak jika terjadi kegagalan.
- **100% rollout + bug** = 100% user potentially affected.
- **10% rollout + bug** = blast radius jauh lebih kecil (hanya 10% user terganggu).

Inilah alasan feature flag sangat penting bagi risk management.

## Feature Flag Lifecycle

Feature flag yang dibiarkan permanen menjadi technical debt. Siklus ideal:

1. Create Flag
2. Deploy OFF
3. Internal
4. 10%
5. 50%
6. 100%
7. Stable
8. Delete Flag (Cleanup)
9. Delete Legacy Code

Sebelum cleanup:
```go
if feature("online_booking") { newFlow() } else { legacyFlow() }
```
Sesudah cleanup:
```go
newFlow()
```

## Cleanup / Technical Debt

Jangan biarkan flag lama menumpuk. Setelah fitur mencapai 100% rollout dan terbukti stabil di production, hapus flag dan kode legacy yang tidak terpakai untuk mengurangi kompleksitas (`if/else` hell).

## Exercises

Jalankan perintah `make lab-12-demo` (yang menjalankan `cmd/simulate/main.go`) dan perhatikan outputnya pada setiap langkah berikut:

- **Exercise 1**: Deploy dengan `online_booking = OFF`. Pastikan 100% user menggunakan legacy flow.
- **Exercise 2**: Aktifkan untuk internal user. Amati bagaimana hanya user spesifik yang mendapat fitur.
- **Exercise 3**: Rollout 10%. Lihat distribusi traffic.
- **Exercise 4**: Naikkan menjadi 50%.
- **Exercise 5**: Naikkan menjadi 100%.
- **Exercise 6**: Simulasikan error pada new booking flow. Amati metrics yang error rate-nya melonjak. Kemudian sistem menggunakan *Kill Switch*.
- **Exercise 7**: Jelaskan, kenapa rollback feature menggunakan feature flag dapat lebih cepat dan lebih kecil risikonya dibanding emergency code deployment?

## Senior Engineer Takeaways

Pertanyaan Refleksi:
1. Kenapa feature flag tidak sama dengan environment variable biasa?
2. Kenapa percentage rollout tidak sebaiknya menggunakan random setiap request?
3. Kenapa rollout 10% dapat memperkecil blast radius?
4. Kapan engineer sebaiknya meningkatkan rollout dari 10% ke 50%?
5. Signal apa yang membuat rollout harus dihentikan?
6. Kapan kill switch harus digunakan?
7. Kenapa feature flag bukan authorization?
8. Apa risiko menyimpan feature flag terlalu lama?
9. Kenapa jalur ON dan OFF harus sama-sama diuji?
10. Setelah feature mencapai 100% dan stabil, apa yang harus dilakukan terhadap flag dan legacy implementation?

### Expected Learning Outcome

Memahami: **Develop ≠ Deploy ≠ Release**

**Success Workflow:**
Code Complete → Deploy → Feature OFF → Internal → 10% → Observe → 50% → Observe → 100% → Stable → Remove Flag

**Failure Path:**
10% Rollout → Metrics Memburuk → STOP ROLLOUT → Kill Switch → Legacy Flow → Investigate
