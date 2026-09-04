# Lab 12 — Feature Flags: Cara Merilis Fitur Baru Tanpa Takut Production Rusak

Lab ini mengajarkan bahwa **Deploy ≠ Release**. Feature Flag adalah mekanisme untuk mengendalikan risiko ketika merilis perubahan ke production, bukan sekadar if/else.

## Tujuan

Developer harus memahami:
- Feature OFF saat deployment.
- Internal rollout untuk testing di production.
- Percentage rollout secara deterministic.
- Kill switch untuk fallback cepat.

## Skenario: Booking Service Online

Fitur baru: `online_booking`.
Masalah: Belum tahu aman digunakan atau tidak.
Solusi: Deploy dengan flag OFF. Uji internal, lalu persentase, lalu 100%. Fallback menggunakan legacy flow sederhana.

## Implementasi Inti

- `FeatureService`: abstraction untuk evaluasi flag (hindari pengecekan tersebar).
- Deterministic bucketing: `hash(user_id + feature_key) % 100`.
- Safe fallback: fitur OFF != crash.

## Navigasi

- **Previous**: [Lab 11 — Lock Contention](../11-lock-contention/)
- **Next**: [Lab 13 — Deadlock](../13-deadlock/)
