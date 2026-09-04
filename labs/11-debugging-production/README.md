# Lab 11 — Debugging Sistem Production: Hypothesis-Driven Debugging

> Debugging bukan mencari error secara acak. Debugging adalah proses membuktikan atau membantah hipotesis berdasarkan data.

## Problem

Aplikasi memiliki satu endpoint (`GET /purchase-orders?range=1y`) yang mengalami masalah performa (lambat / timeout). 
Sebuah `504 Gateway Timeout` di load balancer adalah *symptom* (gejala), namun *root cause* (akar masalah) terletak pada bottleneck di salah satu layer aplikasi.

Lab ini menunjukkan proses menemukan root cause menggunakan *hypothesis-driven debugging*, memisahkan fakta dan asumsi, serta menguji hipotesis dengan *evidence* (waktu eksekusi tiap layer) tanpa mengandalkan *observability stack* yang kompleks.

## Reproduce

Jalankan server dan akses endpoint.

```bash
go test -v -run TestReproduceBug
```

Perhatikan bahwa response menjadi sangat lambat (> 4 detik) pada query dengan rentang waktu besar.

## Known Facts

- Endpoint `GET /purchase-orders?range=1m` cepat (< 20ms).
- Endpoint `GET /purchase-orders?range=1y` lambat (> 4000ms).
- Flow aplikasi: `Handler` -> `Service` -> `Repository`.

## Hypotheses

Sebelum mengumpulkan bukti lebih lanjut, beberapa kemungkinan hipotesis:

1. **H1**: Bottleneck ada di `Handler` (misal: JSON marshaling sangat lambat karena payload besar).
2. **H2**: Bottleneck ada di `Service` (misal: business logic lambat).
3. **H3**: Bottleneck ada di `Repository` (misal: query database melambat karena rentang waktu besar).

## Collect Evidence

Tambahkan measurement sederhana menggunakan `time.Since()` di tiap layer.

Jalankan test pengumpulan evidence:

```bash
go test -v -run TestCollectEvidence
```

Contoh output:
```text
request_id=abc123
handler_duration=5ms
service_duration=12ms
repository_duration=4200ms
total_duration=4220ms
```

## Test Hypotheses

Menggunakan data dari log:
- `handler_duration = 5ms` -> **H1 rejected** (Handler cepat)
- `service_duration = 12ms` -> **H2 rejected** (Service cepat)
- `repository_duration = 4200ms` -> **H3 supported** (Repository sangat lambat)

## Root Cause

Berdasarkan *evidence*, bottleneck berada di layer `Repository`. 
Masalah bukan di network atau parsing data, melainkan data access yang sengaja dibuat lambat (simulasi missing index / inefficient query).

## Fix

Lakukan optimasi terisolasi pada `Repository` dengan mengubah satu variabel saja (One Variable at a Time). 

Ubah konstanta lambat menjadi cepat di implementasi fix.

## Verification

Verifikasi perbaikan dengan memastikan endpoint yang sama kini berjalan cepat.

```bash
go test -v -run TestVerifyFix
```

## What We Learned

1. **Memisahkan fakta dan asumsi**: Masalah lambat tidak langsung berarti "database error" atau "network down".
2. **Membuat hipotesis yang dapat diuji**: "Mungkin layer X yang lambat".
3. **Mengumpulkan evidence (timing/log)**: Menggunakan `time.Since` untuk mengukur durasi setiap layer.
4. **Mengeliminasi hipotesis**: Menyisihkan H1 dan H2 karena waktu eksekusi singkat.
5. **Mengubah satu hal dalam satu waktu**: Hanya memperbaiki `Repository` untuk melihat apakah metrik berubah.
6. **Memverifikasi fix**: Memastikan bahwa root cause yang telah diperbaiki menyelesaikan masalah performa secara nyata.
